package email

import (
	"unibee/api/bean/detail"

	"github.com/gogf/gf/v2/frame/g"
)

type GatewaySetupReq struct {
	g.Meta      `path:"/gateway_setup" tags:"Email" method:"post" summary:"Email Gateway Setup"`
	GatewayName string `json:"gatewayName"  dc:"The name of email gateway, 'sendgrid' or other for future updates" v:"required"`
	Data        string `json:"data" dc:"The setup data of email gateway" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"Whether setup the gateway as default or not, default is true" `
}

type GatewaySetupRes struct {
	Data string `json:"data" dc:"Data" dc:"The hide star data"`
}

type SendTemplateEmailToUserReq struct {
	g.Meta          `path:"/send_template_email_to_user" tags:"Email" method:"post" summary:"Send Template Email To User"`
	TemplateName    string                 `json:"templateName" dc:"The name of email template"       v:"required"`
	UserId          int64                  `json:"userId" dc:"UserId" v:"required" `
	Variables       map[string]interface{} `json:"variables" dc:"Variables，Map"`
	AttachInvoiceId string                 `json:"attachInvoiceId" dc:"AttachInvoiceId"`
}

type SendTemplateEmailToUserRes struct {
}

type SendEmailToUserReq struct {
	g.Meta            `path:"/send_email_to_user" tags:"Email" method:"post" summary:"Send Email To User"`
	Email             string                 `json:"email" dc:"Email" v:"required" `
	GatewayTemplateId string                 `json:"gatewayTemplateId" dc:"GatewayTemplateId" `
	Variables         map[string]interface{} `json:"variables" dc:"Variables，Map"`
	Subject           string                 `json:"subject"`
	Content           string                 `json:"content"`
	AttachInvoiceId   string                 `json:"attachInvoiceId" dc:"AttachInvoiceId"`
}

type SendEmailToUserRes struct {
}

type SenderSetupReq struct {
	g.Meta  `path:"/email_sender_setup" tags:"Email" method:"post" summary:"Email Sender Setup"`
	Name    string `json:"name"  dc:"The name of email sender, like 'no-reply'" v:"required"`
	Address string `json:"address" dc:"The address of email sender, like 'no-reply@unibee.dev'" v:"required"`
}

type SenderSetupRes struct {
}

type HistoryListReq struct {
	g.Meta          `path:"/history_list" tags:"Email" method:"get" summary:"Get Email History List" dc:"Get email send history list"`
	SearchKey       string `json:"searchKey" dc:"Search Key, email or title"  `
	Email           string `json:"email" dc:"Filter Email"  `
	Status          []int  `json:"status" dc:"status, 0-pending, 1-success, 2-failure" `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
}

type HistoryListRes struct {
	EmailHistoryStatistics *detail.MerchantEmailHistoryStatistics `json:"emailHistoryStatistics"`
	EmailHistories         []*detail.MerchantEmailHistoryDetail   `json:"emailHistories" dc:"Email History Object List"`
	Total                  int                                    `json:"total" dc:"Total"`
}
