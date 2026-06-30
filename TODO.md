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
