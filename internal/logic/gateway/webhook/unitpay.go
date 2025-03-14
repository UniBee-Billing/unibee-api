package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strings"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/user/sub_update"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type JsonRes struct {
	Data      interface{} `json:"data"`
	Redirect  string      `json:"redirect"`
	RequestId string      `json:"requestId"`
}

type UnitpayWebhook struct {
}

func (c UnitpayWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (c UnitpayWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	_ = query.UpdateGatewayWebhookSecret(ctx, gateway.Id, gateway.GatewaySecret)
	return nil
}

func (c UnitpayWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	var err error
	method := fmt.Sprintf("%s", r.GetQuery("method"))
	var params = make(map[string]interface{})
	if r.GetQuery("params") != nil && r.GetQuery("params").IsMap() {
		params = r.GetQuery("params").Map()
	}
	if (method != "check" && method != "pay") || len(params) == 0 {
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Invalid request\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	params["method"] = method
	jsonData := gjson.New(params)
	g.Log().Infof(r.Context(), "Unitpay_GatewayWebhook url: %s method:%s params:%s body:%s", r.GetUrl(), method, params, jsonData)
	g.Log().Info(r.Context(), "Receive_Webhook_Channel:", gateway.GatewayName, " hook:", jsonData.String())
	var responseBack = http.StatusOK
	var response = make(map[string]interface{})
	if method == "check" && jsonData.Contains("account") &&
		(strings.Contains(jsonData.Get("account").String(), "test") ||
			query.GetPaymentByPaymentId(r.Context(), jsonData.Get("account").String()) != nil ||
			(jsonData.Contains("test") && strings.Contains(jsonData.Get("test").String(), "1"))) {
		response["result"] = map[string]interface{}{
			"message": "success",
		}
	} else if method == "pay" && jsonData.Contains("unitpayId") && jsonData.Contains("account") {
		err = ProcessPaymentWebhook(r.Context(), jsonData.Get("account").String(), jsonData.Get("unitpayId").String(), gateway)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error ProcessPaymentWebhook: %s\n", gateway.GatewayName, err.Error())
		} else {
			{
				if jsonData.Contains("subscriptionId") && len(jsonData.Get("subscriptionId").String()) > 0 {
					go func() {
						backgroundCtx := context.Background()
						defer func() {
							if exception := recover(); exception != nil {
								fmt.Printf("Unitpay Save GatewayPaymentMethod panic error:%s\n", exception)
								return
							}
						}()
						paymentMethod := jsonData.Get("subscriptionId").String()
						//update gateway payment method
						payment := query.GetPaymentByPaymentId(backgroundCtx, jsonData.Get("account").String())
						if payment != nil {
							sub_update.UpdateUserDefaultGatewayPaymentMethod(backgroundCtx, payment.UserId, payment.GatewayId, paymentMethod, "")
						}
					}()
				}
			}
		}
		response["result"] = map[string]interface{}{
			"message": "success",
		}
	} else {
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Payment not found paymentId:%s\n", gateway.GatewayName, jsonData.Get("account").String())
		responseBack = http.StatusBadRequest
		response["error"] = map[string]interface{}{
			"message": "invalid payment",
		}
	}
	log.SaveChannelHttpLog("GatewayWebhook", jsonData, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	r.Response.WriteHeader(responseBack)
	r.Response.WriteJson(response)

	return
}

// redirect https://url/?paymentId=443597283&account=test
func (c UnitpayWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	payIdStr := r.Get("account").String()
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
						err = handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
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
						err = handler2.HandlePayFailure(r.Context(), &handler2.HandlePayReq{
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
