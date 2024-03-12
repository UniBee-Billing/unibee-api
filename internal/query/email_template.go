package query

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func convertMerchantEmailTemplateToVo(emailTemplate *entity.MerchantEmailTemplate) *bean.EmailTemplateVo {
	var status = "Active"
	if emailTemplate.Status != 0 {
		status = "InActive"
	}
	return &bean.EmailTemplateVo{
		Id:                  emailTemplate.Id,
		MerchantId:          emailTemplate.MerchantId,
		TemplateName:        emailTemplate.TemplateName,
		TemplateDescription: "",
		TemplateTitle:       emailTemplate.TemplateTitle,
		TemplateContent:     emailTemplate.TemplateContent,
		TemplateAttachName:  emailTemplate.TemplateAttachName,
		CreateTime:          emailTemplate.CreateTime,
		UpdateTime:          emailTemplate.GmtModify.Timestamp(),
		Status:              status,
	}
}

func convertEmailDefaultTemplateToVo(emailTemplate *entity.EmailDefaultTemplate) *bean.EmailTemplateVo {
	return &bean.EmailTemplateVo{
		Id:                  emailTemplate.Id,
		MerchantId:          0,
		TemplateName:        emailTemplate.TemplateName,
		TemplateDescription: emailTemplate.TemplateDescription,
		TemplateTitle:       emailTemplate.TemplateTitle,
		TemplateContent:     emailTemplate.TemplateContent,
		TemplateAttachName:  emailTemplate.TemplateAttachName,
		CreateTime:          emailTemplate.CreateTime,
		UpdateTime:          emailTemplate.GmtModify.Timestamp(),
		Status:              "Active",
	}
}

func GetEmailDefaultTemplateByTemplateName(ctx context.Context, templateName string) *bean.EmailTemplateVo {
	if len(templateName) == 0 {
		return nil
	}
	var one *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	if err != nil {
		one = nil
	}
	if one != nil {
		return convertEmailDefaultTemplateToVo(one)
	} else {
		return nil
	}
}

func GetMerchantEmailTemplateByTemplateName(ctx context.Context, merchantId uint64, templateName string) *bean.EmailTemplateVo {
	if len(templateName) == 0 || merchantId <= 0 {
		return nil
	}
	var one *entity.MerchantEmailTemplate
	err := dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	if one == nil {
		// No Setting, Get Default
		return GetEmailDefaultTemplateByTemplateName(ctx, templateName)
	}
	return convertMerchantEmailTemplateToVo(one)
}
