# SaveAnyBot Plugins

SaveAnyBot 可透過插件擴充功能, 目前僅支援 Parser (解析器) 插件.

## Parser

解析器為 SaveAnyBot 提供了處理非 Telegram 檔案的能力, 例如下載其他網站的圖片或影片.

當前解析器介面定義如下:

```go
type Parser interface {
	CanHandle(url string) bool // 判斷是否能處理給定的 URL
	Parse(ctx context.Context, url string) (*Item, error) // 解析 URL, 返回 Item
}

// Resource is a single downloadable resource with metadata.
type Resource struct {
	URL       string            `json:"url"`
	Filename  string            `json:"filename"` // with ext
	MimeType  string            `json:"mime_type"`
	Extension string            `json:"extension"`
	Size      int64             `json:"size"`    // 0 when unknown
	Hash      map[string]string `json:"hash"`    // {"md5": "...", "sha256": "..."}
	Headers   map[string]string `json:"headers"` // HTTP headers when downloading
	Extra     map[string]any    `json:"extra"`
}

type Item struct {
	Site        string         `json:"site"`
	URL         string         `json:"url"` // original URL of the item
	Title       string         `json:"title"`
	Author      string         `json:"author"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags"`
	Resources   []Resource     `json:"resources"`
	Extra       map[string]any `json:"extra"`
}
```

### Write a Parser Plugin

解析器插件可使用 JavaScript 撰寫, SaveAnyBot 使用 [goja](https://github.com/dop251/goja) 提供執行環境, 並向其中注入了以下全域函式或物件:

- **registerParser**: 用於註冊解析器, 每個插件必須呼叫此函式以進行註冊
- **console.log**: 呼叫 go 端的 logger 印出日誌
- **ghttp**: 提供 HTTP 請求功能
- **playwright**: 提供基於 Playwright 的瀏覽器自動化請求功能

插件需要提供元資料 `metadata` 並實作 `canHandle` 和 `parse` 兩個函式, 最後呼叫 `registerParser` 註冊解析器.

#### Plugin Metadata

插件元資料是一個 JavaScript 物件:

```js
const metadata = {
    version: "1.0.0", // 插件相容版本號, 必須提供, 其他欄位可選
    name: "Example Parser", // 插件名稱
    description: "A parser for example links", // 插件描述
    author: "Krau", // 插件作者
}
```

#### canHandle Function

`canHandle`: `canHandle(url: string): boolean` , 用於判斷當前解析器能否解析給定的 URL, 返回布林值, 例如:

```js
const canHandle = function (url) {
	return url.includes("youtube.com/watch?v");
};
```

這將讓 SaveAnyBot 在遇到包含 `youtube.com/watch?v` 的 url 時呼叫當前解析器的 `parse`.

#### parse Function

`parse`: `parse(url: string): Item` , 是核心解析函式, 用於解析給定的 url, 返回一個 `Item` 物件, 例:

```js
const parse = function (url) {
    var result = {
        // 元資訊
        site: "YouTube",
        url: url,
        title: "測試 YouTube 影片",
        author: "某影片作者",
        description: "這是一個測試影片",
        tags: ["test", "youtube"],
        // 資源(可下載的檔案)清單
        resources: [
            {
                url: "https://example.com/video1.mp4", // 檔案直鏈
                filename: "somevideo.mp4", // 檔案名稱
                mime_type: "video/mp4", // 檔案 MIME 類型, 可選
                extension: "mp4", // 檔案副檔名, 可選
                size: 100 * 1024 * 1024, // 檔案大小, 單位為位元組, 未知可設為 0
                hash: {}, // 檔案雜湊, 可選, 格式為 {"md5": "xxx", "sha256": "xxx"} 等
                headers: {}, // 下載檔案時所需的 HTTP 標頭, 可選, 例如 {"User-Agent": "Mozilla/5.0"}
                extra: {} // 額外資訊, 可選, 可以包含任何自訂資料
            },
            {
                url: "https://example.com/picture1.png",
                filename: "picture1.png",
                mime_type: "image/png",
                extension: "png",
                size: 1 * 1024 * 1024,
                hash: {},
                headers: {},
                extra: {}
            }
        ],
        extra: {}
    };
    return result;
}
```

#### HTTP Requests

使用 `ghttp` 物件以發起 HTTP 請求.

**ghttp.get(url: string)** 發起 GET 請求, 當成功時返回回應內容字串, 失敗時或回應狀態碼不為 200 時返回一個包含 `error` 欄位的物件:

```js
const response = ghttp.get("https://example.com/someapi");
if (response.error) {
	console.log("Request failed:", response.error);
}
if (response.status) {
	console.log("Response status:", response.status);
}
```

**ghttp.getJSON(url: string)** 發起 GET 請求並將回應內容解析為 JSON 物件, 始終返回以下物件:

```js
{
	data?: any, // 當請求成功且回應內容為合法 JSON 時包含解析後的資料
	error?: string, // 當請求失敗或回應狀態碼不為 200 時包含錯誤訊息
	status?: number, // 回應狀態碼, 僅當回應狀態碼不為 200 時包含
}
```

#### Playwright

使用 `playwright` 物件以發起基於瀏覽器的請求.

**playwright.get(url: string)** 發起基於瀏覽器的 GET 請求, 當成功時返回回應內容字串, 失敗時或回應狀態碼不為 200 時返回一個包含 `error` 欄位的物件:

```js
const response = playwright.get("https://example.com/somepage");
if (response.error) {
    console.log("Request failed:", response.error);
}
```

---

最後別忘了呼叫 `registerParser` 註冊解析器:

```js
registerParser({
	metadata,
	canHandle,
	parse
});
```

### Examples

請先查看 [example_parser_basic.js](./example_parser_basic.js) 瞭解最簡示範解析器插件的實作.

然後查看 [example_parser_danbooru.js](./example_parser_danbooru.js) , 這是一個可直接使用的插件, 用於解析 Danbooru 圖片頁面並提取圖片資源.
