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

- **備註**：2026-06 的工作階段將其安裝於暫存（scratchpad）目錄，工作階段結束後就消失了。
- **2026-11-01 重新安裝**：改裝在使用者目錄 `~/.local/go`（不需 root，會保留下來）。`GOPATH` 使用預設的 `~/go`，模組快取已在那裡。

  ```bash
  curl -sSL -o go.tgz https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
  tar -C ~/.local -xzf go.tgz
  export PATH=$HOME/.local/go/bin:$PATH
  ```
