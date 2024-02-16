package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/nacos-group/nacos-sdk-go/util"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func CreateOrUpdateSubscriptionTimeline(ctx context.Context, sub *entity.Subscription, source string) error {
	utility.Assert(sub != nil, "subscription is null ")
	uniqueKey := fmt.Sprintf("%s-%d-%d-%s-%d-%d-%s", sub.SubscriptionId, sub.PlanId, sub.Quantity, sub.AddonData, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, source)
	uniqueId := util.Md5(uniqueKey)
	one := query.GetSubscriptionTimeLineByUniqueId(ctx, uniqueId)
	if one == nil {
		var periodStart = sub.CurrentPeriodStart
		if sub.LastUpdateTime > sub.CurrentPeriodStart {
			periodStart = sub.LastUpdateTime
		}
		var periodEnd = gtime.Now().Timestamp()
		if periodEnd > sub.CurrentPeriodEnd {
			//表示已经过了当前周期, 部分通道可能会提前支付并生成发票 todo mark
			periodEnd = sub.CurrentPeriodEnd
		}
		//创建
		one = &entity.SubscriptionTimeline{
			MerchantId:      sub.MerchantId,
			UserId:          sub.UserId,
			SubscriptionId:  sub.SubscriptionId,
			InvoiceId:       "", // todo mark
			UniqueId:        uniqueId,
			UniqueKey:       uniqueKey,
			Currency:        sub.Currency,
			PlanId:          sub.PlanId,
			Quantity:        sub.Quantity,
			AddonData:       sub.AddonData,
			GatewayId:       sub.GatewayId,
			PeriodStart:     periodStart,
			PeriodEnd:       periodEnd,
			PeriodStartTime: gtime.NewFromTimeStamp(periodStart),
			PeriodEndTime:   gtime.NewFromTimeStamp(periodEnd),
			CreateTime:      gtime.Now().Timestamp(),
		}

		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateSubscriptionTimeline record insert failure %s`, err.Error())
			return err
		}
	}
	return nil
}
