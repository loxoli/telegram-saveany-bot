package strutil_test

import (
	"reflect"
	"testing"

	"github.com/krau/SaveAny-Bot/common/utils/strutil"
)

func TestExtractTagsFromText(t *testing.T) {
	tests := []struct {
		text     string
		expected []string
	}{
		{
			text: `初音ミクHappy 16th Birthday -Dear Creators-
			✨エンドイラスト公開！✨
			https://piapro.net/miku16thbd/
			#初音ミク #miku16th`,
			expected: []string{"初音ミク", "miku16th"},
		},
		{
			text: `ひっつきむし
			#創作百合`,
			expected: []string{"創作百合"},
		},
		{
			text:     `#創作百合 #原創`,
			expected: []string{"創作百合", "原創"},
		},
		{
			text:     `プラニャ　#ブルアカ`,
			expected: []string{"ブルアカ"},
		},
		{
			text:     `原神是一款#開放世界#冒險遊戲，由中國著名遊戲公司#miHoYo開發。`,
			expected: []string{},
		},
	}

	for _, test := range tests {
		result := strutil.ExtractTagsFromText(test.text)
		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("ExtractTagsFromText(%s) = %v, expected %v", test.text, result, test.expected)
		}
	}
}

func TestParseIntStrRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		sep     string
		wantMin int64
		wantMax int64
		wantErr bool
	}{
		{
			name:    "normal range",
			input:   "10-20",
			sep:     "-",
			wantMin: 10,
			wantMax: 20,
		},
		{
			name:    "reverse order",
			input:   "30 - 10",
			sep:     "-",
			wantMin: 10,
			wantMax: 30,
		},
		{
			name:    "invalid format",
			input:   "10",
			sep:     "-",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "a-b",
			sep:     "-",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, err := strutil.ParseIntStrRange(tt.input, tt.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntStrRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if min != tt.wantMin || max != tt.wantMax {
					t.Errorf("ParseIntStrRange(%q) = (%d, %d), want (%d, %d)", tt.input, min, max, tt.wantMin, tt.wantMax)
				}
			}
		})
	}
}

func TestParseArgsRespectQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple split",
			input: `/rule add FILENAME-REGEX (?i)\.(mp4|mkv)$ "我的 Alist" /影片`,
			want:  []string{"/rule", "add", "FILENAME-REGEX", "(?i)\\.(mp4|mkv)$", "我的 Alist", "/影片"},
		},
		{
			name:  "escaped quotes",
			input: `/rule add "My \"Awesome\" Folder"`,
			want:  []string{"/rule", "add", `My "Awesome" Folder`},
		},
		{
			name:  "escaped backslash",
			input: `/cmd "C:\\Users\\Admin" test`,
			want:  []string{"/cmd", `C:\Users\Admin`, "test"},
		},
		{
			name:  "multiple quoted parts",
			input: `"Hello World" "你好 世界"`,
			want:  []string{"Hello World", "你好 世界"},
		},
		{
			name:  "unquoted words",
			input: "a b c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "mixed quotes and plain",
			input: `cmd "quoted arg" plain`,
			want:  []string{"cmd", "quoted arg", "plain"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strutil.ParseArgsRespectQuotes(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgsRespectQuotes(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsMeaninglessFileName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"", true},
		{"123456.png", true},
		{"123456", true},
		{"photo_123456.jpg", true},
		{"123456_ck8b1q2a.mp4", true},
		{"987654321.mp4", true},
		{".jpg", true},
		{"my_report.pdf", false},
		{"IMG_2024.jpg", false},
		{"初音ミク.png", false},
		{"report123.pdf", false},
		{"123abc.png", false},
	}
	for _, tt := range tests {
		if got := strutil.IsMeaninglessFileName(tt.name); got != tt.expected {
			t.Errorf("IsMeaninglessFileName(%q) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestGenTextFileNameBase(t *testing.T) {
	tests := []struct {
		text     string
		maxRunes int
		expected string
	}{
		{"這是一段測試文字內容範例會被截斷的部分", 16, "這是一段測試文字內容範例會被截斷"},
		{"  Hello, World! 你好 ", 16, "HelloWorld你好"},
		{"#初音ミク #miku16th", 16, "初音ミクmiku16th"},
		{"😀😀😀", 16, ""},
		{"", 16, ""},
		{"abc", 0, ""},
		{"a b c d e", 3, "abc"},
	}
	for _, tt := range tests {
		if got := strutil.GenTextFileNameBase(tt.text, tt.maxRunes); got != tt.expected {
			t.Errorf("GenTextFileNameBase(%q, %d) = %q, want %q", tt.text, tt.maxRunes, got, tt.expected)
		}
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"⚫️2025最新資訊，最新AI發展&商機!.mp4", "2025最新資訊_最新AI發展_商機.mp4"},
		{"my_report.pdf", "my_report.pdf"},
		{"可愛貓咪-1.jpg", "可愛貓咪-1.jpg"},
		{"Hello, World!.txt", "Hello_World.txt"},
		{"  spaced   name .png", "spaced_name.png"},
		{"😀😀😀.gif", "😀😀😀.gif"}, // 清理後為空 → 保留原檔名
		{"123456.png", "123456.png"},
		{"a___b.mp4", "a_b.mp4"},
		{"純文字檔名.mp4", "純文字檔名.mp4"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := strutil.SanitizeFileName(tt.name); got != tt.expected {
			t.Errorf("SanitizeFileName(%q) = %q, want %q", tt.name, got, tt.expected)
		}
	}
}
