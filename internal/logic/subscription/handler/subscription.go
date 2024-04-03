package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/payment/method"
	subscription2 "unibee/internal/logic/subscription"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

func ChangeSubscriptionGateway(ctx context.Context, subscriptionId string, gatewayId uint64, paymentMethodId string) error {
	utility.Assert(gatewayId > 0, "gatewayId is nil")
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionFirstPaymentSuccess sub not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway.MerchantId == sub.MerchantId, "merchant not match:"+strconv.FormatUint(gatewayId, 10))
	if gateway.GatewayType != consts.GatewayTypeCrypto {
		utility.Assert(len(paymentMethodId) > 0, "paymentMethodId invalid")
		paymentMethod := method.QueryPaymentMethod(ctx, sub.MerchantId, sub.UserId, gatewayId, paymentMethodId)
		// todo mark user attach check
		utility.Assert(paymentMethod != nil, "paymentMethodId not found")
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().GatewayId:                   gatewayId,
		dao.Subscription.Columns().GatewayDefaultPaymentMethod: paymentMethodId,
		dao.Subscription.Columns().GmtModify:                   gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func HandleSubscriptionFirstPaymentSuccess(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	utility.Assert(payment != nil, "HandleSubscriptionFirstPaymentSuccess payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "HandleSubscriptionFirstPaymentSuccess payment subId is nil")
	sub = query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionFirstPaymentSuccess sub not found")
	if sub.Status == consts.SubStatusActive {
		return nil
	}
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "HandleSubscriptionFirstPaymentSuccess invoice not found payment:"+payment.PaymentId)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, invoice.PeriodEnd, sub.PlanId)
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPaidTime:          payment.PaidTime,
		dao.Subscription.Columns().TrialEnd:               invoice.PeriodStart - 1,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	SubscriptionNewTimeline(ctx, invoice)
	return nil
}

func HandleSubscriptionNextBillingCyclePaymentSuccess(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	utility.Assert(payment != nil, "UpdateSubscriptionBillingCycleWithPayment payment is nil")
	utility.Assert(payment.Status == consts.PaymentSuccess, "payment not success")
	utility.Assert(len(payment.SubscriptionId) > 0, "UpdateSubscriptionBillingCycleWithPayment payment subId is nil")

	sub = query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "UpdateSubscriptionBillingCycleWithPayment sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "UpdateSubscriptionBillingCycleWithPayment invoice not found payment:"+payment.PaymentId)
	if sub.CurrentPeriodEnd >= invoice.PeriodEnd {
		// sub cycle never go back time
		return nil
	}
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(invoice.PeriodEnd, sub.TrialEnd), sub.PlanId)
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().TrialEnd:               invoice.PeriodStart - 1,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	SubscriptionNewTimeline(ctx, invoice)
	return nil
}

func HandleSubscriptionIncomplete(ctx context.Context, subscriptionId string, nowTimeStamp int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < nowTimeStamp, "subscription not incomplete base on time now")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:    consts.SubStatusIncomplete,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionIncomplete.Topic,
		Tag:   redismq2.TopicSubscriptionIncomplete.Tag,
		Body:  sub.SubscriptionId,
	})
	return nil
}

func UpdateSubscriptionDefaultPaymentMethod(ctx context.Context, subscriptionId string, paymentMethod string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	if len(paymentMethod) == 0 {
		return nil
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().GatewayDefaultPaymentMethod: paymentMethod,
		dao.Subscription.Columns().GmtModify:                   gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}
