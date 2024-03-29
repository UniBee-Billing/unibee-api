package method

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	"unibee/api/bean"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/query"
	"unibee/utility"
)

type NewPaymentMethodInternalReq struct {
	MerchantId uint64      `json:"merchantId" dc:"MerchantId" `
	UserId     uint64      `json:"userId" dc:"UserId" `
	GatewayId  uint64      `json:"gatewayId" dc:"GatewayId" `
	Type       string      `json:"type"`
	Data       *gjson.Json `json:"data"`
}

func NewPaymentMethod(ctx context.Context, req *NewPaymentMethodInternalReq) *bean.PaymentMethod {
	merchant := query.GetMerchantById(ctx, req.MerchantId)
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	utility.Assert(len(req.Type) > 0 && strings.Compare(req.Type, "card") == 0, "invalid type, should be card")
	createResult, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserCreateAndBindPaymentMethod(ctx, gateway, req.UserId, req.Data)
	utility.AssertError(err, "Server Error")
	return createResult.PaymentMethod
}

type PaymentMethodListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" `
	UserId     uint64 `json:"userId" dc:"UserId" `
	GatewayId  uint64 `json:"gatewayId" dc:"GatewayId" `
	PaymentId  string `json:"paymentId" dc:"PaymentId"  `
}

func QueryPaymentMethodList(ctx context.Context, req *PaymentMethodListInternalReq) []*bean.PaymentMethod {
	merchant := query.GetMerchantById(ctx, req.MerchantId)
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	var gatewayPaymentId string
	if len(req.PaymentId) > 0 {
		one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
		if one != nil {
			gatewayPaymentId = one.GatewayPaymentId
		}
	}
	listQuery, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
		UserId:           req.UserId,
		GatewayPaymentId: gatewayPaymentId,
	})
	if err != nil {
		return nil
	}
	return listQuery.PaymentMethods
}

func QueryPaymentMethod(ctx context.Context, merchantId uint64, userId uint64, gatewayId uint64, gatewayPaymentMethodId string) *bean.PaymentMethod {
	merchant := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	listQuery, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
		UserId:                 userId,
		GatewayPaymentMethodId: gatewayPaymentMethodId,
	})
	if err != nil {
		return nil
	}
	if listQuery != nil && len(listQuery.PaymentMethods) == 1 {
		return listQuery.PaymentMethods[0]
	}
	return nil
}
