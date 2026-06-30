---
title: "監聽聊天室"
weight: 4
---

# 監聽聊天室

{{< hint warning >}}
該功能需開啟 UserBot 整合.
{{< /hint >}}

監聽指定聊天室的訊息, 並自動儲存到預設儲存空間中, 遵從儲存規則, 並且可以設定過濾器來只儲存符合條件的訊息.

監聽聊天室:

```
/watch <chat_id/username> [filter] 
```

取消監聽:

```
/unwatch <chat_id/username>
```

過濾器類型:

## msgre

正則比對訊息文字, 例如:

```
/watch 12345678 msgre:.*hello.*
```

這將會監聽 ID 為 12345678 的聊天室, 並且只儲存訊息文字中包含 "hello" 的訊息.
