package merchant

import (
	"context"
	"unibee/api/bean"
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
	template := query.GetMerchantEmailTemplateByTemplateName(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	utility.Assert(template != nil, "template not found")

	var one *bean.MerchantLocalizationVersion
	for _, v := range template.LocalizationVersions {
		if req.VersionId == v.VersionId {
			one = v
		}
	}
	utility.Assert(one != nil, "Invalid localization versionId")
	one.Localizations = req.Localizations
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, template.LocalizationVersions)
	if err != nil {
		return nil, err
	}

	return &email.EditLocalizationVersionRes{LocalizationVersion: one}, nil
}
