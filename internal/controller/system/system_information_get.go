package system

import (
	"context"
	"unibee/api/system/information"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/logic/currency"
	"unibee/time"
	"unibee/utility"
)

// GetSupportedLanguagesList returns the list of supported languages
func GetSupportedLanguagesList() []*information.Language {
	var languages []*information.Language
	for code, fullName := range consts.SupportedLanguages {
		languages = append(languages, &information.Language{
			Code:     code,
			FullName: fullName,
		})
	}
	return languages
}

func (c *ControllerInformation) Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error) {
	res = &information.GetRes{}

	res.SupportTimeZone = time.GetTimeZoneList()
	res.Env = config.GetConfigInstance().Env
	res.IsProd = config.GetConfigInstance().IsProd()
	res.SupportCurrency = currency.GetMerchantCurrencies()
	res.Mode = config.GetConfigInstance().Mode
	res.BuildVersion = utility.ReadBuildVersionInfo(ctx)

	// Set supported languages
	res.SupportLanguage = GetSupportedLanguagesList()

	// Set supported country codes
	res.SupportCountryCode = utility.GetCountryCodeList()

	oauthConfig := config.GetConfigInstance().OAuth
	res.OAuth = &information.OAuthConfig{
		TokenSecret:    oauthConfig.TokenSecret,
		GoogleClientId: oauthConfig.GoogleClientId,
		//GoogleClientSecret: oauthConfig.GoogleClientSecret,
		GithubClientId: oauthConfig.GithubClientId,
		//GithubClientSecret: oauthConfig.GithubClientSecret,
	}

	return res, nil
}
