package stripe

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
	"github.com/stripe/stripe-go/v76/balance"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/invoiceitem"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/refund"
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
	handler2 "go-oversea-pay/internal/logic/payment/handler"
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

func parseStripeRefund(item *stripe.Refund) *ro.OutPayRefundRo {
	var channelPaymentId string
	if item.PaymentIntent != nil {
		channelPaymentId = item.PaymentIntent.ID
	}
	var status = consts.REFUND_ING
	if strings.Compare(string(item.Status), "succeeded") == 0 {
		status = consts.REFUND_SUCCESS
	} else if strings.Compare(string(item.Status), "canceled") == 0 || strings.Compare(string(item.Status), "failed") == 0 {
		status = consts.REFUND_FAILED
	}
	return &ro.OutPayRefundRo{
		MerchantId:       "",
		ChannelRefundId:  item.ID,
		ChannelPaymentId: channelPaymentId,
		Status:           status,
		Reason:           string(item.Reason),
		RefundFee:        item.Amount,
		Currency:         strings.ToUpper(string(item.Currency)),
		RefundTime:       gtime.NewFromTimeStamp(item.Created),
	}
}

func parseStripePayment(item *stripe.PaymentIntent) *ro.OutPayRo {
	var channelInvoiceId string
	if item.Invoice != nil {
		channelInvoiceId = item.Invoice.ID
	}
	var channelUserId string
	if item.Customer != nil {
		channelUserId = item.Customer.ID
	}
	var status = consts.TO_BE_PAID
	if strings.Compare(string(item.Status), "succeeded") == 0 {
		status = consts.PAY_SUCCESS
	} else if strings.Compare(string(item.Status), "canceled") == 0 {
		status = consts.PAY_FAILED
	}
	var captureStatus = consts.WAITING_AUTHORIZED
	if strings.Compare(string(item.Status), "requires_capture") == 0 {
		captureStatus = consts.AUTHORIZED
	} else if strings.Compare(string(item.Status), "requires_confirmation") == 0 {
		captureStatus = consts.CAPTURE_REQUEST
	}
	return &ro.OutPayRo{
		ChannelInvoiceId: channelInvoiceId,
		ChannelUserId:    channelUserId,
		ChannelPaymentId: item.ID,
		Status:           status,
		CaptureStatus:    captureStatus,
		PayFee:           item.Amount,
		ReceiptFee:       item.AmountReceived,
		Currency:         strings.ToUpper(string(item.Currency)),
		PayTime:          gtime.NewFromTimeStamp(item.Created),
		CreateTime:       gtime.NewFromTimeStamp(item.Created),
		CancelTime:       gtime.NewFromTimeStamp(item.CanceledAt),
		CancelReason:     string(item.CancellationReason),
	}
}

func parseStripeSubscription(subscription *stripe.Subscription) *ro.ChannelDetailSubscriptionInternalResp {
	var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	if strings.Compare(string(subscription.Status), "trialing") == 0 ||
		strings.Compare(string(subscription.Status), "active") == 0 {
		status = consts.SubStatusActive
	} else if strings.Compare(string(subscription.Status), "incomplete") == 0 ||
		strings.Compare(string(subscription.Status), "unpaid") == 0 {
		status = consts.SubStatusCreate
	} else if strings.Compare(string(subscription.Status), "incomplete_expired") == 0 {
		status = consts.SubStatusExpired
	} else if strings.Compare(string(subscription.Status), "pass_due") == 0 {
		status = consts.SubStatusPendingInActive
	} else if strings.Compare(string(subscription.Status), "paused") == 0 {
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
		BillingCycleAnchor:     subscription.BillingCycleAnchor,
		TrialEnd:               subscription.TrialEnd,
	}
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

func (s Stripe) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.OverseaPayChannel, listReq *ro.ChannelPaymentListReq) (res []*ro.OutPayRo, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.PaymentIntentListParams{}
	params.Customer = stripe.String(listReq.ChannelUserId)
	params.Limit = stripe.Int64(200)
	paymentList := paymentintent.List(params)
	log.SaveChannelHttpLog("DoRemoteChannelPaymentList", params, paymentList, err, "", nil, payChannel)
	var list []*ro.OutPayRo
	for _, item := range paymentList.PaymentIntentList().Data {
		list = append(list, parseStripePayment(item))
	}

	return list, nil
}

func (s Stripe) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.RefundListParams{}
	params.PaymentIntent = stripe.String(channelPaymentId)
	params.Limit = stripe.Int64(100)
	refundList := refund.List(params)
	log.SaveChannelHttpLog("DoRemoteChannelRefundList", params, refundList, err, "", nil, payChannel)
	var list []*ro.OutPayRefundRo
	for _, item := range refundList.RefundList().Data {
		list = append(list, parseStripeRefund(item))
	}

	return list, nil
}

func (s Stripe) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res *ro.OutPayRo, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PaymentIntentParams{}
	response, err := paymentintent.Get(channelPaymentId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPaymentDetail", params, response, err, "", nil, payChannel)
	if err != nil {
		return nil, err
	}

	return parseStripePayment(response), nil
}

func (s Stripe) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(payChannel != nil, "支付渠道异常 gateway not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{}
	response, err := refund.Get(channelRefundId, params)
	log.SaveChannelHttpLog("DoRemoteChannelRefundDetail", params, response, err, "", nil, payChannel)
	if err != nil {
		return nil, err
	}
	return parseStripeRefund(response), nil
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

func (s Stripe) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		TrialEndNow:       stripe.Bool(true),
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
	return details, nil
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
			result, err := dao.SubscriptionVatRateChannel.Ctx(ctx).Data(channelVatRate).OmitNil().Insert(channelVatRate)
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
		//SubscriptionProrationBehavior: stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice)),// 设置了就只会输出 Proration 账单
	}
	params.SubscriptionProrationDate = stripe.Int64(updateUnixTime)
	detail, err := invoice.Upcoming(params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdateProrationPreview", params, detail, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}

	// 拆开 invoice Proration into invoice,nextPeriodInvoice
	var currentInvoiceItems []*ro.ChannelDetailInvoiceItem
	var nextInvoiceItems []*ro.ChannelDetailInvoiceItem
	var currentSubAmount int64 = 0
	var currentSubAmountExcludingTax int64 = 0
	var nextSubAmount int64 = 0
	var nextSubAmountExcludingTax int64 = 0
	for _, line := range detail.Lines.Data {
		if line.Proration {
			currentInvoiceItems = append(currentInvoiceItems, &ro.ChannelDetailInvoiceItem{
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
				Description:            line.Description,
				Proration:              line.Proration,
				Quantity:               line.Quantity,
				Currency:               strings.ToUpper(string(line.Currency)),
			})
			currentSubAmount = currentSubAmount + line.Amount
			currentSubAmountExcludingTax = currentSubAmountExcludingTax + line.AmountExcludingTax
		} else {
			nextInvoiceItems = append(nextInvoiceItems, &ro.ChannelDetailInvoiceItem{
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
				Description:            line.Description,
				Proration:              line.Proration,
				Quantity:               line.Quantity,
				Currency:               strings.ToUpper(string(line.Currency)),
			})
			nextSubAmount = nextSubAmount + line.Amount
			nextSubAmountExcludingTax = nextSubAmountExcludingTax + line.AmountExcludingTax
		}
	}

	currentInvoice := &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    currentSubAmount,
		TotalAmountExcludingTax:        currentSubAmountExcludingTax,
		TaxAmount:                      currentSubAmount - currentSubAmountExcludingTax,
		SubscriptionAmount:             currentSubAmount,
		SubscriptionAmountExcludingTax: currentSubAmountExcludingTax,
		Lines:                          currentInvoiceItems,
		ChannelSubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		ChannelId:                      int64(channelEntity.Id),
		ChannelUserId:                  detail.Customer.ID,
	}

	nextPeriodInvoice := &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    nextSubAmount,
		TotalAmountExcludingTax:        nextSubAmountExcludingTax,
		TaxAmount:                      nextSubAmount - nextSubAmountExcludingTax,
		SubscriptionAmount:             nextSubAmount,
		SubscriptionAmountExcludingTax: nextSubAmountExcludingTax,
		Lines:                          nextInvoiceItems,
		ChannelSubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		ChannelId:                      int64(channelEntity.Id),
		ChannelUserId:                  detail.Customer.ID,
	}

	return &ro.ChannelUpdateSubscriptionPreviewInternalResp{
		Data:          utility.FormatToJsonString(detail),
		TotalAmount:   currentInvoice.TotalAmount,
		Currency:      strings.ToUpper(string(detail.Currency)),
		ProrationDate: updateUnixTime,
		Invoice:       currentInvoice,
		//Invoice: parseStripeInvoice(detail, int64(channelEntity.Id)),
		NextPeriodInvoice: nextPeriodInvoice,
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
		newMap := make(map[string]*ro.SubscriptionPlanAddonRo)
		for _, addon := range subscriptionRo.AddonPlans {
			newMap[addon.AddonPlanChannel.ChannelPlanId] = addon
		}
		//匹配
		var replace = false
		for _, item := range stripeSubscriptionItems {
			if strings.Compare(item.Price.ID, subscriptionRo.OldPlanChannel.ChannelPlanId) == 0 {
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Quantity),
				})
				replace = true
			} else if addon, ok := newMap[item.Price.ID]; ok {
				//替换
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
				})
				delete(newMap, item.Price.ID)
			} else {
				//删除之前全部，新增 Plan 和 Addons 方式
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Quantity: stripe.Int64(0),
				})
			}
		}
		if !replace {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
				Quantity: stripe.Int64(subscriptionRo.Quantity),
			})
		}
		//新增剩余的Addons
		for channelPlanId, addon := range newMap {
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

	if subscriptionRo.EffectImmediate {
		queryParams := &stripe.InvoiceParams{}
		newInvoice, err := invoice.Get(updateSubscription.LatestInvoice.ID, queryParams)
		log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", queryParams, newInvoice, err, "GetInvoice", nil, channelEntity)
		g.Log().Infof(ctx, "query new invoice:", newInvoice)

		return &ro.ChannelUpdateSubscriptionInternalResp{
			Data:            utility.FormatToJsonString(updateSubscription),
			ChannelUpdateId: newInvoice.ID,
			Link:            newInvoice.HostedInvoiceURL,
			Paid:            newInvoice.Paid,
		}, nil
	} else {
		//EffectImmediate=false 不需要支付 获取的发票是之前最新的发票
		return &ro.ChannelUpdateSubscriptionInternalResp{
			Data: utility.FormatToJsonString(updateSubscription),
			Paid: true,
		}, nil
	}
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

	return parseStripeSubscription(response), nil
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

func (s Stripe) processPaymentWebhook(ctx context.Context, eventType string, payment stripe.PaymentIntent, payChannel *entity.OverseaPayChannel) error {
	details, err := s.DoRemoteChannelPaymentDetail(ctx, payChannel, payment.ID)
	if err != nil {
		return err
	}
	details.ChannelId = int64(payChannel.Id)
	utility.Assert(len(details.ChannelUserId) > 0, "invalid channelUserId")
	details.ChannelUser = query.GetUserChannelByChannelUserId(ctx, details.ChannelUserId, details.ChannelId)
	utility.Assert(details.ChannelUser != nil, "channelUser not found")
	if payment.Invoice != nil {
		//可能来自 SubPendingUpdate 流程，需要补充 Invoice 信息获取 ChannelUpdateId
		invoiceDetails, err := s.DoRemoteChannelInvoiceDetails(ctx, payChannel, payment.Invoice.ID)
		if err != nil {
			return err
		}
		details.ChannelInvoiceDetail = invoiceDetails
		details.ChannelInvoiceId = payment.Invoice.ID
		details.ChannelUpdateId = invoiceDetails.ChannelInvoiceId
		oneSub := query.GetSubscriptionByChannelSubscriptionId(ctx, invoiceDetails.ChannelSubscriptionId)
		if oneSub != nil {
			plan := query.GetPlanById(ctx, oneSub.PlanId)
			planChannel := query.GetPlanChannel(ctx, oneSub.PlanId, oneSub.ChannelId)
			subDetails, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, oneSub)
			if err != nil {
				return err
			}
			details.ChannelSubscriptionDetail = subDetails
			details.Subscription = oneSub
		}
	}
	err = handler2.HandlePaymentWebhookEvent(ctx, eventType, details)
	if err != nil {
		return err
	}
	return nil
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
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %s\n", payChannel.Channel, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
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
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			err = s.processSubscriptionWebhook(r.Context(), string(event.Type), subscription)
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
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
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
			// Then define and call a func to handle the successful attachment of a PaymentMethod.
			// handleSubscriptionTrialWillEnd(subscription)
			err = s.processPaymentWebhook(r.Context(), string(event.Type), stripePayment, payChannel)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandlePaymentWebhookEvent: %s\n", payChannel.Channel, err.Error())
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
						detail := parseStripeSubscription(result.SubscriptionSearchResult().Data[0])
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

func (s Stripe) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, one *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, pay.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 gateway not found")
	params := &stripe.RefundParams{PaymentIntent: stripe.String(pay.ChannelPaymentId)}
	params.Reason = stripe.String(one.RefundComment)
	params.Amount = stripe.Int64(one.RefundFee)
	params.Currency = stripe.String(strings.ToLower(one.Currency))
	result, err := refund.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelRefund", params, result, err, "refund", nil, channelEntity)
	utility.Assert(err == nil, fmt.Sprintf("call stripe refund error %s", err))
	utility.Assert(result != nil, "Stripe refund failed, result is nil")
	return &ro.OutPayRefundRo{
		ChannelRefundId: result.ID,
		Status:          consts.REFUND_ING,
	}, nil
}
