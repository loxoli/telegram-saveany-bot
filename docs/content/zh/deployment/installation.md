---
title: "安裝與更新"
---

# 安裝與更新

## 從預編譯檔案部署（推薦）

在 [Release](https://github.com/krau/SaveAny-Bot/releases) 頁面下載對應平台的二進位檔案.

在解壓縮後的目錄中新增 `config.toml` 檔案, 參考 [設定說明](../configuration) 編輯設定檔

執行:

```bash
chmod +x saveany-bot
./saveany-bot
```

### 程序守護

{{< tabs "daemon" >}}
{{< tab "systemd (常規 Linux)" >}}

建立檔案 <code>/etc/systemd/system/saveany-bot.service</code> 並寫入以下內容:

{{< codeblock >}}
[Unit]
Description=SaveAnyBot
After=systemd-user-sessions.service

[Service]
Type=simple
WorkingDirectory=/yourpath/
ExecStart=/yourpath/saveany-bot
Restart=always

[Install]
WantedBy=multi-user.target
{{< /codeblock >}}

設為開機啟動並啟動服務:

{{< codeblock >}}
systemctl enable --now saveany-bot
{{< /codeblock >}}

{{< /tab >}}

{{< tab "procd (OpenWrt)" >}}

<h4>新增開機自動啟動服務</h4>

建立檔案 <code>/etc/init.d/saveanybot</code> ，參考 <a href="https://github.com/krau/SaveAny-Bot/blob/main/docs/confs/wrt_init" target="_blank">wrt_init</a> 並自行修改:

{{< codeblock >}}
#!/bin/sh /etc/rc.common

#This is the OpenWRT init.d script for SaveAnyBot

START=99 
STOP=10
description="SaveAnyBot"

WORKING_DIR="/mnt/mmc1-1/SaveAnyBot"
EXEC_PATH="$WORKING_DIR/saveany-bot"
start() {
    echo "Starting SaveAnyBot..."
    cd $WORKING_DIR
    $EXEC_PATH &
}
stop() {
    echo "Stopping SaveAnyBot..."
    killall saveany-bot
}
reload() {
    stop
    start
}

{{< /codeblock >}}

賦予權限:

{{< codeblock >}}
chmod +x /etc/init.d/saveanybot
{{< /codeblock >}}

然後將檔案複製到 <code>/etc/rc.d</code> 並重新命名為 <code>S99saveanybot</code>, 同樣賦予權限:

{{< codeblock >}}
chmod +x /etc/rc.d/S99saveanybot
{{< /codeblock >}}

<h4>新增快捷指令</h4>

建立檔案 <code>/usr/bin/sabot</code> ，參考 <a href="https://github.com/krau/SaveAny-Bot/blob/main/docs/confs/wrt_bin" target="_blank">wrt_bin</a>  並自行修改，注意此處檔案編碼僅支援 ANSI 936 .

隨後賦予權限:

{{< codeblock >}}
chmod +x /usr/bin/sabot
{{< /codeblock >}}

使用: <code>sudo sabot start|stop|restart|status|enable|disable</code>

{{< /tab >}}
{{< /tabs >}}


## 使用 Docker 部署

### Docker Compose

下載 [docker-compose.yml](https://github.com/krau/SaveAny-Bot/blob/main/docker-compose.yml) 檔案, 在同目錄下新增 `config.toml` 檔案, 參考 [config.example.toml](https://github.com/krau/SaveAny-Bot/blob/main/config.example.toml) 編輯設定檔.

啟動:

```bash
docker compose up -d
```

### Docker

```shell
docker run -d --name saveany-bot \
    -v /path/to/config.toml:/app/config.toml \
    -v /path/to/downloads:/app/downloads \
    ghcr.io/krau/saveany-bot:latest
```

{{< hint info >}}
關於 docker 映像的變體版本
<br />
<ul>
<li>預設版本: 包含所有功能和依賴, 體積較大. 如果沒有特殊需求, 請使用此版本</li>
<li>micro: 精簡版本, 去除部分選用依賴, 體積較小</li>
<li>pico: 極簡版本, 僅包含核心功能, 體積最小</li>
</ul>
您可以根據需要, 透過指定不同的標籤來拉取合適的版本, 例如: <code>ghcr.io/krau/saveany-bot:micro</code>
<br />
關於變體版本更詳細的差異, 請參考專案根目錄下的 Dockerfile 檔案.
{{< /hint >}}

## 更新

若使用預編譯二進位檔案部署, 使用以下 CLI 命令更新:

```bash
./saveany-bot up
```

如果是 Docker 部署, 使用以下命令更新:

docker:

```bash
docker pull ghcr.io/krau/saveany-bot:latest
docker restart saveany-bot
```

docker compose:

```bash
docker compose pull
docker compose restart
```
