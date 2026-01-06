package i18n

import (
	"context"
	"strings"
	"unibee/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

func IsLangAvailable(lang string) bool {
	if len(lang) == 0 {
		return false
	}
	return consts.IsSupportedLanguage(strings.ToLower(strings.TrimSpace(lang)))
}

func LocalizationFormat(ctx context.Context, format string, values ...interface{}) string {
	// First try to translate using the current context language
	localize := gi18n.Tf(ctx, format, values...)

	// If the translation result contains {# marker, it means the translation key doesn't exist or translation failed
	if strings.Contains(localize, "{#") {
		// Keep hardcoded English fallback, but use original context
		return g.I18n().Tf(
			gi18n.WithLanguage(ctx, `en`),
			format, values...,
		)
	} else {
		return localize
	}
}
