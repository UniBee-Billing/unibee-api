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

func (c *ControllerEmail) DeleteLocalizationVersion(ctx context.Context, req *email.DeleteLocalizationVersionReq) (res *email.DeleteLocalizationVersionRes, err error) {
	utility.Assert(len(req.TemplateName) > 0, "Invalid template name")
	utility.Assert(len(req.VersionId) > 0, "Invalid versionId")
	template := query.GetMerchantEmailTemplateByTemplateName(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	utility.Assert(template != nil, "template not found")

	var one *bean.MerchantLocalizationVersion
	var list []*bean.MerchantLocalizationVersion
	for _, v := range template.LocalizationVersions {
		if req.VersionId == v.VersionId {
			one = v
		} else {
			list = append(list, v)
		}
	}
	utility.Assert(one != nil, "Invalid localization versionId")
	utility.Assert(!one.Activate, "Cannot delete active localization version")
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, list)
	if err != nil {
		return nil, err
	}

	return &email.DeleteLocalizationVersionRes{}, nil
}
