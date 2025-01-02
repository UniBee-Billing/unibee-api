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
	metric2 "unibee/internal/logic/metric"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/timeline"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
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
	} else {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionPendingUpdateCancel.Topic,
			Tag:        redismq2.TopicSubscriptionPendingUpdateCancel.Tag,
			Body:       pendingUpdateId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
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
	} else {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionPendingUpdateSuccess.Topic,
			Tag:        redismq2.TopicSubscriptionPendingUpdateSuccess.Tag,
			Body:       pendingUpdateId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	var billingCycleAnchor = invoice.BillingCycleAnchor
	if billingCycleAnchor <= 0 {
		billingCycleAnchor = sub.BillingCycleAnchor
	}

	periodEnd := utility.MaxInt64(invoice.PeriodEnd, sub.CurrentPeriodEnd)
	if one.EffectImmediate == 1 {
		periodEnd = invoice.PeriodEnd
	}
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, periodEnd, sub.PlanId)

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
		dao.Subscription.Columns().Amount:                 invoice.TotalAmount,
		dao.Subscription.Columns().Currency:               one.UpdateCurrency,
		dao.Subscription.Columns().GatewayId:              invoice.GatewayId,
		dao.Subscription.Columns().LastUpdateTime:         gtime.Now().Timestamp(),
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().PendingUpdateId:        "", //clear PendingUpdateId
		dao.Subscription.Columns().TrialEnd:               invoice.PeriodStart - 1,
		dao.Subscription.Columns().MetaData:               invoice.MetaData,
		dao.Subscription.Columns().TaxPercentage:          invoice.TaxPercentage,
		dao.Subscription.Columns().DiscountCode:           invoice.DiscountCode,
		dao.Subscription.Columns().CancelAtPeriodEnd:      0,
		dao.Subscription.Columns().LatestInvoiceId:        invoice.InvoiceId,
		dao.Subscription.Columns().Data:                   fmt.Sprintf("UpgradedBy-%v", pendingUpdateId),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return false, err
	}

	{
		if vat_gateway.GetDefaultVatGateway(ctx, one.MerchantId) == nil {
			sub_update.UpdateUserTaxPercentageOnly(ctx, one.UserId, one.TaxPercentage)
		}
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
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         sub.UserId,
				SubscriptionId: sub.SubscriptionId,
				Description:    "SubscriptionUpdateSuccess",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	return true, nil
}
