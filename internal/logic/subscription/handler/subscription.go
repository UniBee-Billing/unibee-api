package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	email2 "unibee/internal/logic/email"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/logic/payment/method"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/timeline"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func ChangeSubscriptionGateway(ctx context.Context, subscriptionId string, gatewayId uint64, paymentMethodId string) (*entity.Subscription, error) {
	utility.Assert(gatewayId > 0, "gatewayId is nil")
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionFirstInvoicePaid sub not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway.MerchantId == sub.MerchantId, "merchant not match:"+strconv.FormatUint(gatewayId, 10))
	var newPaymentMethodId = ""
	if gateway.GatewayType == consts.GatewayTypeCard && len(paymentMethodId) > 0 {
		paymentMethod := method.QueryPaymentMethod(ctx, sub.MerchantId, sub.UserId, gatewayId, paymentMethodId)
		utility.Assert(paymentMethod != nil, "card not found")
		newPaymentMethodId = paymentMethodId
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().GatewayId:                   gatewayId,
		dao.Subscription.Columns().GatewayDefaultPaymentMethod: newPaymentMethodId,
		dao.Subscription.Columns().GmtModify:                   gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "UpdateUserDefaultGatewayPaymentMethod subscriptionId:%d gatewayId:%d, paymentMethodId:%s error:%s", subscriptionId, gatewayId, paymentMethodId, err.Error())
		return nil, err
	} else {
		g.Log().Errorf(ctx, "UpdateUserDefaultGatewayPaymentMethod subscriptionId:%d gatewayId:%d, paymentMethodId:%s success", subscriptionId, gatewayId, paymentMethodId)
		sub_update.UpdateUserDefaultGatewayPaymentMethod(ctx, sub.UserId, gatewayId, paymentMethodId)
	}
	return sub, nil
}

func HandleSubscriptionFirstInvoicePaid(ctx context.Context, sub *entity.Subscription, invoice *entity.Invoice) error {
	utility.Assert(invoice != nil, "HandleSubscriptionFirstInvoicePaid invoice is nil")
	sub = query.GetSubscriptionBySubscriptionId(ctx, sub.SubscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionFirstInvoicePaid sub not found")
	if sub.Status == consts.SubStatusActive {
		return nil
	}
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, invoice.PeriodEnd, sub.PlanId)
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodPaid:      1,
		dao.Subscription.Columns().BillingCycleAnchor:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPaidTime:          gtime.Now().Timestamp(),
		dao.Subscription.Columns().TrialEnd:               invoice.TrialEnd,
		dao.Subscription.Columns().LastUpdateTime:         gtime.Now().Timestamp(),
		dao.Subscription.Columns().CancelAtPeriodEnd:      0,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "HandleSubscriptionFirstInvoicePaid update sub error:%s", err.Error())
		return err
	}
	timeline.SubscriptionFirstPaidTimeline(ctx, invoice)
	if invoice.TrialEnd > 0 && invoice.TrialEnd > invoice.PeriodStart {
		// trial start
		oneUser := query.GetUserAccountById(ctx, sub.UserId)
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
		if oneUser != nil && plan != nil && merchant != nil {
			err = email2.SendTemplateEmail(ctx, sub.MerchantId, oneUser.Email, oneUser.TimeZone, oneUser.Language, email2.TemplateSubscriptionTrialStart, "", &email2.TemplateVariable{
				UserName:              oneUser.FirstName + " " + oneUser.LastName,
				MerchantProductName:   plan.PlanName,
				MerchantCustomerEmail: merchant.Email,
				MerchantName:          query.GetMerchantCountryConfigName(ctx, sub.MerchantId, oneUser.CountryCode),
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail TemplateSubscriptionTrialStart:%s", err.Error())
			}
		}
	}
	if utility.TryLock(ctx, fmt.Sprintf("HandleSubscriptionFirstInvoicePaid_%s", invoice.InvoiceId), 60) {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionActive.Topic,
			Tag:        redismq2.TopicSubscriptionActive.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionUpdate.Topic,
			Tag:        redismq2.TopicSubscriptionUpdate.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         sub.UserId,
				SubscriptionId: sub.SubscriptionId,
				Description:    "SubscriptionFirstActivate",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	return nil
}

func HandleSubscriptionNextBillingCyclePaymentSuccess(ctx context.Context, sub *entity.Subscription, paymentInvoice *entity.Invoice) error {
	utility.Assert(sub != nil, "sub is nil")
	utility.Assert(len(paymentInvoice.SubscriptionId) > 0, "UpdateSubscriptionBillingCycleWithPayment payment subId is nil")

	sub = query.GetSubscriptionBySubscriptionId(ctx, sub.SubscriptionId)
	utility.Assert(sub != nil, "UpdateSubscriptionBillingCycleWithPayment sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, paymentInvoice.InvoiceId)
	utility.Assert(invoice != nil, "UpdateSubscriptionBillingCycleWithPayment invoice not found payment:"+paymentInvoice.PaymentId)
	utility.Assert(invoice.Status == consts.InvoiceStatusPaid || invoice.Status == consts.InvoiceStatusReversed, fmt.Sprintf("invoice not success:%v", invoice.Status))
	if sub.CurrentPeriodEnd > invoice.PeriodEnd && sub.Status == consts.SubStatusActive {
		// sub cycle never go back time
		return nil
	}
	//var recurringDiscountCode *string
	//if len(invoice.DiscountCode) > 0 {
	//	discount := query.GetDiscountByCode(ctx, invoice.MerchantId, invoice.DiscountCode)
	//	if discount.BillingType == consts.DiscountBillingTypeRecurring {
	//		recurringDiscountCode = &invoice.DiscountCode
	//	}
	//}
	var billingCycleAnchor = invoice.BillingCycleAnchor
	if billingCycleAnchor <= 0 {
		billingCycleAnchor = sub.BillingCycleAnchor
	}
	periodEnd := utility.MaxInt64(invoice.PeriodEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, periodEnd, sub.PlanId)
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().BillingCycleAnchor:     billingCycleAnchor,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       periodEnd,
		dao.Subscription.Columns().Amount:                 invoice.TotalAmount,
		dao.Subscription.Columns().CurrentPeriodPaid:      1,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(periodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().TrialEnd:               invoice.PeriodStart - 1,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().TaxPercentage:          invoice.TaxPercentage,
		dao.Subscription.Columns().DiscountCode:           invoice.DiscountCode,
		dao.Subscription.Columns().LastUpdateTime:         gtime.Now().Timestamp(),
		dao.Subscription.Columns().Data:                   fmt.Sprintf("AutoChargeBy-%v", invoice.InvoiceId),
		dao.Subscription.Columns().CancelAtPeriodEnd:      0,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	{
		if vat_gateway.GetDefaultVatGateway(ctx, invoice.MerchantId) == nil {
			sub_update.UpdateUserTaxPercentageOnly(ctx, invoice.UserId, invoice.TaxPercentage)
		}
	}
	timeline.SubscriptionNewTimeline(ctx, invoice)
	if utility.TryLock(ctx, fmt.Sprintf("HandleSubscriptionNextBillingCyclePaymentSuccess_%s", invoice.InvoiceId), 60) {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionUpdate.Topic,
			Tag:        redismq2.TopicSubscriptionUpdate.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		if sub.Status != consts.SubStatusIncomplete && sub.Status != consts.SubStatusActive {
			_, _ = redismq.Send(&redismq.Message{
				Topic:      redismq2.TopicSubscriptionActive.Topic,
				Tag:        redismq2.TopicSubscriptionActive.Tag,
				Body:       sub.SubscriptionId,
				CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
			})
		}
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         sub.UserId,
				SubscriptionId: sub.SubscriptionId,
				Description:    "SubscriptionPaymentSuccess",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	// need cancel payment、 invoice and send invoice email
	return nil
}

func HandleSubscriptionIncomplete(ctx context.Context, subscriptionId string, nowTimeStamp int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < nowTimeStamp, "subscription not incomplete base on time now")
	err := MakeSubscriptionIncomplete(ctx, subscriptionId)
	if err != nil {
		return err
	}
	return nil
}

func MakeSubscriptionIncomplete(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:         consts.SubStatusIncomplete,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicSubscriptionIncomplete.Topic,
		Tag:        redismq2.TopicSubscriptionIncomplete.Tag,
		Body:       subscriptionId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionIncomplete",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	return nil
}

func UpdateSubscriptionDefaultPaymentMethod(ctx context.Context, subscriptionId string, paymentMethod string) error {
	g.Log().Infof(ctx, "UpdateSubscriptionDefaultPaymentMethod subscriptionId:%s paymentMethod:%s", subscriptionId, paymentMethod)
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
