package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

//https://cryptadium.com/
//https://cryptadium.gitbook.io/cryptadium-api/basics/editor/fiat-payment

type Cryptadium struct {
}

func (c Cryptadium) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Cryptadium",
		Description:                   "Use public and private keys to secure the crypto payment.",
		DisplayName:                   "Cryptadium",
		GatewayWebsiteLink:            "https://cryptadium.com/",
		GatewayWebhookIntegrationLink: "https://cryptadium.gitbook.io/cryptadium-api/webhooks",
		GatewayLogo:                   "https://api.unibee.top/oss/file/d76q5bxsotbt0uzajb.png",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d6yhnz0wty7w6m7zhd.svg", "https://api.unibee.top/oss/file/d6yho8slal03ywl65c.svg", "https://api.unibee.top/oss/file/d6yhoilcikizou9ztk.svg", "https://api.unibee.top/oss/file/d6yhotsmefitw0cav1.svg"},
		GatewayType:                   consts.GatewayTypeCrypto,
	}
}

func (c Cryptadium) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   0,
		CryptoCurrency: "USDT",
		Rate:           0,
	}, nil
}

func (c Cryptadium) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	urlPath := "/api/v1/payment/pages/fiat"
	param := map[string]interface{}{
		"FiatCurrency": "USD",
		"AmountFiat":   999,
		"Currency":     "USDT",
		"Amount":       0,
		"Name":         "Cristianot",
		"Lastname":     "Ronaldo",
		"BillingId":    uuid.New().String(),
		"ClientId":     "201",
		"Email":        "mail@hotmail.com",
	}
	responseJson, err := SendCryptadiumPaymentRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("id"), "invalid keys, id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c Cryptadium) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/api/v1/payment/pages/fiat"
	var gasPayer string
	if createPayContext.Pay.GasPayer == "merchant" {
		gasPayer = "MERCHANT"
	} else {
		gasPayer = "CUSTOMER"
	}
	var name = ""
	var description = ""
	if len(createPayContext.Invoice.Lines) > 0 {
		var line = createPayContext.Invoice.Lines[0]
		if len(line.Name) == 0 {
			name = line.Description
		} else {
			name = line.Name
			description = line.Description
		}
	}

	param := map[string]interface{}{
		"Currency":             createPayContext.Pay.CryptoCurrency,
		"Amount":               0,
		"FiatCurrency":         createPayContext.Pay.Currency,
		"AmountFiat":           utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		"CheckBilling":         true,
		"Name":                 name,
		"Lastname":             description,
		"BillingId":            createPayContext.Pay.PaymentId,
		"ClientId":             strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"Email":                createPayContext.Email,
		"success_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"failure_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		"fees_payer":           gasPayer, // who pay the fee
		"payment_data":         createPayContext.Metadata,
		"pending_deadline_at":  time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	responseJson, err := SendCryptadiumPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "CryptadiumNewPayment", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("id") {
		return nil, gerror.New("invalid request, id is nil")
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := ""
	if responseJson.Contains("id") {
		gatewayPaymentId = responseJson.Get("id").String()
	}
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   responseJson.Get("paymentLink").String(),
	}, nil
}

func (c Cryptadium) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c Cryptadium) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api/v1/PaymentByBillingId/" + payment.PaymentId
	param := map[string]interface{}{}
	responseJson, err := SendCryptadiumPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "CryptadiumPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return parseCryptadiumPayment(responseJson), nil
}

func (c Cryptadium) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Cryptadium) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (c Cryptadium) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:       strconv.FormatUint(payment.MerchantId, 10),
		GatewayRefundId:  refund.GatewayRefundId,
		GatewayPaymentId: payment.GatewayPaymentId,
		Status:           consts.RefundCancelled,
		Reason:           refund.RefundComment,
		RefundAmount:     refund.RefundAmount,
		Currency:         refund.Currency,
		RefundTime:       gtime.Now(),
	}, nil
}

func parseCryptadiumPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if strings.Compare(item.Get("state").String(), "Success") == 0 && item.Get("Insufficient").String() == "false" {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("state").String(), "Canceled") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("state").String(), "Failed") == 0 {
		status = consts.PaymentFailed
	}
	var authorizeReason = ""
	var paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("Amount").String(), item.Get("Currency").String())

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: item.Get("PaymentId").String(),
		Status:           status,
		AuthorizeStatus:  authorizeStatus,
		AuthorizeReason:  authorizeReason,
		CancelReason:     "",
		PaymentData:      item.String(),
		TotalAmount:      utility.ConvertDollarStrToCent(item.Get("AmountFiat").String(), item.Get("FiatCurrency").String()),
		PaymentAmount:    paymentAmount,
		PaidTime:         gtime.Now(),
	}
}

func SendCryptadiumPaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")
	if !config.GetConfigInstance().IsProd() {
		urlPath = strings.ReplaceAll(urlPath, "payment/pages", "sandbox")
	}
	param["ShopId"] = publicKey
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nCryptadium_Start %s %s %s %s\n", method, "https://dashboard.cryptadium.com"+urlPath, privateKey, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type": "application/json",
		"Key":          privateKey,
	}

	response, err := utility.SendRequest("https://dashboard.cryptadium.com"+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nCryptadium_End %s %s response: %s error %s\n", method, "https://dashboard.cryptadium.com"+urlPath, response, err)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}
