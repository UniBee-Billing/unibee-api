package method

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/query"
	"unibee/utility"
)

type NewPaymentMethodInternalReq struct {
	MerchantId     uint64                 `json:"merchantId" dc:"MerchantId" `
	UserId         uint64                 `json:"userId" dc:"UserId" `
	GatewayId      uint64                 `json:"gatewayId" dc:"GatewayId" `
	Currency       string                 `json:"currency" dc:""`
	RedirectUrl    string                 `json:"redirectUrl" dc:"Redirect Url"`
	SubscriptionId string                 `json:"subscriptionId" dc:"SubscriptionId"`
	Type           string                 `json:"type"`
	Metadata       map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func NewPaymentMethod(ctx context.Context, req *NewPaymentMethodInternalReq) (url string, paymentMethod *bean.PaymentMethod) {
	merchant := query.GetMerchantById(ctx, req.MerchantId)
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	req.Currency = strings.ToUpper(req.Currency)
	if req.Metadata == nil {
		req.Metadata = map[string]interface{}{}
	}
	if len(req.RedirectUrl) > 0 && len(req.SubscriptionId) > 0 {
		if !strings.Contains(req.RedirectUrl, "?") {
			req.RedirectUrl = fmt.Sprintf("%s?subId=%s", req.RedirectUrl, req.SubscriptionId)
		} else {
			req.RedirectUrl = fmt.Sprintf("%s&subId=%s", req.RedirectUrl, req.SubscriptionId)
		}
	}
	req.Metadata["RedirectUrl"] = req.RedirectUrl
	req.Metadata["SubscriptionId"] = req.SubscriptionId
	req.Metadata["MerchantId"] = req.MerchantId
	createResult, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserCreateAndBindPaymentMethod(ctx, gateway, req.UserId, req.Currency, req.Metadata)
	utility.AssertError(err, "Server Error")
	return createResult.Url, &bean.PaymentMethod{
		Id:        createResult.PaymentMethod.Id,
		Type:      createResult.PaymentMethod.Type,
		IsDefault: createResult.PaymentMethod.IsDefault,
		Data:      createResult.PaymentMethod.Data,
	}
}

func DeletePaymentMethod(ctx context.Context, merchantId uint64, userId uint64, gatewayId uint64, paymentMethodId string) error {
	merchant := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(len(paymentMethodId) > 0, "invalid paymentMethodId")
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	_, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayUserDeAttachPaymentMethodQuery(ctx, gateway, userId, paymentMethodId)
	return err
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
	if req.GatewayId <= 0 {
		return make([]*bean.PaymentMethod, 0)
	}
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
		g.Log().Errorf(ctx, "GatewayUserPaymentMethodListQuery error:%s", err.Error())
		return nil
	}
	if req.UserId > 0 {
		user := query.GetUserAccountById(ctx, req.UserId)
		if user != nil && listQuery != nil && len(listQuery.PaymentMethods) > 0 {
			for _, one := range listQuery.PaymentMethods {
				if one.Id == user.PaymentMethod && user.GatewayId == fmt.Sprintf("%d", req.GatewayId) {
					one.IsDefault = true
				}
			}
		}
	}
	var list []*bean.PaymentMethod
	for _, one := range listQuery.PaymentMethods {
		list = append(list, &bean.PaymentMethod{
			Id:        one.Id,
			Type:      one.Type,
			IsDefault: one.IsDefault,
			Data:      one.Data,
		})
	}
	return list
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
		return &bean.PaymentMethod{
			Id:        listQuery.PaymentMethods[0].Id,
			Type:      listQuery.PaymentMethods[0].Type,
			IsDefault: listQuery.PaymentMethods[0].IsDefault,
			Data:      listQuery.PaymentMethods[0].Data,
		}
	}
	return nil
}
