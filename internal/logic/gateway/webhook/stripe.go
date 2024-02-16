package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	sub "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
	"net/http"
	"strings"
	"unibee-api/internal/consts"
	_gateway "unibee-api/internal/logic/gateway"
	"unibee-api/internal/logic/gateway/api"
	"unibee-api/internal/logic/gateway/api/log"
	"unibee-api/internal/logic/gateway/ro"
	handler2 "unibee-api/internal/logic/payment/handler"
	"unibee-api/internal/logic/subscription/handler"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type StripeWebhook struct {
}

// 测试数据
// 付款成功
// 4242 4242 4242 4242
// 付款需要验证
// 4000 0025 0000 3155
// 付款被拒绝
// 4000 0000 0000 9995
func (s StripeWebhook) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

// GatewayCheckAndSetupWebhook https://stripe.com/docs/billing/subscriptions/webhooks  https://stripe.com/docs/api/events/types
func (s StripeWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	utility.Assert(gateway != nil, "gateway is nil")
	stripe.Key = gateway.GatewaySecret
	params := &stripe.WebhookEndpointListParams{}
	params.Limit = stripe.Int64(10)
	result := webhookendpoint.List(params)
	if len(result.WebhookEndpointList().Data) > 1 {
		return gerror.New("webhook endpoints count > 1")
	}
	//过滤不可用
	if len(result.WebhookEndpointList().Data) == 0 {
		//创建
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				stripe.String("customer.subscription.deleted"),
				stripe.String("customer.subscription.updated"),
				stripe.String("customer.subscription.created"),
				stripe.String("customer.subscription.trial_will_end"),
				stripe.String("customer.subscription.paused"),
				stripe.String("customer.subscription.resumed"),
				stripe.String("invoice.upcoming"),
				stripe.String("invoice.created"),
				stripe.String("invoice.updated"), //todo mark 并发所有发票都会产生支付，并发所有订阅更新都会产生支付，可能从贷方余额支付（需确认）或者更新会产生退款从情况，所有 invoice.paid 可能必须要接
				stripe.String("invoice.paid"),
				stripe.String("invoice.voided"),
				stripe.String("invoice.will_be_due"),
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.created"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("payment_intent.canceled"),
				stripe.String("payment_intent.partially_funded"),
				stripe.String("payment_intent.payment_failed"),
				stripe.String("payment_intent.requires_action"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			URL: stripe.String(_gateway.GetPaymentWebhookEntranceUrl(int64(gateway.Id))),
		}
		result, err := webhookendpoint.New(params)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", params, result, err, "", nil, gateway)
		if err != nil {
			return nil
		}
		//更新 secret
		utility.Assert(len(result.Secret) > 0, "secret is nil")
		err = query.UpdateGatewayWebhookSecret(ctx, int64(gateway.Id), result.Secret)
		if err != nil {
			return err
		}
	} else {
		utility.Assert(len(result.WebhookEndpointList().Data) == 1, "internal webhook update, count is not 1")
		//检查并更新, todo mark 优化检查逻辑，如果 evert 一致不用发起更新
		webhook := result.WebhookEndpointList().Data[0]
		utility.Assert(strings.Compare(webhook.Status, "enabled") == 0, "webhook not status enabled")
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				//订阅相关 webhook
				stripe.String("customer.subscription.deleted"),
				stripe.String("customer.subscription.updated"),
				stripe.String("customer.subscription.created"),
				stripe.String("customer.subscription.trial_will_end"),
				stripe.String("customer.subscription.paused"),
				stripe.String("customer.subscription.resumed"),
				stripe.String("invoice.upcoming"),
				stripe.String("invoice.created"),
				stripe.String("invoice.updated"),
				stripe.String("invoice.paid"),
				stripe.String("invoice.voided"),
				stripe.String("invoice.will_be_due"),
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.created"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("payment_intent.canceled"),
				stripe.String("payment_intent.partially_funded"),
				stripe.String("payment_intent.payment_failed"),
				stripe.String("payment_intent.requires_action"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			URL: stripe.String(_gateway.GetPaymentWebhookEntranceUrl(int64(gateway.Id))),
		}
		result, err := webhookendpoint.Update(webhook.ID, params)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", params, result, err, webhook.ID, nil, gateway)
		if err != nil {
			return err
		}
		utility.Assert(strings.Compare(result.Status, "enabled") == 0, "webhook not status enabled after updated")
	}

	return nil
}

func (s StripeWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	endpointSecret := gateway.WebhookSecret
	signatureHeader := r.Header.Get("Stripe-Signature")
	var event stripe.Event
	var err error
	if !consts.GetConfigInstance().IsServerDev() {
		event, err = webhook.ConstructEvent(r.GetBody(), signatureHeader, endpointSecret)
		if err != nil {
			g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook signature verification failed. %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
			return
		}
	} else {
		if err := json.Unmarshal(r.GetBody(), &event); err != nil {
			g.Log().Errorf(r.Context(), "Failed to parse webhook body json: %s", err.Error())
			r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
			return
		}
	}

	data, _ := gjson.Marshal(event)
	g.Log().Info(r.Context(), "Receive_Webhook_Channel: ", gateway.GatewayName, " hook:", string(data))

	var responseBack = http.StatusOK
	var requestId = ""
	switch event.Type {
	case "customer.subscription.deleted", "customer.subscription.created", "customer.subscription.updated", "customer.subscription.trial_will_end":
		var subscription stripe.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error parsing webhook JSON: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Gateway:%s, Event %s for Sub %s\n", gateway.GatewayName, string(event.Type), subscription.ID)
			// Then define and call a func to handle the successful attachment of a GatewayDefaultPaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			requestId = subscription.ID
			err = s.processSubscriptionWebhook(r.Context(), string(event.Type), subscription, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandleSubscriptionWebhookEvent: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "invoice.upcoming", "invoice.created", "invoice.updated", "invoice.paid", "invoice.payment_failed", "invoice.payment_action_required", "invoice.voided", "invoice.will_be_due":
		var stripeInvoice stripe.Invoice
		err = json.Unmarshal(event.Data.Raw, &stripeInvoice)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error parsing webhook JSON: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Gateway:%s, Event %s for Invoice %s\n", gateway.GatewayName, string(event.Type), stripeInvoice.ID)
			requestId = stripeInvoice.ID
			// Then define and call a func to handle the successful attachment of a GatewayDefaultPaymentMethod.
			err = s.processInvoiceWebhook(r.Context(), string(event.Type), stripeInvoice, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandleInvoiceWebhookEvent: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "payment_intent.created", "payment_intent.succeeded", "payment_intent.canceled", "payment_intent.partially_funded", "payment_intent.payment_failed", "payment_intent.requires_action":
		var stripePayment stripe.PaymentIntent
		err = json.Unmarshal(event.Data.Raw, &stripePayment)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error parsing webhook JSON: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Gateway:%s, Event %s for Payment %s\n", gateway.GatewayName, string(event.Type), stripePayment.ID)
			// Then define and call a func to handle the successful attachment of a GatewayDefaultPaymentMethod.
			requestId = stripePayment.ID

			err = s.processPaymentWebhook(r.Context(), string(event.Type), stripePayment, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandlePaymentWebhookEvent: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "charge.refund.updated":
		var stripeRefund stripe.Refund
		err = json.Unmarshal(event.Data.Raw, &stripeRefund)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error parsing webhook JSON: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Gateway:%s, Event %s for Refund %s\n", gateway.GatewayName, string(event.Type), stripeRefund.ID)
			requestId = stripeRefund.ID
			// Then define and call a func to handle the successful attachment of a GatewayDefaultPaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)

			err = s.processRefundWebhook(r.Context(), string(event.Type), stripeRefund, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandlePaymentWebhookEvent: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "checkout.session.completed":
		var stripeCheckoutSession stripe.CheckoutSession
		err = json.Unmarshal(event.Data.Raw, &stripeCheckoutSession)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error parsing webhook JSON: %s\n", gateway.GatewayName, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Gateway:%s, Event %s for Refund %s\n", gateway.GatewayName, string(event.Type), stripeCheckoutSession.ID)
			// Then define and call a func to handle the successful attachment of a GatewayDefaultPaymentMethod.

			err = s.processCheckoutSessionWebhook(r.Context(), string(event.Type), stripeCheckoutSession, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandlePaymentWebhookEvent: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	default:
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Unhandled event type: %s\n", gateway.GatewayName, event.Type)
		r.Response.WriteHeader(http.StatusBadRequest)
		responseBack = http.StatusBadRequest
	}
	log.SaveChannelHttpLog("GatewayWebhook", event, responseBack, err, string(event.Type), requestId, gateway)
	r.Response.WriteHeader(responseBack)
}

func (s StripeWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *ro.GatewayRedirectInternalResp, err error) {
	params, err := r.GetJson()
	if err != nil {
		g.Log().Printf(r.Context(), "StripeNotifyController redirect params:%s err:%s", params, err.Error())
		r.Response.Writeln(err)
		return
	}
	payIdStr := r.Get("paymentId").String()
	SubIdStr := r.Get("subId").String()
	var response string
	var status = false
	var returnUrl = ""
	if len(payIdStr) > 0 {
		response = "not implement"
		//Payment Redirect
		if r.Get("success").Bool() {
			stripe.Key = gateway.GatewaySecret
			s.setUnibeeAppInfo()
			payment := query.GetPaymentByPaymentId(r.Context(), payIdStr)
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PAY_SUCCESS {
				returnUrl = payment.ReturnUrl
				response = "success"
				status = true
			} else {
				//需要去检索
				returnUrl = payment.ReturnUrl
				params := &stripe.CheckoutSessionParams{}
				result, err := session.Get(
					payment.GatewayPaymentIntentId,
					params,
				)
				if err != nil {
					response = "payment not match"
				}
				gatewayUser := query.GetGatewayUser(r.Context(), payment.UserId, int64(gateway.Id))
				if gatewayUser != nil && result != nil {
					//find
					if strings.Compare(result.Customer.ID, gatewayUser.GatewayUserId) != 0 {
						response = "user not match"
					} else if strings.Compare(string(result.Status), "complete") == 0 && result.PaymentIntent != nil && len(result.PaymentIntent.ID) > 0 {
						paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), int64(gateway.Id)).GatewayPaymentDetail(r.Context(), gateway, result.PaymentIntent.ID)
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							if paymentIntentDetail.Status == consts.PAY_SUCCESS {
								err := handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
									PaymentId:                        payment.PaymentId,
									GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
									GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
									TotalAmount:                      paymentIntentDetail.TotalAmount,
									PayStatusEnum:                    consts.PAY_SUCCESS,
									PaidTime:                         paymentIntentDetail.PayTime,
									PaymentAmount:                    paymentIntentDetail.PaymentAmount,
									CaptureAmount:                    0,
									Reason:                           paymentIntentDetail.Reason,
									ChannelDefaultPaymentMethod:      paymentIntentDetail.GatewayPaymentMethod,
									ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
								})
								if err != nil {
									response = fmt.Sprintf("%v", err)
								} else {
									response = "payment success"
									status = true
								}
							} else if paymentIntentDetail.Status == consts.PAY_FAILED {
								err := handler2.HandlePayFailure(r.Context(), &handler2.HandlePayReq{
									PaymentId:                        payment.PaymentId,
									GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
									GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
									TotalAmount:                      paymentIntentDetail.TotalAmount,
									PayStatusEnum:                    consts.PAY_FAILED,
									PaidTime:                         paymentIntentDetail.PayTime,
									PaymentAmount:                    paymentIntentDetail.PaymentAmount,
									CaptureAmount:                    0,
									Reason:                           paymentIntentDetail.Reason,
									ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
								})
								if err != nil {
									response = fmt.Sprintf("%v", err)
								}
							}
						}
					} else {
						response = "not complete"
					}
				} else {
					//not found
					response = "payment not found"
				}
			}
		} else {
			response = "user cancelled"
		}
	} else if len(SubIdStr) > 0 {
		//subscription redirect
		if r.Get("success").Bool() {
			stripe.Key = gateway.GatewaySecret
			s.setUnibeeAppInfo()
			unibSub := query.GetSubscriptionBySubscriptionId(r.Context(), SubIdStr)
			if unibSub == nil {
				response = "subId invalid"
			} else if len(unibSub.GatewaySubscriptionId) > 0 && unibSub.Status == consts.SubStatusActive {
				returnUrl = unibSub.ReturnUrl
				response = "active"
				status = true
			} else {
				//search
				returnUrl = unibSub.ReturnUrl
				params := &stripe.SubscriptionSearchParams{
					SearchParams: stripe.SearchParams{
						Query: "metadata['SubId']:'" + SubIdStr + "'",
					},
				}
				result := sub.Search(params)
				gatewayUser := query.GetGatewayUser(r.Context(), unibSub.UserId, int64(gateway.Id))
				if gatewayUser != nil && result.SubscriptionSearchResult().Data != nil && len(result.SubscriptionSearchResult().Data) == 1 {
					//找到
					if strings.Compare(result.SubscriptionSearchResult().Data[0].Customer.ID, gatewayUser.GatewayUserId) != 0 {
						response = "customId not match"
					} else {
						detail := parseStripeSubscription(result.SubscriptionSearchResult().Data[0])
						err := handler.UpdateSubWithGatewayDetailBack(r.Context(), unibSub, detail)
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else if detail.Status == consts.SubStatusActive {
							response = "subscription active"
							status = true
						} else {
							response = "not complete"
						}
					}
				} else {
					//not find
					response = "subscription not paid"
				}
			}
		} else {
			response = "user cancelled"
		}
	}
	log.SaveChannelHttpLog("GatewayRedirect", params, response, err, "", nil, gateway)
	return &ro.GatewayRedirectInternalResp{
		Status:    status,
		Message:   response,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}

func (s StripeWebhook) processRefundWebhook(ctx context.Context, eventType string, refund stripe.Refund, gateway *entity.MerchantGateway) error {
	refundDetail, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayRefundDetail(ctx, gateway, refund.ID)
	if err != nil {
		return err
	}
	//details.Id = int64(gateway.Id)
	//utility.Assert(len(details.GatewayUserId) > 0, "invalid gatewayUserId")
	//if payment.Invoice != nil {
	//	//可能来自 SubPendingUpdate 流程，需要补充 Invoice 信息获取 GatewaySubscriptionUpdateId
	//	invoiceDetails, err := s.GatewayInvoiceDetails(ctx, gateway, payment.Invoice.ID)
	//	if err != nil {
	//		return err
	//	}
	//	details.GatewayInvoiceDetail = invoiceDetails
	//	details.GatewayInvoiceId = payment.Invoice.ID
	//	details.GatewaySubscriptionUpdateId = invoiceDetails.GatewayInvoiceId
	//	oneSub := query.GetSubscriptionByGatewaySubscriptionId(ctx, invoiceDetails.GatewaySubscriptionId)
	//	if oneSub != nil {
	//		plan := query.GetPlanById(ctx, oneSub.PlanId)
	//		gatewayPlan := query.GetGatewayPlan(ctx, oneSub.PlanId, oneSub.Id)
	//		subDetails, err := s.GatewaySubscriptionDetails(ctx, plan, gatewayPlan, oneSub)
	//		if err != nil {
	//			return err
	//		}
	//		details.GatewaySubscriptionDetail = subDetails
	//	}
	//}
	//details.UniqueId = details.GatewayPaymentIntentId
	err = handler2.HandleRefundWebhookEvent(ctx, refundDetail)
	if err != nil {
		return err
	}

	return nil
}

func (s StripeWebhook) processPaymentWebhook(ctx context.Context, eventType string, stripePayment stripe.PaymentIntent, gateway *entity.MerchantGateway) error {
	if paymentId, ok := stripePayment.Metadata["PaymentId"]; ok {
		payment := query.GetPaymentByPaymentId(ctx, paymentId)
		if payment != nil {
			paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayPaymentDetail(ctx, gateway, stripePayment.ID)
			if err != nil {
				return err
			}
			if len(paymentIntentDetail.PaymentData) == 0 && stripePayment.NextAction != nil {
				paymentIntentDetail.PaymentData = utility.MarshalToJsonString(stripePayment.NextAction)
			}
			if paymentIntentDetail.Status == consts.PAY_SUCCESS {
				err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_SUCCESS,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.Reason,
					ChannelDefaultPaymentMethod:      paymentIntentDetail.GatewayPaymentMethod,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PAY_FAILED {
				err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_FAILED,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.CancelReason,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PAY_CANCEL {
				err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           paymentIntentDetail.GatewayPaymentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_CANCEL,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.CancelReason,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return err
				}
			} else if paymentIntentDetail.AuthorizeStatus == consts.WAITING_AUTHORIZED {
				err := handler2.HandlePayNeedAuthorized(ctx, payment, paymentIntentDetail.AuthorizeReason, paymentIntentDetail.PaymentData)
				if err != nil {
					return err
				}
			}
		} else {
			return gerror.New("Payment Not Found")
		}
	} else {
		return gerror.New("No PaymentId Metadata")
	}
	return nil
}

func (s StripeWebhook) processInvoiceWebhook(ctx context.Context, eventType string, invoice stripe.Invoice, gateway *entity.MerchantGateway) error {
	utility.Assert(len(invoice.ID) > 0, "processInvoiceWebhook gatewayInvoiceId Invalid")
	invoiceDetails, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayInvoiceDetails(ctx, gateway, invoice.ID)
	if err != nil {
		return err
	}
	// require_action status not deal here, use subscription_update webhook
	//if strings.Compare("invoice.payment_action_required", eventType) == 0 {
	//	return gerror.New("require_action status not deal processInvoiceWebhook, use processSubscriptionWebhook webhook")
	//}

	var status = consts.TO_BE_PAID
	var authorizeStatus = consts.AUTHORIZED
	var authorizeReason = ""
	var cancelReason = ""
	var paymentData = ""
	if invoiceDetails.Status == consts.InvoiceStatusPaid {
		status = consts.PAY_SUCCESS
		authorizeStatus = consts.CAPTURE_REQUEST
	} else if invoiceDetails.Status == consts.InvoiceStatusFailed {
		status = consts.PAY_FAILED
	} else if invoiceDetails.Status == consts.InvoiceStatusCancelled {
		status = consts.PAY_CANCEL
	} else if strings.Compare("invoice.payment_action_required", eventType) == 0 {
		authorizeStatus = consts.WAITING_AUTHORIZED
	}

	if len(invoiceDetails.GatewayPaymentId) > 0 {
		paymentIntentDetail, _ := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayPaymentDetail(ctx, gateway, invoiceDetails.GatewayPaymentId)
		if paymentIntentDetail != nil {
			authorizeReason = paymentIntentDetail.AuthorizeReason
			cancelReason = paymentIntentDetail.CancelReason
			paymentData = paymentIntentDetail.PaymentData
		}
	}

	var gatewaySubscriptionDetail *ro.GatewayDetailSubscriptionInternalResp
	if len(invoiceDetails.GatewaySubscriptionId) > 0 {
		unibeeSub := query.GetSubscriptionByGatewaySubscriptionId(ctx, invoiceDetails.GatewaySubscriptionId)
		var subNeedUpdate = false
		if unibeeSub == nil && len(invoiceDetails.SubscriptionId) > 0 {
			unibeeSub = query.GetSubscriptionBySubscriptionId(ctx, invoiceDetails.SubscriptionId)
			unibeeSub.GatewaySubscriptionId = invoiceDetails.GatewaySubscriptionId
			subNeedUpdate = true
		}
		if unibeeSub != nil {
			plan := query.GetPlanById(ctx, unibeeSub.PlanId)
			gatewayPlan := query.GetGatewayPlan(ctx, unibeeSub.PlanId, unibeeSub.GatewayId)
			gatewaySubscriptionDetail, err = api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewaySubscriptionDetails(ctx, plan, gatewayPlan, unibeeSub)
			if subNeedUpdate {
				err = handler.HandleSubscriptionWebhookEvent(ctx, unibeeSub, eventType, gatewaySubscriptionDetail)
				if err != nil {
					return err
				}
			}
		}
	}

	err = handler2.HandlePaymentWebhookEvent(ctx, &ro.GatewayPaymentRo{
		MerchantId:                  gateway.MerchantId,
		Status:                      status,
		AuthorizeStatus:             authorizeStatus,
		AuthorizeReason:             authorizeReason,
		Currency:                    invoiceDetails.Currency,
		TotalAmount:                 invoiceDetails.TotalAmount,
		PaymentAmount:               invoiceDetails.PaymentAmount,
		BalanceAmount:               invoiceDetails.BalanceAmount,
		BalanceStart:                invoiceDetails.BalanceStart,
		BalanceEnd:                  invoiceDetails.BalanceEnd,
		Reason:                      invoiceDetails.Reason,
		CancelReason:                cancelReason,
		PaymentData:                 paymentData,
		UniqueId:                    invoiceDetails.GatewayInvoiceId,
		PayTime:                     gtime.NewFromTimeStamp(invoiceDetails.PaymentTime),
		CreateTime:                  gtime.NewFromTimeStamp(invoiceDetails.CreateTime),
		CancelTime:                  gtime.NewFromTimeStamp(invoiceDetails.CancelTime),
		GatewayId:                   int64(gateway.Id),
		GatewayUserId:               invoiceDetails.GatewayUserId,
		GatewayPaymentId:            invoiceDetails.GatewayPaymentId,
		GatewayPaymentMethod:        invoiceDetails.GatewayDefaultPaymentMethod,
		GatewayInvoiceId:            invoiceDetails.GatewayInvoiceId,
		GatewaySubscriptionId:       invoiceDetails.GatewaySubscriptionId,
		GatewaySubscriptionUpdateId: invoiceDetails.GatewayInvoiceId,
		GatewayInvoiceDetail:        invoiceDetails,
		GatewaySubscriptionDetail:   gatewaySubscriptionDetail,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s StripeWebhook) processSubscriptionWebhook(ctx context.Context, eventType string, subscription stripe.Subscription, gateway *entity.MerchantGateway) error {
	unibeeSub := query.GetSubscriptionByGatewaySubscriptionId(ctx, subscription.ID)
	if unibeeSub == nil {
		if unibSubId, ok := subscription.Metadata["SubId"]; ok {
			unibeeSub = query.GetSubscriptionBySubscriptionId(ctx, unibSubId)
			unibeeSub.GatewaySubscriptionId = subscription.ID
		}
	}
	if unibeeSub != nil {
		plan := query.GetPlanById(ctx, unibeeSub.PlanId)
		gatewayPlan := query.GetGatewayPlan(ctx, unibeeSub.PlanId, unibeeSub.GatewayId)
		details, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewaySubscriptionDetails(ctx, plan, gatewayPlan, unibeeSub)
		if err != nil {
			return err
		}

		err = handler.HandleSubscriptionWebhookEvent(ctx, unibeeSub, eventType, details)
		if err != nil {
			return err
		}
		if details.Status == consts.SubStatusIncomplete && len(details.GatewayLatestInvoiceId) > 0 {
			//处理支付需要授权事件
			invoiceDetails, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayInvoiceDetails(ctx, gateway, details.GatewayLatestInvoiceId)
			if err != nil {
				return err
			}
			if invoiceDetails.Status != consts.InvoiceStatusPaid {
				//有支付授权 todo mark
				var gatewaySubscriptionDetail *ro.GatewayDetailSubscriptionInternalResp
				if len(invoiceDetails.GatewaySubscriptionId) > 0 {
					oneSub := query.GetSubscriptionByGatewaySubscriptionId(ctx, invoiceDetails.GatewaySubscriptionId)
					if oneSub != nil {
						plan := query.GetPlanById(ctx, oneSub.PlanId)
						gatewayPlan := query.GetGatewayPlan(ctx, oneSub.PlanId, oneSub.GatewayId)
						gatewaySubscriptionDetail, err = api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewaySubscriptionDetails(ctx, plan, gatewayPlan, oneSub)
					}
				}

				err = handler2.HandlePaymentWebhookEvent(ctx, &ro.GatewayPaymentRo{
					MerchantId:                  gateway.MerchantId,
					Status:                      consts.TO_BE_PAID,
					AuthorizeStatus:             consts.WAITING_AUTHORIZED,
					Currency:                    invoiceDetails.Currency,
					TotalAmount:                 invoiceDetails.TotalAmount,
					PaymentAmount:               invoiceDetails.PaymentAmount,
					BalanceAmount:               invoiceDetails.BalanceAmount,
					BalanceStart:                invoiceDetails.BalanceStart,
					BalanceEnd:                  invoiceDetails.BalanceEnd,
					Reason:                      invoiceDetails.Reason,
					UniqueId:                    invoiceDetails.GatewayInvoiceId,
					PayTime:                     gtime.NewFromTimeStamp(invoiceDetails.PaymentTime),
					CreateTime:                  gtime.NewFromTimeStamp(invoiceDetails.CreateTime),
					CancelTime:                  gtime.NewFromTimeStamp(invoiceDetails.CancelTime),
					GatewayId:                   int64(gateway.Id),
					GatewayUserId:               invoiceDetails.GatewayUserId,
					GatewayPaymentId:            invoiceDetails.GatewayPaymentId,
					GatewayPaymentMethod:        invoiceDetails.GatewayDefaultPaymentMethod,
					GatewayInvoiceId:            invoiceDetails.GatewayInvoiceId,
					GatewaySubscriptionId:       invoiceDetails.GatewaySubscriptionId,
					GatewaySubscriptionUpdateId: invoiceDetails.GatewayInvoiceId,
					GatewayInvoiceDetail:        invoiceDetails,
					GatewaySubscriptionDetail:   gatewaySubscriptionDetail,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return gerror.New("subscription not found on gatewaySubId:" + subscription.ID)
	}
}

func (s StripeWebhook) processCheckoutSessionWebhook(ctx context.Context, event string, checkoutSession stripe.CheckoutSession, gateway *entity.MerchantGateway) error {
	if paymentId, ok := checkoutSession.Metadata["PaymentId"]; ok {
		payment := query.GetPaymentByPaymentId(ctx, paymentId)
		if checkoutSession.PaymentIntent != nil {
			paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayPaymentDetail(ctx, gateway, checkoutSession.PaymentIntent.ID)
			if err != nil {
				return gerror.New(fmt.Sprintf("%s", err.Error()))
			}
			if paymentIntentDetail.Status == consts.PAY_SUCCESS {
				err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_SUCCESS,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.Reason,
					ChannelDefaultPaymentMethod:      paymentIntentDetail.GatewayPaymentMethod,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PAY_FAILED {
				err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           payment.GatewayPaymentIntentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_FAILED,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.Reason,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PAY_CANCEL {
				err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					GatewayPaymentIntentId:           paymentIntentDetail.GatewayPaymentId,
					GatewayPaymentId:                 paymentIntentDetail.GatewayPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_CANCEL,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.CancelReason,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.GatewayInvoiceDetail,
				})
				if err != nil {
					return err
				}
			} else if paymentIntentDetail.AuthorizeStatus == consts.WAITING_AUTHORIZED {
				err := handler2.HandlePayNeedAuthorized(ctx, payment, paymentIntentDetail.AuthorizeReason, paymentIntentDetail.PaymentData)
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			return gerror.New("no PaymentIntent")
		}
	} else {
		return gerror.New("No PaymentId Metadata")
	}
}

func parseStripeSubscription(subscription *stripe.Subscription) *ro.GatewayDetailSubscriptionInternalResp {
	//https://stripe.com/docs/billing/subscriptions/overview
	/**
	trialing	订阅目前处于试用期，可以安全地为您的客户配置您的产品。订阅会自动转换到active首次付款时。
	active	订阅信誉良好，最近一次付款成功。为您的客户配置您的产品是安全的。
	incomplete	需要在23小时内成功付款才能激活订阅。或者付款需要采取行动，例如客户身份验证。incomplete如果有待付款并且 PaymentIntent 状态为 ，则订阅也可以为processing。
	incomplete_expired	订阅的首次付款失败，并且在创建订阅后 23 小时内未成功付款。这些订阅不会向客户收取费用。存在此状态是为了让您可以跟踪未能激活订阅的客户。
	past_due	最新最终发票的付款失败或未尝试。订阅将继续创建发票。您的订阅设置决定了订阅的下一个状态。如果在尝试所有智能重试后发票仍未支付，您可以将订阅配置为移至canceled、unpaid，或保留为past_due。要将订阅转移到active，请在到期日之前支付最新的发票。
	canceled	订阅已被取消。取消期间，将禁用所有未付发票的自动收取 ( auto_advance=false)。这是无法更新的最终状态。
	unpaid	最新的发票尚未支付，但订阅仍然有效。最新发票仍处于打开状态，并且继续生成发票，但不会尝试付款。您应该在订阅时撤销对产品的访问权限，unpaid因为已尝试付款并在订阅时重试past_due。要将订阅转移到active，请在到期日之前支付最新的发票。
	paused	订阅已结束试用期，没有默认付款方式，并且trial_settings.end_behavior.missing_payment_method设置为pause。将不再为订阅创建发票。为客户附加默认付款方式后，您可以恢复订阅。
	*/
	var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	if strings.Compare(string(subscription.Status), "trialing") == 0 ||
		strings.Compare(string(subscription.Status), "active") == 0 {
		status = consts.SubStatusActive
	} else if strings.Compare(string(subscription.Status), "unpaid") == 0 {
		status = consts.SubStatusCreate
	} else if strings.Compare(string(subscription.Status), "incomplete_expired") == 0 {
		status = consts.SubStatusExpired
	} else if strings.Compare(string(subscription.Status), "incomplete") == 0 ||
		strings.Compare(string(subscription.Status), "pass_due") == 0 {
		status = consts.SubStatusIncomplete
	} else if strings.Compare(string(subscription.Status), "paused") == 0 {
		status = consts.SubStatusSuspended
	} else if strings.Compare(string(subscription.Status), "canceled") == 0 {
		status = consts.SubStatusCancelled
	}
	var latestChannelPaymentId = ""
	if subscription.LatestInvoice != nil && subscription.LatestInvoice.PaymentIntent != nil {
		latestChannelPaymentId = subscription.LatestInvoice.PaymentIntent.ID
	}

	return &ro.GatewayDetailSubscriptionInternalResp{
		Status:                 status,
		GatewaySubscriptionId:  subscription.ID,
		GatewayStatus:          string(subscription.Status),
		Data:                   utility.FormatToJsonString(subscription),
		GatewayItemData:        utility.MarshalToJsonString(subscription.Items.Data),
		GatewayLatestInvoiceId: subscription.LatestInvoice.ID,
		GatewayLatestPaymentId: latestChannelPaymentId,
		CancelAtPeriodEnd:      subscription.CancelAtPeriodEnd,
		CurrentPeriodStart:     subscription.CurrentPeriodStart,
		CurrentPeriodEnd:       subscription.CurrentPeriodEnd,
		BillingCycleAnchor:     subscription.BillingCycleAnchor,
		TrialEnd:               subscription.TrialEnd,
	}
}
