package bean

import (
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type MerchantEmailTemplate struct {
	Id                   int64                          `json:"id"                 description:""`                //
	MerchantId           uint64                         `json:"merchantId"         description:""`                //
	TemplateName         string                         `json:"templateName"       description:""`                //
	TemplateDescription  string                         `json:"templateDescription" description:""`               //
	TemplateTitle        string                         `json:"templateTitle"      description:""`                //
	TemplateContent      string                         `json:"templateContent"    description:""`                //
	TemplateAttachName   string                         `json:"templateAttachName" description:""`                //
	CreateTime           int64                          `json:"createTime"         description:"create utc time"` // create utc time
	UpdateTime           int64                          `json:"updateTime"         description:"update utc time"` // create utc time
	Status               string                         `json:"status"             description:""`                //
	GatewayTemplateId    string                         `json:"gatewayTemplateId"  description:""`                //
	LanguageData         []*EmailLocalizationTemplate   `json:"languageData"       description:""`                //
	LocalizationVersions []*MerchantLocalizationVersion `json:"localizationVersions" description:""`
}

func (t *MerchantEmailTemplate) LocalizationSubject(language string) (subject string) {
	if len(language) == 0 {
		language = "en" // default language
	}
	if len(t.LanguageData) == 0 || len(language) == 0 {
		return t.TemplateTitle
	}
	for _, one := range t.LanguageData {
		if one.Language == language {
			return one.Title
		}
	}
	return t.TemplateTitle
}

func (t *MerchantEmailTemplate) LocalizationContent(language string) (content string) {
	if len(language) == 0 {
		language = "en" // default language
	}
	if len(t.LanguageData) == 0 || len(language) == 0 {
		return t.TemplateContent
	}
	for _, one := range t.LanguageData {
		if one.Language == language {
			return one.Content
		}
	}
	return t.TemplateContent
}

func SimplifyMerchantEmailTemplate(emailTemplate *entity.MerchantEmailTemplate) *MerchantEmailTemplate {
	var status = "Active"
	if emailTemplate.Status != 0 {
		status = "InActive"
	}
	var languageData = make([]*EmailLocalizationTemplate, 0)
	if len(emailTemplate.LanguageData) > 0 {
		_ = utility.UnmarshalFromJsonString(emailTemplate.LanguageData, &languageData)
	}
	var localizationVersions = make([]*MerchantLocalizationVersion, 0)
	if len(emailTemplate.LanguageVersionData) > 0 {
		_ = utility.UnmarshalFromJsonString(emailTemplate.LanguageVersionData, &localizationVersions)
	}
	return &MerchantEmailTemplate{
		Id:                   emailTemplate.Id,
		MerchantId:           emailTemplate.MerchantId,
		TemplateName:         emailTemplate.TemplateName,
		TemplateDescription:  "",
		TemplateTitle:        emailTemplate.TemplateTitle,
		TemplateContent:      emailTemplate.TemplateContent,
		TemplateAttachName:   emailTemplate.TemplateAttachName,
		CreateTime:           emailTemplate.CreateTime,
		UpdateTime:           emailTemplate.GmtModify.Timestamp(),
		Status:               status,
		GatewayTemplateId:    emailTemplate.GatewayTemplateId,
		LanguageData:         languageData,
		LocalizationVersions: localizationVersions,
	}
}

type MerchantLocalizationVersion struct {
	VersionId     string                       `json:"versionId"       description:""`
	Activate      bool                         `json:"activate"       description:""`
	Localizations []*EmailLocalizationTemplate `json:"localizations" description:""`
}

type MerchantLocalizationEmailTemplate struct {
	TemplateName        string                       `json:"templateName"       description:""`
	TemplateDescription string                       `json:"templateDescription" description:""`
	Attach              string                       `json:"attach"       description:""`
	Activate            bool                         `json:"activate"       description:""`
	Localizations       []*EmailLocalizationTemplate `json:"localizations" description:""`
}

type EmailLocalizationTemplate struct {
	Language string `json:"language"       description:""`
	Title    string `json:"title"       description:""`
	Content  string `json:"content"       description:""`
}
