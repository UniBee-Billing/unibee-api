package bean

import (
	"context"
	"strings"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type EmailTemplateVariable struct {
	InvoiceId             string      `json:"InvoiceId" group:"Invoice Information"`
	UserName              string      `json:"User name" key:"UserName" group:"User Information"`
	MerchantProductName   string      `json:"Merchant Product Name" key:"ProductName" group:"Subscription Information"`
	MerchantCustomerEmail string      `json:"Merchant’s customer support email address" key:"SupportEmail" group:"Company Information"`
	MerchantName          string      `json:"Merchant Name" key:"MerchantName" group:"Company Information"`
	DateNow               *gtime.Time `json:"DateNow" layout:"2006-01-02" group:"Subscription Information"`
	PeriodEnd             *gtime.Time `json:"PeriodEnd" layout:"2006-01-02" group:"Subscription Information"`
	PaymentAmount         string      `json:"Payment Amount" key:"PaymentAmount" group:"Invoice Information"`
	RefundAmount          string      `json:"Refund Amount" key:"RefundAmount" group:"Invoice Information"`
	Currency              string      `json:"Currency" key:"Currency" group:"Invoice Information"`
	TokenExpireMinute     string      `json:"TokenExpireMinute" group:"Authentication Information"`
	CodeExpireMinute      string      `json:"CodeExpireMinute" group:"Authentication Information"`
	Code                  string      `json:"Code" group:"Authentication Information"`
	Link                  string      `json:"Link" group:"Invoice Information"`
	HttpLink              string      `json:"HttpLink" group:"Invoice Information"`
	AccountHolder         string      `json:"Account Holder" key:"WireTransferAccountHolder"`
	Address               string      `json:"Address" key:"WireTransferAddress" group:"Company Information"`
	BIC                   string      `json:"BIC" key:"WireTransferBIC" group:"Company Information"`
	IBAN                  string      `json:"IBAN" key:"WireTransferIBAN" group:"Company Information"`
	BankData              string      `json:"Bank Data" key:"WireTransferBankData" group:"Company Information"`
}

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
	VariableGroups       []*TemplateVariableGroup       `json:"VariableGroups" description:""`
}

type TemplateVariableGroup struct {
	GroupName string              `json:"groupName"                 description:""`
	Variables []*TemplateVariable `json:"variables"            description:""`
}

type TemplateVariable struct {
	VariableName string `json:"variableName"                 description:""`
}

func (t *MerchantEmailTemplate) LocalizationSubject(language string, languageData *[]*EmailLocalizationTemplate) (subject string) {
	if len(language) == 0 {
		language = "en" // default language
	}
	targetLanguageData := make([]*EmailLocalizationTemplate, 0)
	if languageData != nil {
		targetLanguageData = *languageData
	} else {
		targetLanguageData = t.LanguageData
	}
	if len(targetLanguageData) == 0 || len(language) == 0 {
		return t.TemplateTitle
	}
	for _, one := range targetLanguageData {
		if one.Language == language {
			return one.Title
		} else if strings.ToLower(one.Language) == "cn" && strings.ToLower(language) == "zh" {
			return one.Title
		} else if strings.ToLower(one.Language) == "zh" && strings.ToLower(language) == "cn" {
			return one.Title
		}
	}
	return t.TemplateTitle
}

func (t *MerchantEmailTemplate) LocalizationContent(language string, languageData *[]*EmailLocalizationTemplate) (content string) {
	if len(language) == 0 {
		language = "en" // default language
	}
	targetLanguageData := make([]*EmailLocalizationTemplate, 0)
	if languageData != nil {
		targetLanguageData = *languageData
	} else {
		targetLanguageData = t.LanguageData
	}
	if len(targetLanguageData) == 0 || len(language) == 0 {
		return t.TemplateContent
	}
	for _, one := range targetLanguageData {
		if one.Language == language {
			return one.Content
		} else if strings.ToLower(one.Language) == "cn" && strings.ToLower(language) == "zh" {
			return one.Content
		} else if strings.ToLower(one.Language) == "zh" && strings.ToLower(language) == "cn" {
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
		err := utility.UnmarshalFromJsonString(emailTemplate.LanguageData, &languageData)
		if err != nil {
			g.Log().Errorf(context.Background(), "error:%s", err.Error())
		}
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
	VersionName   string                       `json:"versionName"       description:""`
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
