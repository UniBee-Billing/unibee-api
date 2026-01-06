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

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) EditLocalizationVersion(ctx context.Context, req *email.EditLocalizationVersionReq) (res *email.EditLocalizationVersionRes, err error) {
	utility.Assert(len(req.TemplateName) > 0, "Invalid template name")
	utility.Assert(len(req.VersionId) > 0, "Invalid versionId")
	utility.Assert(req.Localizations != nil, "Invalid localizations")

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

	// Find and update the specific version
	var targetVersion *bean.MerchantLocalizationVersion
	for _, version := range template.LocalizationVersions {
		if version.VersionId == req.VersionId {
			targetVersion = version
			break
		}
	}
	utility.Assert(targetVersion != nil, "version not found")

	// Update version properties
	if req.VersionName != nil {
		targetVersion.VersionName = *req.VersionName
	}
	targetVersion.Localizations = req.Localizations

	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, template.LocalizationVersions)
	if err != nil {
		return nil, err
	}

	return &email.EditLocalizationVersionRes{LocalizationVersion: targetVersion}, nil
}
