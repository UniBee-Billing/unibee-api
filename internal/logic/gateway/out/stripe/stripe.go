package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/balance"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/invoiceitem"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	sub "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/taxrate"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/out"
	"go-oversea-pay/internal/logic/gateway/out/log"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/gateway/util"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Stripe struct {
}

func (s Stripe) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.BalanceParams{}
	response, err := balance.Get(params)
	if err != nil {
		return nil, err
	}

	var availableBalances []*ro.ChannelBalance
	for _, item := range response.Available {
		availableBalances = append(availableBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var connectReservedBalances []*ro.ChannelBalance
	for _, item := range response.ConnectReserved {
		connectReservedBalances = append(connectReservedBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var pendingBalances []*ro.ChannelBalance
	for _, item := range response.ConnectReserved {
		pendingBalances = append(pendingBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	return &ro.ChannelMerchantBalanceQueryInternalResp{
		AvailableBalance:       availableBalances,
		ConnectReservedBalance: connectReservedBalances,
		PendingBalance:         pendingBalances,
	}, nil
}

func (s Stripe) DoRemoteChannelUserBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, customerId string) (res *ro.ChannelUserBalanceQueryInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.CustomerParams{}
	response, err := customer.Get(customerId, params)
	if err != nil {
		return nil, err
	}
	var cashBalances []*ro.ChannelBalance
	if response.CashBalance != nil {
		for currency, amount := range response.CashBalance.Available {
			cashBalances = append(cashBalances, &ro.ChannelBalance{
				Amount:   amount,
				Currency: strings.ToUpper(currency),
			})
		}
	}

	var invoiceCreditBalances []*ro.ChannelBalance
	for currency, amount := range response.InvoiceCreditBalance {
		invoiceCreditBalances = append(invoiceCreditBalances, &ro.ChannelBalance{
			Amount:   amount,
			Currency: strings.ToUpper(currency),
		})
	}
	return &ro.ChannelUserBalanceQueryInternalResp{
		Balance: &ro.ChannelBalance{
			Amount:   response.Balance,
			Currency: strings.ToUpper(string(response.Currency)),
		},
		CashBalance:          cashBalances,
		InvoiceCreditBalance: invoiceCreditBalances,
		Description:          response.Description,
		Email:                response.Email,
	}, nil
}

func parseStripeInvoice(detail *stripe.Invoice, channelId int64) *ro.ChannelDetailInvoiceInternalResp {
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
	var invoiceItems []*ro.ChannelDetailInvoiceItem
	for _, line := range detail.Lines.Data {
		var start int64 = 0
		var end int64 = 0
		if line.Period != nil {
			start = line.Period.Start
			end = line.Period.End
		}
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
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

	var channelPaymentId string
	if detail.PaymentIntent != nil {
		channelPaymentId = detail.PaymentIntent.ID
	}

	return &ro.ChannelDetailInvoiceInternalResp{
		ChannelSubscriptionId:          detail.Subscription.ID,
		TotalAmount:                    detail.Total,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		ChannelId:                      channelId,
		Status:                         status,
		ChannelUserId:                  detail.Customer.ID,
		Link:                           detail.HostedInvoiceURL,
		ChannelStatus:                  string(detail.Status),
		ChannelInvoiceId:               detail.ID,
		ChannelInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
		ChannelPaymentId:               channelPaymentId,
	}
}

func (s Stripe) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.OverseaPayChannel, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoiceMarkUncollectibleParams{}
	response, err := invoice.MarkUncollectible(cancelInvoiceInternalReq.ChannelInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCancel", params, response, err, "", nil, payChannel)
	return parseStripeInvoice(response, int64(payChannel.Id)), nil
}

func (s Stripe) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	if len(createInvoiceInternalReq.Invoice.ChannelUserId) == 0 {
		params := &stripe.CustomerParams{
			//Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
			Email: stripe.String(createInvoiceInternalReq.Invoice.SendEmail),
		}

		createCustomResult, err := customer.New(params)
		if err != nil {
			g.Log().Printf(ctx, "customer.New: %v", err.Error())
			return nil, err
		}
		createInvoiceInternalReq.Invoice.ChannelUserId = createCustomResult.ID
	}

	params := &stripe.InvoiceParams{
		Currency: stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
		Customer: stripe.String(createInvoiceInternalReq.Invoice.ChannelUserId)}
	if createInvoiceInternalReq.PayMethod == 1 {
		params.CollectionMethod = stripe.String("charge_automatically")
	} else {
		params.CollectionMethod = stripe.String("send_invoice")
		if createInvoiceInternalReq.DaysUtilDue > 0 {
			params.DaysUntilDue = stripe.Int64(int64(createInvoiceInternalReq.DaysUtilDue))
		}
		// todo mark tax 设置
	}
	result, err := invoice.New(params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCancel", params, result, err, "New", nil, payChannel)

	for _, line := range createInvoiceInternalReq.InvoiceLines {
		ItemParams := &stripe.InvoiceItemParams{
			Invoice:     stripe.String(result.ID),
			Currency:    stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
			UnitAmount:  stripe.Int64(line.UnitAmountExcludingTax),
			Description: stripe.String(line.Description),
			Quantity:    stripe.Int64(line.Quantity),
			Customer:    stripe.String(createInvoiceInternalReq.Invoice.ChannelUserId)}
		// todo mark tax 设置
		_, err = invoiceitem.New(ItemParams)
		if err != nil {
			return nil, err
		}
	}

	finalizeInvoiceParam := &stripe.InvoiceFinalizeInvoiceParams{}
	if createInvoiceInternalReq.PayMethod == 1 {
		finalizeInvoiceParam.AutoAdvance = stripe.Bool(true)
	} else {
		finalizeInvoiceParam.AutoAdvance = stripe.Bool(false)
	}

	// todo mark 总价格验证

	detail, err := invoice.FinalizeInvoice(result.ID, finalizeInvoiceParam)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCancel", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, payChannel)
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
	var invoiceItems []*ro.ChannelDetailInvoiceItem
	for _, line := range detail.Lines.Data {
		var start int64 = 0
		var end int64 = 0
		if line.Period != nil {
			start = line.Period.Start
			end = line.Period.End
		}
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
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

	return &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    detail.Total,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		ChannelId:                      int64(payChannel.Id),
		Status:                         status,
		ChannelUserId:                  detail.Customer.ID,
		Link:                           detail.HostedInvoiceURL,
		ChannelStatus:                  string(detail.Status),
		ChannelInvoiceId:               detail.ID,
		ChannelInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
	}, nil
}

func (s Stripe) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoicePayParams{}
	response, err := invoice.Pay(payInvoiceInternalReq.ChannelInvoiceId, params)
	log.SaveChannelHttpLog("DoRemoteChannelInvoicePay", params, response, err, "", nil, payChannel)
	return parseStripeInvoice(response, int64(payChannel.Id)), nil
}

func (s Stripe) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.InvoiceParams{}
	detail, err := invoice.Get(channelInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceDetails", params, detail, err, "", nil, payChannel)
	return parseStripeInvoice(detail, int64(payChannel.Id)), nil
}

// DoRemoteChannelSubscriptionNewTrialEnd https://stripe.com/docs/billing/subscriptions/billing-cycle#add-a-trial-to-change-the-billing-cycle
func (s Stripe) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		TrialEnd:          stripe.Int64(newTrialEnd),
		ProrationBehavior: stripe.String("none"),
	}
	_, err = sub.Update(subscription.ChannelSubscriptionId, params)
	if err != nil {
		return nil, err
	}

	details, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, subscription)
	if err != nil {
		return nil, err
	}
	if details.TrialEnd != newTrialEnd {
		return nil, gerror.New("update new trial end error")
	}
	return details, nil
}

// 测试数据
// 付款成功
// 4242 4242 4242 4242
// 付款需要验证
// 4000 0025 0000 3155
// 付款被拒绝
// 4000 0000 0000 9995
func (s Stripe) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func (s Stripe) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	{
		if len(subscriptionRo.Subscription.ChannelUserId) == 0 {
			params := &stripe.CustomerParams{
				//Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
				Email: stripe.String(subscriptionRo.Subscription.CustomerEmail),
			}

			createCustomResult, err := customer.New(params)
			if err != nil {
				g.Log().Printf(ctx, "customer.New: %v", err.Error())
				return nil, err
			}
			subscriptionRo.Subscription.ChannelUserId = createCustomResult.ID
		}
		//税率创建并处理

		channelVatRate := query.GetSubscriptionVatRateChannel(ctx, subscriptionRo.VatCountryRate.Id, channelEntity.Id)
		if channelVatRate == nil {
			params := &stripe.TaxRateParams{
				DisplayName: stripe.String("VAT"),
				Description: stripe.String(subscriptionRo.VatCountryRate.CountryName),
				Percentage:  stripe.Float64(utility.ConvertTaxPercentageToPercentageFloat(subscriptionRo.VatCountryRate.StandardTaxPercentage)),
				Country:     stripe.String(subscriptionRo.VatCountryRate.CountryCode),
				Active:      stripe.Bool(true),
				//Jurisdiction: stripe.String("DE"),
				Inclusive: stripe.Bool(false),
			}
			vatCreateResult, err := taxrate.New(params)
			if err != nil {
				g.Log().Printf(ctx, "taxrate.New: %v", err.Error())
				return nil, err
			}
			channelVatRate = &entity.SubscriptionVatRateChannel{
				VatRateId:        int64(subscriptionRo.VatCountryRate.Id),
				ChannelId:        int64(channelEntity.Id),
				ChannelVatRateId: vatCreateResult.ID,
			}
			result, err := dao.SubscriptionVatRateChannel.Ctx(ctx).Data(channelVatRate).OmitEmpty().Insert(channelVatRate)
			if err != nil {
				err = gerror.Newf(`SubscriptionVatRateChannel record insert failure %s`, err.Error())
				return nil, err
			}
			id, _ := result.LastInsertId()
			channelVatRate.Id = uint64(uint(id))
		}

		taxInclusive := true
		if subscriptionRo.Plan.TaxInclusive == 0 {
			//税费不包含
			taxInclusive = false
		}

		var checkoutMode = true
		if checkoutMode {
			items := []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.CheckoutSessionLineItemParams{
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
				})
			}
			subscriptionParams := &stripe.CheckoutSessionParams{
				Customer:  stripe.String(subscriptionRo.Subscription.ChannelUserId),
				Currency:  stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
				LineItems: items,
				AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
					Enabled: stripe.Bool(!taxInclusive), //默认值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_tates
				},
				Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
				SuccessURL: stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, true)),
				CancelURL:  stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, false)),
				SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
					Metadata: map[string]string{
						"SubId": subscriptionRo.Subscription.SubscriptionId,
					},
					DefaultTaxRates: []*string{stripe.String(channelVatRate.ChannelVatRateId)},
				},
			}
			createSubscription, err := session.New(subscriptionParams)
			log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCreateSession", subscriptionParams, createSubscription, err, "", nil, channelEntity)
			if err != nil {
				return nil, err
			}
			return &ro.ChannelCreateSubscriptionInternalResp{
				ChannelUserId: subscriptionRo.Subscription.ChannelUserId,
				Link:          createSubscription.URL,
				//ChannelSubscriptionId:     createSubscription.Subscription.ID,
				//ChannelSubscriptionStatus: string(createSubscription.Subscription.Status),
				Data:   utility.FormatToJsonString(createSubscription),
				Status: 0, //todo mark
				//Paid:                      createSubscription.Subscription.LatestInvoice.Paid,
			}, nil
		} else {
			items := []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Main",
						"BillingPlanId":   strconv.FormatInt(subscriptionRo.PlanChannel.PlanId, 10),
					},
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.SubscriptionItemsParams{
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Addon",
						"BillingPlanId":   strconv.FormatInt(addon.AddonPlanChannel.PlanId, 10),
					},
				})
			}
			subscriptionParams := &stripe.SubscriptionParams{
				Customer: stripe.String(subscriptionRo.Subscription.ChannelUserId),
				Currency: stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
				Items:    items,
				AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
					Enabled: stripe.Bool(!taxInclusive), //默认值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_tates
				},
				PaymentBehavior:  stripe.String("default_incomplete"),   // todo mark https://stripe.com/docs/api/subscriptions/create
				CollectionMethod: stripe.String("charge_automatically"), //默认行为 charge_automatically，自动扣款
				Metadata: map[string]string{
					"SubId": subscriptionRo.Subscription.SubscriptionId,
				},
				DefaultTaxRates: []*string{stripe.String(channelVatRate.ChannelVatRateId)},
			}
			subscriptionParams.AddExpand("latest_invoice.payment_intent")
			createSubscription, err := sub.New(subscriptionParams)
			log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCreate", subscriptionParams, createSubscription, err, "", nil, channelEntity)
			if err != nil {
				return nil, err
			}

			return &ro.ChannelCreateSubscriptionInternalResp{
				ChannelUserId:             subscriptionRo.Subscription.ChannelUserId,
				Link:                      createSubscription.LatestInvoice.HostedInvoiceURL,
				ChannelSubscriptionId:     createSubscription.ID,
				ChannelSubscriptionStatus: string(createSubscription.Status),
				Data:                      utility.FormatToJsonString(createSubscription),
				Status:                    0, //todo mark
				Paid:                      createSubscription.LatestInvoice.Paid,
			}, nil
		}
	}
	//{
	//	//付款链接方式可能存在多次重复付款问题
	//	params := &stripe.PaymentLinkParams{
	//		LineItems: []*stripe.PaymentLinkLineItemParams{
	//			{
	//				Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
	//				Quantity: stripe.Int64(1),
	//			},
	//			//{
	//			//	Price: stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
	//			//},
	//		},
	//		AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
	//			Type: stripe.String(string(stripe.PaymentLinkAfterCompletionTypeRedirect)),
	//			Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
	//				URL: stripe.String("https://www.baidu.com"),
	//			},
	//		},
	//		//不启用试用期
	//		//SubscriptionData: &stripe.PaymentLinkSubscriptionDataParams{
	//		//	TrialPeriodDays: stripe.Int64(7),
	//		//},
	//	}
	//	createSubscription, err := paymentlink.New(params)
	//	if err != nil {
	//		return nil, err
	//	}
	//	jsonData, _ := gjson.Marshal(createSubscription)
	//	return &ro.ChannelCreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     createSubscription.ID,
	//		ChannelSubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}
	//{
	//	checkoutParams := &stripe.CheckoutSessionParams{
	//		Metadata: map[string]string{
	//			"orderId": subscriptionRo.Subscription.SubscriptionId,
	//		},
	//		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	//		LineItems: []*stripe.CheckoutSessionLineItemParams{
	//			{
	//				Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
	//				Quantity: stripe.Int64(1),
	//			},
	//		},
	//		SuccessURL: stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, true)),
	//		CancelURL:  stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, false)),
	//	}
	//
	//	result, err := session.New(checkoutParams)
	//	if err != nil {
	//		return nil, err
	//	}
	//	jsonData, _ := gjson.Marshal(result)
	//	return &ro.ChannelCreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     result.ID,
	//		ChannelSubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}

}

// DoRemoteChannelSubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel?dashboard-or-api=api
func (s Stripe) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionCancelInternalReq.Subscription.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionCancelInternalReq.Subscription.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionCancelParams{}
	params.InvoiceNow = stripe.Bool(subscriptionCancelInternalReq.InvoiceNow)
	params.Prorate = stripe.Bool(subscriptionCancelInternalReq.Prorate)
	response, err := sub.Cancel(subscriptionCancelInternalReq.Subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancel", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelSubscriptionInternalResp{}, nil
}

// DoRemoteChannelSubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel
func (s Stripe) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.ChannelSubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancelAtPeriodEnd", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

func (s Stripe) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.ChannelSubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

var usePendingUpdate = true

func (s Stripe) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	// Set the proration date to this moment:
	updateUnixTime := time.Now().Unix()
	if subscriptionRo.ProrationDate > 0 {
		updateUnixTime = subscriptionRo.ProrationDate
	}
	items, err := s.makeSubscriptionUpdateItems(subscriptionRo)
	if err != nil {
		return nil, err
	}
	params := &stripe.InvoiceUpcomingParams{
		Customer:          stripe.String(subscriptionRo.Subscription.ChannelUserId),
		Subscription:      stripe.String(subscriptionRo.Subscription.ChannelSubscriptionId),
		SubscriptionItems: items,
	}
	params.SubscriptionProrationDate = stripe.Int64(updateUnixTime)
	result, err := invoice.Upcoming(params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdateProrationPreview", params, result, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}

	//var invoiceItems []*ro2.SubscriptionInvoiceItemRo
	//for _, line := range result.Lines.Data {
	//	invoiceItems = append(invoiceItems, &ro2.SubscriptionInvoiceItemRo{
	//		Currency:    strings.ToUpper(string(line.Currency)),
	//		Amount:      line.Amount,
	//		Description: line.Description,
	//		Proration:   line.Proration,
	//	})
	//}

	return &ro.ChannelUpdateSubscriptionPreviewInternalResp{
		Data:          utility.FormatToJsonString(result),
		TotalAmount:   result.Total,
		Currency:      strings.ToUpper(string(result.Currency)),
		ProrationDate: updateUnixTime,
		Invoice:       parseStripeInvoice(result, int64(channelEntity.Id)),
	}, nil
}

func (s Stripe) makeSubscriptionUpdateItems(subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) ([]*stripe.SubscriptionItemsParams, error) {

	var items []*stripe.SubscriptionItemsParams

	var stripeSubscriptionItems []*stripe.SubscriptionItem
	if !subscriptionRo.EffectImmediate {
		if len(subscriptionRo.Subscription.ChannelItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.ChannelItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.ChannelSubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//方案 1 遍历并删除，下周期生效，不支持 PendingUpdate
		for _, item := range stripeSubscriptionItems {
			//删除之前全部，新增 Plan 和 Addons 方式
			items = append(items, &stripe.SubscriptionItemsParams{
				ID:      stripe.String(item.ID),
				Deleted: stripe.Bool(true),
			})
		}
		//新增新的项目
		items = append(items, &stripe.SubscriptionItemsParams{
			Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
			Quantity: stripe.Int64(subscriptionRo.Quantity),
			Metadata: map[string]string{
				"BillingPlanType": "Main",
				"BillingPlanId":   strconv.FormatInt(subscriptionRo.PlanChannel.PlanId, 10),
			},
		})
		for _, addon := range subscriptionRo.AddonPlans {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
				Quantity: stripe.Int64(addon.Quantity),
				Metadata: map[string]string{
					"BillingPlanType": "Addon",
					"BillingPlanId":   strconv.FormatInt(addon.AddonPlanChannel.PlanId, 10),
				},
			})
		}
	} else {
		//使用PendingUpdate
		if len(subscriptionRo.Subscription.ChannelItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.ChannelItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.ChannelSubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//方案 2 EffectImmediate=true, 使用PendingUpdate，对于删除的 Plan 和 Addon，修改 Quantity 为 0
		newAddonMap := make(map[string]*ro.SubscriptionPlanAddonRo)
		for _, addon := range subscriptionRo.AddonPlans {
			newAddonMap[addon.AddonPlanChannel.ChannelPlanId] = addon
		}
		//匹配
		for _, item := range stripeSubscriptionItems {
			if strings.Compare(item.Price.ID, subscriptionRo.OldPlanChannel.ChannelPlanId) == 0 {
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Quantity),
				})
			} else if addon, ok := newAddonMap[item.Price.ID]; ok {
				//替换
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
				})
				delete(newAddonMap, item.Price.ID)
			} else {
				//删除之前全部，新增 Plan 和 Addons 方式
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Quantity: stripe.Int64(0),
				})
			}
		}
		//新增剩余的Addons
		for channelPlanId, addon := range newAddonMap {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(channelPlanId),
				Quantity: stripe.Int64(addon.Quantity),
			})
		}
	}

	return items, nil
}

// DoRemoteChannelSubscriptionUpdate 需保证同一个 Price 在 Items 中不能出现两份
func (s Stripe) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	items, err := s.makeSubscriptionUpdateItems(subscriptionRo)
	if err != nil {
		return nil, err
	}

	params := &stripe.SubscriptionParams{
		Items: items,
	}
	if subscriptionRo.EffectImmediate {
		params.ProrationDate = stripe.Int64(subscriptionRo.ProrationDate)
		params.PaymentBehavior = stripe.String("pending_if_incomplete") //pendingIfIncomplete 只有部分字段可以更新 Price Quantity
		params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice))
	} else {
		params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorNone))
	}
	updateSubscription, err := sub.Update(subscriptionRo.Subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", params, updateSubscription, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}

	////todo mark EffectImmediate=false 获取的发票是之前最新的发票
	queryParams := &stripe.InvoiceParams{}
	queryParamsResult, err := invoice.Get(updateSubscription.LatestInvoice.ID, queryParams)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", queryParams, queryParamsResult, err, "GetInvoice", nil, channelEntity)
	g.Log().Infof(ctx, "query invoice:", queryParamsResult)

	return &ro.ChannelUpdateSubscriptionInternalResp{
		ChannelSubscriptionId:     queryParamsResult.ID,
		ChannelSubscriptionStatus: string(updateSubscription.Status),
		ChannelInvoiceId:          queryParamsResult.ID,
		Data:                      utility.FormatToJsonString(updateSubscription),
		LatestInvoiceLink:         queryParamsResult.HostedInvoiceURL,
		Status:                    0, //todo mark
		Paid:                      queryParamsResult.Paid,
	}, nil
}

// DoRemoteChannelSubscriptionDetails 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
func (s Stripe) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.SubscriptionParams{}
	response, err := sub.Get(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionDetails", params, response, err, subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}
	//var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	//if strings.Compare(string(response.Status), "trialing") == 0 ||
	//	strings.Compare(string(response.Status), "active") == 0 {
	//	status = consts.SubStatusActive
	//} else if strings.Compare(string(response.Status), "incomplete") == 0 ||
	//	strings.Compare(string(response.Status), "incomplete_expired") == 0 {
	//	status = consts.SubStatusCreate
	//} else if strings.Compare(string(response.Status), "past_due") == 0 ||
	//	strings.Compare(string(response.Status), "unpaid") == 0 ||
	//	strings.Compare(string(response.Status), "paused") == 0 {
	//	status = consts.SubStatusSuspended
	//} else if strings.Compare(string(response.Status), "canceled") == 0 {
	//	status = consts.SubStatusCancelled
	//}

	return parseStripeSubscriptionDetail(response), nil
}

func parseStripeSubscriptionDetail(subscription *stripe.Subscription) *ro.ChannelDetailSubscriptionInternalResp {
	var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	if strings.Compare(string(subscription.Status), "trialing") == 0 ||
		strings.Compare(string(subscription.Status), "active") == 0 {
		status = consts.SubStatusActive
	} else if strings.Compare(string(subscription.Status), "incomplete") == 0 ||
		strings.Compare(string(subscription.Status), "unpaid") == 0 {
		status = consts.SubStatusCreate
	} else if strings.Compare(string(subscription.Status), "incomplete_expired") == 0 {
		status = consts.SubStatusExpired
	} else if strings.Compare(string(subscription.Status), "past_due") == 0 ||
		strings.Compare(string(subscription.Status), "paused") == 0 {
		status = consts.SubStatusSuspended
	} else if strings.Compare(string(subscription.Status), "canceled") == 0 {
		status = consts.SubStatusCancelled
	}

	return &ro.ChannelDetailSubscriptionInternalResp{
		Status:                 status,
		ChannelSubscriptionId:  subscription.ID,
		ChannelStatus:          string(subscription.Status),
		Data:                   utility.FormatToJsonString(subscription),
		ChannelItemData:        utility.MarshalToJsonString(subscription.Items.Data),
		ChannelLatestInvoiceId: subscription.LatestInvoice.ID,
		CancelAtPeriodEnd:      subscription.CancelAtPeriodEnd,
		CurrentPeriodStart:     subscription.CurrentPeriodStart,
		CurrentPeriodEnd:       subscription.CurrentPeriodEnd,
		TrialEnd:               subscription.TrialEnd,
	}
}

// DoRemoteChannelCheckAndSetupWebhook https://stripe.com/docs/billing/subscriptions/webhooks
func (s Stripe) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
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
				stripe.String("invoice.updated"),
				stripe.String("invoice.paid"),
				stripe.String("invoice.payment_failed"),
				stripe.String("invoice.payment_action_required"),
				stripe.String("payment_intent.created"),
				stripe.String("payment_intent.succeeded"),
			},
			URL: stripe.String(out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
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
			},
			URL: stripe.String(out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
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

// DoRemoteChannelPlanActive 使用 price 代替 plan  https://stripe.com/docs/api/plans
func (s Stripe) DoRemoteChannelPlanActive(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(true) // todo mark 使用这种方式可能不能用
	result, err := price.Update(planChannel.ChannelPlanId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanActive", params, result, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelPlanDeactivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(false) // todo mark 使用这种方式可能不能用
	result, err := price.Update(planChannel.ChannelPlanId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanDeactivate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.ProductParams{
		Active:      stripe.Bool(true),
		Description: stripe.String(plan.ChannelProductDescription), // todo mark 暂时不确定 description 如果为空会怎么样
		Name:        stripe.String(plan.ChannelProductName),
	}
	if len(plan.ImageUrl) > 0 {
		params.Images = stripe.StringSlice([]string{plan.ImageUrl})
	}
	if len(plan.HomeUrl) > 0 {
		params.URL = stripe.String(plan.HomeUrl)
	}
	result, err := product.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelProductCreate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	//Prod 创建好了之后似乎并不是Active 状态 todo mark
	return &ro.ChannelCreateProductInternalResp{
		ChannelProductId:     result.ID,
		ChannelProductStatus: fmt.Sprintf("%v", result.Active),
	}, nil
}

func (s Stripe) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//params := &stripe.PlanParams{
	//	//todo mark
	//	Active:   stripe.Bool(true),
	//	Amount:   stripe.Int64(1200),
	//	Currency: stripe.String(string(stripe.CurrencyUSD)),
	//	Interval: stripe.String(string(stripe.PlanIntervalMonth)),
	//	Product:  &stripe.PlanProductParams{ID: stripe.String("prod_NjpI7DbZx6AlWQ")},
	//}
	//result, err := plan.New(params)
	// 使用 Price 代替 Plan https://stripe.com/docs/api/plans
	params := &stripe.PriceParams{
		Currency:   stripe.String(strings.ToLower(targetPlan.Currency)),
		UnitAmount: stripe.Int64(targetPlan.Amount), //todo mark 小数点可能不用处理
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(targetPlan.IntervalUnit),
			IntervalCount: stripe.Int64(int64(targetPlan.IntervalCount)),
		},
		Product: stripe.String(planChannel.ChannelProductId),

		//ProductData: &stripe.PriceProductDataParams{
		//	ID:   stripe.String(planChannel.ChannelProductId),
		//	Name: stripe.String(targetPlan.PlanName),
		//},//这里是创建的意思
	}
	result, err := price.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanCreateAndActivate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCreatePlanInternalResp{
		ChannelPlanId:     result.ID,
		ChannelPlanStatus: fmt.Sprintf("%v", result.Active),
		Data:              utility.FormatToJsonString(result),
		Status:            consts.PlanChannelStatusActive,
	}, nil
}

func (s Stripe) processInvoiceWebhook(ctx context.Context, eventType string, invoice stripe.Invoice, payChannel *entity.OverseaPayChannel) error {
	details, err := s.DoRemoteChannelInvoiceDetails(ctx, payChannel, invoice.ID)
	if err != nil {
		return err
	}
	err = handler.HandleInvoiceWebhookEvent(ctx, eventType, details)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) processSubscriptionWebhook(ctx context.Context, eventType string, subscription stripe.Subscription) error {
	unibSub := query.GetSubscriptionByChannelSubscriptionId(ctx, subscription.ID)
	if unibSub == nil {
		if unibSubId, ok := subscription.Metadata["SubId"]; ok {
			unibSub = query.GetSubscriptionBySubscriptionId(ctx, unibSubId)
			unibSub.ChannelSubscriptionId = subscription.ID
		}
	}
	if unibSub != nil {
		plan := query.GetPlanById(ctx, unibSub.PlanId)
		planChannel := query.GetPlanChannel(ctx, unibSub.PlanId, unibSub.ChannelId)
		details, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, unibSub)
		if err != nil {
			return err
		}

		err = handler.HandleSubscriptionWebhookEvent(ctx, unibSub, eventType, details)
		if err != nil {
			return err
		}
		return nil
	} else {
		return gerror.New("subscription not found on channelSubId:" + subscription.ID)
	}
}

func (s Stripe) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	endpointSecret := payChannel.WebhookSecret
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(r.GetBody(), signatureHeader, endpointSecret)
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %v\n", payChannel.Channel, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	data, _ := gjson.Marshal(event)
	g.Log().Info(r.Context(), "Receive_Webhook_Channel: ", payChannel.Channel, " hook:", string(data))

	var responseBack = http.StatusOK
	switch event.Type {
	case "customer.subscription.deleted", "customer.subscription.created", "customer.subscription.updated", "customer.subscription.trial_will_end":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription trial will end for %d.", payChannel.Channel, subscription.ID)
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			err := s.processSubscriptionWebhook(r.Context(), string(event.Type), subscription)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "invoice.upcoming", "invoice.created", "invoice.updated", "invoice.paid", "invoice.payment_failed", "invoice.payment_action_required":
		var stripeInvoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &stripeInvoice)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Invoice %s for %d.", payChannel.Channel, string(event.Type), stripeInvoice.ID)
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			err := s.processInvoiceWebhook(r.Context(), string(event.Type), stripeInvoice, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleInvoiceWebhookEvent: %v\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	case "payment_intent.created", "payment_intent.succeeded":
		var stripePayment stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &stripePayment)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			responseBack = http.StatusBadRequest
		} else {
			g.Log().Infof(r.Context(), "Webhook Channel:%s, Payment %s for %d.", payChannel.Channel, string(event.Type), stripePayment.ID)
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			//err := s.processInvoiceWebhook(r.Context(), string(event.Type), stripePayment, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandlePaymentWebhookEvent: %v\n", payChannel.Channel, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	default:
		g.Log().Errorf(r.Context(), "Webhook Channel:%s, Unhandled event type: %s\n", payChannel.Channel, event.Type)
	}
	log.SaveChannelHttpLog("DoRemoteChannelWebhook", event, responseBack, err, string(event.Type), nil, payChannel)
	r.Response.WriteHeader(http.StatusOK)
}

func (s Stripe) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	params, err := r.GetJson()
	g.Log().Printf(r.Context(), "StripeNotifyController redirect params:%s err:%s", params, err.Error())
	if err != nil {
		r.Response.Writeln(err)
		return
	}
	payIdStr := r.Get("payId").String()
	SubIdStr := r.Get("subId").String()
	var response string
	var status bool = false
	var returnUrl string = ""
	if len(payIdStr) > 0 {
		response = "not implement"
	} else if len(SubIdStr) > 0 {
		//订阅回跳
		if r.Get("success").Bool() {
			stripe.Key = payChannel.ChannelSecret
			s.setUnibeeAppInfo()
			unibSub := query.GetSubscriptionBySubscriptionId(r.Context(), SubIdStr)
			if unibSub == nil || len(unibSub.ChannelUserId) == 0 {
				response = "subId invalid or customId empty"
			} else if len(unibSub.ChannelSubscriptionId) > 0 && unibSub.Status == consts.SubStatusActive {
				returnUrl = unibSub.ReturnUrl
				response = "active"
				status = true
			} else {
				//需要去检索
				returnUrl = unibSub.ReturnUrl
				params := &stripe.SubscriptionSearchParams{
					SearchParams: stripe.SearchParams{
						Query: "metadata['SubId']:'" + SubIdStr + "'",
					},
				}
				result := sub.Search(params)
				if result.SubscriptionSearchResult().Data != nil && len(result.SubscriptionSearchResult().Data) == 1 {
					//找到
					if strings.Compare(result.SubscriptionSearchResult().Data[0].Customer.ID, unibSub.ChannelUserId) != 0 {
						response = "customId not match"
					} else {
						detail := parseStripeSubscriptionDetail(result.SubscriptionSearchResult().Data[0])
						err := handler.UpdateSubWithChannelDetailBack(r.Context(), unibSub, detail)
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							response = "subscription active"
							status = true
						}
					}
				} else {
					//找不到
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

func (s Stripe) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
