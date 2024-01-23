package handler

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func CreateOrUpdateSubscriptionTimeline(ctx context.Context, sub *entity.Subscription, channelPaymentId string) error {
	utility.Assert(len(channelPaymentId) > 0, "channelPaymentId is null")
	utility.Assert(sub != nil, "subscription is null ")
	one := query.GetSubscriptionTimeLineByUniqueId(ctx, channelPaymentId)
	if one == nil {
		var periodEnd = gtime.Now().Timestamp()
		if periodEnd > sub.CurrentPeriodEnd {
			//表示已经过了当前周期, 部分通道可能会提前支付并生成发票 todo mark
			periodEnd = sub.CurrentPeriodEnd
		}
		//创建
		one = &entity.SubscriptionTimeline{
			MerchantId:       sub.MerchantId,
			UserId:           sub.UserId,
			SubscriptionId:   sub.SubscriptionId,
			InvoiceId:        "", // todo mark
			UniqueId:         channelPaymentId,
			Currency:         sub.Currency,
			PlanId:           sub.PlanId,
			Quantity:         sub.Quantity,
			AddonData:        sub.AddonData,
			ChannelId:        sub.ChannelId,
			PeriodStart:      gtime.Now().Timestamp(),
			PeriodEnd:        periodEnd,
			PeriodStartTime:  gtime.Now(),
			PeriodEndTime:    gtime.NewFromTimeStamp(periodEnd),
			ChannelPaymentId: channelPaymentId,
		}

		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateSubscriptionTimeline record insert failure %s`, err.Error())
			return err
		}
	}
	return nil
}
