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
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

//https://api.pay.changelly.com/
//https://pay.changelly.com/

type Changelly struct {
}

func (c Changelly) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Changelly",
		Description:                   "Use public and private keys to secure the crypto payment.",
		DisplayName:                   "Crypto",
		GatewayWebsiteLink:            "https://pay.changelly.com/",
		GatewayWebhookIntegrationLink: "https://app.pay.changelly.com/integrations",
		GatewayLogo:                   "https://api.unibee.top/oss/file/d6z44ghlnhkglgo76y.svg",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d6yhnz0wty7w6m7zhd.svg", "https://api.unibee.top/oss/file/d6yho8slal03ywl65c.svg", "https://api.unibee.top/oss/file/d6yhoilcikizou9ztk.svg", "https://api.unibee.top/oss/file/d6yhotsmefitw0cav1.svg"},
		GatewayType:                   consts.GatewayTypeCrypto,
		Sort:                          5,
	}
}

func (c Changelly) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	from.CountryCode = "FR" // not all countryCode contain data
	urlPath := "/v1/offers?currencyFrom=" + from.Currency + "&currencyTo=USDT20&amountFrom=100&country=" + from.CountryCode
	//urlPath := "/v1/offers?currencyFrom=" + from.Currency + "&currencyTo=USD&amountFrom=100&country=" + from.CountryCode
	param := map[string]interface{}{}
	responseJson, err := SendChangellyFiatRequest(ctx, "105ee736f1d01918ffd7c794ebc4b6e94169454e10fbcf094a347f0a370c7f08", "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRRERrQnpuakV0NVRJUnYKMDhmSVUvYmt6NjdSbUc2ZDNVSXlhL0liRmRjY0QzREdyWHRKTDJJaG5ROGhyV0RzaDFaUWV1K1dSaWF5VkI5NApEM01Ia0o3WWsxb05SaVNyV1FiSXA5eTZrbWdiMnpQeDhwRWJ0R2txRlVScWoxS2JNYTdEQmtwTjFxUnUxMTNsCmNHQTliWnBLL0lsUlVKOEpZUFBieU5kZnV3d2hPNDlqakg3M1FmVmVYY2RFNUd0WkpBK0R1R3pybTJsc1Z3N2cKQ3lhZlMyRGdLaGwzdjN3N3ZrOTRRVHJRaDFOR0hyM0lIaG5mSmxnQUlJOHFTNWtIUkh4SGZFbG1sREgwcDYxbgpRQWhJK0ljN3FiRzVKMnBDdk95S084L2FIY3BiZ1dpc2EzanV0cDFuYWNEa2lzdkhOcWpjRHRXQ1Y3K0pJYWk5CnlJZGw0OEVuQWdNQkFBRUNnZ0VBQVdlTGpaalRBZkFCYU5HRjRja1lsaWx0QVZHdG1iSDVHUkR3RUM4UlRjRHoKOTdsUE9yQ29iUmdLVlJ6cSsvS000Q1JzakJFU1BoTDBsdWJRSzFjU2pQajNSZTUzQjU2cXVjNHdKa3cwTmQ4QgpKSCt3SXhsM1FvekhONnl5ZnkyQUtua1B3amRadnZkZnFVSXNCQnBCYUJYSzRWYWU0eHV3ZWZ6d21iTm5odVdQCllvR0J0S2I2Mkx3blYvdXhGdWNwMTQra3QrUWsrRU5QbUwyVWE2emJscjBCTmFJMytIejJKSzd5L3JpcXZoOXEKbUNCektPVG5nNVdkV01BazhPc3ZmdnJ3WW1jU3FDbE9mcjhQbWhrbC9UeUw1d0F3clJGdEcyZW8rZmtFRG9zYgo1SnBQaXI0NmhBOENBNHBXTGdrUkcvbm54YmptYmphQlk4NUFtbnFad1FLQmdRRHdnVFFkRGVrcG1JRWpFbmMwCnM1WkxOZUY4SzlDWDRKbVVTOFhROHBEM3RQMkd0VDBraURkUWlWNTVWL2cxdXR2ZXdWT0xlbjFmdExub1ExR00KeGJaMytIOExWUEZLVVdsYno3SFhQOUxzMll5R1d4STJmU1R1R3pOTU1kR01ZVnZJUGFZMDlMSmRGdTRwM3F6bgptZHFBZHJtdC9ZTExpRy85Q2dZeFpYSVlmd0tCZ1FEUUthZ0tiOTBPcGNWa1RjaTdabEtlajhCT2pnSXhEV2FICm11SVMvMjZwN3FpYzdadVQzMEllckV3V0l5cjNDbk1TQ2paR2wyUGFTNE1uYmNmUU5lSHFlUzF6dXVmeTdxOTEKR1Q5UkNqQnZCaFZNUDVMVy9FcXAyMlhiK0hBZndnaStVcjAxbjlvNnVCV29mZys0ZGtwK0lQa3diSWk1QklsRwpYYWRIMHpCRFdRS0JnRE1yeW1MRUt1L211dE16Z3BsNy9HWlVPSDJxOU43YnN0R1NyYXdmY0NqRUlZMGYwcnFMCklQbkp3SWdnNTNiSEl6RHFBVlNUNDBrUnN0eHdObEcxWDNWM01kQy9hZmRlQ3dTMTFDandNM2loY1B6Rk04TFUKTFo5YnVqWmtBeW5UTFN3VnNkOWlrUEN0aUU3d1NlbWRHcGhxcW1jU29WbWMxZmNJd3ZpUGxROFJBb0dCQUszbwpIMXVZMlRYRGlJV2o5bStackt5THJEMzBwaUFVOGZPWWtnY05INGNZdkFWZS9Qc3RLakEyQWRyOHhvaGRVb1ZmCndyaDNBaFQ5d1RUUG9uOXdoSzAvVDVuQUxNZm9ZTzJUaWpKS01PeVFTSHJMSWdJNkJLYWpoUldoR1F0dkw1N1IKd1FGcjZ3WGpoVFNmSE1NZkVGMFBieC9sak5RRjFpblRWRTNOUWlVQkFvR0FNbk1EUGpxVWVkc2h0b2UySER5bwo4cG42S0ZiUVVIcDlnU2JXZk9XZnZqMkF0TzhUanlYUzhZVko5bVltNkpPZlo5cEVlMHBZamxEY1Z2MkRZWDcwCjBCMEIyWmkvdFZUOVU1ZUpVa2hTUEVHSnlwRHMwR1VVV3hXWFV0bllUU2tRcWZLTnNBS3V2WVpwbHg2NDU3Z3kKeXpCMFA3Q1gvdlgvc1dwZWx1WnhBQ3M9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K", "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayCryptoFiatTrans", urlPath, responseJson, err, "ChangelyFiatTrans", nil, from.Gateway)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	var invertedRate float64 = 0
	for _, a := range responseJson.Array() {
		if method, ok := a.(map[string]interface{}); ok {
			if method["invertedRate"] != nil {
				rateFloat, err := strconv.ParseFloat(method["invertedRate"].(string), 64)
				if err == nil && rateFloat > 0 && (invertedRate == 0 || invertedRate > rateFloat) {
					invertedRate = rateFloat
				}
			}
		}
	}
	if invertedRate == 0 {
		return nil, gerror.New("rate not found")
	}
	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   utility.RoundUp(float64(from.Amount) / invertedRate),
		CryptoCurrency: "USDT",
		Rate:           invertedRate,
	}, nil
}

func (c Changelly) GatewayGetCurrency(ctx context.Context, key string, secret string) (responseJson *gjson.Json, err error) {
	urlPath := "/api/payment/v1/currencies"
	param := map[string]interface{}{}
	responseJson, err = SendChangellyPaymentRequest(ctx, key, secret, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return responseJson, err
}

func (c Changelly) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	urlPath := "/api/payment/v1/payments"
	param := map[string]interface{}{
		"nominal_currency": "BNB",
		"nominal_amount":   "0.00001",
		"title":            "test crypto payment",
		"description":      "test crypto payment description",
		"order_id":         uuid.New().String(),
		"customer_id":      "17",
		"customer_email":   "jack.fu@wowow.io",
	}
	responseJson, err := SendChangellyPaymentRequest(ctx, key, secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call changelly error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("id"), "invalid keys, id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c Changelly) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	utility.Assert(len(req.GatewayPaymentId) > 0, "gatewayPaymentId is nil")
	urlPath := "/api/payment/v1/payments/" + req.GatewayPaymentId + "/payment_methods"
	param := map[string]interface{}{}
	responseJson, err := SendChangellyPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentMethodList", param, responseJson, err, "ChangelyPaymentMethodList", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	var paymentMethods []*gateway_bean.PaymentMethod
	for _, a := range responseJson.Array() {
		if method, ok := a.(map[string]interface{}); ok {
			if method["code"] != nil && method["networks"] != nil {
				currencyCode := method["code"].(string)
				for _, network := range method["networks"].([]interface{}) {
					if network.(map[string]interface{})["code"] != nil && len(network.(map[string]interface{})["code"].(string)) > 0 {
						paymentMethods = append(paymentMethods, &gateway_bean.PaymentMethod{
							Id: currencyCode + "|" + network.(map[string]interface{})["code"].(string),
						})
					}
				}
			}
		}
	}
	return &gateway_bean.GatewayUserPaymentMethodListResp{PaymentMethods: paymentMethods}, nil
}

func (c Changelly) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
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
		"nominal_currency":     createPayContext.Pay.CryptoCurrency,
		"nominal_amount":       utility.ConvertCentToDollarStr(createPayContext.Pay.CryptoAmount, createPayContext.Pay.CryptoCurrency),
		"title":                name,
		"description":          description,
		"order_id":             createPayContext.Pay.PaymentId,
		"customer_id":          strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"customer_email":       createPayContext.Email,
		"success_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"failure_redirect_url": webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		"fees_payer":           gasPayer, // who pay the fee
		"payment_data":         createPayContext.Metadata,
		"pending_deadline_at":  time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	responseJson, err := SendChangellyPaymentRequest(ctx, createPayContext.Gateway.GatewayKey, createPayContext.Gateway.GatewaySecret, "POST", urlPath, param)
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

func (c Changelly) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api/payment/v1/payments/" + gatewayPaymentId
	param := map[string]interface{}{}
	responseJson, err := SendChangellyPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "ChangelyPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return parseChangellyPayment(responseJson), nil
}

func (c Changelly) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Changelly) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (c Changelly) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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
	if item.Contains("selected_payment_method") && item.GetJson("selected_payment_method").Contains("payins") {
		//for _, payin := range item.GetJson("selected_payment_method").GetJsons("payins") {
		//	paymentAmount = paymentAmount + utility.ConvertDollarStrToCent(payin.Get("amount").String(), item.Get("nominal_currency").String())
		//	paymentAmount = paymentAmount - utility.ConvertDollarStrToCent(payin.Get("fee").String(), item.Get("nominal_currency").String())
		//}
		//paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("selected_payment_method").Get("expected_payin_amount").String(), item.Get("nominal_currency").String())
		paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("nominal_amount").String(), item.Get("nominal_currency").String())
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
		PaidTime:             paidTime,
	}
}

func SendChangellyFiatRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")
	datetime := getExpirationDateTime(1)

	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nChangelly_Start %s %s %s %s %s\n", method, urlPath, publicKey, jsonString, datetime)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type":    "application/json",
		"X-Api-Signature": signForFiat(method, "https://fiat-api.changelly.com"+urlPath, datetime, privateKey, body),
		"X-Api-Key":       publicKey,
	}
	response, err := utility.SendRequest("https://fiat-api.changelly.com"+urlPath, method, body, headers)
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

func SendChangellyPaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
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

func signForFiat(method string, urlPath string, dateTime string, purePrivateKey string, postJson []byte) (sign string) {
	var builder strings.Builder

	builder.WriteString(urlPath)
	builder.WriteString(string(postJson))
	payload := builder.String()
	decodedBytes, err := base64.StdEncoding.DecodeString(purePrivateKey)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}
	purePrivateKey = string(decodedBytes)
	block, _ := pem.Decode([]byte(purePrivateKey))
	utility.Assert(block != nil, "rsa encrypt error")
	prv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	utility.AssertError(err, "rsa encrypt error")
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(payload))
	utility.AssertError(err, "sha256 hash encrypt error")
	result, err := rsa.SignPKCS1v15(rand.Reader, prv.(*rsa.PrivateKey), crypto.SHA256, msgHash.Sum(nil))
	//result, err := utility.RsaEncrypt([]byte(key), []byte(sha256Encoding(builder.String())))
	utility.AssertError(err, "rsa encrypt error")
	return base64Encoding(result)
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
