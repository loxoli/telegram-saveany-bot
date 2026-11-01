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
  - ⚠️ **編譯前先確認磁碟空間**：沒有快取時，完整編譯約需 3～4 GB。原因有兩個：`FROM golang:alpine` 沒有固定版本，抓到新版就會整個重編；模組在兩個快取掛載裡各下載一次。2026-11-01 就曾因此把主機磁碟塞滿，導致編譯失敗。空間吃緊時，建議先單獨執行 `docker compose -f docker-compose.local.yml build`，邊編譯邊用 `df` 監控，編譯成功後再 `up -d`。
- **敏感檔案**:
  - `config.toml`：保密，不可入版本庫（已在 `.gitignore`），應備份
  - `data/`：含 SQLite 資料庫（已在 `.gitignore`），應規劃備份

## 查詢部署主機的 log

容器 log 用 json-file 並設定輪替（`docker-compose.local.yml` 的 `logging`：單檔 10MB、保留 3 檔），最多保留約 30MB，更舊的會被丟掉。

```bash
# 最新 N 行（先用這個確認最新一筆的時間）
ssh deploy@10.0.10.11 'docker logs --timestamps --tail 20 saveany-bot 2>&1'

# 最近一段時間 / 某時間點（UTC）之後的 log，略過下載進度
ssh deploy@10.0.10.11 'docker logs --timestamps --since 30m saveany-bot 2>&1 | grep -v "Progress update"'
ssh deploy@10.0.10.11 'docker logs --timestamps --since 2026-11-01T06:50:00Z saveany-bot 2>&1 | grep -v "Progress update"'

# 等新訊息進來再印出（請使用者傳訊息測試時用），最多 60 秒
ssh deploy@10.0.10.11 'base=$(docker logs --tail 1 saveany-bot 2>&1); for i in $(seq 12); do
  [ "$(docker logs --tail 1 saveany-bot 2>&1)" != "$base" ] && break; sleep 5; done
  docker logs --timestamps --tail 30 saveany-bot 2>&1 | grep -v "Progress update"'
```

- log 時間是 **UTC**（台北時間減 8 小時）。加上 `--timestamps` 才看得到日期；程式自己印的時間只有時:分:秒。
- **結果是 0 行時，不要直接判定 bot 沒收到訊息。** 先用 `--tail 20` 看最新一筆的時間，確認查詢方式本身沒問題。
- 訊息處理會留 log 的只有連結或檔案（`Got media`、`Parsed item…`）；`/start` 這類純指令不一定會留 log。
- 重建容器（`up -d` 時設定有變更）會換新的 log 檔，舊 log 就讀不到了；重建前要先備份。
- deploy 帳號沒有權限讀 `/var/lib/docker`，無法直接讀 log 原始檔。

### 2026-11-01 之前的舊容器（紀錄）

舊容器的 log 從頭往後讀時，會在 `2026-07-26T02:41:00` 那行靜默停止，所以 `--since` 回傳 0 行、大 `--tail` 只讀到 7/26。這曾讓我們誤判成「bot 沒收到訊息」。2026-11-01 加上 log 輪替並重建容器之後已經解決。舊 log 的備份在主機的 `/home/deploy/saveany-log-backup-20261101/`：

- `part1-0714-0726.log`
- `part2-0729-now.log`（7/29 到 11/01 07:00）

## 語言

- 所有回應與 git commit 訊息須使用正體中文（zh-TW）。
