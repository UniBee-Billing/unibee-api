package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"unibee/internal/consts"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Coinbase struct {
}

func (c Coinbase) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	urlPath := "/charges/"
	param := map[string]interface{}{
		"pricing_type": "fixed_price",
		"name":         "test crypto payment",
		"description":  "test crypto payment description",
		"checkout_id":  uuid.New().String(),
		"local_price": map[string]string{
			"amount":   "10.00",
			"currency": "USD",
		},
	}
	responseJson, err := SendCoinbasePaymentRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("id"), "invalid keys, id is nil")
	return "", consts.GatewayTypeCrypto, nil
}

func (c Coinbase) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/charges/"
	var gasPayer string
	if createPayContext.Pay.GasPayer == "merchant" {
		return nil, gerror.New("Not Support")
	} else {
		gasPayer = "CUSTOMER"
		createPayContext.Metadata["gasPayer"] = gasPayer
	}
	param := map[string]interface{}{
		"pricing_type": "fixed_price",
		"name":         "crypto payment",
		"description":  "crypto payment description",
		"checkout_id":  createPayContext.Pay.PaymentId,
		"local_price": map[string]string{
			"amount":   utility.ConvertCentToDollarStr(createPayContext.Pay.CryptoAmount, createPayContext.Pay.CryptoCurrency),
			"currency": createPayContext.Pay.CryptoCurrency,
		},
		"redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"cancel_url\n": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		"metadata":     createPayContext.Metadata,
	}
	responseJson, err := SendCoinbasePaymentRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "CoinbaseNewPayment", nil, createPayContext.Gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("id") {
		return nil, gerror.New("invalid request, id is nil")
	}
	if err != nil {
		return nil, err
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := responseJson.Get("id").String()
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   responseJson.Get("hosted_url").String(),
	}, nil
}

func (c Coinbase) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c Coinbase) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	utility.Assert(from.Currency == "USD", "only USD support")
	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   from.Amount,
		CryptoCurrency: "USD",
		Rate:           1,
	}, nil
}

func (c Coinbase) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Coinbase) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundSuccess,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (c Coinbase) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func SendCoinbasePaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")
	datetime := getExpirationDateTime(1)

	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nCoinbase_Start %s %s %s %s %s\n", method, urlPath, publicKey, jsonString, datetime)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-CC-Api-Key": publicKey,
	}
	response, err := utility.SendRequest("https://api.commerce.coinbase.com"+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nCoinbase_End %s %s response: %s error %s\n", method, urlPath, response, err)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}
