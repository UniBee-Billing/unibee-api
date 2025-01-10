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
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

//https://unitpay.ru/
//https://help.unitpay.ru/other/test-api
//https://help.unitpay.ru/master/payments/create-payment
//https://help.unitpay.ru/master/book-of-reference/payment-system-codes
//How to use test mode https://help.unitpay.ru/master/book-of-reference/test-api

type UnitPay struct {
}

func (c UnitPay) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (c UnitPay) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	urlPath := "/api?method=initPayment"
	param := map[string]interface{}{
		"currency":    "RUB",
		"sum":         100,
		"paymentType": "card",
		"desc":        "test_payment_description",
		"projectId":   key,
		"account":     "test_user",
		//"locale":      "en",
		//"preauth":     0,
		//"ip": "15155-ae12d",
	}
	param["signature"] = getUnitPayFormSignature(fmt.Sprintf("%v", param["account"]), fmt.Sprintf("%v", param["currency"]), fmt.Sprintf("%v", param["desc"]), fmt.Sprintf("%v", param["sum"]), secret)
	responseJson, err := SendUnitPayPaymentRequest(ctx, key, secret, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("paymentId"), "invalid keys, id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c UnitPay) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
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

func (c UnitPay) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/api?method=initPayment"
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

	param := map[string]interface{}{
		"currency":    createPayContext.Pay.Currency,
		"sum":         utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		"paymentType": "card",
		//"title":               name,
		"desc":      description,
		"projectId": createPayContext.Gateway.GatewayKey,
		//"customer_id": strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"account":   createPayContext.Pay.PaymentId,
		"resultUrl": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"backUrl":   webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		//"payment_data":        createPayContext.Metadata,
		//"pending_deadline_at": time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	param["signature"] = getUnitPayFormSignature(fmt.Sprintf("%v", param["account"]), fmt.Sprintf("%v", param["currency"]), fmt.Sprintf("%v", param["desc"]), fmt.Sprintf("%v", param["sum"]), createPayContext.Gateway.GatewaySecret)
	responseJson, err := SendUnitPayPaymentRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "UnitPayNewPayment", nil, createPayContext.Gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("paymentId") {
		return nil, gerror.New("invalid request, paymentId is nil")
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := responseJson.Get("paymentId").String()
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   responseJson.Get("redirectUrl").String(),
	}, nil
}

func (c UnitPay) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
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
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "UnitPayPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return parseUnitPayPayment(responseJson), nil
}

func (c UnitPay) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c UnitPay) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	urlPath := "/api?method=getPayment"
	param := map[string]interface{}{}
	param["paymentId"] = gatewayRefundId
	param["secretKey"] = gateway.GatewaySecret
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "UnitPayPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	/*
		success — successful payment;
		wait — payment is pending;
		error — payment error;
		error_pay — error/failure of the store at the PAY stage; in statistics, it is displayed as "incomplete";
		error_check — error/failure of the store at the CHECK stage, in statistics it is displayed as "rejected";
		refund — refund to the buyer;
		secure — being verified by the Bank's security service.
	*/
	var status consts.RefundStatusEnum = consts.RefundCreated
	if strings.Compare(responseJson.Get("state").String(), "refund") == 0 {
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

// https://help.unitpay.ru/master/payments/payment-refund
func (c UnitPay) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	urlPath := "/api?method=refundPayment"
	param := map[string]interface{}{}
	gateway := util.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	param["paymentId"] = payment.GatewayPaymentId
	param["secretKey"] = gateway.GatewaySecret
	param["sum"] = utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency)
	responseJson, err := SendUnitPayPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayRefund", param, responseJson, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: payment.GatewayPaymentId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (c UnitPay) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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
	if strings.Compare(item.Get("state").String(), "wait") == 0 {
		authorizeStatus = consts.Authorized
	} else if strings.Compare(item.Get("state").String(), "success") == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("state").String(), "error_check") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("state").String(), "error_pay") == 0 {
		status = consts.PaymentFailed
	}

	var authorizeReason = ""
	var paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("orderSum").String(), item.Get("orderCurrency").String())
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

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("paymentId").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          utility.ConvertDollarStrToCent(item.Get("payerSum").String(), item.Get("payerCurrency").String()),
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: paymentMethod,
		PaidTime:             paidTime,
	}
}

func SendUnitPayPaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")

	param["secretKey"] = privateKey
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nUnitPay_Start %s %s %s %s\n", method, urlPath, publicKey, jsonString)
	body := []byte(jsonString)
	paramPath := ""
	for k, v := range param {
		paramPath = fmt.Sprintf("%s&params[%s]=%v", paramPath, k, url.QueryEscape(fmt.Sprintf("%v", v)))
	}

	if !config.GetConfigInstance().IsProd() {
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
