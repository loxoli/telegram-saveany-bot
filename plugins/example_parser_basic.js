// 這是一個最簡範例解析器插件, 用於展示插件所需實作的基本功能
// 此插件將會模擬處理 YouTube 的影片連結

/**
 * 插件元資料
 * 版本號是 saveany-bot 本體支援的插件規範版本號, 必須提供
 */
const metadata = {
    name: "Example Parser", // 插件名稱
    version: "1.0.0", // 插件版本號
    description: "A parser for example links", // 插件描述
    author: "Krau", // 插件作者
}

// 你可以使用 console.log 來在終端中使用 go 的 logger 印出資訊
console.log("Parser loaded", "name", metadata.name);

/**
 * canHandle 函式用於判斷目前解析器能否解析給定的 URL
 */
const canHandle = function (url) {
    // 這裡我們簡單地檢查 URL 是否包含 "youtube.com/watch?v"
    return url.includes("youtube.com/watch?v");
}

/**
 * 解析 url 並回傳一個 Item 物件, 型別定義在 pkg/parser.go 中
 */
const parse = function (url) {
    var result = {
        // 元資訊
        site: "YouTube",
        url: url,
        title: "測試 YouTube 影片",
        author: "某影片作者",
        description: "這是一個測試影片",
        tags: ["test", "youtube"],
        // 資源（可下載的檔案）列表
        resources: [
            {
                url: "https://example.com/video1.mp4", // 檔案直連
                filename: "somevideo.mp4", // 檔名
                mime_type: "video/mp4", // 檔案 MIME 類型, 可選
                extension: "mp4", // 副檔名, 可選
                size: 100 * 1024 * 1024, // 檔案大小, 單位為位元組, 未知可設為 0
                hash: {}, // 檔案雜湊, 可選, 格式為 {"md5": "xxx", "sha256": "xxx"} 等
                headers: {}, // 下載檔案時所需的 HTTP 標頭, 可選, 例如 {"User-Agent": "Mozilla/5.0"}
                extra: {} // 額外資訊, 可選, 可包含任何自訂資料
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

// 最後需要呼叫 registerParser 來註冊這個解析器
registerParser({
    metadata,
    canHandle,
    parse
});

// 更進一步的插件撰寫資訊, 請查看 plugins/example_parser_danbooru.js
