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
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
	"net/http"
	"strconv"
	"strings"
	"unibee/internal/consts"
	_gateway "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/ro"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
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
	if len(result.WebhookEndpointList().Data) == 0 {
		//create
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				stripe.String("invoice.upcoming"),
				stripe.String("invoice.created"),
				stripe.String("invoice.updated"),
				stripe.String("invoice.paid"),
				stripe.String("invoice.voided"),
				stripe.String("invoice.will_be_due"),
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("payment_intent.canceled"),
				stripe.String("payment_intent.partially_funded"),
				stripe.String("payment_intent.payment_failed"),
				stripe.String("payment_intent.requires_action"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			Metadata: map[string]string{"MerchantId": strconv.FormatUint(gateway.MerchantId, 10)},
			URL:      stripe.String(_gateway.GetPaymentWebhookEntranceUrl(gateway.Id)),
		}
		result, err := webhookendpoint.New(params)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", params, result, err, "", nil, gateway)
		if err != nil {
			return nil
		}
		//更新 secret
		utility.Assert(len(result.Secret) > 0, "secret is nil")
		err = query.UpdateGatewayWebhookSecret(ctx, gateway.Id, result.Secret)
		if err != nil {
			return err
		}
	} else {
		utility.Assert(len(result.WebhookEndpointList().Data) == 1, "internal webhook update, count is not 1")
		one := result.WebhookEndpointList().Data[0]
		utility.Assert(strings.Compare(one.Status, "enabled") == 0, "webhook not status enabled")
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				//webhook
				stripe.String("invoice.upcoming"),
				stripe.String("invoice.created"),
				stripe.String("invoice.updated"),
				stripe.String("invoice.paid"),
				stripe.String("invoice.voided"),
				stripe.String("invoice.will_be_due"),
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.succeeded"),
				stripe.String("payment_intent.canceled"),
				stripe.String("payment_intent.partially_funded"),
				stripe.String("payment_intent.payment_failed"),
				stripe.String("payment_intent.requires_action"),
				stripe.String("checkout.session.completed"),
				stripe.String("charge.refund.updated"),
			},
			URL:      stripe.String(_gateway.GetPaymentWebhookEntranceUrl(gateway.Id)),
			Metadata: map[string]string{"MerchantId": strconv.FormatUint(gateway.MerchantId, 10)},
		}
		result, err := webhookendpoint.Update(one.ID, params)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", params, result, err, one.ID, nil, gateway)
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
	g.Log().Debug(r.Context(), "Receive_Webhook_Channel: ", gateway.GatewayName, " hook:", string(data))

	var responseBack = http.StatusOK
	var requestId = ""
	switch event.Type {
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
	var response string
	var status = false
	var returnUrl = ""
	if len(payIdStr) > 0 {
		response = ""
		//Payment Redirect
		if r.Get("success").Bool() {
			stripe.Key = gateway.GatewaySecret
			s.setUnibeeAppInfo()
			payment := query.GetPaymentByPaymentId(r.Context(), payIdStr)
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				returnUrl = payment.ReturnUrl
				response = "success"
				status = true
			} else {
				returnUrl = payment.ReturnUrl
				params := &stripe.CheckoutSessionParams{}
				result, err := session.Get(
					payment.GatewayPaymentIntentId,
					params,
				)
				if err != nil {
					response = "payment not match"
				}
				gatewayUser := query.GetGatewayUser(r.Context(), payment.UserId, gateway.Id)
				if gatewayUser != nil && result != nil {
					//find
					if strings.Compare(result.Customer.ID, gatewayUser.GatewayUserId) != 0 {
						response = "user not match"
					} else if strings.Compare(string(result.Status), "complete") == 0 && result.PaymentIntent != nil && len(result.PaymentIntent.ID) > 0 {
						paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayPaymentDetail(r.Context(), gateway, result.PaymentIntent.ID)
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
									CaptureAmount:          0,
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
									TotalAmount:            paymentIntentDetail.TotalAmount,
									PayStatusEnum:          consts.PaymentFailed,
									PaidTime:               paymentIntentDetail.PayTime,
									PaymentAmount:          paymentIntentDetail.PaymentAmount,
									CaptureAmount:          0,
									Reason:                 paymentIntentDetail.Reason,
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
	refundDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayRefundDetail(ctx, gateway, refund.ID)
	if err != nil {
		return err
	}
	err = handler2.HandleRefundWebhookEvent(ctx, refundDetail)
	if err != nil {
		return err
	}

	return nil
}

func (s StripeWebhook) processPaymentWebhook(ctx context.Context, eventType string, stripePayment stripe.PaymentIntent, gateway *entity.MerchantGateway) error {
	if paymentId, ok := stripePayment.Metadata["PaymentId"]; ok {
		// PaymentIntent Under UniBee Control
		payment := query.GetPaymentByPaymentId(ctx, paymentId)
		if payment != nil {
			paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, stripePayment.ID)
			if err != nil {
				return err
			}
			if len(paymentIntentDetail.PaymentData) == 0 && stripePayment.NextAction != nil {
				paymentIntentDetail.PaymentData = utility.MarshalToJsonString(stripePayment.NextAction)
			}
			if paymentIntentDetail.Status == consts.PaymentSuccess {
				err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentSuccess,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.Reason,
					GatewayPaymentMethod:   paymentIntentDetail.GatewayPaymentMethod,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PaymentFailed {
				err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentFailed,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.CancelReason,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PaymentCancelled {
				err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: paymentIntentDetail.GatewayPaymentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentCancelled,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.CancelReason,
				})
				if err != nil {
					return err
				}
			} else if paymentIntentDetail.AuthorizeStatus == consts.WaitingAuthorized {
				err := handler2.HandlePayNeedAuthorized(ctx, payment, paymentIntentDetail.AuthorizeReason, paymentIntentDetail.PaymentData)
				if err != nil {
					return err
				}
			}
		} else {
			return gerror.New("Payment Not Found")
		}
	} else {
		//Maybe PaymentIntent Create By Invoice
		g.Log().Errorf(ctx, "No PaymentId Metadata PaymentIntentId:%s", stripePayment.ID)
		return nil
	}
	return nil
}

func parseStripeInvoice(detail stripe.Invoice) *GatewayDetailInvoiceInternalResp {
	var status consts.InvoiceStatusEnum = consts.InvoiceStatusInit
	if strings.Compare(string(detail.Status), "draft") == 0 {
		status = consts.InvoiceStatusPending
	} else if strings.Compare(string(detail.Status), "open") == 0 {
		status = consts.InvoiceStatusProcessing
	} else if strings.Compare(string(detail.Status), "paid") == 0 {
		status = consts.InvoiceStatusPaid
	} else if strings.Compare(string(detail.Status), "uncollectible") == 0 {
		status = consts.InvoiceStatusFailed
	} else if strings.Compare(string(detail.Status), "void") == 0 {
		status = consts.InvoiceStatusCancelled
	}
	var invoiceItems []*ro.InvoiceItemDetailRo
	for _, line := range detail.Lines.Data {
		var start int64 = 0
		var end int64 = 0
		if line.Period != nil {
			start = line.Period.Start
			end = line.Period.End
		}
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               strings.ToUpper(string(line.Currency)),
			Amount:                 line.Amount,
			AmountExcludingTax:     line.AmountExcludingTax,
			UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
			Description:            line.Description,
			Proration:              line.Proration,
			Quantity:               line.Quantity,
			PeriodStart:            start,
			PeriodEnd:              end,
		})
	}

	var gatewayPaymentId string
	if detail.PaymentIntent != nil {
		gatewayPaymentId = detail.PaymentIntent.ID
	}
	var subscriptionId string
	if detail.SubscriptionDetails != nil {
		subscriptionId = detail.SubscriptionDetails.Metadata["SubId"]
	}
	var gatewayUserId string
	if detail.Customer != nil {
		gatewayUserId = detail.Customer.ID
	}
	var paymentTime int64
	var cancelTime int64
	if detail.StatusTransitions != nil {
		paymentTime = detail.StatusTransitions.PaidAt
		cancelTime = detail.StatusTransitions.VoidedAt
	}
	var gatewayDefaultPaymentMethod = ""
	if detail.DefaultPaymentMethod != nil {
		gatewayDefaultPaymentMethod = detail.DefaultPaymentMethod.ID
	}

	return &GatewayDetailInvoiceInternalResp{
		GatewayDefaultPaymentMethod:    gatewayDefaultPaymentMethod,
		TotalAmount:                    detail.Total,
		PaymentAmount:                  detail.AmountPaid,
		BalanceAmount:                  -(detail.StartingBalance) - -(detail.EndingBalance),
		BalanceStart:                   -detail.StartingBalance,
		BalanceEnd:                     -detail.EndingBalance,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		Status:                         status,
		Link:                           detail.HostedInvoiceURL,
		GatewayStatus:                  string(detail.Status),
		GatewayInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
		GatewayInvoiceId:               detail.ID,
		GatewayUserId:                  gatewayUserId,
		SubscriptionId:                 subscriptionId,
		GatewayPaymentId:               gatewayPaymentId,
		PaymentTime:                    paymentTime,
		Reason:                         string(detail.BillingReason),
		CreateTime:                     detail.Created,
		CancelTime:                     cancelTime,
	}
}

func (s StripeWebhook) processInvoiceWebhook(ctx context.Context, eventType string, invoice stripe.Invoice, gateway *entity.MerchantGateway) error {
	utility.Assert(len(invoice.ID) > 0, "processInvoiceWebhook gatewayInvoiceId Invalid")
	invoiceDetails := parseStripeInvoice(invoice)

	var status = consts.PaymentCreated
	var authorizeStatus = consts.Authorized
	var authorizeReason = ""
	var cancelReason = ""
	var paymentData = ""
	if invoiceDetails.Status == consts.InvoiceStatusPaid {
		status = consts.PaymentSuccess
		authorizeStatus = consts.CaptureRequest
	} else if invoiceDetails.Status == consts.InvoiceStatusFailed {
		status = consts.PaymentFailed
	} else if invoiceDetails.Status == consts.InvoiceStatusCancelled {
		status = consts.PaymentCancelled
	} else if strings.Compare("invoice.payment_action_required", eventType) == 0 {
		authorizeStatus = consts.WaitingAuthorized
	}

	if len(invoiceDetails.GatewayPaymentId) > 0 {
		paymentIntentDetail, _ := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, invoiceDetails.GatewayPaymentId)
		if paymentIntentDetail != nil {
			authorizeReason = paymentIntentDetail.AuthorizeReason
			cancelReason = paymentIntentDetail.CancelReason
			paymentData = paymentIntentDetail.PaymentData
		}
	}

	err := handler2.HandlePaymentWebhookEvent(ctx, &ro.GatewayPaymentRo{
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		Currency:             invoiceDetails.Currency,
		TotalAmount:          invoiceDetails.TotalAmount,
		PaymentAmount:        invoiceDetails.PaymentAmount,
		BalanceAmount:        invoiceDetails.BalanceAmount,
		BalanceStart:         invoiceDetails.BalanceStart,
		BalanceEnd:           invoiceDetails.BalanceEnd,
		Reason:               invoiceDetails.Reason,
		CancelReason:         cancelReason,
		PaymentData:          paymentData,
		UniqueId:             invoiceDetails.GatewayInvoiceId,
		PayTime:              gtime.NewFromTimeStamp(invoiceDetails.PaymentTime),
		CreateTime:           gtime.NewFromTimeStamp(invoiceDetails.CreateTime),
		CancelTime:           gtime.NewFromTimeStamp(invoiceDetails.CancelTime),
		GatewayUserId:        invoiceDetails.GatewayUserId,
		GatewayPaymentId:     invoiceDetails.GatewayPaymentId,
		GatewayPaymentMethod: invoiceDetails.GatewayDefaultPaymentMethod,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s StripeWebhook) processCheckoutSessionWebhook(ctx context.Context, event string, checkoutSession stripe.CheckoutSession, gateway *entity.MerchantGateway) error {
	if paymentId, ok := checkoutSession.Metadata["PaymentId"]; ok {
		payment := query.GetPaymentByPaymentId(ctx, paymentId)
		if checkoutSession.PaymentIntent != nil {
			paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, checkoutSession.PaymentIntent.ID)
			if err != nil {
				return gerror.New(fmt.Sprintf("%s", err.Error()))
			}
			if paymentIntentDetail.Status == consts.PaymentSuccess {
				err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentSuccess,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.Reason,
					GatewayPaymentMethod:   paymentIntentDetail.GatewayPaymentMethod,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PaymentFailed {
				err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentFailed,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.Reason,
				})
				if err != nil {
					return gerror.New(fmt.Sprintf("%s", err.Error()))
				}
			} else if paymentIntentDetail.Status == consts.PaymentCancelled {
				err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: paymentIntentDetail.GatewayPaymentId,
					GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
					TotalAmount:            paymentIntentDetail.TotalAmount,
					PayStatusEnum:          consts.PaymentCancelled,
					PaidTime:               paymentIntentDetail.PayTime,
					PaymentAmount:          paymentIntentDetail.PaymentAmount,
					CaptureAmount:          0,
					Reason:                 paymentIntentDetail.CancelReason,
				})
				if err != nil {
					return err
				}
			} else if paymentIntentDetail.AuthorizeStatus == consts.WaitingAuthorized {
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

type GatewayDetailInvoiceInternalResp struct {
	GatewayDefaultPaymentMethod    string                    `json:"gatewayDefaultPaymentMethod"`
	SubscriptionId                 string                    `json:"subscriptionId"           `
	TotalAmount                    int64                     `json:"totalAmount"        `
	PaymentAmount                  int64                     `json:"paymentAmount"              `
	BalanceAmount                  int64                     `json:"balanceAmount"              `
	BalanceStart                   int64                     `json:"balanceStart"              `
	BalanceEnd                     int64                     `json:"balanceEnd"              `
	TotalAmountExcludingTax        int64                     `json:"totalAmountExcludingTax"        `
	TaxAmount                      int64                     `json:"taxAmount"          `
	SubscriptionAmount             int64                     `json:"subscriptionAmount" `
	SubscriptionAmountExcludingTax int64                     `json:"subscriptionAmountExcludingTax" `
	Currency                       string                    `json:"currency"           `
	Lines                          []*ro.InvoiceItemDetailRo `json:"lines"              `
	Status                         consts.InvoiceStatusEnum  `json:"status"             `
	Reason                         string                    `json:"reason"             `
	GatewayUserId                  string                    `json:"gatewayUserId"             `
	Link                           string                    `json:"link"               `
	GatewayStatus                  string                    `json:"gatewayStatus"      `
	GatewayInvoiceId               string                    `json:"gatewayInvoiceId"   `
	GatewayInvoicePdf              string                    `json:"GatewayInvoicePdf"   `
	PeriodEnd                      int64                     `json:"periodEnd"`
	PeriodStart                    int64                     `json:"periodStart"`
	GatewayPaymentId               string                    `json:"gatewayPaymentId"`
	PaymentTime                    int64                     `json:"paymentTime"        `
	CreateTime                     int64                     `json:"createTime"        `
	CancelTime                     int64                     `json:"cancelTime"        `
}
