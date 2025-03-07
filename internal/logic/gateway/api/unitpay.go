package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"net/url"
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
	"unibee/utility/unibee"
)

//https://unitpay.ru/
//https://help.unitpay.ru/other/test-api
//https://help.unitpay.ru/master/payments/create-payment
//https://help.unitpay.ru/master/book-of-reference/payment-system-codes
//How to use test mode https://help.unitpay.ru/master/book-of-reference/test-api
//https://help.unitpay.ru/other/test-api
// todo mark 3ds check send email
// todo mark subscription test
// todo mark auto-charge response didn't contain paymentId

type UnitPay struct {
}

func (c UnitPay) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "UnitPay",
		Description:                   "Use Project Id and Secret Key to secure the payment",
		DisplayName:                   "UnitPay",
		GatewayWebsiteLink:            "https://unitpay.ru",
		GatewayWebhookIntegrationLink: "https://unitpay.ru/partner",
		GatewayLogo:                   "https://api.unibee.top/oss/file/d76q4ctiz7jjaemhsr.png",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d76q4ctiz7jjaemhsr.png"},
		GatewayType:                   consts.GatewayTypeCard,
		CurrencyExchangeEnabled:       true,
		QueueForRefund:                true,
		Sort:                          70,
		PublicKeyName:                 "Project Id",
	}
}

func (c UnitPay) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (c UnitPay) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(!strings.Contains(req.Key, "-"), "Please using unitpay 'ProjectId' instead 'publicKey'")
	urlPath := "/api?method=initPayment"
	param := map[string]interface{}{
		"currency":    "RUB",
		"sum":         "100",
		"paymentType": "card",
		"desc":        "test_payment_description",
		"projectId":   req.Key,
		//"account":     "test_unitpay",
		"account": "test",
		//"subscriptionId": 1,
		//"subscription": true,
		//"locale":      "en",
		//"preauth":     0,
		//"ip": "15155-ae12d",
	}
	param["signature"] = getUnitPayFormSignature(fmt.Sprintf("%v", param["account"]), fmt.Sprintf("%v", param["currency"]), fmt.Sprintf("%v", param["desc"]), fmt.Sprintf("%v", param["sum"]), req.Secret)
	responseJson, err := SendUnitPayPaymentRequest(ctx, req.Secret, "GET", urlPath, param, config.GetConfigInstance().IsProd())
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	g.Log().Infof(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("result.paymentId"), "invalid keys, paymentId is nil")
	utility.Assert(responseJson.Contains("result.redirectUrl"), "invalid keys, redirectUrl is nil")
	redirectUrl, err := url.PathUnescape(responseJson.Get("result.redirectUrl").String())
	utility.Assert(err == nil, "invalid keys, invalid redirectUrl")
	g.Log().Debugf(ctx, "redirectUrl: %s", redirectUrl)
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c UnitPay) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GetSubscription(ctx context.Context, secret string, subscriptionId string) (res *gjson.Json, err error) {
	urlPath := "/api?method=getSubscription"
	param := map[string]interface{}{}
	param["subscriptionId"] = subscriptionId
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	return responseJson, err
}

func (c UnitPay) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/api?method=initPayment"
	//var name = ""
	var description = createPayContext.Invoice.ProductName
	if len(createPayContext.Invoice.Lines) > 0 {
		var line = createPayContext.Invoice.Lines[0]
		if len(line.Name) > 0 {
			description = line.Name
		} else if len(line.Description) > 0 {
			description = line.Description
		}
	}

	var currency = createPayContext.Pay.Currency
	var totalAmount = createPayContext.Pay.TotalAmount
	{
		// Currency Exchange
		if createPayContext.GatewayCurrencyExchange != nil && createPayContext.ExchangeAmount > 0 && len(createPayContext.ExchangeCurrency) > 0 {
			currency = createPayContext.ExchangeCurrency
			totalAmount = createPayContext.ExchangeAmount
		}
	}
	param := map[string]interface{}{
		"currency":    currency,
		"sum":         utility.ConvertCentToDollarStr(totalAmount, currency),
		"paymentType": "card",
		//"title":               name,
		"desc":      description,
		"projectId": gateway.GatewayKey,
		//"customer_id": strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"account":   createPayContext.Pay.PaymentId,
		"resultUrl": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"backUrl":   webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		//"payment_data":        createPayContext.Metadata,
		//"pending_deadline_at": time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	if len(createPayContext.Gateway.UniqueProductId) > 0 && createPayContext.CheckoutMode {
		param["subscription"] = true
	}
	if len(gateway.UniqueProductId) > 0 &&
		len(createPayContext.GatewayPaymentMethod) > 0 &&
		!createPayContext.CheckoutMode {
		subscription, err := c.GetSubscription(ctx, gateway.GatewayKey, createPayContext.GatewayPaymentMethod)
		if err == nil && subscription != nil &&
			subscription.Contains("result") &&
			subscription.Contains("result.status") &&
			subscription.Get("result.status").String() == "active" {
			param["subscriptionId"] = createPayContext.GatewayPaymentMethod
		}
	}
	param["signature"] = getUnitPayFormSignature(fmt.Sprintf("%v", param["account"]), fmt.Sprintf("%v", param["currency"]), fmt.Sprintf("%v", param["desc"]), fmt.Sprintf("%v", param["sum"]), gateway.GatewaySecret)
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewaySecret, "GET", urlPath, param, config.GetConfigInstance().IsProd())
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "UnitPayNewPayment", nil, gateway)
	if err != nil {
		return nil, err
	}
	if responseJson.Contains("error") && responseJson.Contains("error.message") {
		return nil, gerror.New(fmt.Sprintf("invalid request, %s", responseJson.Get("error.message")))
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("result.paymentId") {
		return nil, gerror.New(fmt.Sprintf("invalid request, paymentId is nil, %s", responseJson.Get("error.message")))
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := responseJson.Get("result.paymentId").String()
	redirectUrl, err := url.PathUnescape(responseJson.Get("result.redirectUrl").String())
	if err != nil {
		return nil, gerror.New("invalid request, invalid redirectUrl")
	}
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   redirectUrl,
	}, nil
}

func (c UnitPay) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c UnitPay) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api?method=getPayment"
	param := map[string]interface{}{}
	param["paymentId"] = gatewayPaymentId
	param["secretKey"] = gateway.GatewaySecret
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewaySecret, "GET", urlPath, param, config.GetConfigInstance().IsProd())
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "UnitPayPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Infof(ctx, "GatewayPaymentDetail paymentId:%s responseJson :%s", payment.PaymentId, responseJson.String())

	return parseUnitPayPayment(responseJson), nil
}

func (c UnitPay) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	payment := util.GetPaymentByPaymentId(ctx, refund.PaymentId)
	if payment == nil {
		return nil, gerror.New("payment not found")
	}

	detail, err := c.GatewayPaymentDetail(ctx, gateway, payment.GatewayPaymentId, payment)
	if err != nil {
		return nil, err
	}
	if detail.RefundSequence > int64(refund.RefundGatewaySequence) {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: gatewayRefundId,
			Status:          consts.RefundSuccess,
			Reason:          refund.RefundComment,
			RefundAmount:    refund.RefundAmount,
			Currency:        strings.ToUpper(refund.Currency),
			RefundTime:      gtime.Now(),
		}, nil
	} else {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: gatewayRefundId,
			Status:          consts.RefundCreated,
			Reason:          refund.RefundComment,
			RefundAmount:    refund.RefundAmount,
			Currency:        strings.ToUpper(refund.Currency),
			RefundTime:      gtime.Now(),
		}, nil
	}
}

// https://help.unitpay.ru/master/payments/payment-refund
func (c UnitPay) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	detail, err := c.GatewayPaymentDetail(ctx, gateway, createPaymentRefundContext.Payment.GatewayPaymentId, createPaymentRefundContext.Payment)
	if err != nil {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("Get Gateway Refund Sequence Failed:%s", err.Error()),
		}, nil
	}
	urlPath := "/api?method=refundPayment"
	param := map[string]interface{}{}
	param["paymentId"] = createPaymentRefundContext.Payment.GatewayPaymentId
	if createPaymentRefundContext.GatewayCurrencyExchange != nil && createPaymentRefundContext.ExchangeRefundAmount > 0 && len(createPaymentRefundContext.ExchangeRefundCurrency) > 0 {
		param["sum"] = utility.ConvertCentToDollarStr(createPaymentRefundContext.ExchangeRefundAmount, createPaymentRefundContext.ExchangeRefundCurrency)
	} else {
		param["sum"] = utility.ConvertCentToDollarStr(createPaymentRefundContext.Refund.RefundAmount, createPaymentRefundContext.Refund.Currency)
	}
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewaySecret, "GET", urlPath, param, config.GetConfigInstance().IsProd())
	log.SaveChannelHttpLog("GatewayRefund", param, responseJson, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}
	if responseJson.Contains("error.message") {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          responseJson.Get("error.message").String(),
		}, nil
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
		RefundSequence:  unibee.Int64(detail.RefundSequence),
	}, nil
}

func (c UnitPay) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func parseUnitPayPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
	/*
		success — successful payment;
		wait — payment is pending;
		error — payment error;
		error_pay — error/failure of the store at the PAY stage; in statistics, it is displayed as "incomplete";
		error_check — error/failure of the store at the CHECK stage, in statistics it is displayed as "rejected";
		refund — refund to the buyer;
		secure — being verified by the Bank's security service.
	*/
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if strings.Compare(item.Get("result.status").String(), "wait") == 0 {
		authorizeStatus = consts.Authorized
	} else if strings.Compare(item.Get("result.status").String(), "success") == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("result.status").String(), "error_check") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("result.status").String(), "error_pay") == 0 {
		status = consts.PaymentFailed
	}

	var authorizeReason = ""
	var paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("result.orderSum").String(), item.Get("result.orderCurrency").String())
	var paymentMethod = ""
	if item.Contains("selected_payment_method") && item.GetJson("selected_payment_method").Contains("payins") {
		//paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("orderSum").String(), item.Get("orderCurrency").String())
		//paymentMethod = item.Get("payin_currency").String() + "|" + item.Get("payin_network").String()
	}
	var paidTime *gtime.Time
	if item.Contains("date") {
		if t, err := gtime.StrToTime(item.Get("date").String()); err == nil {
			paidTime = t
		}
	}

	var refundSequence int64 = 0
	if item.Contains("refunds") {
		refundSequence = int64(len(item.GetJsons("refunds")))
	}
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("result.paymentId").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          utility.ConvertDollarStrToCent(item.Get("result.payerSum").String(), item.Get("result.payerCurrency").String()),
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: paymentMethod,
		PaidTime:             paidTime,
		RefundSequence:       refundSequence,
	}
}

func SendUnitPayPaymentRequest(ctx context.Context, privateKey string, method string, urlPath string, param map[string]interface{}, isProd bool) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")

	param["secretKey"] = privateKey
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nUnitPay_Start %s %s %s\n", method, urlPath, jsonString)
	body := []byte(jsonString)
	paramPath := ""
	for k, v := range param {
		paramPath = fmt.Sprintf("%s&params[%s]=%v", paramPath, k, url.QueryEscape(fmt.Sprintf("%v", v)))
	}

	if !isProd {
		//paramPath = fmt.Sprintf("%s&params[test]=1&params[login]=senseybiz@gmail.com", paramPath)
		paramPath = fmt.Sprintf("%s&params[test]=1&params[login]=senseybiz@gmail.com", paramPath)
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	urlPath = fmt.Sprintf("%s&%s", urlPath, paramPath)
	response, err := utility.SendRequest("https://unitpay.ru"+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nUnitPay_End %s %s response: %s error %s\n", method, "https://unitpay.ru"+urlPath, string(response), err)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	g.Log().Debugf(ctx, "\nUnitPay_End %s %s decodeResponse: %s error %s\n", method, "https://unitpay.ru"+urlPath, responseJson, err)
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}

func getUnitPayFormSignature(account, currency, desc, sum, secretKey string) string {
	hashStr := account + "{up}" + currency + "{up}" + desc + "{up}" + sum + "{up}" + secretKey
	g.Log().Debugf(context.Background(), "UnitPay_Start before decode signature: %s \n", hashStr)
	hash := sha256.Sum256([]byte(hashStr))
	g.Log().Debugf(context.Background(), "UnitPay_Start after decode signature: %s \n", hex.EncodeToString(hash[:]))
	return hex.EncodeToString(hash[:])
}
