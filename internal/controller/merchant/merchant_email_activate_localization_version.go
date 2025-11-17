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

func (c *ControllerEmail) ActivateLocalizationVersion(ctx context.Context, req *email.ActivateLocalizationVersionReq) (res *email.ActivateLocalizationVersionRes, err error) {
	utility.Assert(len(req.TemplateName) > 0, "Invalid template name")
	utility.Assert(len(req.VersionId) > 0, "Invalid versionId")
	template := query.GetMerchantEmailTemplateByTemplateName(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	utility.Assert(template != nil, "template not found")
	var one *bean.MerchantLocalizationVersion
	for _, v := range template.LocalizationVersions {
		if req.VersionId == v.VersionId {
			one = v
		}
	}
	utility.Assert(one != nil, "Invalid localization versionId")
	one.Activate = true
	for _, v := range template.LocalizationVersions {
		if v != one {
			v.Activate = false
		}
	}
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, template.LocalizationVersions)
	if err != nil {
		return nil, err
	}

	return &email.ActivateLocalizationVersionRes{}, nil
}
