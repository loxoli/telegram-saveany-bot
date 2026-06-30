# 專案任務清單

## 狀態說明
- [ ] 待處理 (Pending)
- [x] 已完成 (Done)
- [?] 需要釐清/卡住 (Blocked/Needs clarification)
- [-] 暫停/取消 (Skipped)

## 任務
- [x] telegram 下載檔案名稱優化
  - 目標：下載檔案時，若檔案名稱無意義並且文章有文字時，取16個純文字作為檔案名稱，並按照順去加上編號，以"-"連接。
  - 驗收：測試一個曾經執行過的下載，確認檔案名稱。
  - 實作：
    - `strutil.IsMeaninglessFileName`（判定純數字/ID 樣式檔名）、`strutil.GenTextFileNameBase`（取 16 個純文字字元）。
    - `tgutil.ExtFromMedia`、`tgutil.GenContentlessFileName`（無文字時的後備名稱）。
    - `mediautil.RefineFileNames`：批次優化檔名；同批次共用文字基底時依訊息順序加上 `-N` 編號。
    - 套用於 save / 媒體相簿 / 批次儲存 / 連結訊息等下載流程。
    - 單元測試：`strutil/string_test.go`、`mediautil/refine_test.go`（全部通過）。
  - 規則：原始檔名具意義者保留；單一檔案不加編號（同批次共用文字才編號）。

- [x] 清理檔名中的 emoji 與特殊符號
  - 目標：有時下載到的檔名包含 emoji 或特殊字元，需移除這些符號後重新命名為一般文字檔名（保留副檔名）。
  - 特殊字元檔名範例：`⚫️2025最新資訊，最新AI發展&商機!.mp4`
  - 驗收：測試一個檔名含 emoji/特殊符號的下載，確認輸出檔名已清理為一般文字。
  - 實作：
    - `strutil.SanitizeFileName`：保留字母（含 CJK）/數字/`-`/`_`，其餘（emoji、標點、符號、空白）以底線取代並收斂、去頭尾底線，保留副檔名；清理後為空則保留原名。
    - 範例輸出：`⚫️2025最新資訊，最新AI發展&商機!.mp4` → `2025最新資訊_最新AI發展_商機.mp4`。
    - 套用：`mediautil.RefineFileNames` 末段統一清理所有檔名（save/相簿/批次/連結）；另於 bot watch 與 user watch 流程清理 → 涵蓋所有下載檔名。
    - 規則：採「底線分隔」風格；編號用的連字號 `-N` 會保留。
    - 單元測試：`strutil/string_test.go`、`mediautil/refine_test.go`（全部通過）。

- [x] bug：Twitter（parsed 任務）下載完成不顯示檔案路徑
  - 現象：下載完成訊息僅顯示「處理完成，資源數量: 1\n儲存路徑: [本機1]:」，`儲存路徑` 後方為空。
  - 原因：單一資源時 `dirPath`/`Task.StorPath` 為空（僅多資源才以 item.Title 建立子目錄），而完成訊息只顯示 `StoragePath()`（目錄），實際檔案存放在 `path.Join(StorPath, resource.Filename)`。
  - 目標：完成訊息顯示有意義的儲存路徑（單一資源時帶上檔名）。
  - 驗收：下載一則單張圖片/影片的 Twitter 連結，確認儲存路徑顯示為 `[本機1]:<檔名>`。
  - 實作：`core/tasks/parsed/taskinfo.go::StoragePath()` 單一資源時回傳 `path.Join(StorPath, resource.Filename)`，多資源維持回傳目錄（`StorPath`）。
