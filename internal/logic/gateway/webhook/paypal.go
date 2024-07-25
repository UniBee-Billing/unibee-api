package webhook

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_gateway "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/api/paypal"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type PaypalWebhook struct {
}

func init() {
	//注册 gateway_webhook_entry
}

// todo mark 确认改造成单例是否可行，不用每次都去获取 accessToken
func NewClient(clientID string, secret string, APIBase string) (*paypal.Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &paypal.Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

func (p PaypalWebhook) GetPaypalHost() string {
	var apiHost = "https://api-m.paypal.com"
	if !config.GetConfigInstance().IsProd() {
		apiHost = "https://api-m.sandbox.paypal.com"
	}
	return apiHost
}

// GatewayCheckAndSetupWebhook https://developer.paypal.com/docs/subscriptions/webhooks/
func (p PaypalWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	utility.Assert(gateway != nil, "gateway is nil")
	client, err := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	if err != nil {
		return err
	}
	_, err = client.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	result, err := client.ListWebhooks(ctx, paypal.AncorTypeApplication)
	if err != nil {
		return err
	}
	webhookUrl := _gateway.GetPaymentWebhookEntranceUrl(gateway.Id)
	var targetEventTypes = []paypal.WebhookEventType{
		{Name: "PAYMENT.SALE.COMPLETED"},
		{Name: "PAYMENT.SALE.REFUNDED"},
		{Name: "PAYMENT.SALE.REVERSED"},
		{Name: "CHECKOUT.ORDER.COMPLETED"},
		{Name: "CHECKOUT.ORDER.APPROVED"},
		{Name: "CHECKOUT.PAYMENT-APPROVAL.REVERSED"},
		{Name: "VAULT.PAYMENT-TOKEN.CREATED"},
		{Name: "VAULT.PAYMENT-TOKEN.DELETED"},
		{Name: "VAULT.PAYMENT-TOKEN.DELETION-INITIATED"},
	}
	var one *paypal.Webhook
	for _, endpoint := range result.Webhooks {
		if strings.Compare(endpoint.URL, webhookUrl) == 0 {
			if len(endpoint.EventTypes) != len(targetEventTypes) {
				err = client.DeleteWebhook(ctx, endpoint.ID)
				if err != nil {
					g.Log().Errorf(ctx, "Delete Paypal Webhook Endpoint error:%s", err.Error())
					return err
				}
			} else {
				one = &endpoint
				break
			}
		}
	}
	if one == nil {
		param := &paypal.CreateWebhookRequest{
			URL:        webhookUrl,
			EventTypes: targetEventTypes,
		}
		response, err := client.CreateWebhook(ctx, param)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", param, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
		if err != nil {
			return err
		}
		utility.Assert(len(response.ID) > 0, "secret is nil")
		err = query.UpdateGatewayWebhookSecret(ctx, gateway.Id, response.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p PaypalWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//params, err := r.GetJson()
	//if err != nil {
	//	g.Log().Printf(r.Context(), "Paypal redirect params:%s err:%s", params, err.Error())
	//	r.Response.Writeln(err)
	//	return
	//}
	payIdStr := r.Get("paymentId").String()
	var response string
	var status = false
	var returnUrl = ""
	var isSuccess = true
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
			if r.Get("PayerID") != nil {
				customerId := r.Get("PayerID").String()
				var gatewayPaymentMethodId = ""
				if r.Get("token") != nil {
					gatewayPaymentMethodId = r.Get("token").String()
				}
				_, _ = query.CreateOrUpdateGatewayUser(r.Context(), payment.UserId, gateway.Id, customerId, gatewayPaymentMethodId)
			}
		}
		if r.Get("success").Bool() {
			if payment == nil || len(payment.GatewayPaymentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				response = "success"
				status = true
			} else {
				//find
				if payment.AuthorizeStatus == consts.Authorized {
					//_, err = api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayCapture(r.Context(), payment)
					err = service.PaymentGatewayCapture(r.Context(), payment)
					if err != nil {
						g.Log().Errorf(r.Context(), "GatewayRedirect paypal GatewayCapture error:%s", err.Error())
					}
				}
				paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayPaymentDetail(r.Context(), gateway, payment.GatewayPaymentId, payment)
				if err != nil {
					response = fmt.Sprintf("GatewayPaymentDetail %v", err)
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
							response = fmt.Sprintf("HandlePaySuccess %v", err)
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
							response = fmt.Sprintf("HandlePayFailure %v", err)
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

func (p PaypalWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook Get PortalJson failed. %v\n", gateway.GatewayName, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = client.GetAccessToken(r.Context())
	if err != nil {
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	signature, err := client.VerifyWebhookSignature(r.Context(), r.Request, gateway.WebhookSecret)
	if err != nil {
		log.SaveChannelHttpLog("GatewayWebhook", jsonData, "VerifyError-400", err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook signature verification err:%s\n", gateway.GatewayName, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	var eventType = ""
	if strings.Compare(signature.VerificationStatus, "SUCCESS") == 0 {
		g.Log().Info(r.Context(), "Receive_Webhook_Channel:", gateway.GatewayName, " hook:", jsonData.String())
		eventType = jsonData.Get("event_type").String()
		var responseBack = http.StatusOK
		switch eventType {
		case "CHECKOUT.ORDER.COMPLETED", "CHECKOUT.ORDER.APPROVED", "CHECKOUT.PAYMENT-APPROVAL.REVERSED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s-%d, Error parsing webhook resource is nil\n", gateway.GatewayName, gateway.Id)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Gateway:%s-%d, Subscription updated for %d.", gateway.GatewayName, gateway.Id, resource.Get("id").String())
				// Then define and call a func to handle the successful attachment of a PaymentMethod.
				gatewayPaymentId := resource.Get("id").String()
				payment := query.GetPaymentByGatewayPaymentId(r.Context(), gatewayPaymentId)
				if payment.MerchantId != gateway.MerchantId {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s-%d, Payment Merchant Not Match error:%s\n", gateway.GatewayName, gateway.Id, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				} else if payment != nil {
					if eventType == "CHECKOUT.ORDER.APPROVED" {
						_, err = api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayCapture(r.Context(), payment)
						if err != nil {
							g.Log().Errorf(r.Context(), "Webhook Gateway paypal GatewayCapture error:%s", err.Error())
						}
					}
					err = ProcessPaymentWebhook(r.Context(), payment.PaymentId, gatewayPaymentId, gateway)
					if err != nil {
						g.Log().Errorf(r.Context(), "Webhook Channel:%s-%d, ProcessPaymentWebhook error:%s\n", gateway.GatewayName, gateway.Id, err.Error())
						r.Response.WriteHeader(http.StatusBadRequest)
						responseBack = http.StatusBadRequest
					}
				} else {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s-%d, Error Payment not match: %v\n", gateway.GatewayName, gateway.Id, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		case "VAULT.PAYMENT-TOKEN.CREATED", "VAULT.PAYMENT-TOKEN.DELETED", "VAULT.PAYMENT-TOKEN.DELETION-INITIATED":

		default:
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Unhandled event type: %s\n", gateway.GatewayName, eventType)
		}
		r.Response.WriteHeader(http.StatusOK)
		log.SaveChannelHttpLog("GatewayWebhook", jsonData, responseBack, err, fmt.Sprintf("%s-%s-%d", eventType, gateway.GatewayName, gateway.Id), nil, gateway)
		return
	} else {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook signature verification failed.\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
}
