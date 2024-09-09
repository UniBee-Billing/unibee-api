package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/email"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/timeline"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func HandlePendingUpdatePaymentFailure(ctx context.Context, pendingUpdateId string) (bool, error) {
	one := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, pendingUpdateId)
	if one == nil {
		return false, gerror.New("HandlePendingUpdatePaymentSuccess PendingUpdate Not Found:" + one.PendingUpdateId)
	}
	if one.Status == consts.PendingSubStatusFinished {
		return true, nil
	}
	if one.Status == consts.PendingSubStatusCancelled {
		return true, nil
	}
	_, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusCancelled,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).OmitNil().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func HandlePendingUpdatePaymentSuccess(ctx context.Context, sub *entity.Subscription, pendingUpdateId string, invoice *entity.Invoice) (bool, error) {
	one := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, pendingUpdateId)
	utility.Assert(one != nil, "HandlePendingUpdatePaymentSuccess PendingUpdate Not Found:"+pendingUpdateId)
	if one.Status == consts.PendingSubStatusFinished {
		return true, nil
	}
	_, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusFinished,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return false, err
	}
	var billingCycleAnchor = invoice.BillingCycleAnchor
	if billingCycleAnchor <= 0 {
		billingCycleAnchor = sub.BillingCycleAnchor
	}

	periodEnd := utility.MaxInt64(invoice.PeriodEnd, sub.CurrentPeriodEnd)

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(invoice.PeriodEnd, sub.TrialEnd), sub.PlanId)
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().BillingCycleAnchor:     billingCycleAnchor,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       periodEnd,
		dao.Subscription.Columns().CurrentPeriodPaid:      1,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(periodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().PlanId:                 one.UpdatePlanId,
		dao.Subscription.Columns().Quantity:               one.UpdateQuantity,
		dao.Subscription.Columns().AddonData:              one.UpdateAddonData,
		dao.Subscription.Columns().Amount:                 one.UpdateAmount,
		dao.Subscription.Columns().Currency:               one.UpdateCurrency,
		dao.Subscription.Columns().GatewayId:              one.GatewayId,
		dao.Subscription.Columns().LastUpdateTime:         gtime.Now().Timestamp(),
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().PendingUpdateId:        "", //clear PendingUpdateId
		dao.Subscription.Columns().TrialEnd:               invoice.PeriodStart - 1,
		dao.Subscription.Columns().MetaData:               one.MetaData,
		dao.Subscription.Columns().TaxPercentage:          one.TaxPercentage,
		dao.Subscription.Columns().DiscountCode:           one.DiscountCode,
		dao.Subscription.Columns().Data:                   fmt.Sprintf("UpgradedBy-%v", pendingUpdateId),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return false, err
	}

	user := query.GetUserAccountById(ctx, sub.UserId)
	merchant := query.GetMerchantById(ctx, sub.MerchantId)

	timeline.SubscriptionNewTimeline(ctx, invoice)
	err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, user.Language, email.TemplateSubscriptionUpdate, "", &email.TemplateVariable{
		UserName:              user.FirstName + " " + user.LastName,
		MerchantProductName:   query.GetPlanById(ctx, one.UpdatePlanId).PlanName,
		MerchantCustomerEmail: merchant.Email,
		MerchantName:          query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
		PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
	})
	if err != nil {
		g.Log().Errorf(ctx, "SendTemplateEmail HandlePendingUpdatePaymentSuccess:%s", err.Error())
	}
	if utility.TryLock(ctx, fmt.Sprintf("HandlePendingUpdatePaymentSuccess_%s", invoice.InvoiceId), 60) {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionUpdate.Topic,
			Tag:        redismq2.TopicSubscriptionUpdate.Tag,
			Body:       one.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		if sub.Status != consts.SubStatusIncomplete && sub.Status != consts.SubStatusActive {
			_, _ = redismq.Send(&redismq.Message{
				Topic:      redismq2.TopicSubscriptionActive.Topic,
				Tag:        redismq2.TopicSubscriptionActive.Tag,
				Body:       one.SubscriptionId,
				CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
			})
		}
	}
	return true, nil
}
