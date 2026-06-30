---
title: "命令列子命令"
weight: 21
---

# 命令列子命令

除了直接執行 `./saveany-bot` (不帶子命令) 啟動 Telegram Bot 外, 這個二進位檔案還提供兩個將本機檔案上傳到儲存後端的輔助子命令: `upload` (一次性) 和 `watch` (持續監聽).

這些子命令會讀取與 Bot 相同的 `config.toml`, 初始化資料庫和快取, 然後執行任務. 它們**不會**啟動 Telegram Bot 本身, 但 `telegram` 類型的儲存會在需要上傳時臨時啟動 Bot 客戶端來執行上傳.

## `upload` — 上傳單一檔案

```
saveany-bot upload -f <檔案> -s <儲存名稱> [-d <目錄>] [--no-progress]
```

參數:

| 參數 | 必填 | 說明 |
|---|---|---|
| `-f, --file` | 是 | 待上傳的本機檔案路徑 |
| `-s, --storage` | 是 | 目標儲存名稱 (必須存在於 `config.toml`) |
| `-d, --dir` | 否 | 儲存中的目標目錄, 預設使用儲存的 `base_path` |
| `--no-progress` | 否 | 關閉終端機進度條 |

範例:

```bash
# 上傳檔案到 "MyAlist" 的預設目錄
./saveany-bot upload -f ./movie.mp4 -s MyAlist

# 上傳到指定子目錄
./saveany-bot upload -f ./movie.mp4 -s MyAlist -d movies/2026

# 透過 Telegram 儲存上傳並關閉進度條
./saveany-bot upload -f ./photo.jpg -s MyChannel --no-progress
```

## `watch` — 監聽目錄並自動上傳

`watch` 子命令持續監聽一個本機目錄, 將新建或修改的檔案上傳到儲存後端, 並保留相對監聽根目錄的子目錄結構.

```
saveany-bot watch -p <路徑> -s <儲存名稱> [-d <目錄>] [選項]
```

參數:

| 參數 | 預設值 | 說明 |
|---|---|---|
| `-p, --path` | *(必填)* | 要監聽的本機目錄 |
| `-s, --storage` | *(必填)* | 目標儲存名稱 |
| `-d, --dir` | 儲存的 `base_path` | 儲存中的目標目錄 |
| `-r, --recursive` | `false` | 是否遞迴監聽子目錄 |
| `--overwrite` | `false` | 覆蓋儲存上已有的檔案, 而非略過 |
| `--initial-scan` | `false` | 啟動時將目錄中已存在的檔案也上傳 |
| `--debounce` | `2s` | 檔案最後一次寫入後, 等待多久再上傳 |
| `--upload-workers` | `config.workers` | 並行上傳數 |
| `--retry-delay` | `3s` | 上傳重試之間的延遲 |

{{< hint info >}}
寫入完成偵測: 監聽器會按檔案做防抖動處理, 僅當檔案大小在一個 debounce 視窗內保持不變時才上傳, 因此不會上傳未寫完的半成品檔案.
<br />
若某檔案在上傳過程中又被修改, 它會在當前上傳完成後再上傳一次, 而不是被重複排入佇列.
{{< /hint >}}

範例:

```bash
# 遞迴監聽 ./inbox 並且把新檔案上傳到 "MyAlist"
./saveany-bot watch -p ./inbox -s MyAlist -r

# 自訂目標目錄並覆蓋已有檔案
./saveany-bot watch -p ./inbox -s MyAlist -d backup --overwrite

# 啟動時把 ./inbox 中已有的內容也一併上傳
./saveany-bot watch -p ./inbox -s MyAlist --initial-scan
```

### 行為說明

- 相對子目錄結構會被保留: 以 `--path ./inbox` 為例, 寫入 `./inbox/sub/file.txt` 的檔案會被上傳到 `<目標目錄>/sub/file.txt`.
- `watch` 會一直執行直到被中斷 (如 `Ctrl-C` / `SIGINT`), 退出前會等待所有進行中的上傳完成.
- 重試次數遵循 `config.toml` 中的全域 `retry` 值, 各次重試之間間隔 `--retry-delay`.
- `telegram` 類型的儲存會自動啟動 Bot 客戶端來執行上傳.

{{< hint warning >}}
`watch` 子命令與 Bot 內的 `/watch` 命令 (監聽 Telegram 聊天室) 無關. 本子命令監聽的是**本機檔案系統目錄**, 不依賴 Telegram.
{{< /hint >}}
