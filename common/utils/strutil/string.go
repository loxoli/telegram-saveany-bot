package strutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/duke-git/lancet/v2/slice"
)

func HashString(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

var TagRe = regexp.MustCompile(`(?:^|[\p{Zs}\s.,!?(){}[\]<>\"\'，。！？（）：；、])#([\p{L}\d_]+)`)

func ExtractTagsFromText(text string) []string {
	matches := TagRe.FindAllStringSubmatch(text, -1)
	tags := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}
	return slice.Compact(tags)
}

// meaninglessNameRe 匹配自動產生、無實際語意的檔名(忽略副檔名)，例如:
//   - "123456"          純數字 ID
//   - "photo_123456"    照片 ID
//   - "123456_ck8b1q2"  id_xid 形式
var meaninglessNameRe = regexp.MustCompile(`^(?:photo_)?\d+(?:_[a-zA-Z0-9]+)?$`)

// IsMeaninglessFileName 判斷檔名是否為「無意義」的自動產生名稱(純數字/ID 樣式)，
// 例如 "123456.png"、"photo_123.jpg"、"123_ck8b1q2.mp4"。副檔名不列入判斷。
// 空字串亦視為無意義。
func IsMeaninglessFileName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return true
	}
	base := strings.TrimSuffix(name, filepath.Ext(name))
	if base == "" {
		return true
	}
	return meaninglessNameRe.MatchString(base)
}

// GenTextFileNameBase 從 text 取出至多 maxRunes 個「純文字」字元(字母與數字，
// 含 CJK)作為檔名基底，會略過空白、標點、emoji 等符號。若找不到可用字元則回傳空字串。
func GenTextFileNameBase(text string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := make([]rune, 0, maxRunes)
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			runes = append(runes, r)
			if len(runes) >= maxRunes {
				break
			}
		}
	}
	return string(runes)
}

func ParseIntStrRange(input string, sep string) (int64, int64, error) {
	parts := strings.Split(input, sep)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format: %s", input)
	}
	min, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minimum value: %s", parts[0])
	}
	max, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid maximum value: %s", parts[1])
	}
	if min > max {
		min, max = max, min
	}
	return min, max, nil
}

func ParseArgsRespectQuotes(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	escaped := false

	for _, r := range input {
		switch {
		case escaped:
			if r == '"' || r == '\\' {
				current.WriteRune(r)
			} else {
				current.WriteRune('\\')
				current.WriteRune(r)
			}
			escaped = false

		case r == '\\':
			escaped = true

		case r == '"':
			inQuotes = !inQuotes

		case r == ' ' || r == '\t':
			if inQuotes {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
