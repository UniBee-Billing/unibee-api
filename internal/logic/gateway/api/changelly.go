package api

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"unibee/internal/consts"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

//https://api.pay.changelly.com/
//https://pay.changelly.com/

type Changelly struct {
}

func (c Changelly) GatewayTest(ctx context.Context, key string, secret string) (gatewayType int64, err error) {
	urlPath := "/api/payment/v1/payments"
	param := map[string]interface{}{
		"nominal_currency": "USDT",
		"nominal_amount":   "1.08",
		"title":            "test crypto payment",
		"description":      "test crypto payment description",
		"order_id":         uuid.New().String(),
		"customer_id":      "17",
		"customer_email":   "jack.fu@wowow.io",
	}
	data, err := SendChangellyRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("invalid keys, json parse error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("id"), "invalid keys, id is nil")
	return consts.GatewayTypeCrypto, nil
}

func (c Changelly) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserAttachPaymentMethodInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserDeAttachPaymentMethodInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *ro.GatewayUserPaymentMethodReq) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	utility.Assert(len(req.GatewayPaymentId) > 0, "gatewayPaymentId is nil")
	urlPath := "/api/payment/v1/payments/" + req.GatewayPaymentId + "/payment_methods"
	param := map[string]interface{}{}
	data, err := SendChangellyRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentMethodList", param, data, err, "ChangelyPaymentMethodList", nil, gateway)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(data))
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if err != nil {
		return nil, err
	}
	var paymentMethods []*ro.PaymentMethod
	for _, method := range responseJson.GetJsons("") {
		if method.Contains("code") && method.Contains("networks") {
			currencyCode := method.Get("code").String()
			for _, network := range method.GetJsons("networks") {
				if network.Contains("code") {
					paymentMethods = append(paymentMethods, &ro.PaymentMethod{
						Id: currencyCode + "|" + network.Get("code").String(),
					})
				}
			}
		}
	}
	return &ro.GatewayUserPaymentMethodListInternalResp{PaymentMethods: paymentMethods}, nil
}

func (c Changelly) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64, data *gjson.Json) (res *ro.GatewayUserPaymentMethodCreateAndBindInternalResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayNewPayment(ctx context.Context, createPayContext *ro.NewPaymentInternalReq) (res *ro.NewPaymentInternalResp, err error) {
	urlPath := "/api/payment/v1/payments"
	var gasPayer string
	if createPayContext.Pay.GasPayer == "merchant" {
		gasPayer = "MERCHANT"
	} else {
		gasPayer = "CUSTOMER"
	}
	param := map[string]interface{}{
		"nominal_currency":     createPayContext.Pay.Currency,
		"nominal_amount":       utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		"title":                "",
		"description":          "",
		"order_id":             createPayContext.Pay.PaymentId,
		"customer_id":          createPayContext.Pay.UserId,
		"customer_email":       createPayContext.Email,
		"success_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"failure_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		"fees_payer":           gasPayer, // who pay the fee
		"payment_data":         createPayContext.Metadata,
		"pending_deadline_at":  "",
	}
	data, err := SendChangellyRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, data, err, "ChangelyNewPayment", nil, createPayContext.Gateway)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(data))
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if responseJson.Contains("id") {
		return nil, gerror.New("invalid request, id is nil")
	}
	if err != nil {
		return nil, err
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := responseJson.Get("id").String()
	return &ro.NewPaymentInternalResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   responseJson.Get("payment_url").String(),
	}, nil
}

func (c Changelly) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	urlPath := "/api/payment/v1/payments/" + gatewayPaymentId
	param := map[string]interface{}{}
	data, err := SendChangellyRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, data, err, "ChangelyPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(data))
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if err != nil {
		return nil, err
	}

	return parseChangellyPayment(responseJson), nil
}

func (c Changelly) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	return nil, gerror.New("Not Support")
}

func parseChangellyPayment(item *gjson.Json) *ro.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if strings.Compare(item.Get("state").String(), "WAITING") == 0 {
		authorizeStatus = consts.Authorized
	} else if strings.Compare(item.Get("state").String(), "COMPLETED") == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("state").String(), "CANCELED") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("state").String(), "FAILED") == 0 {
		status = consts.PaymentFailed
	}

	var authorizeReason = ""
	//var gatewayPaymentMethod string
	//if item.PaymentMethod != nil {
	//	gatewayPaymentMethod = item.PaymentMethod.ID
	//}
	var paymentAmount int64 = 0
	var paymentMethod = ""
	if item.Contains("selected_payment_method") && item.GetJson("selected_payment_method").Contains("expected_payin_amount") {
		paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("selected_payment_method").Get("expected_payin_amount").String(), item.Get("nominal_currency").String())
		paymentMethod = item.Get("payin_currency").String() + "|" + item.Get("payin_network").String()
	}
	var paidTime *gtime.Time
	if item.Contains("completed_at") {
		if t, err := gtime.StrToTime(item.Get("completed_at").String()); err == nil {
			paidTime = t
		}
	}

	return &ro.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("id").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          utility.ConvertDollarStrToCent(item.Get("nominal_amount").String(), item.Get("nominal_currency").String()),
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: paymentMethod,
		PayTime:              paidTime,
	}
}

func SendChangellyRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res []byte, err error) {
	utility.Assert(param != nil, "param is nil")
	datetime := getExpirationDateTime(1)

	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nChangelly_Start %s %s %s %s %s\n", method, urlPath, publicKey, jsonString, datetime)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Signature":  sign(method, urlPath, datetime, privateKey, body),
		"X-Api-Key":    publicKey,
	}
	response, err := utility.SendRequest("https://api.pay.changelly.com"+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nChangelly_End %s %s response: %s error %s\n", method, urlPath, response, err)
	return response, err
}

func sign(method string, urlPath string, dateTime string, purePrivateKey string, postJson []byte) (sign string) {
	var builder strings.Builder
	lineSeparator := lineSeparator()
	builder.WriteString(method)
	builder.WriteString(lineSeparator)
	builder.WriteString(urlPath)
	builder.WriteString(lineSeparator)
	builder.WriteString(base64Encoding(postJson))
	builder.WriteString(lineSeparator)
	builder.WriteString(dateTime)
	payload := builder.String()
	privateKey := purePrivateKey
	if !strings.Contains(privateKey, "BEGIN PRIVATE KEY") {
		privateKey = `
***REMOVED***
` + purePrivateKey + `
***REMOVED***
`
	}
	block, _ := pem.Decode([]byte(privateKey))
	utility.Assert(block != nil, "rsa encrypt error")
	prv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	utility.AssertError(err, "rsa encrypt error")
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(payload))
	utility.AssertError(err, "sha256 hash encrypt error")
	result, err := rsa.SignPKCS1v15(rand.Reader, prv.(*rsa.PrivateKey), crypto.SHA256, msgHash.Sum(nil))
	//result, err := utility.RsaEncrypt([]byte(key), []byte(sha256Encoding(builder.String())))
	utility.AssertError(err, "rsa encrypt error")
	return base64Encoding([]byte(base64Encoding(result) + lineSeparator + dateTime))
}

func getExpirationDateTime(hour int64) (datetime string) {
	return strconv.FormatInt(gtime.Now().Unix()+(hour*3600), 10)
}

func lineSeparator() string {
	return ":"
}

func base64Encoding(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}