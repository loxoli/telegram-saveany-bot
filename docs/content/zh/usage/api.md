---
title: "HTTP API"
weight: 20
---

# HTTP API

SaveAny-Bot 提供了一套 HTTP API，允許您透過程式化方式建立下載/轉存任務、查詢任務狀態、取消任務等，無需透過 Telegram 操作。

## 啟用 API

在 `config.toml` 中新增或修改以下設定：

```toml
[api]
enable = true
host   = "0.0.0.0"   # 監聽位址，預設 0.0.0.0
port   = 8080         # 監聽連接埠，預設 8080
token  = "your-token" # 鑑權 Token，強烈建議設定
```

也可透過環境變數覆蓋（前綴 `SAVEANY_`）：

| 環境變數 | 對應設定項 |
|---|---|
| `SAVEANY_API_ENABLE` | `api.enable` |
| `SAVEANY_API_HOST` | `api.host` |
| `SAVEANY_API_PORT` | `api.port` |
| `SAVEANY_API_TOKEN` | `api.token` |

{{< hint warning >}}
若 `token` 為空，API 服務將**不進行任何鑑權**即可存取，存在安全風險。
{{< /hint >}}

## 鑑權

當設定了 `token` 時，所有 API 請求均需在 HTTP 請求標頭中攜帶 Bearer Token：

```
Authorization: Bearer <your-token>
```

鑑權失敗時返回 `401`：

```json
{ "error": "unauthorized", "message": "invalid token" }
```

## 錯誤回應格式

所有錯誤均使用統一的 JSON 格式：

```json
{
  "error":   "error_code",
  "message": "錯誤說明"
}
```

常見錯誤碼：

| 錯誤碼 | HTTP 狀態 | 含義 |
|---|---|---|
| `unauthorized` | 401 | 鑑權失敗 |
| `method_not_allowed` | 405 | HTTP 方法不正確 |
| `invalid_request` | 400 | 請求內容/參數非法 |
| `task_creation_failed` | 400 | 任務建立失敗 |
| `task_not_found` | 404 | 任務 ID 不存在 |
| `cancel_failed` | 500 | 取消任務失敗 |
| `internal_error` | 500 | 伺服器內部錯誤 |

---

## 介面清單

### GET /health — 健康檢查

無需鑑權。

**回應 `200 OK`：**

```json
{ "status": "ok" }
```

---

### GET /api/v1/storages — 列出儲存空間

返回當前所有已載入的儲存後端。

**回應 `200 OK`：**

```json
{
  "storages": [
    { "name": "local",   "type": "local" },
    { "name": "MyMinio", "type": "s3" }
  ]
}
```

---

### GET /api/v1/task-types — 列出支援的任務類型

**回應 `200 OK`：**

```json
{
  "types": [
    "directlinks",
    "ytdlp",
    "aria2",
    "parseditem",
    "tgfiles",
    "tphpics",
    "transfer"
  ]
}
```

---

### POST /api/v1/tasks — 建立任務

**請求標頭：**

```
Content-Type: application/json
Authorization: Bearer <token>
```

**請求內容：**

```json
{
  "type":    "<任務類型>",
  "storage": "<儲存名稱>",
  "path":    "<子目錄>",
  "webhook": "<回呼URL>",
  "params":  { }
}
```

| 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `type` | string | 是 | 任務類型，見下文 |
| `storage` | string | 是 | 目標儲存名稱，須與設定中的儲存名稱一致 |
| `path` | string | 否 | 儲存內的子目錄路徑 |
| `webhook` | string | 否 | 任務完成/失敗時的回呼位址 |
| `params` | object | 是 | 各任務類型的專屬參數，見下文 |

**回應 `201 Created`：**

```json
{
  "task_id":    "abc123xyz",
  "type":       "directlinks",
  "status":     "queued",
  "created_at": "2026-03-11T10:00:00Z"
}
```

#### 任務類型與 params

##### directlinks — 直接下載連結

下載一個或多個 HTTP/HTTPS 直鏈檔案。

```json
{
  "type":    "directlinks",
  "storage": "local",
  "path":    "downloads",
  "params": {
    "urls": [
      "https://example.com/file.zip",
      "https://example.com/other.zip"
    ]
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `urls` | []string | 是 | 下載位址清單，至少 1 條 |

##### ytdlp — yt-dlp 影片下載

{{< hint warning >}}
需要在系統中安裝 yt-dlp。
{{< /hint >}}

透過 yt-dlp 下載影片/音訊，支援 YouTube、Bilibili 等 1000+ 個網站。

```json
{
  "type":    "ytdlp",
  "storage": "local",
  "path":    "videos",
  "params": {
    "urls":  ["https://www.youtube.com/watch?v=xxx"],
    "flags": ["--extract-audio", "--audio-format", "mp3"]
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `urls` | []string | 是 | 媒體連結清單，至少 1 條 |
| `flags` | []string | 否 | 額外的 yt-dlp 命令列參數 |

##### aria2 — Aria2 下載

{{< hint warning >}}
需要在設定檔中啟用並設定 Aria2 RPC。
{{< /hint >}}

透過 Aria2 下載管理器下載檔案，支援 HTTP/HTTPS、FTP、BitTorrent（磁力連結、種子）等協定。

```json
{
  "type":    "aria2",
  "storage": "local",
  "path":    "downloads",
  "params": {
    "urls":    ["magnet:?xt=urn:btih:..."],
    "options": { "split": "4" }
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `urls` | []string | 是 | 下載位址清單，至少 1 條 |
| `options` | map[string]string | 否 | Aria2 下載選項 |

##### parseditem — 解析器下載

將 URL 交由已註冊的 JS 插件或內建解析器處理後下載。

```json
{
  "type":    "parseditem",
  "storage": "local",
  "path":    "parsed",
  "params": {
    "url": "https://some-site.com/page"
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `url` | string | 是 | 待解析的頁面 URL |

若沒有任何解析器能處理該 URL，則返回 `400 task_creation_failed`。

##### tgfiles — Telegram 訊息檔案下載

透過 Telegram 訊息連結下載檔案。支援以下連結格式：

- `https://t.me/username/123` — 公開頻道/群組
- `https://t.me/c/123456789/123` — 私有頻道（數字 ID）
- `https://t.me/c/123456789/111/456` — 話題訊息
- `https://t.me/username/111/456` — 使用者名稱頻道下的話題訊息

若訊息屬於媒體組（相簿），預設下載整組檔案。在連結末尾追加 `?single` 可強制只下載單條訊息的檔案。

```json
{
  "type":    "tgfiles",
  "storage": "local",
  "path":    "telegram",
  "params": {
    "message_links": [
      "https://t.me/username/123",
      "https://t.me/c/1234567890/456"
    ]
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `message_links` | []string | 是 | Telegram 訊息連結清單，至少 1 條 |

##### tphpics — Telegraph 文章圖片下載

下載 Telegra.ph 文章中的所有圖片。

支援的連結前綴：`https://telegra.ph/`、`http://telegra.ph/`、`https://telegraph.co/`、`http://telegraph.co/`

```json
{
  "type":    "tphpics",
  "storage": "local",
  "path":    "telegraph",
  "params": {
    "telegraph_url": "https://telegra.ph/Some-Article-01-01"
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `telegraph_url` | string | 是 | Telegra.ph 文章 URL |

##### transfer — 儲存間檔案傳輸

在兩個儲存後端之間直接傳輸檔案，無需經過 Telegram。來源儲存須支援列舉（list）和讀取（read）操作。

{{< hint info >}}
`transfer` 任務中，頂層的 `storage` 欄位仍然必須填寫（用於通過參數驗證），但實際使用的儲存由 `params` 中的 `source_storage` 和 `target_storage` 決定。
{{< /hint >}}

```json
{
  "type":    "transfer",
  "storage": "local",
  "params": {
    "source_storage": "MyS3",
    "source_path":    "backups/",
    "target_storage": "LocalDisk",
    "target_path":    "restored/"
  }
}
```

| params 欄位 | 類型 | 必填 | 說明 |
|---|---|---|---|
| `source_storage` | string | 是 | 來源儲存名稱 |
| `source_path` | string | 是 | 來源儲存中的路徑，須包含至少一個檔案 |
| `target_storage` | string | 是 | 目標儲存名稱 |
| `target_path` | string | 是 | 目標儲存中的路徑 |

---

### GET /api/v1/tasks — 列出所有任務

返回所有 API 建立的任務（僅在記憶體中保留，重新啟動後清空）。

**回應 `200 OK`：**

```json
{
  "tasks": [
    {
      "task_id":    "abc123xyz",
      "type":       "directlinks",
      "status":     "running",
      "title":      "file.zip",
      "storage":    "local",
      "path":       "downloads",
      "error":      "",
      "created_at": "2026-03-11T10:00:00Z",
      "updated_at": "2026-03-11T10:00:05Z",
      "progress": {
        "total_bytes":      10485760,
        "downloaded_bytes": 5242880,
        "percent":          50.0
      }
    }
  ],
  "total": 1
}
```

`progress` 欄位僅在 `total_bytes > 0` 時出現。`error` 欄位僅在有錯誤時出現。

---

### GET /api/v1/tasks/{task_id} — 查詢任務

**路徑參數：** `task_id` — 建立任務時返回的 ID。

**回應 `200 OK`：** 同上清單中的單一任務物件。

**錯誤回應：**
- `400 invalid_request` — 路徑中未提供 task_id
- `404 task_not_found` — 任務不存在

---

### DELETE /api/v1/tasks/{task_id} — 取消任務

**路徑參數：** `task_id`

**回應 `200 OK`：**

```json
{ "message": "task cancelled successfully" }
```

**錯誤回應：**
- `400 invalid_request` — 路徑中未提供 task_id
- `404 task_not_found` — 任務不存在
- `500 cancel_failed` — 取消操作失敗

---

## 任務狀態

| 狀態值 | 含義 |
|---|---|
| `queued` | 已入佇列，等待執行 |
| `running` | 正在執行 |
| `completed` | 已成功完成 |
| `failed` | 執行失敗 |
| `cancelled` | 已透過 DELETE 介面取消 |

---

## Webhook 回呼

建立任務時可設定 `webhook` 欄位。當任務進入終態（`completed`、`failed`、`cancelled`）時，Bot 會向該位址傳送一個 `POST` 請求。

**回呼請求標頭：**

```
Content-Type: application/json
User-Agent: SaveAny-Bot/1.0
```

**回呼請求內容：**

```json
{
  "task_id":      "abc123xyz",
  "type":         "directlinks",
  "status":       "completed",
  "storage":      "local",
  "path":         "downloads",
  "completed_at": "2026-03-11T10:01:00Z",
  "error":        ""
}
```

`completed_at` 僅在狀態為 `completed` 或 `failed` 時出現。`error` 僅在有錯誤時出現。

**重試機制：** 最多重試 3 次，重試間隔依次為 1 秒、2 秒、3 秒。每次請求逾時為 30 秒。
