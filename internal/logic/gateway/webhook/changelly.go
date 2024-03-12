package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type ChangellyWebhook struct {
}

func (c ChangellyWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (c ChangellyWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	panic("implement me")
}

func (c ChangellyWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	params, err := r.GetJson()
	if err != nil {
		g.Log().Printf(r.Context(), "ChangellyNotify redirect params:%s err:%s", params, err.Error())
		r.Response.Writeln(err)
		return
	}
	payIdStr := r.Get("paymentId").String()
	var response string
	var status = false
	var returnUrl = ""
	if len(payIdStr) > 0 {
		response = ""
		//Payment Redirect
		if r.Get("success").Bool() {
			payment := query.GetPaymentByPaymentId(r.Context(), payIdStr)
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				returnUrl = payment.ReturnUrl
				response = "success"
				status = true
			} else {
				returnUrl = payment.ReturnUrl
				//find
				paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayPaymentDetail(r.Context(), gateway, payment.GatewayPaymentId)
				if err != nil {
					response = fmt.Sprintf("%v", err)
				} else {
					if paymentIntentDetail.Status == consts.PaymentSuccess {
						err := handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							TotalAmount:            paymentIntentDetail.TotalAmount,
							PayStatusEnum:          consts.PaymentSuccess,
							PaidTime:               paymentIntentDetail.PayTime,
							PaymentAmount:          paymentIntentDetail.PaymentAmount,
							Reason:                 paymentIntentDetail.Reason,
							GatewayPaymentMethod:   paymentIntentDetail.GatewayPaymentMethod,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							response = "payment success"
							status = true
						}
					} else if paymentIntentDetail.Status == consts.PaymentFailed {
						err := handler2.HandlePayFailure(r.Context(), &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							PayStatusEnum:          consts.PaymentFailed,
							Reason:                 paymentIntentDetail.Reason,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						}
					}
				}
			}
		} else {
			response = "user cancelled"
		}
	}
	log.SaveChannelHttpLog("GatewayRedirect", params, response, err, "", nil, gateway)
	return &gateway_bean.GatewayRedirectResp{
		Status:    status,
		Message:   response,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}
