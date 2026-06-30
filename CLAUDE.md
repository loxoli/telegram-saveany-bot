# SaveAny-Bot — Claude Code 提示

## 開始任何任務前，請先讀取以下兩個檔案

1. **`.claude/project-overview.md`** — 技術棧與目錄結構（AGENTS.md 未涵蓋的部分）
2. **`AGENTS.md`** — 開發規範、架構慣例、Build/Test 指令、程式碼風格、常見模式

> Claude Code 不會自動讀取 AGENTS.md，必須透過此指示手動載入。

## 部署

- **伺服器**: `deploy@10.0.10.11`，以 ai-agent 的 SSH private key 連線
- **專案目錄**: `/home/deploy/telegram-saveany-bot`
- **部署方式（從原始碼編譯）**：目前一律以本地原始碼編譯部署，暫不使用官方 registry 映像。
  1. `git fetch origin <branch>` → `git checkout <branch>`（部署 main 以外的 branch 時務必切換）
  2. `docker compose -f docker-compose.local.yml up -d --build`
  - `docker-compose.local.yml` 含 `build: .`，會用 Dockerfile 編譯當前 checkout 的程式碼。
  - ⚠️ **不要使用** `docker-compose.yml`（預設檔）/ `docker compose up -d --build`：該檔直接 pull 官方映像 `ghcr.io/krau/saveany-bot:latest`，`--build` 為空操作，**不會套用本地程式碼**。
- **敏感檔案**:
  - `config.toml`：保密，不可入版本庫（已在 `.gitignore`），應備份
  - `data/`：含 SQLite 資料庫（已在 `.gitignore`），應規劃備份

## 語言

- 所有回應與 git commit 訊息須使用正體中文（zh-TW）。
