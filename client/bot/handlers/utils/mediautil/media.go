package mediautil

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gotd/td/tg"
	"github.com/krau/SaveAny-Bot/common/utils/strutil"
	"github.com/krau/SaveAny-Bot/common/utils/tgutil"
	"github.com/krau/SaveAny-Bot/database"
	"github.com/krau/SaveAny-Bot/pkg/enums/fnamest"
	"github.com/krau/SaveAny-Bot/pkg/tfile"
)

// filenameTextMaxRunes 由訊息文字產生檔名時擷取的純文字字元數上限。
const filenameTextMaxRunes = 16

func IsSupported(media tg.MessageMediaClass) bool {
	switch media.(type) {
	case *tg.MessageMediaDocument, *tg.MessageMediaPhoto:
		return true
	default:
		return false
	}
}

// RefineFileNames 將「無意義」的自動產生檔名(純數字/ID 樣式)優化為由訊息文字
// 產生的名稱：取訊息前 16 個純文字字元作為檔名。當同一批次中有多個檔案產生相同
// 的文字基底名稱時，會依訊息順序加上 "-N" 編號，避免互相覆蓋。
//
// 已具有實際意義之原始檔名(例如使用者上傳時的檔名)會保留語意，但仍會清理其中的
// emoji 與特殊符號(見 strutil.SanitizeFileName)。
func RefineFileNames(files []tfile.TGFileMessage) {
	type entry struct {
		file  tfile.TGFileMessage
		base  string
		ext   string
		msgID int
	}
	groups := make(map[string][]*entry)
	for _, f := range files {
		if f == nil {
			continue
		}
		msg := f.Message()
		if msg == nil {
			continue
		}
		if !strutil.IsMeaninglessFileName(f.Name()) {
			continue
		}
		text := strings.TrimSpace(msg.GetMessage())
		if text == "" {
			continue
		}
		base := strutil.GenTextFileNameBase(text, filenameTextMaxRunes)
		if base == "" {
			continue
		}
		ext := filepath.Ext(f.Name())
		if ext == "" {
			ext = tgutil.ExtFromMedia(msg.Media)
		}
		groups[base] = append(groups[base], &entry{file: f, base: base, ext: ext, msgID: msg.GetID()})
	}
	for _, entries := range groups {
		if len(entries) == 1 {
			e := entries[0]
			e.file.SetName(e.base + e.ext)
			continue
		}
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].msgID < entries[j].msgID
		})
		for i, e := range entries {
			e.file.SetName(fmt.Sprintf("%s-%d%s", e.base, i+1, e.ext))
		}
	}

	// 最後統一清理所有檔名中的 emoji 與特殊符號(含未被改名的原始檔名)
	for _, f := range files {
		if f == nil {
			continue
		}
		name := f.Name()
		if name == "" {
			continue
		}
		if cleaned := strutil.SanitizeFileName(name); cleaned != name {
			f.SetName(cleaned)
		}
	}
}

type FilenameTemplateData struct {
	MsgID    string `json:"msgid,omitempty"`
	MsgTags  string `json:"msgtags,omitempty"`
	MsgGen   string `json:"msggen,omitempty"`
	MsgDate  string `json:"msgdate,omitempty"`
	MsgRaw   string `json:"msgraw,omitempty"`
	OrigName string `json:"origname,omitempty"`
	ChatID   string `json:"chatid,omitempty"`
}

func (f FilenameTemplateData) ToMap() map[string]string {
	return map[string]string{
		"msgid":    f.MsgID,
		"msgtags":  f.MsgTags,
		"msggen":   f.MsgGen,
		"msgraw":   f.MsgRaw,
		"msgdate":  f.MsgDate,
		"origname": f.OrigName,
		"chatid":   f.ChatID,
	}
}

func TfileOptions(ctx context.Context, user *database.User, message *tg.Message) []tfile.TGFileOption {
	opts := make([]tfile.TGFileOption, 0)
	var fnameOpt tfile.TGFileOption
	switch user.FilenameStrategy {
	case fnamest.Message.String():
		fnameOpt = tfile.WithName(tgutil.GenFileNameFromMessage(*message))
	case fnamest.Template.String():
		if user.FilenameTemplate == "" {
			log.FromContext(ctx).Warnf("empty filename template")
			fnameOpt = tfile.WithNameIfEmpty(tgutil.GenContentlessFileName(*message))
			break
		}
		tmpl, err := template.New("filename").Parse(user.FilenameTemplate)
		if err != nil {
			log.FromContext(ctx).Errorf("failed to parse filename template: %s", err)
			fnameOpt = tfile.WithNameIfEmpty(tgutil.GenContentlessFileName(*message))
			break
		}
		data := BuildFilenameTemplateData(message)
		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		if err != nil {
			log.FromContext(ctx).Errorf("failed to execute filename template: %s", err)
			fnameOpt = tfile.WithNameIfEmpty(tgutil.GenContentlessFileName(*message))
			break
		}
		fnameOpt = tfile.WithName(sb.String())
	default:
		// 預設策略：保留具意義的原始檔名，其餘維持無意義名稱，
		// 交由 RefineFileNames 依訊息文字優化。
		fnameOpt = tfile.WithNameIfEmpty(tgutil.GenContentlessFileName(*message))
	}
	opts = append(opts, fnameOpt, tfile.WithMessage(message))
	return opts
}

func BuildFilenameTemplateData(message *tg.Message) map[string]string {
	data := FilenameTemplateData{
		MsgID: func() string {
			id := message.GetID()
			if id == 0 {
				return ""
			}
			return fmt.Sprintf("%d", id)
		}(),
		MsgTags: func() string {
			tags := strutil.ExtractTagsFromText(message.GetMessage())
			if len(tags) == 0 {
				return ""
			}
			return strings.Join(tags, "_")
		}(),
		MsgGen: tgutil.GenFileNameFromMessage(*message),
		OrigName: func() string {
			f, _ := tgutil.GetMediaFileName(message.Media)
			return f
		}(),
		MsgDate: func() string {
			date := message.GetDate()
			if date == 0 {
				return ""
			}
			t := time.Unix(int64(date), 0)
			return t.Format("2006-01-02_15-04-05")
		}(),
		MsgRaw: message.GetMessage(),
		ChatID: func() string {
			// 如果訊息是頻道的(從訊息連結中fetch的) 直接使用其chat id,
			// 無論它是否是從其他來源轉發的
			if message.GetPost() {
				peer := message.GetPeerID()
				switch p := peer.(type) {
				case *tg.PeerChannel:
					return intToStringOmitZero(p.ChannelID)
				default: // impossible case
					return intToStringOmitZero(tgutil.ChatIdFromPeer(peer))
				}
			}
			fwdHeader, ok := message.GetFwdFrom()
			if !ok {
				return intToStringOmitZero(tgutil.ChatIdFromPeer(message.GetPeerID()))
			}
			fwdFrom, ok := fwdHeader.GetFromID()
			if !ok {
				return intToStringOmitZero(tgutil.ChatIdFromPeer(message.GetPeerID()))
			}
			return intToStringOmitZero(tgutil.ChatIdFromPeer(fwdFrom))
		}(),
	}.ToMap()
	return data
}

func intToStringOmitZero(i int64) string {
	if i == 0 {
		return ""
	}
	return fmt.Sprintf("%d", i)
}
