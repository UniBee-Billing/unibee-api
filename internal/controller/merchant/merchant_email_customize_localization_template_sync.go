package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) CustomizeLocalizationTemplateSync(ctx context.Context, req *email.CustomizeLocalizationTemplateSyncReq) (res *email.CustomizeLocalizationTemplateSyncRes, err error) {
	for _, template := range req.TemplateData {
		// check default language
		var enTemplate *bean.EmailLocalizationTemplate
		for _, lang := range template.Localizations {
			if lang.Language == "en" {
				enTemplate = lang
			}
		}
		utility.Assert(enTemplate != nil, "en not found")
		var one *entity.MerchantEmailTemplate
		err = dao.MerchantEmailTemplate.Ctx(ctx).
			Where(dao.MerchantEmailTemplate.Columns().MerchantId, _interface.GetMerchantId(ctx)).
			Where(dao.MerchantEmailTemplate.Columns().TemplateName, template.TemplateName).
			Scan(&one)
		utility.AssertError(err, "Server Error")
		if one == nil {
			//insert
			one = &entity.MerchantEmailTemplate{
				MerchantId:          _interface.GetMerchantId(ctx),
				TemplateName:        template.TemplateName,
				TemplateDescription: template.TemplateDescription,
				TemplateAttachName:  template.Attach,
				TemplateTitle:       enTemplate.Title,
				TemplateContent:     enTemplate.Content,
				LanguageData:        utility.MarshalToJsonString(template.Localizations),
				CreateTime:          gtime.Now().Timestamp(),
				Status:              0,
			}
			result, err := dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
			if err != nil {
				g.Log().Errorf(ctx, `record insert failure %s`, err.Error())
				continue
			}
			id, _ := result.LastInsertId()
			one.Id = id
		} else {
			//update
			_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
				dao.MerchantEmailTemplate.Columns().MerchantId:          _interface.GetMerchantId(ctx),
				dao.MerchantEmailTemplate.Columns().TemplateName:        template.TemplateName,
				dao.MerchantEmailTemplate.Columns().TemplateDescription: template.TemplateDescription,
				dao.MerchantEmailTemplate.Columns().TemplateAttachName:  template.Attach,
				dao.MerchantEmailTemplate.Columns().TemplateTitle:       enTemplate.Title,
				dao.MerchantEmailTemplate.Columns().TemplateContent:     enTemplate.Content,
				dao.MerchantEmailTemplate.Columns().LanguageData:        utility.MarshalToJsonString(template.Localizations),
				dao.MerchantEmailTemplate.Columns().GmtModify:           gtime.Now(),
				dao.MerchantEmailTemplate.Columns().Status:              0,
			}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		}
		//Sync to Gateway
		err = email2.SyncMerchantEmailTemplateToGateway(ctx, one.Id, req.VersionEnable)
		if err != nil {
			g.Log().Errorf(ctx, `SyncMerchantEmailTemplateToGateway error %s`, err.Error())
			continue
		}
	}

	return &email.CustomizeLocalizationTemplateSyncRes{}, nil
}
