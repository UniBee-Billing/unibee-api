package evonet

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/gateway/log"
	"go-oversea-pay/internal/logic/payment/gateway/out"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
	"go-oversea-pay/internal/logic/payment/gateway/util"
	"go-oversea-pay/internal/logic/payment/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//const DEV_ENDPOINT = "https://hkg-online-uat.everonet.com"
//const PROD_ENDPOINT = "https://hkg-online.everonet.com"

type Evonet struct{}

func (e Evonet) DoRemoteChannelCustomerBalanceQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, customerId string) (res *ro.ChannelCustomerBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoiceCreate(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelCreateInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionNewTrailEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrailEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
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

func (e Evonet) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (e Evonet) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	g.Log().Infof(r.Context(), "EvonetNotifyController 收到 channel_webhook_entry 结果通知:%s", r.GetBody())
	notificationJson, err := r.GetJson()
	if err != nil {
		r.Response.Writeln(err)
	}
	g.Log().Infof(r.Context(), "EvonetNotifyController channel_webhook_entry notifications:%s", notificationJson)
	if notificationJson == nil {
		r.Response.Writeln("invalid body")
	}
	eventCode := notificationJson.Get("eventCode").String()
	paymentMethod := notificationJson.GetJson("paymentMethod")
	payment := notificationJson.GetJson("payment")
	capture := notificationJson.GetJson("capture")
	cancel := notificationJson.GetJson("cancel")
	refund := notificationJson.GetJson("refund")
	executeResult := false
	notifyEventCode := ""
	notifyEventDate := gtime.Now()
	notifyReason := ""
	notifyMerchantReference := ""
	notifyIsSuccess := false
	if strings.Compare(eventCode, "Payment") == 0 &&
		payment != nil &&
		payment.Contains("status") &&
		strings.Compare(payment.Get("status").String(), "Authorised") == 0 &&
		payment.Contains("merchantTransInfo") &&
		payment.GetJson("merchantTransInfo").Contains("merchantTransID") &&
		payment.Contains("evoTransInfo") &&
		payment.GetJson("evoTransInfo").Contains("evoTransID") {
		//授权成功
		data := payment
		merchantTradeNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		channelTradeNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetOverseaPayByMerchantOrderNo(r.Context(), merchantTradeNo)
		transAmount := data.GetJson("transAmount")
		if one != nil && transAmount != nil &&
			one.PaymentFee == utility.ConvertYuanStrToFen(transAmount.Get("value").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			one.ChannelTradeNo = channelTradeNo
			err := handler.HandlePayAuthorized(r.Context(), one)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s handlePayAuthorized object:%s hook:%s err:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.MerchantOrderNo
				notifyIsSuccess = true
				notifyEventCode = "AUTHORISATION"
				executeResult = true
			}
		} else {
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s not match object:%s hook:%s", eventCode, one, notificationJson.String())
		}
	} else if (strings.Compare(eventCode, "Cancel") == 0 && cancel != nil &&
		cancel.Contains("status") &&
		strings.Compare(cancel.Get("status").String(), "Success") == 0 &&
		cancel.Contains("merchantTransInfo") &&
		cancel.GetJson("merchantTransInfo").Contains("merchantTransID") &&
		cancel.Contains("evoTransInfo") &&
		cancel.GetJson("evoTransInfo").Contains("evoTransID")) ||
		(strings.Compare(eventCode, "Payment") == 0 &&
			payment != nil &&
			payment.Contains("status") &&
			strings.Compare(payment.Get("status").String(), "Failed") == 0 &&
			payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") &&
			payment.Contains("evoTransInfo") &&
			payment.GetJson("evoTransInfo").Contains("evoTransID")) {
		//取消成功
		data := &gjson.Json{}
		if strings.Compare(eventCode, "Cancel") == 0 {
			data = cancel
		} else if strings.Compare(eventCode, "Payment") == 0 {
			data = payment
		} else {
			utility.Assert(true, fmt.Sprintf("data eventCode error notificationJson:%s", notificationJson.String()))
		}
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantTradeNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		//从 payment 获取更加精确
		if payment != nil && payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") {
			merchantTradeNo = payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
		}
		channelTradeNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetOverseaPayByMerchantOrderNo(r.Context(), merchantTradeNo)
		if one != nil && len(one.ChannelTradeNo) > 0 {
			utility.Assert(strings.Compare(channelTradeNo, one.ChannelTradeNo) == 0, "channelTradeNo not match")
		}

		reason := fmt.Sprintf("from_webhook:%s", data.Get("failureReason").String())
		if one != nil {
			req := &handler.HandlePayReq{
				MerchantOrderNo: merchantTradeNo,
				ChannelTradeNo:  channelTradeNo,
				PayStatusEnum:   consts.PAY_FAILED,
				Reason:          reason,
			}
			err := handler.HandlePayFailure(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.MerchantOrderNo
				notifyIsSuccess = true
				notifyEventCode = "CANCELLATION"
				notifyReason = reason
				executeResult = true
			}
		} else {
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s not match object:%s hook:%s", eventCode, one, notificationJson.String())
		}

	} else if (strings.Compare(eventCode, "Capture") == 0 && capture != nil &&
		capture.Contains("status") &&
		strings.Compare(capture.Get("status").String(), "Success") == 0 &&
		capture.Contains("merchantTransInfo") &&
		capture.GetJson("merchantTransInfo").Contains("merchantTransID") &&
		capture.Contains("evoTransInfo") &&
		capture.GetJson("evoTransInfo").Contains("evoTransID")) ||
		(strings.Compare(eventCode, "Payment") == 0 &&
			payment != nil &&
			payment.Contains("status") &&
			strings.Compare(payment.Get("status").String(), "Captured") == 0 &&
			payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") &&
			payment.Contains("evoTransInfo") &&
			payment.GetJson("evoTransInfo").Contains("evoTransID")) {
		//捕获成功
		data := &gjson.Json{}
		if strings.Compare(eventCode, "Capture") == 0 {
			data = cancel
		} else if strings.Compare(eventCode, "Payment") == 0 {
			data = payment
		} else {
			utility.Assert(true, fmt.Sprintf("data eventCode error notificationJson:%s", notificationJson.String()))
		}
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantTradeNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		//从 payment 获取更加精确
		if payment != nil && payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") {
			merchantTradeNo = payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
		}
		channelTradeNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetOverseaPayByMerchantOrderNo(r.Context(), merchantTradeNo)
		transAmount := data.GetJson("transAmount")
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertYuanStrToFen(transAmount.Get("value").String()) > 0 &&
			one.PaymentFee == utility.ConvertYuanStrToFen(transAmount.Get("value").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			receiveFee := utility.ConvertYuanStrToFen(transAmount.Get("value").String())
			req := &handler.HandlePayReq{
				MerchantOrderNo: merchantTradeNo,
				ChannelTradeNo:  channelTradeNo,
				PayStatusEnum:   consts.PAY_SUCCESS,
				PayFee:          one.PaymentFee,
				ReceiptFee:      receiveFee,
				PaidTime:        gtime.Now(),
			}
			err := handler.HandlePaySuccess(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.MerchantOrderNo
				notifyIsSuccess = true
				notifyEventCode = "CAPTURE"
				executeResult = true
			}
		} else {
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s not match object:%s hook:%s", eventCode, one, notificationJson.String())
		}
	} else if strings.Compare(eventCode, "Refund") == 0 &&
		refund != nil &&
		refund.Contains("status") &&
		strings.Compare(refund.Get("status").String(), "Success") == 0 &&
		refund.Contains("merchantTransInfo") &&
		refund.GetJson("merchantTransInfo").Contains("merchantTransID") &&
		refund.Contains("evoTransInfo") &&
		refund.GetJson("evoTransInfo").Contains("evoTransID") {
		//退款成功
		data := refund
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantRefundNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		channelRefundNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetOverseaRefundByMerchantRefundNo(r.Context(), merchantRefundNo)
		transAmount := data.GetJson("transAmount")
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertYuanStrToFen(transAmount.Get("value").String()) > 0 &&
			one.RefundFee == utility.ConvertYuanStrToFen(transAmount.Get("value").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			req := &handler.HandleRefundReq{
				MerchantRefundNo: merchantRefundNo,
				ChannelRefundNo:  channelRefundNo,
				RefundStatusEnum: consts.REFUND_SUCCESS,
				RefundTime:       gtime.Now(),
			}
			err := handler.HandleRefundSuccess(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantRefundNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.OutRefundNo
				notifyIsSuccess = true
				notifyEventCode = "REFUND"
				executeResult = true
			}
		} else {
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s not match object:%s hook:%s", eventCode, one, notificationJson.String())
		}
	} else if strings.Compare(eventCode, "Refund") == 0 &&
		refund != nil &&
		refund.Contains("status") &&
		strings.Compare(refund.Get("status").String(), "Failed") == 0 &&
		refund.Contains("merchantTransInfo") &&
		refund.GetJson("merchantTransInfo").Contains("merchantTransID") &&
		refund.Contains("evoTransInfo") &&
		refund.GetJson("evoTransInfo").Contains("evoTransID") {
		//退款失败
		data := refund
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantRefundNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		channelRefundNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetOverseaRefundByMerchantRefundNo(r.Context(), merchantRefundNo)
		transAmount := data.GetJson("transAmount")
		reason := fmt.Sprintf("from_webhook:%s", data.Get("failureReason").String())
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertYuanStrToFen(transAmount.Get("value").String()) > 0 &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			req := &handler.HandleRefundReq{
				MerchantRefundNo: merchantRefundNo,
				ChannelRefundNo:  channelRefundNo,
				RefundStatusEnum: consts.REFUND_FAILED,
				RefundTime:       gtime.Now(),
				Reason:           reason,
			}
			err := handler.HandleRefundFailure(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantRefundNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.OutRefundNo
				notifyIsSuccess = true
				notifyEventCode = "REFUND_FAILED"
				executeResult = true
			}
		} else {
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s not match object:%s hook:%s", eventCode, one, notificationJson.String())
		}
	} else {
		requestId := strconv.FormatInt(utility.CurrentTimeMillis(), 10)
		if paymentMethod != nil &&
			paymentMethod.Contains("merchantTransInfo") &&
			paymentMethod.GetJson("merchantTransInfo").Contains("merchantTransID") {
			requestId = paymentMethod.GetJson("merchantTransInfo").Get("merchantTransID").String()
		} else if payment != nil &&
			payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") {
			requestId = payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
		} else if cancel != nil &&
			cancel.Contains("merchantTransInfo") &&
			cancel.GetJson("merchantTransInfo").Contains("merchantTransID") {
			requestId = cancel.GetJson("merchantTransInfo").Get("merchantTransID").String()
		} else if capture != nil &&
			capture.Contains("merchantTransInfo") &&
			capture.GetJson("merchantTransInfo").Contains("merchantTransID") {
			requestId = capture.GetJson("merchantTransInfo").Get("merchantTransID").String()
		} else if refund != nil &&
			refund.Contains("merchantTransInfo") &&
			refund.GetJson("merchantTransInfo").Contains("merchantTransID") {
			requestId = refund.GetJson("merchantTransInfo").Get("merchantTransID").String()
		}

		log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", "", eventCode, requestId, "evonet webhook")
	}
	if executeResult {
		//向商户推送渠道原消息
		//Message message = new Message(MqTopicEnum.ChannelPayV2WebHookReceive,new PaymentNotificationMqItem(notifyIsSuccess,notifyEventDate,notifyEventCode,notifyMerchantReference,notifyReason));
		//boolean sendResult = producerWrapper.send(message);
	}
	g.Log().Infof(r.Context(), "channel_webhook_entry execute result:%s %s %s %s %s ", notifyIsSuccess, notifyEventDate, notifyEventCode, notifyMerchantReference, notifyReason)

	r.Response.Writeln("success")
}

func (e Evonet) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	payIdStr := r.Get("payId").String()
	redirectResult := r.Get("redirectResult").String()
	g.Log().Printf(r.Context(), "EvonetNotifyController evonet_redirect payId:%s redirectResult:%s", payIdStr, redirectResult)
	log.DoSaveChannelLog(r.Context(), payIdStr, "redirect", redirectResult, "redirect", payIdStr, "evonet redirect")
	utility.Assert(len(payIdStr) > 0, "参数错误，payId未传入")
	payId, err := strconv.Atoi(payIdStr)
	utility.Assert(err == nil, "参数错误，payId 需 int 类型 %s")
	overseaPay := query.GetOverseaPayById(r.Context(), int64(payId))
	utility.Assert(overseaPay != nil, fmt.Sprintf("找不到支付单 payId: %s", payIdStr))
	channelEntity := query.GetPaymentTypePayChannelById(r.Context(), overseaPay.ChannelId)
	utility.Assert(channelEntity != nil, fmt.Sprintf("支付渠道异常 payId: %s", payIdStr))
	g.Log().Infof(r.Context(), "DoRemoteChannelRedirect payId:%s notifyUrl:%s", payIdStr, overseaPay.NotifyUrl)
	if len(overseaPay.NotifyUrl) > 0 {
		r.Response.Writeln(fmt.Sprintf("<head>\n<meta http-equiv=\"refresh\" content=\"1;url=%s\">\n</head>", overseaPay.NotifyUrl))
	} else {
		//r.Response.Writeln(r.Get("channelId"))
		r.Response.Writeln(overseaPay)
	}
	return &ro.ChannelRedirectInternalResp{
		Status:  true,
		Message: "",
	}, nil
}

func (e Evonet) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.Pay != nil, "pay is nil")
	utility.Assert(createPayContext.PayChannel != nil, "pay gateway config is nil")

	//其他渠道所需参数校验
	utility.Assert(len(createPayContext.Pay.CountryCode) > 0, "countryCode is nil")
	utility.Assert(createPayContext.Pay.PaymentFee > 0, "paymentFee is nil")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(len(createPayContext.Pay.MerchantOrderNo) > 0, "merchantOrderNo is nil")
	utility.Assert(len(createPayContext.ShopperEmail) > 0, "shopperEmail is nil")
	utility.Assert(len(createPayContext.UserId) > 0, "shopperReference is nil")
	utility.Assert(createPayContext.Items != nil, "lineItems is nil")
	urlPath := "/g2/auth/payment/mer/" + createPayContext.PayChannel.ChannelAccountId + "/evo.e-commerce.payment"
	channelType := createPayContext.PayChannel.SubChannel
	if len(channelType) == 0 {
		channelType = createPayContext.PayChannel.Channel
	}
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":        createPayContext.Pay.MerchantOrderNo,
			"merchantOrderReference": createPayContext.MerchantOrderReference,
			"merchantTransTime":      getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": createPayContext.Pay.Currency,
			"value":    utility.ConvertFenToYuanMinUnitStr(createPayContext.Pay.PaymentFee),
		},
		"paymentMethod": map[string]interface{}{
			"recurringProcessingModel": createPayContext.RecurringProcessingModel,
			"type":                     channelType,
		},
		"userinfo": map[string]interface{}{
			"email":     createPayContext.ShopperEmail,
			"reference": createPayContext.UserId,
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
			"totalQuantity": fmt.Sprintf("%d", len(createPayContext.Items)),
		},
		"returnURL": out.GetPaymentRedirectEntranceUrl(createPayContext.Pay),
		"webhook":   out.GetPaymentWebhookEntranceUrl(createPayContext.Pay.ChannelId),
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
		Action:         responseJson.GetJson("action"),
		AdditionalData: responseJson.GetJson("paymentMethod"),
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.capture" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreateMerchantOrderNo(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": pay.Currency,
			"value":    utility.ConvertFenToYuanMinUnitStr(pay.BuyerPayFee),
		},
		"webhook": out.GetPaymentWebhookEntranceUrl(pay.ChannelId),
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
	pspReference := captureJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "capture", responseJson.String(), "支付捕获", pspReference, channelEntity.Channel)
	res = &ro.OutPayCaptureRo{
		PspReference: pspReference,
		Status:       status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.cancel" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreateMerchantOrderNo(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"webhook": out.GetPaymentWebhookEntranceUrl(pay.ChannelId),
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
	pspReference := cancelJson.GetJson("evoTransInfo").Get("evoTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "cancel", responseJson.String(), "支付取消", pspReference, channelEntity.Channel)
	res = &ro.OutPayCancelRo{
		PspReference: pspReference,
		Status:       status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.payment"
	param := map[string]interface{}{
		"merchantTransID": pay.MerchantOrderNo,
	}
	data, err := sendEvonetRequest(ctx, channelEntity, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay支付查询失败 result is nil")
	resultJson := responseJson.GetJson("result")
	payment := responseJson.GetJson("payment")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		payment != nil &&
		payment.Contains("status") &&
		payment.Contains("evoTransInfo") &&
		payment.GetJson("evoTransInfo").Contains("evoTransID") &&
		payment.GetJson("merchantTransInfo").Contains("merchantTransID"),
		fmt.Sprintf("Evonetpay支付查询失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := payment.Get("status").String()
	pspReference := payment.GetJson("evoTransInfo").Get("evoTransID").String()
	merchantPspReference := payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "payment_query", responseJson.String(), "支付查询", pspReference, channelEntity.Channel)
	utility.Assert(strings.Compare(merchantPspReference, pay.MerchantOrderNo) == 0, "merchantPspReference not match")
	res = &ro.OutPayRo{
		PayFee:    pay.PaymentFee,
		PayStatus: consts.TO_BE_PAID,
	}
	if strings.Compare(status, "Failed") == 0 || strings.Compare(status, "Cancelled") == 0 {
		res.PayStatus = consts.PAY_FAILED
		res.Reason = "from_query:" + payment.Get("failureReason").String()
	} else if strings.Compare(status, "Captured") == 0 {
		res.PayStatus = consts.PAY_SUCCESS
		res.ChannelTradeNo = pspReference
		res.PayTime = gtime.Now()
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.refund" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   refund.OutRefundNo,
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": pay.Currency,
			"value":    utility.ConvertFenToYuanMinUnitStr(refund.RefundFee),
		},
		"webhook": out.GetPaymentWebhookEntranceUrl(pay.ChannelId),
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
		ChannelRefundNo: pspReference,
		RefundStatus:    consts.REFUND_ING,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	urlPath := "/g2/auth/payment/mer/" + channelEntity.ChannelAccountId + "/evo.e-commerce.refund"
	param := map[string]interface{}{
		"merchantTransID": refund.OutRefundNo,
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
	utility.Assert(strings.Compare(merchantPspReference, refund.OutRefundNo) == 0, "merchantPspReference not match")
	log.DoSaveChannelLog(ctx, log.ConvertToStringIgnoreErr(param), "refund_query", responseJson.String(), "退款查询", pspReference, channelEntity.Channel)
	res = &ro.OutPayRefundRo{
		RefundFee:    refund.RefundFee,
		RefundStatus: consts.REFUND_ING,
	}
	if strings.Compare(status, "Failed") == 0 {
		res.RefundStatus = consts.REFUND_FAILED
		res.Reason = "from_query:" + refundJson.Get("failureReason").String()
	} else if strings.Compare(status, "Success") == 0 {
		res.RefundStatus = consts.REFUND_SUCCESS
		res.ChannelRefundNo = pspReference
		res.RefundTime = gtime.Now()
	}
	return res, nil
}

func sendEvonetRequest(ctx context.Context, channelEntity *entity.OverseaPayChannel, method string, urlPath string, param map[string]interface{}) (res []byte, err error) {
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
