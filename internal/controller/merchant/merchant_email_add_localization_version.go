package merchant

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/os/gtime"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) AddLocalizationVersion(ctx context.Context, req *email.AddLocalizationVersionReq) (res *email.AddLocalizationVersionRes, err error) {
	utility.Assert(len(req.TemplateName) > 0, "Invalid template name")
	utility.Assert(req.Localizations != nil, "Invalid localizations")
	utility.Assert(len(req.Localizations) > 0, "please setup at least one template localization version")

	// Validate languages
	for _, loc := range req.Localizations {
		//if !system.IsSupportedLanguage(loc.Language) {
		//	return nil, gerror.Newf("Unsupported language: %s. Supported languages are: en, ru, vi, cn, pt", loc.Language)
		//}
		if loc.Language == "" {
			loc.Language = "en"
		}
		utility.Assert(consts.IsSupportedLanguage(loc.Language), fmt.Sprintf("Unsupported language: %s. Supported languages are: %s", loc.Language, strings.Join(consts.GetSupportedLanguagesList(), ", ")))
	}

	template := query.GetMerchantEmailTemplateByTemplateName(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	utility.Assert(template != nil, "template not found")
	utility.Assert(len(template.LocalizationVersions) < 10, "Reach the maximum number (10) of localization versions")
	one := &bean.MerchantLocalizationVersion{
		VersionId:     fmt.Sprintf("v_%s_%d", utility.JodaTimePrefix(), gtime.Now().Timestamp()),
		VersionName:   req.VersionName,
		Activate:      false,
		Localizations: req.Localizations,
	}
	template.LocalizationVersions = append(template.LocalizationVersions, one)
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, template.LocalizationVersions)
	if err != nil {
		return nil, err
	}
	return &email.AddLocalizationVersionRes{LocalizationVersion: one}, nil
}
