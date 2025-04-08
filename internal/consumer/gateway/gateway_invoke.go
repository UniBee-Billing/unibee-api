package gateway

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean/detail"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/service"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func init() {
	redismq.RegisterInvoke("SetupMerchantGateway", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "SetupMerchantGateway:%s", request)
		var req *SetupGatewayInvokeReq
		err = utility.UnmarshalFromJsonString(request.(string), &req)
		if err != nil {
			return nil, err
		}
		gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, req.GatewayName).GatewayInfo(ctx)
		utility.Assert(gatewayInfo != nil, "gateway not ready")
		var one *entity.MerchantGateway
		if req.GatewayName == "wire_transfer" {
			one = service.SetupWireTransferGateway(ctx, &service.WireTransferSetupReq{
				MerchantId:    req.MerchantId,
				Currency:      req.Currency,
				MinimumAmount: req.MinimumAmount,
				Bank:          req.Bank,
				DisplayName:   req.DisplayName,
				GatewayIcon:   req.GatewayIcons,
				Sort:          req.Sort,
			})
		} else {
			one = service.SetupGateway(ctx, req.MerchantId, req.GatewayName, req.GatewayKey, req.GatewaySecret, req.SubGateway, req.GatewayPaymentTypes, req.DisplayName, req.GatewayIcons, req.Sort, req.CurrencyExchange)
			if one != nil && len(req.WebhookSecret) > 0 {
				utility.Assert(one.MerchantId == req.MerchantId, "merchant not match")
				gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
				if len(gatewayInfo.GatewayWebhookIntegrationLink) > 0 && len(req.WebhookSecret) > 0 {
					err = query.UpdateGatewayWebhookSecret(ctx, one.Id, req.WebhookSecret)
					if err != nil {
						return nil, err
					}
				}
			}
			if one != nil {
				one = query.GetGatewayById(ctx, one.Id)
			}
		}

		return one, nil
	})
}

type SetupGatewayInvokeReq struct {
	MerchantId          uint64                            `json:"merchantId"  dc:"The id of merchant"`
	GatewayName         string                            `json:"gatewayName"  dc:"The name of payment gateway, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
	DisplayName         *string                           `json:"displayName"  dc:"The displayName of payment gateway"`
	GatewayIcons        *[]string                         `json:"gatewayIcons"  dc:"The icons of payment gateway"`
	Sort                *int64                            `json:"sort"  dc:"The sort value of payment gateway, The bigger, the closer to the front"`
	GatewayKey          string                            `json:"gatewayKey"  dc:"The key of payment gateway" `
	GatewaySecret       string                            `json:"gatewaySecret"  dc:"The secret of payment gateway" `
	SubGateway          string                            `json:"subGateway"  dc:"The sub gateway of payment gateway" `
	CurrencyExchange    []*detail.GatewayCurrencyExchange `json:"currencyExchange" dc:"The currency exchange for gateway payment, effect at start of payment creation when currency matched"`
	GatewayPaymentTypes []string                          `json:"gatewayPaymentTypes"  dc:"Selected gateway payment types"`
	WebhookSecret       string                            `json:"webhookSecret"  dc:"The secret of gateway webhook"`
	Currency            string                            `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount       int64                             `json:"minimumAmount"   dc:"The minimum amount of wire transfer" v:"required" `
	Bank                *detail.GatewayBank               `json:"bank"   dc:"The receiving bank of wire transfer " v:"required" `
}
