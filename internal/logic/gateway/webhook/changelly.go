package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type ChangellyWebhook struct {
}

func (c ChangellyWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (c ChangellyWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (c ChangellyWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Webhook Get PortalJson failed. %v\n", gateway.GatewayName, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	//client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	//_, err = client.GetAccessToken(context.Background())
	//if err != nil {
	//	r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
	//	return
	//}
	//signature, err := client.VerifyWebhookSignature(r.Context(), r.Request, jsonData.Get("id").String())
	//if err != nil {
	//	g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Webhook signature verification success\n", gateway.GatewayName)
	//	r.Response.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//if strings.Compare(signature.VerificationStatus, "SUCCESS") == 0 {
	g.Log().Info(r.Context(), "Receive_Webhook_Channel:", gateway.GatewayName, " hook:", jsonData.String())
	var responseBack = http.StatusOK
	if jsonData.Contains("payment_id") {
		err = ProcessPaymentWebhook(r.Context(), jsonData.Get("order_id").String(), jsonData.Get("payment_id").String(), gateway)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error ProcessPaymentWebhook: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		}
	} else {
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Unhandled paymentId\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
	}
	log.SaveChannelHttpLog("GatewayWebhook", jsonData, responseBack, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	r.Response.WriteHeader(responseBack)
	return
	//} else {
	//	g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Webhook signature verification failed.\n", gateway.GatewayName)
	//	r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
	//	return
	//}
}

func (c ChangellyWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//params, err := r.GetJson()
	//if err != nil {
	//	g.Log().Printf(r.Context(), "ChangellyNotify redirect params:%s err:%s", params, err.Error())
	//	r.Response.Writeln(err)
	//	return
	//}
	payIdStr := r.Get("paymentId").String()
	var response string
	var status = false
	var returnUrl = ""
	var isSuccess = false
	if len(payIdStr) > 0 {
		response = ""
		//Payment Redirect
		payment := query.GetPaymentByPaymentId(r.Context(), payIdStr)
		if payment != nil {
			success := r.Get("success")
			if success != nil {
				if success.String() == "true" {
					isSuccess = true
				}
				returnUrl = GetPaymentRedirectUrl(r.Context(), payment, success.String())
			} else {
				returnUrl = GetPaymentRedirectUrl(r.Context(), payment, "")
			}
		}
		if r.Get("success").Bool() {
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				response = "success"
				status = true
			} else {
				//find
				paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayPaymentDetail(r.Context(), gateway, payment.GatewayPaymentId, payment)
				if err != nil {
					response = fmt.Sprintf("%v", err)
				} else {
					if paymentIntentDetail.Status == consts.PaymentSuccess {
						err := handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							GatewayUserId:          paymentIntentDetail.GatewayUserId,
							TotalAmount:            paymentIntentDetail.TotalAmount,
							PayStatusEnum:          consts.PaymentSuccess,
							PaidTime:               paymentIntentDetail.PaidTime,
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
	log.SaveChannelHttpLog("GatewayRedirect", r.URL, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	return &gateway_bean.GatewayRedirectResp{
		Status:    status,
		Message:   response,
		Success:   isSuccess,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}
