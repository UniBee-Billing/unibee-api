package setup

import "github.com/gogf/gf/v2/frame/g"

type NewReq struct {
	g.Meta `path:"/new_webhook" tags:"Webhook-Controller" method:"post" summary:"New Webhook Endpoint"`
	Url    string   `p:"url" dc:"Url" v:"required"`
	Events []string `p:"events" dc:"Events"`
}

type NewRes struct {
}

type UpdateReq struct {
	g.Meta `path:"/update_webhook" tags:"Webhook-Controller" method:"post" summary:"Update Webhook Endpoint"`
	Url    string   `p:"url" dc:"Url" v:"required"`
	Events []string `p:"events" dc:"Events"`
}

type UpdateRes struct {
}
