package fnamest

//go:generate go-enum --values --names --noprefix --flag --nocase

// FnameST
/* ENUM(
default, message, template
) */
type FnameST string

var fnameSTDisplay = map[FnameST]map[string]string{
	Default:  {"zh-CN": "預設", "en": "Default"},
	Message:  {"zh-CN": "優先從訊息產生", "en": "Gen From Msg First"},
	Template: {"zh-CN": "自訂模板", "en": "Template"},
}

func GetDisplay(st FnameST, lang string) string {
	if display, ok := fnameSTDisplay[st]; ok {
		if str, ok := display[lang]; ok {
			return str
		}
	}
	return fnameSTDisplay[st]["en"]
}
