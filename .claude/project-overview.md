# SaveAny-Bot 專案概要

> 開發規範、架構慣例、Build/Test 指令請見 `AGENTS.md`。本檔案只補充 AGENTS.md 未涵蓋的內容。

## 技術棧

| 項目 | 套件 |
|------|------|
| Telegram Bot | `celestix/gotgproto` |
| Telegram UserBot (MTProto) | `gotd/td` |
| CLI | `spf13/cobra` |
| 設定 | `spf13/viper` (config.toml / `SAVEANY_*` 環境變數) |
| ORM | `gorm.io/gorm` + SQLite (`glebarez/sqlite`) |
| JS 執行環境 | `dop251/goja` (Plugin 系統) |
| 瀏覽器自動化 | `playwright-community/playwright-go` |
| 快取 | `dgraph-io/ristretto/v2` |
| 日誌 | `charmbracelet/log` |
| i18n | `nicksnyder/go-i18n/v2` (zh-Hans, en) |

## 目錄結構

```
telegram-saveany-bot/
├── cmd/              # CLI 入口 (Cobra)：root.go, run.go, upload/, watch/, version.go
├── config/           # Viper 設定定義；config/storage/ 為儲存後端設定
├── core/             # 任務佇列引擎
│   ├── core.go       # Executable 介面、worker pool、AddTask/CancelTask
│   ├── hookutil.go   # 執行外部 hook 指令的工具
│   └── tasks/        # aria2dl, batchtfile, directlinks, parsed, telegraph, tfile, transfer, ytdlp
├── client/
│   ├── bot/          # Bot 客戶端 (gotgproto)；handlers/ 含所有指令處理器
│   ├── user/         # UserBot 客戶端 (gotd MTProto)；負責突破禁止保存限制
│   └── middleware/   # FloodWait、Retry、Recovery 中介層
├── storage/          # 儲存後端介面與實作（local, s3, minio, webdav, alist, telegram, rclone）
├── database/         # GORM 模型：User, Dir, Rule, WatchChat；自動 migrate
├── parsers/          # JS 插件載入器 (Goja runtime)
├── api/              # HTTP API 伺服器（Webhook、檔案串流）
├── pkg/              # 通用套件：queue, rule, parser, enums
├── common/           # i18n、Ristretto cache、工具函數
└── plugins/          # JS 插件範例與 README
```
