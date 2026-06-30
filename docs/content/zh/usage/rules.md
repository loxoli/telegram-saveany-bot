---
title: "儲存規則"
weight: 3
---

# 儲存規則

允許您為 Bot 在上傳檔案到儲存空間時設定一些重導向規則, 用於自動整理所儲存的檔案.

見: <a href="https://github.com/krau/SaveAny-Bot/issues/28" target="_blank">#28</a>

目前支援的規則類型:

1. FILENAME-REGEX
2. MESSAGE-REGEX
3. IS-ALBUM

新增規則的基本語法:

"規則類型 規則內容 儲存名稱 路徑"

注意空格的使用, 語法正確 bot 才能解析, 以下是一條合法的新增規則命令:

```
/rule add FILENAME-REGEX (?i)\.(mp4|mkv|ts|avi|flv)$ MyAlist /影片
```

此外, 規則中的儲存名稱若使用 "CHOSEN" , 則表示儲存到點擊按鈕選擇的儲存端的路徑下

您也可以使用 `/rule switch` 來切換規則模式. 關閉規則模式時, 所有檔案都將儲存到預設儲存空間.

## 預設規則

為常見檔案類型手動撰寫正則規則比較繁瑣, 因此 Bot 內建了一組預設分類 (影片、圖片、音訊、文件、壓縮檔), 可以透過一條命令批次匯入:

```
/rule preset <儲存名稱> [基礎路徑]
```

參數:

- `儲存名稱`: 目標儲存名稱 (必須存在且您有權存取)
- `基礎路徑`: 選用. 各預設分類的子目錄會建立在此路徑下; 若不填則直接使用預設分類目錄名稱

範例:

```
# 匯入預設規則到 "MyAlist", 使用預設目錄設定
/rule preset MyAlist

# 在自訂基礎路徑 "downloads/sorted" 下匯入預設規則
/rule preset MyAlist downloads/sorted
```

此命令會為每個分類建立 `FILENAME-REGEX` 規則, 將符合的檔案路由到 `基礎路徑` 下對應的子目錄:

| 分類 | 符合的副檔名 | 預設目錄 |
|---|---|---|
| 影片 | mp4, mkv, ts, avi, flv, mov, webm, wmv, rmvb, m2ts | `影片` |
| 圖片 | jpg, jpeg, png, gif, webp, bmp | `圖片` |
| 音訊 | mp3, flac, wav, aac, m4a, ogg | `音訊` |
| 文件 | pdf, doc, docx, xls, xlsx, ppt, pptx, txt, md, csv, epub, mobi, azw3, chm | `文件` |
| 壓縮檔 | zip, rar, 7z, tar, gz, bz2, xz, ... | `壓縮檔` |

{{< hint info >}}
匯入後的預設規則就是普通的 `FILENAME-REGEX` 規則. 您可以像其他規則一樣透過 `/rule` 查看或用 `/rule del <id>` 單獨刪除/編輯它們.
{{< /hint >}}

規則類型:

## FILENAME-REGEX

根據檔案名稱正則比對, 規則內容要求為一個合法的正則表達式, 如

```
FILENAME-REGEX (?i)\.(mp4|mkv|ts|avi|flv)$ MyAlist /影片
```

表示將檔案名稱副檔名為 mp4,mkv,ts,avi,flv 的檔案放到名為 MyAlist 儲存下的 /影片 目錄內 (同時受設定檔中的 `base_path` 影響)

## MESSAGE-REGEX

同上, 但是是根據訊息本身的文字內容正則比對

## IS-ALBUM

比對相簿訊息 (media group), 規則內容只能為 `true` 或 `false`.

規則中的路徑若使用 "NEW-FOR-ALBUM" , 則表示為該組訊息新建一個資料夾來儲存它們. 見: https://github.com/krau/SaveAny-Bot/issues/87

例如:

```
IS-ALBUM true MyWebdav NEW-FOR-ALBUM
```

這將會把以 media group 形式傳送的訊息儲存到名為 MyWebdav 的儲存下, 並為每個相簿新建一個資料夾（由第一個檔案生成）來儲存它們.
