package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/nacos-group/nacos-sdk-go/util"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func SubscriptionNewTimeline(ctx context.Context, invoice *entity.Invoice) {
	utility.Assert(invoice != nil, "invoice is null ")
	utility.Assert(len(invoice.SubscriptionId) == 0, "not sub invoice")
	utility.Assert(invoice.Status == consts.InvoiceStatusPaid, "invoice not paid")
	utility.Assert(invoice.PeriodStart > 0, "invalid invoice data")
	utility.Assert(invoice.PeriodEnd > 0, "invalid invoice data")
	sub := query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)
	utility.Assert(sub != nil, "sub not found")
	one := query.GetSubscriptionTimeLineByUniqueId(ctx, invoice.InvoiceId)
	if one == nil {
		//finish old one
		var oldOne *entity.SubscriptionTimeline
		_ = dao.SubscriptionTimeline.Ctx(ctx).
			Where(dao.SubscriptionTimeline.Columns().MerchantId, invoice.MerchantId).
			Where(dao.SubscriptionTimeline.Columns().Status, 1).
			Where(dao.SubscriptionTimeline.Columns().SubscriptionId, sub.SubscriptionId).
			OmitEmpty().Scan(&oldOne)
		if oldOne != nil {
			periodEnd := oldOne.PeriodEnd
			if periodEnd > invoice.PeriodStart {
				periodEnd = invoice.PeriodStart
			}
			_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(g.Map{
				dao.SubscriptionTimeline.Columns().Status:    2,
				dao.SubscriptionTimeline.Columns().PeriodEnd: periodEnd,
			}).Where(dao.SubscriptionTimeline.Columns().Id, oldOne.Id).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, `SubscriptionNewTimeline update old one failure %s`, err.Error())
			}
		}

		//create processing one
		one = &entity.SubscriptionTimeline{
			MerchantId:      invoice.MerchantId,
			UserId:          invoice.UserId,
			SubscriptionId:  invoice.SubscriptionId,
			InvoiceId:       invoice.InvoiceId,
			UniqueId:        invoice.InvoiceId,
			UniqueKey:       util.Md5(invoice.InvoiceId),
			Currency:        invoice.Currency,
			PlanId:          sub.PlanId,
			Quantity:        sub.Quantity,
			AddonData:       sub.AddonData,
			GatewayId:       sub.GatewayId,
			PeriodStart:     invoice.PeriodStart,
			PeriodEnd:       invoice.PeriodEnd,
			PeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
			PeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
			CreateTime:      gtime.Now().Timestamp(),
		}

		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			g.Log().Errorf(ctx, `SubscriptionNewTimeline record insert failure %s`, err.Error())
		}
	}
}
