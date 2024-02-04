package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/channel/log"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/payment/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

type EvonetWebhook struct {
}

func (e EvonetWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantChannelConfig) (err error) {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (e EvonetWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) {
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
		channelPaymentId := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetPaymentByPaymentId(r.Context(), merchantTradeNo)
		transAmount := data.GetJson("transAmount")
		if one != nil && transAmount != nil &&
			one.TotalAmount == utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			one.ChannelPaymentId = channelPaymentId
			err := handler.HandlePayAuthorized(r.Context(), one)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s handlePayAuthorized object:%s hook:%s err:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.PaymentId
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
			utility.Assert(false, fmt.Sprintf("data eventCode error notificationJson:%s", notificationJson.String()))
		}
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantTradeNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		//从 payment 获取更加精确
		if payment != nil && payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") {
			merchantTradeNo = payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
		}
		channelTradeNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetPaymentByPaymentId(r.Context(), merchantTradeNo)
		if one != nil && len(one.ChannelPaymentId) > 0 {
			utility.Assert(strings.Compare(channelTradeNo, one.ChannelPaymentId) == 0, "channelPaymentId not match")
		}

		reason := fmt.Sprintf("from_webhook:%s", data.Get("failureReason").String())
		if one != nil {
			req := &handler.HandlePayReq{
				PaymentId:        merchantTradeNo,
				ChannelPaymentId: channelTradeNo,
				PayStatusEnum:    consts.PAY_FAILED,
				Reason:           reason,
			}
			err := handler.HandlePayFailure(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.PaymentId
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
			utility.Assert(false, fmt.Sprintf("data eventCode error notificationJson:%s", notificationJson.String()))
		}
		utility.Assert(data != nil, fmt.Sprintf("data is nil  notificationJson:%s", notificationJson.String()))
		merchantTradeNo := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		//从 payment 获取更加精确
		if payment != nil && payment.Contains("merchantTransInfo") &&
			payment.GetJson("merchantTransInfo").Contains("merchantTransID") {
			merchantTradeNo = payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
		}
		channelTradeNo := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetPaymentByPaymentId(r.Context(), merchantTradeNo)
		transAmount := data.GetJson("transAmount")
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) > 0 &&
			one.TotalAmount == utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			receiveFee := utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String())
			req := &handler.HandlePayReq{
				PaymentId:        merchantTradeNo,
				ChannelPaymentId: channelTradeNo,
				PayStatusEnum:    consts.PAY_SUCCESS,
				TotalAmount:      one.TotalAmount,
				PaymentAmount:    receiveFee,
				PaidTime:         gtime.Now(),
			}
			err := handler.HandlePaySuccess(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantTradeNo, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.PaymentId
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
		merchantRefundId := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		channelRefundId := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetRefundByRefundId(r.Context(), merchantRefundId)
		transAmount := data.GetJson("transAmount")
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) > 0 &&
			one.RefundAmount == utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			req := &handler.HandleRefundReq{
				RefundId:         merchantRefundId,
				ChannelRefundId:  channelRefundId,
				RefundStatusEnum: consts.REFUND_SUCCESS,
				RefundTime:       gtime.Now(),
			}
			err := handler.HandleRefundSuccess(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantRefundId, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.RefundId
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
		merchantRefundId := data.GetJson("merchantTransInfo").Get("merchantTransID").String()
		channelRefundId := data.GetJson("evoTransInfo").Get("evoTransID").String()
		one := query.GetRefundByRefundId(r.Context(), merchantRefundId)
		transAmount := data.GetJson("transAmount")
		reason := fmt.Sprintf("from_webhook:%s", data.Get("failureReason").String())
		if one != nil &&
			transAmount != nil &&
			len(transAmount.Get("value").String()) > 0 &&
			utility.ConvertDollarStrToCent(transAmount.Get("value").String(), transAmount.Get("currency").String()) > 0 &&
			strings.Compare(one.Currency, transAmount.Get("currency").String()) == 0 {
			req := &handler.HandleRefundReq{
				RefundId:         merchantRefundId,
				ChannelRefundId:  channelRefundId,
				RefundStatusEnum: consts.REFUND_FAILED,
				RefundTime:       gtime.Now(),
				Reason:           reason,
			}
			err := handler.HandleRefundFailure(r.Context(), req)
			log.DoSaveChannelLog(r.Context(), notificationJson.String(), "webhook", strconv.FormatBool(err == nil), eventCode, merchantRefundId, "evonet webhook")
			g.Log().Infof(r.Context(), "channel_webhook_entry action:%s do success object:%s hook:%s result:%s", eventCode, one, notificationJson.String(), err)
			if err != nil {
				executeResult = false
			} else {
				notifyMerchantReference = one.RefundId
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

func (e EvonetWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelRedirectInternalResp, err error) {
	payIdStr := r.Get("payId").String()
	redirectResult := r.Get("redirectResult").String()
	g.Log().Printf(r.Context(), "EvonetNotifyController evonet_redirect payId:%s redirectResult:%s", payIdStr, redirectResult)
	log.DoSaveChannelLog(r.Context(), payIdStr, "redirect", redirectResult, "redirect", payIdStr, "evonet redirect")
	utility.Assert(len(payIdStr) > 0, "参数错误，payId未传入")
	payId, err := strconv.Atoi(payIdStr)
	utility.Assert(err == nil, "参数错误，payId 需 int 类型 %s")
	overseaPay := query.GetPaymentById(r.Context(), int64(payId))
	utility.Assert(overseaPay != nil, fmt.Sprintf("找不到支付单 payId: %s", payIdStr))
	channelEntity := query.GetPaymentTypePayChannelById(r.Context(), overseaPay.ChannelId)
	utility.Assert(channelEntity != nil, fmt.Sprintf("payId: %s", payIdStr))
	g.Log().Infof(r.Context(), "DoRemoteChannelRedirect payId:%s notifyUrl:%s", payIdStr, overseaPay.ReturnUrl)
	if len(overseaPay.ReturnUrl) > 0 {
		r.Response.Writeln(fmt.Sprintf("<head>\n<meta http-equiv=\"refresh\" content=\"1;url=%s\">\n</head>", overseaPay.ReturnUrl))
	} else {
		//r.Response.Writeln(r.Get("channelId"))
		r.Response.Writeln(overseaPay)
	}
	return &ro.ChannelRedirectInternalResp{
		Status:  true,
		Message: "",
	}, nil
}
