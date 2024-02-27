package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	subscription2 "unibee/internal/logic/subscription"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

func HandleSubscriptionCreatePaymentSuccess(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	utility.Assert(payment != nil, "HandleSubscriptionCreatePaymentSuccess payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "HandleSubscriptionCreatePaymentSuccess payment subId is nil")
	sub = query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionCreatePaymentSuccess sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "HandleSubscriptionCreatePaymentSuccess invoice not found payment:"+payment.PaymentId)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, invoice.PeriodEnd, uint64(sub.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPaidTime:          payment.PaidTime,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func FinishNextBillingCycleForSubscription(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	// billing-cycle
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("cycle-paymentId-%s", payment.PaymentId))
	if err != nil {
		g.Log().Errorf(ctx, "FinishNextBillingCycleForSubscription error:%s", err.Error())
	}
	err = UpdateSubscriptionBillingCycleWithPayment(ctx, payment)
	return err
}

func ChangeTrialEnd(ctx context.Context, newTrialEnd int64, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "sub cancelled or sub expired")

	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, uint64(sub.PlanId))
	newStatus := sub.Status
	if newTrialEnd > gtime.Now().Timestamp() {
		//automatic change sub status to active
		newStatus = consts.SubStatusActive
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:             newStatus,
		dao.Subscription.Columns().TrialEnd:           newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
		dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
		dao.Subscription.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionIncomplete(ctx context.Context, subscriptionId string, nowTimeStamp int64) error {
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

func UpdateSubscriptionBillingCycleWithPayment(ctx context.Context, payment *entity.Payment) error {
	utility.Assert(payment != nil, "UpdateSubscriptionBillingCycleWithPayment payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "UpdateSubscriptionBillingCycleWithPayment payment subId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "UpdateSubscriptionBillingCycleWithPayment sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "UpdateSubscriptionBillingCycleWithPayment invoice not found payment:"+payment.PaymentId)
	var FirstPaidTime int64 = 0
	if sub.FirstPaidTime == 0 && payment.Status == consts.PaymentSuccess {
		FirstPaidTime = payment.PaidTime
	}
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(invoice.PeriodEnd, sub.TrialEnd), uint64(sub.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPaidTime:          FirstPaidTime,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}
