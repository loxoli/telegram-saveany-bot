---
title: "yt-dlp 影片下載"
weight: 7
---

# yt-dlp 影片下載

{{< hint warning >}}
該功能需要在系統中安裝 yt-dlp 命令列工具.
{{< /hint >}}

使用 `/ytdlp` 命令可以下載支援的影片網站的影片和音訊, 支援 YouTube、Bilibili、Twitter 等 1000+ 個網站.

```bash
/ytdlp <url1> [url2] [flags...]
```

範例:

```bash
# 基本下載
/ytdlp https://www.youtube.com/watch?v=dQw4w9WgXcQ

# 下載多個影片
/ytdlp https://www.youtube.com/watch?v=video1 https://www.youtube.com/watch?v=video2

# 使用自訂參數
/ytdlp https://www.youtube.com/watch?v=dQw4w9WgXcQ -f best
/ytdlp https://www.youtube.com/watch?v=dQw4w9WgXcQ --extract-audio --audio-format mp3
```

常用參數:

- `-f <format>`: 指定下載格式 (如 `best`, `worst`, `bestvideo+bestaudio`)
- `--extract-audio`: 提取音訊
- `--audio-format <format>`: 音訊格式 (如 `mp3`, `m4a`, `wav`)
- `--write-sub`: 下載字幕
- `--write-thumbnail`: 下載縮圖

更多參數請參考 [yt-dlp 文件](https://github.com/yt-dlp/yt-dlp#usage-and-options).
