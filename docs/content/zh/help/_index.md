---
title: "常見問題"
weight: 15
---

# 常見問題

## 上傳 alist 失敗也會顯示成功

在 alist 管理頁面適當調整上傳分片大小, 為 alist 使用更穩定的網路環境部署, 都可以減少這種情況的發生.

## Bot 提示下載成功但是 alist 未顯示

alist 快取了目錄結構, 參考 <a href="https://alist.nn.ci/zh/guide/drivers/common.html#缓存过期" target="_blank">文件</a> 可以調整快取時間

## docker 部署設定了代理後仍無法連接 telegram (初始化客戶端逾時)

docker 不能直接存取宿主機網路, 如果您不熟悉其用法, 請將容器設為 host 模式.
