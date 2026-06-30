---
title: "儲存端設定"
---

# 儲存端設定

請先閱讀 [設定說明](../) 瞭解設定檔的基本格式.

## Alist

`type=alist`

不支援 Stream 模式.

```toml
url = "https://alist.example.com" # Alist 的 URL
username = "your_username"  # Alist 的使用者名稱
password = "your_password" # Alist 的密碼
base_path = "/path/saveanybot" # Alist 中的基礎路徑, 所有檔案將儲存在此路徑下
token_exp = 3600 # Alist 存取令牌的自動更新時間, 單位秒
token = "your_token" 
# Alist 的存取令牌, 選用, 如果不設定則使用使用者名稱和密碼進行身份驗證. 
# 使用 token 驗證時無法自動更新 token
```

## 本機磁碟

`type=local`

```toml
base_path = "./downloads" # 本機儲存的基礎路徑, 所有檔案將儲存在此路徑下
```

## WebDAV
`type=webdav`

```toml
url = "https://webdav.example.com" # WebDAV 的 URL
username = "your_username"  # WebDAV
password = "your_password" # WebDAV 的密碼
base_path = "/path/to/webdav" # WebDAV 中的基礎路徑, 所有檔案將儲存在此路徑下
```

## S3

`type=s3`

```toml
endpoint = "s3.example.com" # S3 的端點, 預設為 aws S3 的端點
region = "us-east-1" # S3 的區域
access_key_id = "your_access_key_id" # S3 的存取金鑰 ID
secret_access_key = "your_secret_access_key" # S3 的私密存取金鑰
bucket_name = "your_bucket_name" # S3 的儲存桶名稱
base_path = "/path/to/s3" # S3 中的基礎路徑, 所有檔案將儲存在此路徑下
virtual_host = false # 使用虛擬主機風格的 URL, 預設為 false
```

虛擬主機風格的 URL 範例:

```
https://your_bucket_name.s3.example.com/path/to/s3/your_file
```

路徑風格（關閉 virtual_host）的 URL 範例:

```
https://s3.example.com/your_bucket_name/path/to/s3/your_file
```

如果您使用的是第三方相容 S3 的服務, 一般使用的是路徑風格的 URL. 而 AWS S3 則通常使用虛擬主機風格的 URL. 詳情請參考您所使用的 S3 相容服務的文件.

## Telegram

`type=telegram`

不支援 Stream 模式.

```toml
# Telegram 聊天室 ID, Bot 將把檔案傳送到這個聊天室
chat_id = "123456789"
# 是否強制使用檔案方式傳送, 預設為 false
force_file = false
# 是否略過大型檔案, 預設為 false. 如果啟用, 超過 Telegram 限制的檔案將不會上傳.
skip_large = false
# 分卷大小, 單位 MB, 預設為 2000 MB (2 GB). 
# 超過該大小的檔案將被分割成多個部分上傳.（使用 zip 格式）
# 當 skip_large 啟用時, 該選項無效.
spilt_size_mb = 2000
```

## Rclone

`type=rclone`

透過 [rclone](https://rclone.org/) 命令列工具支援多種雲端儲存服務. 需要先安裝 rclone 並設定好遠端儲存.

```toml
# rclone 設定的遠端名稱, 可以是任何在 rclone.conf 中設定的遠端
remote = "mydrive"
# 在遠端儲存中的基礎路徑, 所有檔案將儲存在此路徑下
base_path = "/telegram"
# rclone 設定檔的路徑, 選用, 留空使用預設路徑 (~/.config/rclone/rclone.conf)
config_path = ""
# 傳遞給 rclone 命令的額外參數, 選用
flags = ["--transfers", "4", "--checkers", "8"]
```

### 設定 rclone 遠端

首先需要設定 rclone 遠端, 執行 `rclone config` 命令進行互動式設定, 或直接編輯 `rclone.conf` 檔案.

rclone 支援多種雲端儲存服務, 包括但不限於:
- Google Drive
- Dropbox
- OneDrive
- Amazon S3 及相容服務
- SFTP
- FTP
- 更多服務請參考 [rclone 官方文件](https://rclone.org/overview/)

### 使用範例

設定 Google Drive 後, 可以這樣設定儲存:

```toml
[[storages]]
name = "GoogleDrive"
type = "rclone"
enable = true
remote = "gdrive"
base_path = "/SaveAnyBot"
```

如果使用自訂的 rclone 設定檔:

```toml
[[storages]]
name = "MyRemote"
type = "rclone"
enable = true
remote = "myremote"
base_path = "/backup"
config_path = "/path/to/rclone.conf"
flags = ["--progress"]
```
