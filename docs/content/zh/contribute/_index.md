---
title: "參與開發"
weight: 20
---

# 參與開發

在開始之前, 請 Fork 本專案, 並複製到本機, 安裝好 Go 開發環境.

以下是一些貢獻程式碼的指南或建議, 您不必完全遵守, 但將有助於快速 review 並合併您的提交:

- **新功能請先提交 Issue**, 以便討論設計和實作細節, 並避免因與專案設計不符而被拒絕.
- **使用現代開發工具**, 確保提交前格式化程式碼, 並保持風格一致.
- **使用[語義化提交](https://www.conventionalcommits.org/zh-hans/v1.0.0/)**, 避免提交訊息模糊或過於簡單.

## 貢獻新儲存端

1. 在 `pkg/enums/storage/storages.go` 中新增新的儲存端類型, 並執行程式碼生成
2. 在 `config/storage` 目錄下定義儲存端設定, 並新增到 `config/storage/factory.go` 中
3. 在 `storage` 目錄下新建一個套件, 撰寫儲存端實作, 然後在 `storage/storage.go` 中匯入並新增它
4. 更新文件, 新增設定說明

## 貢獻新解析器

您可以選擇使用 Go 撰寫原生的解析器實作（推薦）, 或是使用 JavaScript 以插件的方式實作.

如果使用 Go 撰寫, 請:

1. 在 `parsers` 目錄下新建一個套件, 撰寫解析器實作
2. 在 `parsers/parser.go` 的 `init` 中註冊解析器

如果使用 JavaScript 撰寫, 請參考 `plugins/example_parser_basic.js` 的實作, 並在該資料夾下新建一個 js 檔案, 實作您的解析邏輯.

需要注意, `plugins` 目錄下解析器預設不會被編譯到二進位檔案中, 使用者需要手動下載它們並放到本機指定目錄下以啟用它們.
