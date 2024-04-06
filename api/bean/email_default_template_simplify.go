package bean

import entity "unibee/internal/model/entity/oversea_pay"

type EmailDefaultTemplateSimplify struct {
	Id                  int64  `json:"id"                 description:""`
	TemplateName        string `json:"templateName"       description:""`
	TemplateDescription string `json:"templateDescription" description:""`
	TemplateTitle       string `json:"templateTitle"      description:""`
	TemplateContent     string `json:"templateContent"    description:""`
	TemplateAttachName  string `json:"templateAttachName" description:""`
}

func SimplifyEmailDefaultTemplate(emailTemplate *entity.EmailDefaultTemplate) *EmailDefaultTemplateSimplify {
	if emailTemplate == nil {
		return nil
	}
	return &EmailDefaultTemplateSimplify{
		Id:                  emailTemplate.Id,
		TemplateName:        emailTemplate.TemplateName,
		TemplateDescription: emailTemplate.TemplateDescription,
		TemplateTitle:       emailTemplate.TemplateTitle,
		TemplateContent:     emailTemplate.TemplateContent,
		TemplateAttachName:  emailTemplate.TemplateAttachName,
	}
}
