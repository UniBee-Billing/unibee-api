package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	merchantPaymentApi "unibee/api/merchant/payment"
)

type NewReq struct {
	g.Meta      `path:"/new" tags:"User-Payment" method:"post" summary:"NewPayment"`
	Currency    string                     `json:"currency" dc:"Currency, either Currency&TotalAmount or PlanId needed" `
	TotalAmount int64                      `json:"totalAmount" dc:"Total PaymentAmount, Cent, either TotalAmount&Currency or PlanId needed"`
	GatewayId   uint64                     `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl string                     `json:"redirectUrl" dc:"Redirect Url"`
	CountryCode string                     `json:"countryCode" dc:"CountryCode"`
	Name        string                     `json:"name" dc:"Name"`
	Description string                     `json:"description" dc:"Description"`
	Items       []*merchantPaymentApi.Item `json:"items" dc:"Items"`
	Metadata    map[string]interface{}     `json:"metadata" dc:"Metadata，Map"`
}

type NewRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"The unique id of payment"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"The external unique id of payment"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}
