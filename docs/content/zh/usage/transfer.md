---
title: "儲存間傳輸"
weight: 8
---

# 儲存間傳輸

使用 `/transfer` 命令可以在不同儲存之間直接傳輸檔案, 無需經過 Telegram.

```bash
/transfer <source_storage>:/<source_path> [filter]
```

參數說明:

- `source_storage`: 來源儲存名稱
- `source_path`: 來源路徑
- `filter`: 選用的正則表達式過濾器, 只傳輸匹配的檔案

範例:

```bash
# 傳輸整個目錄
/transfer local1:/downloads

# 傳輸指定路徑的檔案
/transfer alist1:/media/photos

# 只傳輸 mp4 檔案
/transfer webdav1:/videos ".*\.mp4$"

# 傳輸圖片檔案
/transfer local1:/pictures "(?i)\.(jpg|png|gif)$"
```

Bot 會:

1. 列出來源路徑下的所有檔案
2. 套用過濾器 (如果提供)
3. 顯示檔案數量和總大小
4. 讓您選擇目標儲存
5. 讓您選擇目標目錄 (如果該儲存設定了目錄)
6. 開始傳輸任務

注意:

- 來源儲存必須支援列舉和讀取功能
- 目標儲存必須支援寫入功能
- 傳輸過程顯示即時進度
- 支援取消正在進行的傳輸任務
