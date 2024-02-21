package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type EmailTemplateVo struct {
	Id                  int64  `json:"id"                 description:""`                //
	MerchantId          int64  `json:"merchantId"         description:""`                //
	TemplateName        string `json:"templateName"       description:""`                //
	TemplateDescription string `json:"templateDescription" description:""`               //
	TemplateTitle       string `json:"templateTitle"      description:""`                //
	TemplateContent     string `json:"templateContent"    description:""`                //
	TemplateAttachName  string `json:"templateAttachName" description:""`                //
	CreateTime          int64  `json:"createTime"         description:"create utc time"` // create utc time
	UpdateTime          int64  `json:"updateTime"         description:"update utc time"` // create utc time
	Status              string `json:"status"             description:""`                //
}

func convertMerchantEmailTemplateToVo(emailTemplate *entity.MerchantEmailTemplate) *EmailTemplateVo {
	var status = "Active"
	if emailTemplate.Status != 0 {
		status = "InActive"
	}
	return &EmailTemplateVo{
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

func convertEmailDefaultTemplateToVo(emailTemplate *entity.EmailDefaultTemplate) *EmailTemplateVo {
	return &EmailTemplateVo{
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

func GetEmailDefaultTemplateByTemplateName(ctx context.Context, templateName string) *EmailTemplateVo {
	if len(templateName) == 0 {
		return nil
	}
	var one *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(entity.EmailDefaultTemplate{TemplateName: templateName}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	if one != nil {
		return convertEmailDefaultTemplateToVo(one)
	} else {
		return nil
	}
}

func GetMerchantEmailTemplateByTemplateName(ctx context.Context, merchantId int64, templateName string) *EmailTemplateVo {
	if len(templateName) == 0 || merchantId <= 0 {
		return nil
	}
	var one *entity.MerchantEmailTemplate
	err := dao.MerchantEmailTemplate.Ctx(ctx).
		Where(entity.MerchantEmailTemplate{TemplateName: templateName}).
		Where(entity.MerchantEmailTemplate{MerchantId: merchantId}).
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
