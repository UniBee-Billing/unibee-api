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
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/channel"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/out/log"
	"go-oversea-pay/internal/logic/channel/ro"
	handler2 "go-oversea-pay/internal/logic/payment/handler"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"net/http"
	"strings"
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

// DoRemoteChannelCheckAndSetupWebhook https://stripe.com/docs/billing/subscriptions/webhooks
func (s StripeWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantChannelConfig) (err error) {
	utility.Assert(payChannel != nil, "payChannel is nil")
	stripe.Key = payChannel.ChannelSecret
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
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.created"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			URL: stripe.String(channel.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
		}
		result, err := webhookendpoint.New(params)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", params, result, err, "", nil, payChannel)
		if err != nil {
			return nil
		}
		//更新 secret
		utility.Assert(len(result.Secret) > 0, "secret is nil")
		err = query.UpdatePayChannelWebhookSecret(ctx, int64(payChannel.Id), result.Secret)
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
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.created"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			URL: stripe.String(channel.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
		}
		result, err := webhookendpoint.Update(webhook.ID, params)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", params, result, err, webhook.ID, nil, payChannel)
		if err != nil {
			return err
		}
		utility.Assert(strings.Compare(result.Status, "enabled") == 0, "webhook not status enabled after updated")
	}

	return nil
}

func (s StripeWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) {
	endpointSecret := payChannel.WebhookSecret
	signatureHeader := r.Header.Get("Stripe-Signature")
	var event stripe.Event
	var err error
	if !consts.GetConfigInstance().IsServerDev() {
		event, err = webhook.ConstructEvent(r.GetBody(), signatureHeader, endpointSecret)
		if err != nil {
			g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %s\n", payChannel.Channel, err.Error())
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
	g.Log().Info(r.Context(), "Receive_Webhook_Channel: ", payChannel.Channel, " hook:", string(data))

	var responseBack = http.StatusOK
	switch event.Type {
	case "customer.subscription.deleted", "customer.subscription.created", "customer.subscription.updated", "customer.subscription.trial_will_end":
		var subscription stripe.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %s\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Event %s for Sub %s\n", payChannel.Channel, string(event.Type), subscription.ID)
			// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			err = s.processSubscriptionWebhook(r.Context(), string(event.Type), subscription, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %s\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "invoice.upcoming", "invoice.created", "invoice.updated", "invoice.paid", "invoice.payment_failed", "invoice.payment_action_required":
		var stripeInvoice stripe.Invoice
		err = json.Unmarshal(event.Data.Raw, &stripeInvoice)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %s\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Event %s for Invoice %s\n", payChannel.Channel, string(event.Type), stripeInvoice.ID)
			// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.
			err = s.processInvoiceWebhook(r.Context(), string(event.Type), stripeInvoice, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleInvoiceWebhookEvent: %s\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "payment_intent.created", "payment_intent.succeeded":
		var stripePayment stripe.PaymentIntent
		err = json.Unmarshal(event.Data.Raw, &stripePayment)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %s\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Event %s for Payment %s\n", payChannel.Channel, string(event.Type), stripePayment.ID)
			// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.

			//err = s.processPaymentWebhook(r.Context(), string(event.Type), stripePayment, payChannel)
			//if err != nil {
			//	g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandlePaymentWebhookEvent: %s\n", payChannel.Channel, err.Error())
			//	r.Response.WriteHeader(http.StatusBadRequest)
			//	responseBack = http.StatusBadRequest
			//}
		}
	case "charge.refund.updated":
		var stripeRefund stripe.Refund
		err = json.Unmarshal(event.Data.Raw, &stripeRefund)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %s\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Event %s for Refund %s\n", payChannel.Channel, string(event.Type), stripeRefund.ID)
			// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)

			//err = s.processPaymentWebhook(r.Context(), string(event.Type), stripePayment, payChannel)
			//if err != nil {
			//	g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandlePaymentWebhookEvent: %s\n", payChannel.Channel, err.Error())
			//	r.Response.WriteHeader(http.StatusBadRequest)
			//	responseBack = http.StatusBadRequest
			//}
		}
	case "checkout.session.completed":
		var stripeCheckoutSession stripe.CheckoutSession
		err = json.Unmarshal(event.Data.Raw, &stripeCheckoutSession)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %s\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Event %s for Refund %s\n", payChannel.Channel, string(event.Type), stripeCheckoutSession.ID)
			// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.

			err = s.processCheckoutSessionWebhook(r.Context(), string(event.Type), stripeCheckoutSession, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandlePaymentWebhookEvent: %s\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	default:
		g.Log().Errorf(r.Context(), "Webhook Channel:%s, Unhandled event type: %s\n", payChannel.Channel, event.Type)
		r.Response.WriteHeader(http.StatusBadRequest)
		responseBack = http.StatusBadRequest
	}
	log.SaveChannelHttpLog("DoRemoteChannelWebhook", event, responseBack, err, string(event.Type), nil, payChannel)
	r.Response.WriteHeader(responseBack)
}

func (s StripeWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelRedirectInternalResp, err error) {
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
			stripe.Key = payChannel.ChannelSecret
			s.setUnibeeAppInfo()
			payment := query.GetPaymentByPaymentId(r.Context(), payIdStr)
			if payment == nil || len(payment.ChannelPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.ChannelPaymentId) > 0 && payment.Status == consts.PAY_SUCCESS {
				returnUrl = payment.ReturnUrl
				response = "success"
				status = true
			} else {
				//需要去检索
				returnUrl = payment.ReturnUrl
				params := &stripe.CheckoutSessionParams{}
				result, err := session.Get(
					payment.ChannelPaymentIntentId,
					params,
				)
				if err != nil {
					response = "payment not match"
				}
				channelUser := query.GetUserChannel(r.Context(), payment.UserId, int64(payChannel.Id))
				if channelUser != nil && result != nil {
					//find
					if strings.Compare(result.Customer.ID, channelUser.ChannelUserId) != 0 {
						response = "user not match"
					} else if strings.Compare(string(result.Status), "complete") == 0 && result.PaymentIntent != nil && len(result.PaymentIntent.ID) > 0 {
						paymentIntentDetail, err := out.GetPayChannelServiceProvider(r.Context(), int64(payChannel.Id)).DoRemoteChannelPaymentDetail(r.Context(), payChannel, result.PaymentIntent.ID)
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							if paymentIntentDetail.Status == consts.PAY_SUCCESS {
								err := handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
									PaymentId:                        payment.PaymentId,
									ChannelPaymentIntentId:           payment.ChannelPaymentIntentId,
									ChannelPaymentId:                 paymentIntentDetail.ChannelPaymentId,
									TotalAmount:                      paymentIntentDetail.TotalAmount,
									PayStatusEnum:                    consts.PAY_SUCCESS,
									PaidTime:                         paymentIntentDetail.PayTime,
									PaymentAmount:                    paymentIntentDetail.PaymentAmount,
									CaptureAmount:                    0,
									Reason:                           paymentIntentDetail.Reason,
									ChannelDefaultPaymentMethod:      paymentIntentDetail.ChannelPaymentMethod,
									ChannelDetailInvoiceInternalResp: paymentIntentDetail.ChannelInvoiceDetail,
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
									ChannelPaymentIntentId:           payment.ChannelPaymentIntentId,
									ChannelPaymentId:                 paymentIntentDetail.ChannelPaymentId,
									TotalAmount:                      paymentIntentDetail.TotalAmount,
									PayStatusEnum:                    consts.PAY_FAILED,
									PaidTime:                         paymentIntentDetail.PayTime,
									PaymentAmount:                    paymentIntentDetail.PaymentAmount,
									CaptureAmount:                    0,
									Reason:                           paymentIntentDetail.Reason,
									ChannelDetailInvoiceInternalResp: paymentIntentDetail.ChannelInvoiceDetail,
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
			stripe.Key = payChannel.ChannelSecret
			s.setUnibeeAppInfo()
			unibSub := query.GetSubscriptionBySubscriptionId(r.Context(), SubIdStr)
			if unibSub == nil {
				response = "subId invalid"
			} else if len(unibSub.ChannelSubscriptionId) > 0 && unibSub.Status == consts.SubStatusActive {
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
				channelUser := query.GetUserChannel(r.Context(), unibSub.UserId, int64(payChannel.Id))
				if channelUser != nil && result.SubscriptionSearchResult().Data != nil && len(result.SubscriptionSearchResult().Data) == 1 {
					//找到
					if strings.Compare(result.SubscriptionSearchResult().Data[0].Customer.ID, channelUser.ChannelUserId) != 0 {
						response = "customId not match"
					} else {
						detail := parseStripeSubscription(result.SubscriptionSearchResult().Data[0])
						err := handler.UpdateSubWithChannelDetailBack(r.Context(), unibSub, detail)
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
	log.SaveChannelHttpLog("DoRemoteChannelRedirect", params, response, err, "", nil, payChannel)
	return &ro.ChannelRedirectInternalResp{
		Status:    status,
		Message:   response,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}

func (s StripeWebhook) processRefundWebhook(ctx context.Context, eventType string, refund stripe.Refund, payChannel *entity.MerchantChannelConfig) error {
	refundDetail, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelRefundDetail(ctx, payChannel, refund.ID)
	if err != nil {
		return err
	}
	//details.ChannelId = int64(payChannel.Id)
	//utility.Assert(len(details.ChannelUserId) > 0, "invalid channelUserId")
	//if payment.Invoice != nil {
	//	//可能来自 SubPendingUpdate 流程，需要补充 Invoice 信息获取 ChannelSubscriptionUpdateId
	//	invoiceDetails, err := s.DoRemoteChannelInvoiceDetails(ctx, payChannel, payment.Invoice.ID)
	//	if err != nil {
	//		return err
	//	}
	//	details.ChannelInvoiceDetail = invoiceDetails
	//	details.ChannelInvoiceId = payment.Invoice.ID
	//	details.ChannelSubscriptionUpdateId = invoiceDetails.ChannelInvoiceId
	//	oneSub := query.GetSubscriptionByChannelSubscriptionId(ctx, invoiceDetails.ChannelSubscriptionId)
	//	if oneSub != nil {
	//		plan := query.GetPlanById(ctx, oneSub.PlanId)
	//		planChannel := query.GetPlanChannel(ctx, oneSub.PlanId, oneSub.ChannelId)
	//		subDetails, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, oneSub)
	//		if err != nil {
	//			return err
	//		}
	//		details.ChannelSubscriptionDetail = subDetails
	//	}
	//}
	//details.UniqueId = details.ChannelPaymentIntentId
	err = handler2.HandleRefundWebhookEvent(ctx, refundDetail)
	if err != nil {
		return err
	}

	return nil
}

//
//func (s StripeWebhook) processPaymentWebhook(ctx context.Context, eventType string, payment stripe.PaymentIntent, payChannel *entity.MerchantChannelConfig) error {
//	details, err := s.DoRemoteChannelPaymentDetail(ctx, payChannel, payment.ID)
//	if err != nil {
//		return err
//	}
//	details.ChannelId = int64(payChannel.Id)
//	utility.Assert(len(details.ChannelUserId) > 0, "invalid channelUserId")
//	if payment.Invoice != nil {
//		//可能来自 SubPendingUpdate 流程，需要补充 Invoice 信息获取 ChannelSubscriptionUpdateId
//		invoiceDetails, err := s.DoRemoteChannelInvoiceDetails(ctx, payChannel, payment.Invoice.ID)
//		if err != nil {
//			return err
//		}
//		details.ChannelInvoiceDetail = invoiceDetails
//		details.ChannelInvoiceId = payment.Invoice.ID
//		details.ChannelSubscriptionUpdateId = invoiceDetails.ChannelInvoiceId
//		oneSub := query.GetSubscriptionByChannelSubscriptionId(ctx, invoiceDetails.ChannelSubscriptionId)
//		if oneSub != nil {
//			plan := query.GetPlanById(ctx, oneSub.PlanId)
//			planChannel := query.GetPlanChannel(ctx, oneSub.PlanId, oneSub.ChannelId)
//			subDetails, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, oneSub)
//			if err != nil {
//				return err
//			}
//			details.ChannelSubscriptionDetail = subDetails
//		}
//	}
//	details.UniqueId = details.ChannelPaymentIntentId
//	err = handler2.HandlePaymentWebhookEvent(ctx, details)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (s StripeWebhook) processInvoiceWebhook(ctx context.Context, eventType string, invoice stripe.Invoice, payChannel *entity.MerchantChannelConfig) error {
	invoiceDetails, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelInvoiceDetails(ctx, payChannel, invoice.ID)
	if err != nil {
		return err
	}
	// require_action status not deal here, use subscription_update webhook
	if strings.Compare("invoice.payment_action_required", eventType) == 0 {
		return gerror.New("require_action status not deal processInvoiceWebhook, use processSubscriptionWebhook webhook")
	}

	var status = consts.TO_BE_PAID
	var captureStatus = consts.AUTHORIZED
	if invoiceDetails.Status == consts.InvoiceStatusPaid {
		status = consts.PAY_SUCCESS
		captureStatus = consts.CAPTURE_REQUEST
	} else if invoiceDetails.Status == consts.InvoiceStatusFailed || invoiceDetails.Status == consts.InvoiceStatusCancelled {
		status = consts.PAY_FAILED
	}

	var channelSubscriptionDetail *ro.ChannelDetailSubscriptionInternalResp
	if len(invoiceDetails.ChannelSubscriptionId) > 0 {
		unibeeSub := query.GetSubscriptionByChannelSubscriptionId(ctx, invoiceDetails.ChannelSubscriptionId)
		var subNeedUpdate = false
		if unibeeSub == nil && len(invoiceDetails.SubscriptionId) > 0 {
			unibeeSub = query.GetSubscriptionBySubscriptionId(ctx, invoiceDetails.SubscriptionId)
			unibeeSub.ChannelSubscriptionId = invoiceDetails.ChannelSubscriptionId
			subNeedUpdate = true
		}
		if unibeeSub != nil {
			plan := query.GetPlanById(ctx, unibeeSub.PlanId)
			planChannel := query.GetPlanChannel(ctx, unibeeSub.PlanId, unibeeSub.ChannelId)
			channelSubscriptionDetail, err = out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, unibeeSub)
			if subNeedUpdate {
				err = handler.HandleSubscriptionWebhookEvent(ctx, unibeeSub, eventType, channelSubscriptionDetail)
				if err != nil {
					return err
				}
			}
		}
	}

	err = handler2.HandlePaymentWebhookEvent(ctx, &ro.ChannelPaymentRo{
		MerchantId:                  payChannel.MerchantId,
		Status:                      status,
		CaptureStatus:               captureStatus,
		Currency:                    invoiceDetails.Currency,
		TotalAmount:                 invoiceDetails.TotalAmount,
		PaymentAmount:               invoiceDetails.PaymentAmount,
		BalanceAmount:               invoiceDetails.BalanceAmount,
		BalanceStart:                invoiceDetails.BalanceStart,
		BalanceEnd:                  invoiceDetails.BalanceEnd,
		Reason:                      invoiceDetails.Reason,
		UniqueId:                    invoiceDetails.ChannelInvoiceId,
		PayTime:                     gtime.NewFromTimeStamp(invoiceDetails.PaymentTime),
		CreateTime:                  gtime.NewFromTimeStamp(invoiceDetails.CreateTime),
		CancelTime:                  gtime.NewFromTimeStamp(invoiceDetails.CancelTime),
		ChannelId:                   int64(payChannel.Id),
		ChannelUserId:               invoiceDetails.ChannelUserId,
		ChannelPaymentId:            invoiceDetails.ChannelPaymentId,
		ChannelPaymentMethod:        invoiceDetails.ChannelDefaultPaymentMethod,
		ChannelInvoiceId:            invoiceDetails.ChannelInvoiceId,
		ChannelSubscriptionId:       invoiceDetails.ChannelSubscriptionId,
		ChannelSubscriptionUpdateId: invoiceDetails.ChannelInvoiceId,
		ChannelInvoiceDetail:        invoiceDetails,
		ChannelSubscriptionDetail:   channelSubscriptionDetail,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s StripeWebhook) processSubscriptionWebhook(ctx context.Context, eventType string, subscription stripe.Subscription, payChannel *entity.MerchantChannelConfig) error {
	unibeeSub := query.GetSubscriptionByChannelSubscriptionId(ctx, subscription.ID)
	if unibeeSub == nil {
		if unibSubId, ok := subscription.Metadata["SubId"]; ok {
			unibeeSub = query.GetSubscriptionBySubscriptionId(ctx, unibSubId)
			unibeeSub.ChannelSubscriptionId = subscription.ID
		}
	}
	if unibeeSub != nil {
		plan := query.GetPlanById(ctx, unibeeSub.PlanId)
		planChannel := query.GetPlanChannel(ctx, unibeeSub.PlanId, unibeeSub.ChannelId)
		details, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, unibeeSub)
		if err != nil {
			return err
		}

		err = handler.HandleSubscriptionWebhookEvent(ctx, unibeeSub, eventType, details)
		if err != nil {
			return err
		}
		if details.Status == consts.SubStatusIncomplete && len(details.ChannelLatestInvoiceId) > 0 {
			//处理支付需要授权事件
			invoiceDetails, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelInvoiceDetails(ctx, payChannel, details.ChannelLatestInvoiceId)
			if err != nil {
				return err
			}
			if invoiceDetails.Status != consts.InvoiceStatusPaid {
				//有支付授权 todo mark
				var channelSubscriptionDetail *ro.ChannelDetailSubscriptionInternalResp
				if len(invoiceDetails.ChannelSubscriptionId) > 0 {
					oneSub := query.GetSubscriptionByChannelSubscriptionId(ctx, invoiceDetails.ChannelSubscriptionId)
					if oneSub != nil {
						plan := query.GetPlanById(ctx, oneSub.PlanId)
						planChannel := query.GetPlanChannel(ctx, oneSub.PlanId, oneSub.ChannelId)
						channelSubscriptionDetail, err = out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, oneSub)
					}
				}

				err = handler2.HandlePaymentWebhookEvent(ctx, &ro.ChannelPaymentRo{
					MerchantId:                  payChannel.MerchantId,
					Status:                      consts.TO_BE_PAID,
					CaptureStatus:               consts.WAITING_AUTHORIZED,
					Currency:                    invoiceDetails.Currency,
					TotalAmount:                 invoiceDetails.TotalAmount,
					PaymentAmount:               invoiceDetails.PaymentAmount,
					BalanceAmount:               invoiceDetails.BalanceAmount,
					BalanceStart:                invoiceDetails.BalanceStart,
					BalanceEnd:                  invoiceDetails.BalanceEnd,
					Reason:                      invoiceDetails.Reason,
					UniqueId:                    invoiceDetails.ChannelInvoiceId,
					PayTime:                     gtime.NewFromTimeStamp(invoiceDetails.PaymentTime),
					CreateTime:                  gtime.NewFromTimeStamp(invoiceDetails.CreateTime),
					CancelTime:                  gtime.NewFromTimeStamp(invoiceDetails.CancelTime),
					ChannelId:                   int64(payChannel.Id),
					ChannelUserId:               invoiceDetails.ChannelUserId,
					ChannelPaymentId:            invoiceDetails.ChannelPaymentId,
					ChannelPaymentMethod:        invoiceDetails.ChannelDefaultPaymentMethod,
					ChannelInvoiceId:            invoiceDetails.ChannelInvoiceId,
					ChannelSubscriptionId:       invoiceDetails.ChannelSubscriptionId,
					ChannelSubscriptionUpdateId: invoiceDetails.ChannelInvoiceId,
					ChannelInvoiceDetail:        invoiceDetails,
					ChannelSubscriptionDetail:   channelSubscriptionDetail,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return gerror.New("subscription not found on channelSubId:" + subscription.ID)
	}
}

func (s StripeWebhook) processCheckoutSessionWebhook(ctx context.Context, event string, checkoutSession stripe.CheckoutSession, payChannel *entity.MerchantChannelConfig) error {
	if paymentId, ok := checkoutSession.Metadata["PaymentId"]; ok {
		payment := query.GetPaymentByPaymentId(ctx, paymentId)
		if checkoutSession.PaymentIntent != nil {
			paymentIntentDetail, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPaymentDetail(ctx, payChannel, checkoutSession.PaymentIntent.ID)
			if err != nil {
				return gerror.New(fmt.Sprintf("%s", err.Error()))
			}
			if paymentIntentDetail.Status == consts.PAY_SUCCESS {
				err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					ChannelPaymentIntentId:           payment.ChannelPaymentIntentId,
					ChannelPaymentId:                 paymentIntentDetail.ChannelPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_SUCCESS,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.Reason,
					ChannelDefaultPaymentMethod:      paymentIntentDetail.ChannelPaymentMethod,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.ChannelInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PAY_FAILED {
				err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
					PaymentId:                        payment.PaymentId,
					ChannelPaymentIntentId:           payment.ChannelPaymentIntentId,
					ChannelPaymentId:                 paymentIntentDetail.ChannelPaymentId,
					TotalAmount:                      paymentIntentDetail.TotalAmount,
					PayStatusEnum:                    consts.PAY_FAILED,
					PaidTime:                         paymentIntentDetail.PayTime,
					PaymentAmount:                    paymentIntentDetail.PaymentAmount,
					CaptureAmount:                    0,
					Reason:                           paymentIntentDetail.Reason,
					ChannelDetailInvoiceInternalResp: paymentIntentDetail.ChannelInvoiceDetail,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			}
			return nil
		} else {
			return gerror.New("no PaymentIntent")
		}
	} else {
		return gerror.New("no PaymentId Metadata")
	}
}

func parseStripeSubscription(subscription *stripe.Subscription) *ro.ChannelDetailSubscriptionInternalResp {
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

	return &ro.ChannelDetailSubscriptionInternalResp{
		Status:                 status,
		ChannelSubscriptionId:  subscription.ID,
		ChannelStatus:          string(subscription.Status),
		Data:                   utility.FormatToJsonString(subscription),
		ChannelItemData:        utility.MarshalToJsonString(subscription.Items.Data),
		ChannelLatestInvoiceId: subscription.LatestInvoice.ID,
		ChannelLatestPaymentId: latestChannelPaymentId,
		CancelAtPeriodEnd:      subscription.CancelAtPeriodEnd,
		CurrentPeriodStart:     subscription.CurrentPeriodStart,
		CurrentPeriodEnd:       subscription.CurrentPeriodEnd,
		BillingCycleAnchor:     subscription.BillingCycleAnchor,
		TrialEnd:               subscription.TrialEnd,
	}
}
