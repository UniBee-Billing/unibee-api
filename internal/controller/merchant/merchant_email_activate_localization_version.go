package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) ActivateLocalizationVersion(ctx context.Context, req *email.ActivateLocalizationVersionReq) (res *email.ActivateLocalizationVersionRes, err error) {
	utility.Assert(len(req.TemplateName) > 0, "Invalid template name")
	template := query.GetMerchantEmailTemplateByTemplateName(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	utility.Assert(template != nil, "template not found")
	var one *bean.MerchantLocalizationVersion
	if len(req.VersionId) > 0 {
		for _, v := range template.LocalizationVersions {
			if req.VersionId == v.VersionId {
				one = v
			}
		}
		utility.Assert(one != nil, "Invalid localization versionId")
		one.Activate = true
		if one != nil && one.Localizations != nil && len(one.Localizations) > 0 {
			for _, v := range one.Localizations {
				utility.Assert(len(v.Title) > 0, fmt.Sprintf("Empty Subject For Language:%s", v.Language))
				utility.Assert(len(v.Content) > 0, fmt.Sprintf("Empty Content For Language:%s", v.Language))
			}
		}
	}
	for _, v := range template.LocalizationVersions {
		if one == nil || v != one {
			v.Activate = false
		}
	}
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, template.LocalizationVersions)
	if err != nil {
		return nil, err
	}

	return &email.ActivateLocalizationVersionRes{}, nil
}
