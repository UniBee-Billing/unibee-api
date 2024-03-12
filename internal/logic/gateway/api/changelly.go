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
	"time"
	"unibee/api/bean"
	"unibee/internal/consts"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

//https://api.pay.changelly.com/
//https://pay.changelly.com/

type Changelly struct {
}

func (c Changelly) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
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
	responseJson, err := SendChangellyRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("id"), "invalid keys, id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c Changelly) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	utility.Assert(len(req.GatewayPaymentId) > 0, "gatewayPaymentId is nil")
	urlPath := "/api/payment/v1/payments/" + req.GatewayPaymentId + "/payment_methods"
	param := map[string]interface{}{}
	responseJson, err := SendChangellyRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentMethodList", param, responseJson, err, "ChangelyPaymentMethodList", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	var paymentMethods []*bean.PaymentMethod
	for _, a := range responseJson.Array() {
		if method, ok := a.(map[string]interface{}); ok {
			if method["code"] != nil && method["networks"] != nil {
				currencyCode := method["code"].(string)
				for _, network := range method["networks"].([]interface{}) {
					if network.(map[string]interface{})["code"] != nil && len(network.(map[string]interface{})["code"].(string)) > 0 {
						paymentMethods = append(paymentMethods, &bean.PaymentMethod{
							Id: currencyCode + "|" + network.(map[string]interface{})["code"].(string),
						})
					}
				}
			}
		}
	}
	return &gateway_bean.GatewayUserPaymentMethodListResp{PaymentMethods: paymentMethods}, nil
}

func (c Changelly) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64, data *gjson.Json) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
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
		"customer_id":          strconv.FormatInt(createPayContext.Pay.UserId, 10),
		"customer_email":       createPayContext.Email,
		"success_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"failure_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		"fees_payer":           gasPayer, // who pay the fee
		"payment_data":         createPayContext.Metadata,
		"pending_deadline_at":  time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	responseJson, err := SendChangellyRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "ChangelyNewPayment", nil, createPayContext.Gateway)
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
		Link:                   responseJson.Get("payment_url").String(),
	}, nil
}

func (c Changelly) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c Changelly) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api/payment/v1/payments/" + gatewayPaymentId
	param := map[string]interface{}{}
	responseJson, err := SendChangellyRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "ChangelyPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if err != nil {
		return nil, err
	}

	return parseChangellyPayment(responseJson), nil
}

func (c Changelly) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Changelly) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func parseChangellyPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
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

	return &gateway_bean.GatewayPaymentRo{
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

func SendChangellyRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
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
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	return responseJson, nil
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
-----BEGIN PRIVATE KEY-----
` + purePrivateKey + `
-----END PRIVATE KEY-----
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
