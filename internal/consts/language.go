package consts

// Supported languages constants
const (
	LangEnglish    = "en"
	LangRussian    = "ru"
	LangVietnamese = "vi"
	LangChinese    = "cn"
	LangPortuguese = "pt"
	LangJapanese   = "ja"
	LangGerman     = "de"
	LangFrench     = "fr"
	LangItalian    = "it"
	LangSpanish    = "es"
	LangSwedish    = "sv"
	LangNorwegian  = "no"
	LangHindi      = "hi"
	LangThai       = "th"
	LangKorean     = "ko"
	LangArabic     = "ar"
	LangTurkish    = "tr"
	LangDutch      = "nl"
)

// SupportedLanguages defines all supported languages
var SupportedLanguages = map[string]string{
	LangEnglish:    "English",
	LangRussian:    "Russian",
	LangVietnamese: "Vietnamese",
	LangChinese:    "Chinese",
	LangPortuguese: "Portuguese",
	LangJapanese:   "Japanese",
	LangGerman:     "German",
	LangFrench:     "French",
	LangItalian:    "Italian",
	LangSpanish:    "Spanish",
	LangSwedish:    "Swedish",
	LangNorwegian:  "Norwegian",
	LangHindi:      "Hindi",
	LangThai:       "Thai",
	LangKorean:     "Korean",
	LangArabic:     "Arabic",
	LangTurkish:    "Turkish",
	LangDutch:      "Dutch",
}

// GetSupportedLanguagesList returns the list of supported language codes
func GetSupportedLanguagesList() []string {
	var languages []string
	for code := range SupportedLanguages {
		languages = append(languages, code)
	}
	return languages
}

// IsSupportedLanguage checks if the given language code is supported
func IsSupportedLanguage(lang string) bool {
	_, exists := SupportedLanguages[lang]
	return exists
}
