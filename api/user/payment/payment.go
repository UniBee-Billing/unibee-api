package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type NewReq struct {
	g.Meta      `path:"/new" tags:"User-Payment" method:"post" summary:"NewPayment"`
	PlanId      uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity    int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId   uint64                 `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl string                 `json:"redirectUrl" dc:"Redirect Url"`
	Metadata    map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type NewRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"The unique id of payment"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"The external unique id of payment"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}
