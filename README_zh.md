<div align="center">

# <img src="docs/static/logo.png" width="45" align="center"> Save Any Bot

> **將 Telegram 上的檔案轉存到多種儲存端**

[![Release Date](https://img.shields.io/github/release-date/krau/saveany-bot?label=release)](https://github.com/krau/saveany-bot/releases)
[![tag](https://img.shields.io/github/v/tag/krau/saveany-bot.svg)](https://github.com/krau/saveany-bot/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/krau/saveany-bot/build-release.yml)](https://github.com/krau/saveany-bot/actions/workflows/build-release.yml)
[![Stars](https://img.shields.io/github/stars/krau/saveany-bot?style=flat)](https://github.com/krau/saveany-bot/stargazers)
[![Downloads](https://img.shields.io/github/downloads/krau/saveany-bot/total)](https://github.com/krau/saveany-bot/releases)
[![Issues](https://img.shields.io/github/issues/krau/saveany-bot)](https://github.com/krau/saveany-bot/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/krau/saveany-bot?label=pr)](https://github.com/krau/saveany-bot/pulls)
[![License](https://img.shields.io/github/license/krau/saveany-bot)](./LICENSE)

</div>

## 🎯 特性

- 支援文件/影片/圖片/貼圖…甚至還有 [Telegraph](https://telegra.ph/)
- 破解禁止儲存的檔案
- 批次下載
- 串流傳輸
- 多使用者使用
- 基於儲存規則的自動整理
- 監聽並自動轉存指定聊天室的訊息, 支援過濾
- 在不同儲存端之間轉存檔案
- 整合 yt-dlp, 從所支援的網站下載並轉存媒體檔案
- 整合 Aria2, 支援直鏈/磁力下載和轉存
- 使用 js 撰寫解析器插件以轉存任意網站的檔案
- 儲存端支援:
  - Alist
  - S3
  - WebDAV
  - 本機磁碟
  - Rclone
  - Telegram (重傳至指定聊天室)

## 快速開始

建立檔案 `config.toml` 並填入以下內容:

```toml
[telegram]
token = "" # 您的 Bot Token, 在 @BotFather 取得
[telegram.proxy]
# 啟用代理連接 telegram
enable = false
url = "socks5://127.0.0.1:7890"

[[storages]]
name = "本機磁碟"
type = "local"
enable = true
base_path = "./downloads"

[[users]]
id = 114514 # 您的 Telegram 帳號 id
storages = []
blacklist = true
```

使用 Docker 執行 Save Any Bot:

```bash
docker run -d --name saveany-bot \
    -v ./config.toml:/app/config.toml \
    -v ./downloads:/app/downloads \
    ghcr.io/krau/saveany-bot:latest
```

請 [**查看文件**](https://sabot.unv.app/) 以取得更多設定選項和使用方法.

## 贊助

本專案受到 [YxVM](https://yxvm.com/) 與 [NodeSupport](https://github.com/NodeSeekDev/NodeSupport) 的支持.

若此專案對您有幫助, 您可以考慮透過以下方式贊助:

- [愛發電](https://afdian.com/a/unvapp)

## 鳴謝

- [gotd](https://github.com/gotd/td)
- [TG-FileStreamBot](https://github.com/EverythingSuckz/TG-FileStreamBot)
- [gotgproto](https://github.com/celestix/gotgproto)
- [tdl](https://github.com/iyear/tdl)
- All the dependencies, contributors, sponsors and users.

## 社群與關於作者

- [![通知群組](https://img.shields.io/badge/ProjectSaveAny-Group-blue)](https://t.me/ProjectSaveAny)
- [![討論區](https://img.shields.io/badge/Github-Discussion-white)](https://github.com/krau/saveany-bot/discussions)
- [![個人頻道](https://img.shields.io/badge/Krau-PersonalChannel-cyan)](https://t.me/acherkrau)
