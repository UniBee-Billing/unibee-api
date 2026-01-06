package checkout

import (
	"context"
	"fmt"
	"strings"
	"unibee/utility"

	"unibee/internal/consts"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"

	"unibee/api/checkout/translater"
)

func (c *ControllerTranslater) Translate(ctx context.Context, req *translater.TranslateReq) (res *translater.TranslateRes, err error) {
	utility.Assert(len(req.TargetLang) > 0, "Invalid target Lang")
	utility.Assert(consts.IsSupportedLanguage(req.TargetLang), fmt.Sprintf("Unsupported target Lang: %s; supported: %s", req.TargetLang, strings.Join(consts.GetSupportedLanguagesList(), ",")))
	utility.Assert(len(req.Texts) > 0, "Invalid texts")
	// Source language is fixed to English
	const defaultSourceLang = "en"
	// Cache TTL: 30 days in seconds
	const cacheTTLSeconds = 30 * 24 * 60 * 60
	// TODO: Replace with real Google API Key (placeholder)
	const googleApiKey = "REPLACE_ME_WITH_REAL_GOOGLE_API_KEY"
	// TODO: Replace with real DeepL API Key (placeholder)
	const deepLApiKey = "REPLACE_ME_WITH_REAL_DEEPL_API_KEY"

	// Deduplicate while preserving input order
	texts := make([]string, 0, len(req.Texts))
	seen := make(map[string]struct{})
	for _, t := range req.Texts {
		tt := strings.TrimSpace(t)
		if tt == "" {
			continue
		}
		if _, ok := seen[tt]; ok {
			continue
		}
		seen[tt] = struct{}{}
		texts = append(texts, tt)
	}
	utility.Assert(len(texts) > 0, "Invalid texts")

	result := make(map[string]string, len(texts))
	successAll := true

	// Query cache first
	missList := make([]string, 0)
	for _, t := range texts {
		cacheKey := getTranslateCacheKey(defaultSourceLang, req.TargetLang, t)
		cached, cerr := g.Redis().Get(ctx, cacheKey)
		if cerr == nil && cached != nil && cached.String() != "" {
			result[t] = cached.String()
			continue
		}
		missList = append(missList, t)
	}

	// Call DeepL first for cache misses, fallback to Google on failure
	if len(missList) > 0 {
		// Normalize languages to DeepL codes
		srcDL := mapToDeepLLang(defaultSourceLang)
		tgtDL := mapToDeepLLang(req.TargetLang)
		var translatedMap map[string]string
		ok := false
		if srcDL != "" && tgtDL != "" && deepLApiKey != "" {
			translatedMap, ok = translateWithDeepL(ctx, deepLApiKey, srcDL, tgtDL, missList)
		}
		if !ok {
			// Fallback to Google
			translatedMap, ok = translateWithGoogle(ctx, googleApiKey, defaultSourceLang, req.TargetLang, missList)
		}
		if ok {
			// Write results to cache
			for orig, trans := range translatedMap {
				result[orig] = trans
				cacheKey := getTranslateCacheKey(defaultSourceLang, req.TargetLang, orig)
				// Do not fail the main flow if cache set fails
				if _, e := g.Redis().Set(ctx, cacheKey, trans); e == nil {
					_, _ = g.Redis().Expire(ctx, cacheKey, cacheTTLSeconds)
				}
			}
			// Fallback to original for any missing results (edge cases)
			for _, t := range missList {
				if _, ok := result[t]; !ok {
					result[t] = t
					successAll = false
				}
			}
		} else {
			// If Google fails or quota exceeded, fallback to original for all misses
			for _, t := range missList {
				result[t] = t
			}
			successAll = false
		}
	}

	// Complete mapping for all texts (cached ones already filled)
	for _, t := range texts {
		if _, ok := result[t]; !ok {
			result[t] = t
			successAll = false
		}
	}

	return &translater.TranslateRes{
		Translates: result,
		Success:    successAll,
	}, nil
}

// getTranslateCacheKey generates cache key
func getTranslateCacheKey(src, dst, text string) string {
	// Use MD5 to avoid overly long keys
	sum := utility.MD5(text)
	return fmt.Sprintf("UniBee#Translate#%s#%s#%s", src, dst, sum)
}

// translateWithGoogle calls Google Translate batch API
// Returns map[original]->translated and overall success flag
func translateWithGoogle(ctx context.Context, apiKey string, source, target string, texts []string) (map[string]string, bool) {
	// Google Translate API v2 endpoint (placeholder). You may switch to v3/Advanced later.
	// POST https://translation.googleapis.com/language/translate/v2?key=API_KEY
	// body: multiple q=..., source=en, target=xx, format=text
	if apiKey == "" {
		return nil, false
	}

	endpoint := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", apiKey)
	form := make([]string, 0, len(texts)+3)
	for _, t := range texts {
		// q 可以重复出现，使用标准 form 序列化规则
		form = append(form, fmt.Sprintf("q=%s", gjson.New(t).MustToJsonString()))
	}
	// Note: JSON-based encoding may mis-escape; build x-www-form-urlencoded manually
	// Use manual escaping instead
	form = form[:0]
	for _, t := range texts {
		form = append(form, "q="+urlEncode(t))
	}
	form = append(form, "source="+urlEncode(source))
	form = append(form, "target="+urlEncode(target))
	form = append(form, "format=text")
	payload := []byte(strings.Join(form, "&"))

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	body, err := utility.SendRequest(endpoint, "POST", payload, headers)
	if err != nil {
		g.Log().Warningf(ctx, "Google translate request error: %v", err)
		return nil, false
	}

	// Response schema based on Google Translate v2
	type translateText struct {
		TranslatedText string `json:"translatedText"`
		DetectedSource string `json:"detectedSourceLanguage"`
	}
	type dataWrapper struct {
		Translations []translateText `json:"translations"`
	}
	type respWrapper struct {
		Data  dataWrapper `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	var resp respWrapper
	if e := gjson.Unmarshal(body, &resp); e != nil {
		g.Log().Warningf(ctx, "Google translate unmarshal error: %v", e)
		return nil, false
	}
	if resp.Error != nil {
		g.Log().Warningf(ctx, "Google translate api error: %s", resp.Error.Message)
		return nil, false
	}

	if len(resp.Data.Translations) != len(texts) {
		// If sizes mismatch, do best-effort mapping
		g.Log().Warningf(ctx, "Google translate size mismatch, got=%d want=%d", len(resp.Data.Translations), len(texts))
	}

	out := make(map[string]string, len(texts))
	for i := range texts {
		if i < len(resp.Data.Translations) {
			out[texts[i]] = resp.Data.Translations[i].TranslatedText
		} else {
			out[texts[i]] = texts[i]
		}
	}
	return out, true
}

// translateWithDeepL calls DeepL translate API (primary provider)
// Returns map[original]->translated and overall success flag
func translateWithDeepL(ctx context.Context, apiKey string, source, target string, texts []string) (map[string]string, bool) {
	if apiKey == "" {
		return nil, false
	}
	// DeepL API endpoint (paid): https://api.deepl.com/v2/translate
	// For free plan: https://api-free.deepl.com/v2/translate
	// We use the paid endpoint as placeholder.
	endpoint := "https://api.deepl.com/v2/translate"

	// Build x-www-form-urlencoded body
	form := make([]string, 0, len(texts)+3)
	for _, t := range texts {
		form = append(form, "text="+urlEncode(t))
	}
	form = append(form, "source_lang="+urlEncode(strings.ToUpper(source)))
	form = append(form, "target_lang="+urlEncode(strings.ToUpper(target)))
	payload := []byte(strings.Join(form, "&"))

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "DeepL-Auth-Key " + apiKey,
	}

	body, err := utility.SendRequest(endpoint, "POST", payload, headers)
	if err != nil {
		g.Log().Warningf(ctx, "DeepL translate request error: %v", err)
		return nil, false
	}

	// Response schema per DeepL docs
	// { "translations": [ { "detected_source_language": "EN", "text": "Hallo, Welt!" } ] }
	type deeplTranslation struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	}
	type deeplResp struct {
		Translations []deeplTranslation `json:"translations"`
		Message      string             `json:"message"`
	}

	var resp deeplResp
	if e := gjson.Unmarshal(body, &resp); e != nil {
		g.Log().Warningf(ctx, "DeepL translate unmarshal error: %v", e)
		return nil, false
	}
	if len(resp.Translations) == 0 {
		// Errors may appear as { "message": "..." }
		if resp.Message != "" {
			g.Log().Warningf(ctx, "DeepL translate api error: %s", resp.Message)
		}
		return nil, false
	}
	if len(resp.Translations) != len(texts) {
		g.Log().Warningf(ctx, "DeepL translate size mismatch, got=%d want=%d", len(resp.Translations), len(texts))
	}
	out := make(map[string]string, len(texts))
	for i := range texts {
		if i < len(resp.Translations) {
			out[texts[i]] = resp.Translations[i].Text
		} else {
			out[texts[i]] = texts[i]
		}
	}
	return out, true
}

// mapToDeepLLang maps internal codes to DeepL accepted language codes.
// DeepL generally expects upper-case codes like EN, DE, FR, ES, IT, NL, PL, PT-PT, PT-BR, RU, JA, ZH, etc.
// We map our supported set conservatively; unmapped returns empty string.
func mapToDeepLLang(code string) string {
	switch code {
	case "en":
		return "EN"
	case "de":
		return "DE"
	case "fr":
		return "FR"
	case "es":
		return "ES"
	case "it":
		return "IT"
	case "nl":
		return "NL"
	case "pt":
		// Default to generic PT; caller may refine to PT-PT / PT-BR later if needed
		return "PT"
	case "ru":
		return "RU"
	case "ja":
		return "JA"
	case "zh", "cn":
		return "ZH"
	case "sv":
		return "SV"
	case "no":
		return "NB" // DeepL supports Norwegian Bokmål as NB
	case "hi":
		return "HI"
	case "th":
		return "TH"
	case "ko":
		return "KO"
	case "ar":
		return "AR"
	case "tr":
		return "TR"
	default:
		return ""
	}
}

// urlEncode performs x-www-form-urlencoded encoding
func urlEncode(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
		} else if c == ' ' {
			b.WriteByte('+')
		} else {
			b.WriteString(fmt.Sprintf("%%%02X", c))
		}
	}
	return b.String()
}
