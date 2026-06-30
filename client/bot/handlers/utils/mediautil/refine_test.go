package mediautil_test

import (
	"testing"

	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"github.com/krau/SaveAny-Bot/client/bot/handlers/utils/mediautil"
	"github.com/krau/SaveAny-Bot/pkg/tfile"
)

type fakeFile struct {
	name string
	msg  *tg.Message
}

func (f *fakeFile) Location() tg.InputFileLocationClass { return nil }
func (f *fakeFile) Dler() downloader.Client             { return nil }
func (f *fakeFile) Size() int64                         { return 0 }
func (f *fakeFile) Name() string                        { return f.name }
func (f *fakeFile) SetName(n string)                    { f.name = n }
func (f *fakeFile) Message() *tg.Message                { return f.msg }

func newFile(name string, id int, text string) *fakeFile {
	return &fakeFile{
		name: name,
		msg: &tg.Message{
			ID:      id,
			Message: text,
			Media:   &tg.MessageMediaPhoto{},
		},
	}
}

func TestRefineFileNames(t *testing.T) {
	t.Run("meaningless name with text -> 16 pure-text chars, no number", func(t *testing.T) {
		f := newFile("123456.jpg", 1, "這是一段測試文字內容範例會被截斷的部分")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f})
		if want := "這是一段測試文字內容範例會被截斷.jpg"; f.name != want {
			t.Errorf("got %q, want %q", f.name, want)
		}
	})

	t.Run("shared base in batch -> sequential numbering by msg order", func(t *testing.T) {
		f1 := newFile("100.jpg", 12, "可愛貓咪")
		f2 := newFile("200.jpg", 10, "可愛貓咪")
		f3 := newFile("300.jpg", 11, "可愛貓咪")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f1, f2, f3})
		// ordered by msgID: f2(10)->1, f3(11)->2, f1(12)->3
		if f2.name != "可愛貓咪-1.jpg" {
			t.Errorf("f2 got %q, want 可愛貓咪-1.jpg", f2.name)
		}
		if f3.name != "可愛貓咪-2.jpg" {
			t.Errorf("f3 got %q, want 可愛貓咪-2.jpg", f3.name)
		}
		if f1.name != "可愛貓咪-3.jpg" {
			t.Errorf("f1 got %q, want 可愛貓咪-3.jpg", f1.name)
		}
	})

	t.Run("meaningful original name kept", func(t *testing.T) {
		f := newFile("my_report.pdf", 1, "這是一段測試文字")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f})
		if f.name != "my_report.pdf" {
			t.Errorf("got %q, want my_report.pdf", f.name)
		}
	})

	t.Run("meaningless name without text kept", func(t *testing.T) {
		f := newFile("123456.jpg", 1, "")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f})
		if f.name != "123456.jpg" {
			t.Errorf("got %q, want 123456.jpg", f.name)
		}
	})

	t.Run("ext inferred from media when name has none", func(t *testing.T) {
		f := newFile("123456", 1, "純文字內容")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f})
		if f.name != "純文字內容.jpg" {
			t.Errorf("got %q, want 純文字內容.jpg", f.name)
		}
	})

	t.Run("different texts are not numbered together", func(t *testing.T) {
		f1 := newFile("100.jpg", 1, "貓咪")
		f2 := newFile("200.jpg", 2, "狗狗")
		mediautil.RefineFileNames([]tfile.TGFileMessage{f1, f2})
		if f1.name != "貓咪.jpg" {
			t.Errorf("f1 got %q, want 貓咪.jpg", f1.name)
		}
		if f2.name != "狗狗.jpg" {
			t.Errorf("f2 got %q, want 狗狗.jpg", f2.name)
		}
	})
}
