#!/usr/bin/env bash
# 檢查 deploy 主機上的 SaveAny-Bot 運作狀態

HOST="deploy@10.0.10.11"
PROJECT_DIR="/home/deploy/telegram-saveany-bot"
LOG_LINES="${1:-80}"   # 預設顯示最後 80 行 log，可傳入參數覆蓋

ssh "$HOST" bash <<EOF
set -euo pipefail
cd "$PROJECT_DIR"

echo "====== 容器狀態 ======"
docker compose ps

echo ""
echo "====== 資源使用 ======"
docker stats --no-stream --format "table {{.Name}}\tCPU: {{.CPUPerc}}\tMEM: {{.MemUsage}}"

echo ""
echo "====== 磁碟空間 ======"
df -h "$PROJECT_DIR"

echo ""
echo "====== 最近 $LOG_LINES 行 Log ======"
docker compose logs --no-log-prefix --tail="$LOG_LINES"

echo ""
echo "====== ERROR / WARN 摘要 ======"
docker compose logs --no-log-prefix --tail=500 2>&1 | grep -iE "(error|warn|fatal|panic)" | tail -20 || echo "（無錯誤）"
EOF
