---
title: "Aria2 下載"
weight: 6
---

# Aria2 下載

{{< hint warning >}}
該功能需要在設定檔中啟用 Aria2 並設定 RPC 連接.
{{< /hint >}}

使用 `/aria2dl` 命令可以透過 Aria2 下載管理器下載檔案, 支援 HTTP/HTTPS、FTP、BitTorrent 等多種協定.

```bash
/aria2dl <uri1> [uri2] [uri3] ...
```

範例:

```bash
# 下載 HTTP 連結
/aria2dl https://example.com/file.zip

# 下載磁力連結
/aria2dl magnet:?xt=urn:btih:...

# 下載種子檔案 (需要先上傳 .torrent 檔案)
/aria2dl https://example.com/file.torrent
```

設定 Aria2:

在 `config.toml` 中新增:

```toml
[aria2]
enable = true
url = "http://localhost:6800/jsonrpc"
secret = "your-rpc-secret"  # 如果設定了 rpc-secret
remove_after_transfer = true  # 轉存完成後刪除本機檔案
```
