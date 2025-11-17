package email

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetMerchantEmailTemplateList(ctx context.Context, merchantId uint64) ([]*bean.MerchantEmailTemplate, int) {
	var list = make([]*bean.MerchantEmailTemplate, 0)
	if merchantId > 0 {
		var defaultTemplateList []*entity.EmailDefaultTemplate
		err := dao.EmailDefaultTemplate.Ctx(ctx).
			Scan(&defaultTemplateList)
		if err == nil && len(defaultTemplateList) > 0 {
			for _, emailTemplate := range defaultTemplateList {
				var merchantEmailTemplate *entity.MerchantEmailTemplate
				err = dao.MerchantEmailTemplate.Ctx(ctx).
					Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
					Where(dao.MerchantEmailTemplate.Columns().TemplateName, emailTemplate.TemplateName).
					Scan(&merchantEmailTemplate)
				vo := &bean.MerchantEmailTemplate{
					Id:                  emailTemplate.Id,
					MerchantId:          merchantId,
					TemplateName:        emailTemplate.TemplateName,
					TemplateDescription: emailTemplate.TemplateDescription,
					TemplateTitle:       emailTemplate.TemplateTitle,
					TemplateContent:     emailTemplate.TemplateContent,
					TemplateAttachName:  "", //pdf not customised here
					CreateTime:          emailTemplate.CreateTime,
					UpdateTime:          emailTemplate.GmtModify.Timestamp(),
					Status:              "Active", // default template status should be active
				}
				var languageData = make([]*bean.EmailLocalizationTemplate, 0)
				var languageVersionData = make([]*bean.MerchantLocalizationVersion, 0)
				if err == nil && merchantEmailTemplate != nil {
					if merchantEmailTemplate.Status == 0 {
						vo.Status = "Active"
					} else {
						vo.Status = "InActive"
					}
					if len(merchantEmailTemplate.TemplateDescription) > 0 {
						vo.TemplateDescription = merchantEmailTemplate.TemplateDescription
					}
					if len(merchantEmailTemplate.TemplateTitle) > 0 {
						vo.TemplateTitle = merchantEmailTemplate.TemplateTitle
					}
					if len(merchantEmailTemplate.TemplateContent) > 0 {
						vo.TemplateContent = merchantEmailTemplate.TemplateContent
					}
					if len(merchantEmailTemplate.LanguageData) > 0 {
						_ = utility.UnmarshalFromJsonString(merchantEmailTemplate.LanguageData, &languageData)
					}

					if len(merchantEmailTemplate.LanguageVersionData) > 0 {
						_ = utility.UnmarshalFromJsonString(merchantEmailTemplate.LanguageVersionData, &languageVersionData)
					}
					vo.CreateTime = merchantEmailTemplate.CreateTime
					vo.UpdateTime = merchantEmailTemplate.GmtModify.Timestamp()
				}
				vo.LanguageData = languageData
				vo.LocalizationVersions = languageVersionData
				list = append(list, vo)
			}
		}
	}
	return list, len(list)
}

func UpdateMerchantEmailTemplate(ctx context.Context, merchantId uint64, templateName string, languageVersionData []*bean.MerchantLocalizationVersion) error {
	utility.Assert(merchantId > 0, "Invalid MerchantId")
	utility.Assert(len(templateName) > 0, "Invalid TemplateName")
	var defaultTemplate *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&defaultTemplate)
	utility.AssertError(err, "Server Error")
	utility.Assert(defaultTemplate != nil, "Default Template Not Found")
	var one *entity.MerchantEmailTemplate
	err = dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	utility.AssertError(err, "Server Error")
	var languageData *bean.MerchantLocalizationVersion
	if languageVersionData != nil {
		for _, v := range languageVersionData {
			if v.Activate {
				languageData = v
				break
			}
		}
		if languageData != nil {
			for _, v := range languageVersionData {
				if v != languageData {
					v.Activate = false
				}
			}
		}
	}
	if one == nil {
		//insert
		one = &entity.MerchantEmailTemplate{
			MerchantId:          merchantId,
			TemplateName:        defaultTemplate.TemplateName,
			TemplateDescription: defaultTemplate.TemplateDescription,
			TemplateTitle:       defaultTemplate.TemplateTitle,
			TemplateContent:     defaultTemplate.TemplateContent,
			TemplateAttachName:  defaultTemplate.TemplateAttachName,
			LanguageData:        utility.MarshalToJsonString(languageData),
			LanguageVersionData: utility.MarshalToJsonString(languageVersionData),
			CreateTime:          gtime.Now().Timestamp(),
			Status:              0,
		}
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
		return err
	} else {
		//update
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailTemplate.Columns().MerchantId:          merchantId,
			dao.MerchantEmailTemplate.Columns().TemplateName:        defaultTemplate.TemplateName,
			dao.MerchantEmailTemplate.Columns().TemplateDescription: defaultTemplate.TemplateDescription,
			dao.MerchantEmailTemplate.Columns().TemplateTitle:       defaultTemplate.TemplateTitle,
			dao.MerchantEmailTemplate.Columns().TemplateContent:     defaultTemplate.TemplateContent,
			dao.MerchantEmailTemplate.Columns().TemplateAttachName:  defaultTemplate.TemplateAttachName,
			dao.MerchantEmailTemplate.Columns().LanguageData:        utility.MarshalToJsonString(languageData),
			dao.MerchantEmailTemplate.Columns().LanguageVersionData: utility.MarshalToJsonString(languageVersionData),
			dao.MerchantEmailTemplate.Columns().GmtModify:           gtime.Now(),
			dao.MerchantEmailTemplate.Columns().Status:              0,
		}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		return err
	}
}
