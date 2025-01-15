package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

//https://payssion.com/cn/docs/#api-reference-payment-request

type Payssion struct {
}

func (c Payssion) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Payssion",
		Description:                   "Need to replace Use ClientId and Secret to secure the payment",
		DisplayName:                   "Payssion",
		GatewayWebsiteLink:            "https://payssion.com",
		GatewayWebhookIntegrationLink: "https://www.payssion.com/account/app",
		GatewayLogo:                   "https://api.unibee.top/oss/file/d6yhr3m8mzbbgqla37.png",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d6yhr3m8mzbbgqla37.png"},
		GatewayType:                   consts.GatewayTypeCard,
	}
}

func (c Payssion) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (c Payssion) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	urlPath := "/api/v1/payments"
	pmID := "alipay_cn"
	param := map[string]interface{}{
		"currency":    "EUR",
		"pm_id":       pmID,
		"amount":      100,
		"description": "test payment description",
		"order_id":    uuid.New().String(),
		"payer_email": "jack.fu@wowow.io",
	}
	if !config.GetConfigInstance().IsProd() {
		param["pm_id"] = "payssion_test"
	}
	param["api_key"] = key
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v|%v", param["api_key"], param["pm_id"], param["amount"], param["currency"], param["order_id"], secret))
	responseJson, err := SendPayssionPaymentRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("paymentId"), "invalid keys, id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c Payssion) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/api/v1/payments"
	//var name = ""
	var description = ""
	if len(createPayContext.Invoice.Lines) > 0 {
		var line = createPayContext.Invoice.Lines[0]
		if len(line.Name) == 0 {
			//name = line.Description
		} else {
			//name = line.Name
			description = line.Description
		}
	}
	pmID := "alipay_cn"
	param := map[string]interface{}{
		"currency": createPayContext.Pay.Currency,
		"amount":   utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		"pm_id":    pmID,
		//"title":               name,
		"description": description,
		"order_id":    createPayContext.Pay.PaymentId,
		//"customer_id": strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"payer_email": createPayContext.Email,
		"return_url":  webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		//"backUrl":     webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		//"payment_data": createPayContext.Metadata,
		//"pending_deadline_at": time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	if !config.GetConfigInstance().IsProd() {
		param["pm_id"] = "payssion_test"
	}
	param["api_key"] = createPayContext.Gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v|%v", param["api_key"], param["pm_id"], param["amount"], param["currency"], param["order_id"], createPayContext.Gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "PayssionNewPayment", nil, createPayContext.Gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	//transaction := responseJson.Get("transaction")

	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       responseJson.Get("transaction.transaction_id").String(),
		GatewayPaymentIntentId: responseJson.Get("transaction.transaction_id").String(),
		Link:                   responseJson.Get("redirect_url").String(),
	}, nil
}

func (c Payssion) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c Payssion) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api/v1/payment/getDetail"
	param := map[string]interface{}{}
	param["transaction_id"] = gatewayPaymentId
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v", param["api_key"], param["transaction_id"], param["order_id"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "PayssionPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return parsePayssionPayment(responseJson), nil
}

func (c Payssion) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

// https://payssion.com/cn/docs/#api-reference-payment-details
func (c Payssion) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	urlPath := "/api/v1/payment/getDetail"
	param := map[string]interface{}{}
	param["transaction_id"] = gatewayRefundId
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v", param["api_key"], param["transaction_id"], param["order_id"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "PayssionPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	var status consts.RefundStatusEnum = consts.RefundCreated
	if strings.Compare(responseJson.Get("transaction.state").String(), "refunded") == 0 {
		status = consts.RefundSuccess
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          status,
		Reason:          refund.RefundComment,
		RefundAmount:    refund.RefundAmount,
		Currency:        strings.ToUpper(refund.Currency),
		RefundTime:      gtime.Now(),
	}, nil
}

func (c Payssion) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	urlPath := "/api/v1/refunds"
	param := map[string]interface{}{}
	gateway := util.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	param["transaction_id"] = payment.GatewayPaymentId
	param["amount"] = utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency)
	param["currency"] = strings.ToUpper(refund.Currency)
	param["track_id"] = refund.RefundId
	param["description"] = refund.RefundComment
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v", param["api_key"], param["transaction_id"], param["amount"], param["currency"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayRefund", param, responseJson, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: responseJson.Get("transaction.transaction_id").String(),
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (c Payssion) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func parsePayssionPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if strings.Compare(item.Get("transaction.state").String(), "pending") == 0 {
		authorizeStatus = consts.Authorized
	} else if strings.Compare(item.Get("transaction.state").String(), "completed") == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("transaction.state").String(), "cancelled") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("transaction.state").String(), "failed") == 0 {
		status = consts.PaymentFailed
	}

	var authorizeReason = ""
	var paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("transaction.amount").String(), item.Get("transaction.currency").String())
	var paymentMethod = ""
	if item.Contains("selected_payment_method") && item.GetJson("selected_payment_method").Contains("payins") {
		//paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("orderSum").String(), item.Get("orderCurrency").String())
		//paymentMethod = item.Get("payin_currency").String() + "|" + item.Get("payin_network").String()
	}
	var paidTime *gtime.Time
	if item.Contains("transaction.updated") {
		if t, err := gtime.StrToTime(item.Get("transaction.updated").String()); err == nil {
			paidTime = t
		}
	}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("transaction.transaction_id").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          utility.ConvertDollarStrToCent(item.Get("transaction.amount").String(), item.Get("transaction.currency").String()),
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: paymentMethod,
		PaidTime:             paidTime,
	}
}

func SendPayssionPaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")
	domain := "http://www.payssion.com"
	if !config.GetConfigInstance().IsProd() {
		domain = "http://sandbox.payssion.com"
		param["pm_id"] = "payssion_test"
	}
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nPayssion_Start %s %s %s %s\n", method, urlPath, publicKey, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := utility.SendRequest(domain+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nPayssion_End %s %s response: %s error %s\n", method, urlPath, response, err)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}
