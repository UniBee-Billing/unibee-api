package out

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/channel"
	"go-oversea-pay/internal/logic/channel/log"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/channel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"
)

//const DEV_ENDPOINT = "https://hkg-online-uat.everonet.com"
//const PROD_ENDPOINT = "https://hkg-online.everonet.com"

type Evonet struct{}

func (e Evonet) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantChannelConfig, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantChannelConfig, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantChannelConfig, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantChannelConfig, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantChannelConfig, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.Pay != nil, "payment  is nil")
	utility.Assert(createPayContext.PayChannel != nil, "payment channel config is nil")

	//其他渠道所需参数校验
	utility.Assert(len(createPayContext.Pay.CountryCode) > 0, "countryCode is nil")
	utility.Assert(createPayContext.Pay.TotalAmount > 0, "TotalAmount is nil")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(len(createPayContext.Pay.PaymentId) > 0, "paymentId is nil")
	utility.Assert(len(createPayContext.ShopperEmail) > 0, "shopperEmail is nil")
	utility.Assert(len(createPayContext.ShopperUserId) > 0, "shopperUserId is nil")
	utility.Assert(createPayContext.Invoice.Lines != nil, "lineItems is nil")
	urlPath := "/g2/auth/payment/mer/" + createPayContext.PayChannel.ChannelAccountId + "/evo.e-commerce.payment"
	channelType := createPayContext.PayChannel.SubChannel
	if len(channelType) == 0 {
		channelType = createPayContext.PayChannel.Channel
	}
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":        createPayContext.Pay.PaymentId,
			"merchantOrderReference": createPayContext.MerchantOrderReference,
			"merchantTransTime":      getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": createPayContext.Pay.Currency,
			"value":    utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		},
		"paymentMethod": map[string]interface{}{
			"recurringProcessingModel": createPayContext.RecurringProcessingModel,
			"type":                     channelType,
		},
		"userinfo": map[string]interface{}{
			"email":     createPayContext.ShopperEmail,
			"reference": createPayContext.ShopperUserId,
		},
		"transInitiator": map[string]interface{}{
			"browserInfo": map[string]interface{}{
				"language": createPayContext.ShopperLocale,
			},
			"platform":   createPayContext.Platform,
			"deviceType": createPayContext.DeviceType,
		},
		"tradeInfo": map[string]string{
			"tradeType":     "Sale of goods",
			"totalQuantity": fmt.Sprintf("%d", len(createPayContext.Invoice.Lines)),
		},
		"returnURL": channel.GetPaymentRedirectEntranceUrl(createPayContext.Pay),
		"webhook":   channel.GetPaymentWebhookEntranceUrl(createPayContext.Pay.ChannelId),
	}

	if len(createPayContext.PayChannel.BrandData) > 0 {
		var data map[string]interface{}
		// 使用 json.Unmarshal 解析 JSON 字符串
		err := json.Unmarshal([]byte(createPayContext.PayChannel.BrandData), &data)
		if err == nil {
			//_, ok := data[channelType]
			//if ok {
			//
			//}
			// payeasy 和 payeasy 渠道支持 todo mark
			for key, value := range data {
				param["paymentMethod"].(map[string]interface{})[key] = value
			}
		}
	}

	if createPayContext.ShopperName != nil && len(createPayContext.ShopperName.FirstName) > 0 && len(createPayContext.ShopperName.LastName) > 0 {
		param["userinfo"].(map[string]interface{})["name"] = createPayContext.ShopperName.FirstName + " " + createPayContext.ShopperName.LastName
	}
	match, _ := regexp.MatchString(createPayContext.Mobile, "[0-9]+")
	if len(createPayContext.Mobile) > 0 && match {
		param["authentication"] = map[string]interface{}{
			"securePlus": map[string]string{
				"mobilePhone": createPayContext.Mobile,
			},
		}
	}

	if createPayContext.Pay.CaptureDelayHours > 0 {
		param["captureAfterHours"] = createPayContext.Pay.CaptureDelayHours
	}

	data, err := sendEvonetRequest(ctx, createPayContext.PayChannel, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay捕获失败 result is nil")
	resultJson := responseJson.GetJson("result")
	paymentJson := responseJson.GetJson("payment")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		paymentJson != nil &&
		paymentJson.Contains("evoTransInfo") &&
		paymentJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay字符失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	//status := paymentJson.Get("status").String()
	pspReference := paymentJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "payments", responseJson.String(), "支付", pspReference, createPayContext.PayChannel.Channel)
	res = &ro.CreatePayInternalResp{
		Status: consts.TO_BE_PAID,
		Action: responseJson.GetJson("action"),
		Link:   responseJson.GetJson("action").Get("url").String(),
		//AdditionalData: responseJson.GetJson("paymentMethod"),
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	utility.Assert(payment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, payment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.capture" + "?merchantTransID=" + payment.PaymentId
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreatePaymentId(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": payment.Currency,
			"value":    utility.ConvertCentToDollarStr(payment.PaymentAmount, payment.Currency),
		},
		"webhook": channel.GetPaymentWebhookEntranceUrl(payment.ChannelId),
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay捕获失败 result is nil")
	resultJson := responseJson.GetJson("result")
	captureJson := responseJson.GetJson("capture")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		captureJson != nil &&
		captureJson.Contains("evoTransInfo") &&
		captureJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay捕获失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := captureJson.Get("status").String()
	channelCaptureId := captureJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "capture", responseJson.String(), "支付捕获", channelCaptureId, channelEntity.Channel)
	res = &ro.OutPayCaptureRo{
		ChannelCaptureId: channelCaptureId,
		Status:           status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	utility.Assert(payment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, payment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.cancel" + "?merchantTransID=" + payment.PaymentId
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreatePaymentId(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"webhook": channel.GetPaymentWebhookEntranceUrl(payment.ChannelId),
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay取消失败 result is nil")
	resultJson := responseJson.GetJson("result")
	cancelJson := responseJson.GetJson("cancel")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		cancelJson != nil &&
		cancelJson.Contains("evoTransInfo") &&
		cancelJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay取消失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := cancelJson.Get("status").String()
	channelCancelId := cancelJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "cancel", responseJson.String(), "支付取消", channelCancelId, channelEntity.Channel)
	res = &ro.OutPayCancelRo{
		ChannelCancelId: channelCancelId,
		Status:          status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	utility.Assert(payment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, payment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.payment"
	param := map[string]interface{}{
		"merchantTransID": payment.PaymentId,
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay支付查询失败 result is nil")
	resultJson := responseJson.GetJson("result")
	channelPayment := responseJson.GetJson("payment")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		channelPayment != nil &&
		channelPayment.Contains("status") &&
		channelPayment.Contains("evoTransInfo") &&
		channelPayment.GetJson("evoTransInfo").Contains("evoTransID") &&
		channelPayment.GetJson("merchantTransInfo").Contains("merchantTransID"),
		fmt.Sprintf("Evonetpay支付查询失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := channelPayment.Get("status").String()
	pspReference := channelPayment.GetJson("evoTransInfo").Get("evoTransID").String()
	merchantPspReference := channelPayment.GetJson("merchantTransInfo").Get("merchantTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "payment_query", responseJson.String(), "支付查询", pspReference, channelEntity.Channel)
	utility.Assert(strings.Compare(merchantPspReference, payment.PaymentId) == 0, "merchantPspReference not match")
	res = &ro.ChannelPaymentRo{
		TotalAmount: payment.TotalAmount,
		Status:      consts.TO_BE_PAID,
	}
	if strings.Compare(status, "Failed") == 0 || strings.Compare(status, "Cancelled") == 0 {
		res.Status = consts.PAY_FAILED
		res.Reason = "from_query:" + channelPayment.Get("failureReason").String()
	} else if strings.Compare(status, "Captured") == 0 {
		res.Status = consts.PAY_SUCCESS
		res.ChannelPaymentId = pspReference
		res.PayTime = gtime.Now()
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefund(ctx context.Context, channelPayment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(channelPayment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, channelPayment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.refund" + "?merchantTransID=" + channelPayment.PaymentId
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   refund.RefundId,
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": channelPayment.Currency,
			"value":    utility.ConvertCentToDollarStr(refund.RefundAmount, channelPayment.Currency),
		},
		"webhook": channel.GetPaymentWebhookEntranceUrl(channelPayment.ChannelId),
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay退款失败 result is nil")
	resultJson := responseJson.GetJson("result")
	refundJson := responseJson.GetJson("refund")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		refundJson != nil &&
		refundJson.Contains("evoTransInfo") &&
		refundJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay取消失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	pspReference := refundJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "refund", responseJson.String(), "支付退款", pspReference, channelEntity.Channel)
	res = &ro.OutPayRefundRo{
		ChannelRefundId: pspReference,
		Status:          consts.REFUND_ING,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefundStatusCheck(ctx context.Context, channelPayment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(channelPayment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, channelPayment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.refund"
	param := map[string]interface{}{
		"merchantTransID": refund.RefundId,
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay退款查询失败 result is nil")
	resultJson := responseJson.GetJson("result")
	refundJson := responseJson.GetJson("refund")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		refundJson != nil &&
		refundJson.Contains("status") &&
		refundJson.Contains("evoTransInfo") &&
		refundJson.GetJson("evoTransInfo").Contains("evoTransID") &&
		refundJson.GetJson("merchantTransInfo").Contains("merchantTransID"),
		fmt.Sprintf("Evonetpay退款查询失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := refundJson.Get("status").String()
	pspReference := refundJson.GetJson("evoTransInfo").Get("evoTransID").String()
	merchantPspReference := refundJson.GetJson("merchantTransInfo").Get("merchantTransID").String()
	utility.Assert(strings.Compare(merchantPspReference, refund.RefundId) == 0, "merchantPspReference not match")
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "refund_query", responseJson.String(), "退款查询", pspReference, channelEntity.Channel)
	res = &ro.OutPayRefundRo{
		RefundFee: refund.RefundAmount,
		Status:    consts.REFUND_ING,
	}
	if strings.Compare(status, "Failed") == 0 {
		res.Status = consts.REFUND_FAILED
		res.Reason = "from_query:" + refundJson.Get("failureReason").String()
	} else if strings.Compare(status, "Success") == 0 {
		res.Status = consts.REFUND_SUCCESS
		res.ChannelRefundId = pspReference
		res.RefundTime = gtime.Now()
	}
	return res, nil
}

func sendEvonetRequest(ctx context.Context, channelEntity *entity.MerchantChannelConfig, method string, urlPath string, param map[string]interface{}) (res []byte, err error) {
	utility.Assert(param != nil, "param is nil")
	// 定义自定义的头部信息
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Infof(ctx, "\nEvonet_Start %s %s %s %s\n", method, urlPath, channelEntity.ChannelKey, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Channel": "application/json",
		"Msgid":           msgId,
		"Datetime":        datetime,
		"Authorization":   sign("POST", urlPath, msgId, datetime, channelEntity.ChannelKey, body),
		"Signtype":        "SHA256",
	}
	response, err := sendRequest(channelEntity.Host+urlPath, method, body, headers)
	g.Log().Infof(ctx, "\nEvonet_End %s %s response: %s error %s\n", method, urlPath, response, err)
	return response, nil
}

func sendRequest(url string, method string, data []byte, headers map[string]string) ([]byte, error) {
	// 创建一个字节数组读取器，用于将数据传递给请求体
	bodyReader := bytes.NewReader(data)

	// 创建一个POST请求
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置自定义头部信息
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// 关闭响应体
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func sign(method string, urlPath string, msgId string, dateTime string, key string, postJson []byte) (sign string) {
	var builder strings.Builder
	lineSeparator := lineSeparator()
	builder.WriteString(method)
	builder.WriteString(lineSeparator)
	builder.WriteString(urlPath)
	builder.WriteString(lineSeparator)
	builder.WriteString(dateTime)
	builder.WriteString(lineSeparator)
	builder.WriteString(key)
	builder.WriteString(lineSeparator)
	builder.WriteString(msgId)
	if postJson != nil {
		builder.WriteString(lineSeparator)
		builder.Write(postJson)
	}
	return sha256Encoding(builder.String())
}

func generateMsgId() (msgId string) {
	return fmt.Sprintf("%s%s%d", utility.JodaTimePrefix(), utility.GenerateRandomAlphanumeric(5), utility.CurrentTimeMillis())
}

func getCurrentDateTime() (datetime string) {
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}

func lineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func sha256Encoding(data string) (hash string) {
	hasher := sha256.New()

	// 添加数据到散列器
	hasher.Write([]byte(data))

	// 计算散列值
	hashValue := hasher.Sum(nil)

	// 将散列值转换为十六进制字符串
	return hex.EncodeToString(hashValue)
}
