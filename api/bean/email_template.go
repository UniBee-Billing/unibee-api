package bean

import entity "unibee/internal/model/entity/default"

type MerchantEmailTemplate struct {
	Id                  int64  `json:"id"                 description:""`                //
	MerchantId          uint64 `json:"merchantId"         description:""`                //
	TemplateName        string `json:"templateName"       description:""`                //
	TemplateDescription string `json:"templateDescription" description:""`               //
	TemplateTitle       string `json:"templateTitle"      description:""`                //
	TemplateContent     string `json:"templateContent"    description:""`                //
	TemplateAttachName  string `json:"templateAttachName" description:""`                //
	CreateTime          int64  `json:"createTime"         description:"create utc time"` // create utc time
	UpdateTime          int64  `json:"updateTime"         description:"update utc time"` // create utc time
	Status              string `json:"status"             description:""`                //
	GatewayTemplateId   string `json:"gatewayTemplateId"  description:""`                //
	LanguageData        string `json:"languageData"       description:""`                //
}

func SimplifyMerchantEmailTemplate(emailTemplate *entity.MerchantEmailTemplate) *MerchantEmailTemplate {
	var status = "Active"
	if emailTemplate.Status != 0 {
		status = "InActive"
	}
	return &MerchantEmailTemplate{
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
		GatewayTemplateId:   emailTemplate.GatewayTemplateId,
		LanguageData:        emailTemplate.LanguageData,
	}
}
