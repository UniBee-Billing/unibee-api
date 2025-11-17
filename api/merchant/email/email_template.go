package email

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TemplateListReq struct {
	g.Meta `path:"/template_list" tags:"Email Template" method:"get" summary:"Get Email Template List"`
}

type TemplateListRes struct {
	EmailTemplateList []*bean.MerchantEmailTemplate `json:"emailTemplateList" description:"Email Template Object List" `
	Total             int                           `json:"total" dc:"Total"`
}

type AddLocalizationVersionReq struct {
	g.Meta        `path:"/template_add_localization_version" tags:"Email Template" method:"post" summary:"Add Email Template Localization Version"`
	TemplateName  string                            `json:"templateName" dc:"Template Name" required:"true"`
	Localizations []*bean.EmailLocalizationTemplate `json:"localizations" description:"" required:"true"`
}

type AddLocalizationVersionRes struct {
	LocalizationVersion *bean.MerchantLocalizationVersion `json:"localizationVersion" description:""`
}

type EditLocalizationVersionReq struct {
	g.Meta        `path:"/template_edit_localization_version" tags:"Email Template" method:"post" summary:"Edit Email Template Localization Version"`
	TemplateName  string                            `json:"templateName" dc:"Template Name" required:"true"`
	VersionId     string                            `json:"versionId" description:"" required:"true"`
	Localizations []*bean.EmailLocalizationTemplate `json:"localizations" description:"" required:"true"`
}

type EditLocalizationVersionRes struct {
	LocalizationVersion *bean.MerchantLocalizationVersion `json:"localizationVersion" description:""`
}

type ActivateLocalizationVersionReq struct {
	g.Meta       `path:"/template_activate_localization_version" tags:"Email Template" method:"post" summary:"Activate Email Template Localization Version"`
	TemplateName string `json:"templateName" dc:"Template Name" required:"true"`
	VersionId    string `json:"versionId" description:""`
}

type ActivateLocalizationVersionRes struct {
}

type DeleteLocalizationVersionReq struct {
	g.Meta       `path:"/template_delete_localization_version" tags:"Email Template" method:"post" summary:"Delete Email Template Localization Version"`
	TemplateName string `json:"templateName" dc:"Template Name" required:"true"`
	VersionId    string `json:"versionId" description:""`
}

type DeleteLocalizationVersionRes struct {
}

type CustomizeLocalizationTemplateSyncReq struct {
	g.Meta        `path:"/custom_localization_template_sync" tags:"Email Template" method:"post" summary:"Customize Localization Template Sync" dc:"Sync the custom localization email template to gateway (sendgrid)"`
	TemplateData  []bean.MerchantLocalizationEmailTemplate `json:"templateData" dc:"TemplateData" v:"required"`
	VersionEnable bool                                     `json:"versionEnable" dc:"VersionEnable"`
}

type CustomizeLocalizationTemplateSyncRes struct {
}
