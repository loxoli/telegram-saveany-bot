---
title: "設定說明"
---

# 設定說明

SaveAnyBot 的設定檔使用 toml 格式, 您可以在 [TOML 官方網站](https://toml.io/) 上瞭解更多關於 toml 的語法.

SaveAnyBot 需要讀取工作目錄下的 `config.toml` 檔案作為設定檔, 若缺少該檔案則會建立預設檔案, 並嘗試從環境變數中載入設定.

以下是一個最簡的設定檔範例:

```toml
[telegram]
token = "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ"

[[users]]
# telegram user id
id = 777000
blacklist = true

[[storages]]
name = "本機儲存"
type = "local"
enable = true
base_path = "./downloads"
```

## 詳細設定

### 全域設定

- `lang`: Bot 使用的語言, 預設為 `zh-CN` (簡體中文), 設為 `en` 則使用英語.
- `stream`: 是否啟用 Stream 模式, 預設為 `false`. 啟用後 Bot 將直接將檔案串流傳輸到儲存端（若儲存端支援）, 不需要下載到本機
{{< hint warning >}}
Stream 模式對於磁碟空間有限的部署環境十分有用, 但也有一些缺點:
<br />
<ul>
<li>無法使用多執行緒從 Telegram 下載檔案, 速度較慢.</li>
<li>網路不穩定時, 任務失敗率高.</li>
<li>無法在中間層對檔案進行處理, 例如自動檔案類型識別.</li>
<li>並非支援所有儲存端, 不支援的儲存端可能會降級為普通模式或無法上傳.</li>
</ul>
{{< /hint >}}
- `workers`: 同時處理任務數量, 預設為 3
- `threads`: 下載檔案時使用的執行緒數, 預設為 4. 僅在未啟用 Stream 模式時生效.
- `retry`: 任務失敗時的重試次數, 預設為 3.
- `proxy`: 全域代理設定, 設定後程式內一切網路連接將會嘗試使用該代理, 選用.

```toml
lang = "zh-CN"
stream = false
workers = 3
threads = 4
retry = 3
proxy = "socks5://127.0.0.1:7890"
```

### Telegram 設定

- `token`: 您的 Telegram Bot Token, 可以透過 [BotFather](https://t.me/botfather) 建立 Bot 並取得 Token.
- `app_id`, `app_hash`: Telegram API ID & Hash, 在 [Telegram API](https://my.telegram.org/apps) 建立應用程式取得, 若不提供則使用預設值.
- `flood_retry`: Flood 控制重試次數, 預設為 5.
- `rpc_retry`: RPC 請求重試次數, 預設為 5.
- `proxy`: 代理設定, 選用.
  - `enable`: 是否啟用代理.
  - `url`: 代理位址
- `userbot`: userbot 設定, 選用.
  - `enable`: 啟用 userbot 整合, 需要登入使用者帳號, 此時請務必使用自己的 api id & hash.
  - `session`: userbot 會話檔案路徑, 預設為 `data/usersession.db`.

{{< hint warning >}}
啟用 userbot 整合後, bot 可以下載私密頻道和群組的檔案, 但具有無法避免的帳號被封禁的風險.
<br />
開啟 userbot 整合後第一次啟動 bot 時需要透過終端機互動輸入手機號碼、2FA 和驗證碼.
<br />
如果您使用 docker 部署, 請使用 -it 參數為容器提供互動式環境, 然後執行登入操作.
{{< /hint >}}

```toml
[telegram]
token = "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ"
app_id = 1025907
app_hash = "452b0359b988148995f22ff0f4229750"
flood_retry = 5
rpc_retry = 5
[telegram.proxy]
enable = false
url = "socks5://127.0.0.1:7890"
[telegram.userbot]
enable = false
session = "data/usersession.db"
```

### Aria2 設定

Aria2 是一個強大的下載管理器，支援 HTTP/HTTPS、FTP、BitTorrent 等多種協定。啟用後，Bot 可以使用 `/aria2dl` 命令透過 Aria2 下載檔案。

- `enable`: 是否啟用 Aria2 支援，預設為 `false`
- `url`: Aria2 RPC 位址，通常為 `http://localhost:6800/jsonrpc`
- `secret`: Aria2 RPC 金鑰，如果您在 Aria2 中設定了 `rpc-secret`，需要在此填寫
- `remove_after_transfer`: 轉存完成後是否刪除 Aria2 下載的本機檔案，預設為 `true`

{{< hint info >}}
Aria2 需要單獨安裝和執行。您可以參考 [Aria2 官方文件](https://aria2.github.io/) 瞭解如何安裝和設定 Aria2。
{{< /hint >}}

```toml
[aria2]
enable = true
url = "http://localhost:6800/jsonrpc"
secret = "your-rpc-secret"
remove_after_transfer = true
```

### yt-dlp 設定

用於設定 `/ytdlp` 命令以及 HTTP API 中 `ytdlp` 任務類型在未傳入自訂參數時的預設行為.

- `max_height`: 預設下載的最高影片解析度 (按高度限制), 如 `1080`, `720`, `480`; `0` 表示不限制 (下載最佳畫質). 當設定了 `format` 時此項被忽略.
- `format`: 直接指定 yt-dlp format 選擇表達式, 設定後優先順序高於 `max_height`, 例如 `bv*[height<=720]+ba/b`.
- `recode`: 下載後轉封裝的影片容器格式 (如 `mp4`), 留空則不轉封裝.

{{< hint info >}}
這些預設值僅在使用 `/ytdlp` 命令且未傳入任何自訂參數時生效. 在命令中傳入自訂參數 (或在 API 中傳入 `flags`) 會覆蓋這些預設值.
{{< /hint >}}

```toml
[ytdlp]
max_height = 1080
format = ""        # 例如 "bv*[height<=720]+ba/b"
recode = "mp4"     # 留空則不轉封裝
```

### HTTP API 設定

啟用後, SaveAny-Bot 會暴露一套 HTTP API, 用於以程式化方式建立/查詢/取消任務. 完整的介面說明見 [HTTP API](../../usage/api).

- `enable`: 是否啟用 HTTP API 服務, 預設為 `false`.
- `host`: 監聽位址, 預設 `0.0.0.0`.
- `port`: 監聽連接埠, 預設 `8080`.
- `token`: 鑑權 Token, **強烈建議設定** — 若為空, API 將在無任何鑑權的情況下暴露.

```toml
[api]
enable = false
host = "0.0.0.0"
port = 8080
token = "your-token"
```

### 日誌設定

- `level`: 日誌層級, 可選 `debug`, `info`, `warn`, `error`, `fatal`. 預設為 `info`.

```toml
[log]
level = "info"
```

### 儲存端清單

儲存端清單用於定義 Bot 支援的儲存位置, 每個儲存端需要指定名稱、類型和相關設定, 使用雙中括號語法 `[[storages]]` 定義.

每一個儲存端至少需要以下欄位:

- `name`: 儲存端名稱, 用於在 Bot 中識別, 需要唯一
- `enable`: 是否啟用該儲存端, 預設為 `true`
- `type`: 儲存端類型, 目前支援以下類型:
  - `local`: 本機磁碟
  - `alist`: Alist
  - `webdav`: WebDAV
  - `s3`: aws S3 及其他相容 S3 的服務
  - `rclone`: 呼叫 rclone 實現上傳
  - `telegram`: 上傳到 Telegram

範例, 這是一個包含本機儲存和 webdav 儲存的設定:

```toml
[[storages]]
name = "本機儲存"
type = "local"
enable = true
# 以下是 local 類型儲存的自訂設定
base_path = "./downloads"

[[storages]]
name = "WebDAV"
type = "webdav"
enable = true
# 以下是 webdav 類型儲存的自訂設定
url = "https://example.com/webdav"
base_path = "/path/to/webdav"
username = "your_username"
password = "your_password"
```

所有儲存端的自訂設定項可查看 [儲存端設定](./storages) 

### 使用者清單

使用者清單用於定義對儲存端的存取控制, 每個使用者需要指定 Telegram 上的使用者 ID, 使用雙中括號語法 `[[users]]` 定義.

- `id`: 使用者的 Telegram User ID
- `storages`: 過濾的儲存端清單, 使用儲存端名稱定義, 預設為白名單模式 (即只允許存取清單中的儲存端)
- `blacklist`: 是否啟用黑名單模式, 預設為 `false`. 若啟用黑名單模式, 則僅允許存取**不在**清單中的儲存端.

範例, 這是一個包含三個使用者的設定, 使用者 `123123` 只能存取本機儲存, 使用者 `456456` 只能存取除 WebDAV 以外的儲存, 使用者 `789789` 啟用黑名單模式但沒有指定儲存端, 因此可以存取所有儲存:

```toml
[[users]]
id = 123123
storages = ["本機儲存"]

[[users]]
id = 456456
storages = ["WebDAV"]
blacklist = true

[[users]]
id = 789789
storages = []
blacklist = true
```

### 事件觸發

事件觸發提供了在 Bot 處理任務時根據任務狀態執行自訂操作的能力, 目前僅支援任意命令執行. 使用 `[hook.exec]` 設定.

目前具有以下幾種事件類型:

- `task_before_start`: 任務即將開始前
- `task_success`: 任務成功完成後
- `task_fail`: 任務失敗後
- `task_cancel`: 任務被取消後

提供的設定值需要為完整的命令列命令, Bot 會在事件發生時執行該命令. 範例:

```toml
[hook.exec]
task_before_start = "echo '任務即將開始'"
task_success = "bash /path/to/success_script.sh"
task_fail = "curl -X POST https://example.com/api/notify -d '任務失敗'"
task_cancel = "bash /path/to/cancel_script.sh"
```

### 解析器

解析器為 Bot 提供了處理非 Telegram 檔案的能力, 例如從其他網站下載檔案. 使用 `[parsers]` 設定.

```toml
[parsers]
plugin_enable = true # 是否啟用解析器插件
plugin_dirs = ["./plugins"] # 插件目錄, 可以是多個目錄
```

上述兩個設定項只用於控制以 JavaScript 撰寫的解析器插件, Bot 還有內建的使用 Go 實作的解析器, 目前預設開啟.

### 雜項

```toml
no_clean_cache = false # 是否在退出時不清空快取資料夾
# 暫存下載資料夾設定
[temp]
base_path = "./cache"
```
