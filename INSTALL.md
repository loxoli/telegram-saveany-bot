# 環境安裝紀錄

本檔案記錄為了建置 / 測試本專案而在開發環境額外安裝的系統套件。

## Go 工具鏈

- **原因**：環境中缺少 `go`，無法執行 `go build` / `go test`。
- **版本**：Go 1.25.0（對應 `go.mod` 的 `go 1.25.0`）。
- **安裝方式**：

  ```bash
  curl -sSL -o go.tgz https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
  tar -C <安裝目錄> -xzf go.tgz
  export GOROOT=<安裝目錄>/go
  export GOPATH=<工作目錄>/gopath
  export PATH=$GOROOT/bin:$PATH
  ```

- **備註**：本工作階段將其安裝於暫存（scratchpad）目錄，非系統全域安裝。若要長期使用，建議改裝至 `/usr/local/go` 並將 `PATH` 寫入 shell 設定。
