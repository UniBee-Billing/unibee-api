package webhook

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/webhook"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type EventListReq struct {
	g.Meta `path:"/webhook_event_list" tags:"Merchant-Webhook-Controller" method:"get" summary:"Webhook Event list"`
}

type EventListRes struct {
	Events []string
}

type EndpointListReq struct {
	g.Meta     `path:"/webhook_endpoint_list" tags:"Merchant-Webhook-Controller" method:"get" summary:"Merchant Webhook Endpoint list"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}

type EndpointListRes struct {
	EndpointList []*webhook.MerchantWebhookEndpointVo
}

type EndpointLogListReq struct {
	g.Meta     `path:"/webhook_endpoint_log_list" tags:"Merchant-Webhook-Controller" method:"get" summary:"Merchant Webhook Endpoint Log list"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
	EndpointId int64 `p:"endpointId" dc:"EndpointId" v:"required"`
	Page       int   `p:"page" dc:"Page, Start WIth 0" `
	Count      int   `p:"count" dc:"Count Of Page" `
}

type EndpointLogListRes struct {
	EndpointLogList []*entity.MerchantWebhookLog
}

type NewEndpointReq struct {
	g.Meta     `path:"/new_webhook_endpoint" tags:"Merchant-Webhook-Controller" method:"post" summary:"Merchant New Webhook Endpoint"`
	MerchantId int64    `p:"merchantId" dc:"MerchantId" v:"required"`
	Url        string   `p:"url" dc:"Url" v:"required"`
	Events     []string `p:"events" dc:"Events"`
}

type NewEndpointRes struct {
}

type UpdateEndpointReq struct {
	g.Meta     `path:"/update_webhook_endpoint" tags:"Merchant-Webhook-Controller" method:"post" summary:"Merchant Update Webhook Endpoint"`
	MerchantId int64    `p:"merchantId" dc:"MerchantId" v:"required"`
	EndpointId int64    `p:"endpointId" dc:"EndpointId" v:"required"`
	Url        string   `p:"url" dc:"Url To Update" v:"required"`
	Events     []string `p:"events" dc:"Events To Update"`
}

type UpdateEndpointRes struct {
}

type DeleteEndpointReq struct {
	g.Meta     `path:"/delete_webhook_endpoint" tags:"Merchant-Webhook-Controller" method:"post" summary:"Merchant Delete Webhook Endpoint"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
	EndpointId int64 `p:"endpointId" dc:"EndpointId" v:"required"`
}

type DeleteEndpointRes struct {
}
