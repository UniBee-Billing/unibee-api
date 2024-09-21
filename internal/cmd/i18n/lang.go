package i18n

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"strings"
	"unibee/utility"
)

func IsLangAvailable(lang string) bool {
	if len(lang) == 0 {
		return false
	}
	availableLangs := []string{"cn", "en", "pt", "ru", "vi"}
	if utility.IsStringInArray(availableLangs, strings.ToLower(strings.TrimSpace(lang))) {
		return true
	} else {
		return false
	}
}

func LocalizationFormat(ctx context.Context, format string, values ...interface{}) string {
	localize := gi18n.Tf(ctx, format, values...)
	if strings.Contains(localize, "{#") {
		return g.I18n().Tf(
			gi18n.WithLanguage(context.TODO(), `en`),
			format, values...,
		)
	} else {
		return localize
	}
}
