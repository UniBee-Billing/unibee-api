package webhook

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type EventListReq struct {
	g.Meta `path:"/event_list" tags:"Webhook" method:"get" summary:"Webhook Event list"`
}

type EventListRes struct {
	EventList []string `json:"eventList" dc:"EventList"`
}

type EndpointListReq struct {
	g.Meta `path:"/endpoint_list" tags:"Webhook" method:"get" summary:"Merchant Webhook Endpoint list"`
}

type EndpointListRes struct {
	EndpointList []*bean.MerchantWebhookEndpointSimplify `json:"endpointList" dc:"EndpointList"`
}

type EndpointLogListReq struct {
	g.Meta     `path:"/endpoint_log_list" tags:"Webhook" method:"get" summary:"Merchant Webhook Endpoint Log list"`
	EndpointId int64 `json:"endpointId" dc:"EndpointId" v:"required"`
	Page       int   `json:"page" dc:"Page, Start WIth 0" `
	Count      int   `json:"count" dc:"Count Of Page" `
}

type EndpointLogListRes struct {
	EndpointLogList []*bean.MerchantWebhookLogSimplify `json:"endpointLogList" dc:"EndpointLogList"`
}

type NewEndpointReq struct {
	g.Meta `path:"/new_endpoint" tags:"Webhook" method:"post" summary:"Merchant New Webhook Endpoint"`
	Url    string   `json:"url" dc:"Url" v:"required"`
	Events []string `json:"events" dc:"Events"`
}

type NewEndpointRes struct {
}

type UpdateEndpointReq struct {
	g.Meta     `path:"/update_endpoint" tags:"Webhook" method:"post" summary:"Merchant Update Webhook Endpoint"`
	EndpointId uint64   `json:"endpointId" dc:"EndpointId" v:"required"`
	Url        string   `json:"url" dc:"Url To Update" v:"required"`
	Events     []string `json:"events" dc:"Events To Update"`
}

type UpdateEndpointRes struct {
}

type DeleteEndpointReq struct {
	g.Meta     `path:"/delete_endpoint" tags:"Webhook" method:"post" summary:"Merchant Delete Webhook Endpoint"`
	EndpointId uint64 `json:"endpointId" dc:"EndpointId" v:"required"`
}

type DeleteEndpointRes struct {
}
