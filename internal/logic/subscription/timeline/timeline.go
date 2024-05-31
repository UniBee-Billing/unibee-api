package timeline

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/nacos-group/nacos-sdk-go/util"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func FinishOldTimelineBySubEnd(ctx context.Context, subscriptionId string, endSubStatus consts.SubscriptionStatusEnum) {
	utility.Assert(len(subscriptionId) > 0, "invalid subscriptionId")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "sub not found")
	var oldOne *entity.SubscriptionTimeline
	_ = dao.SubscriptionTimeline.Ctx(ctx).
		Where(dao.SubscriptionTimeline.Columns().MerchantId, sub.MerchantId).
		WhereIn(dao.SubscriptionTimeline.Columns().Status, []int{consts.SubTimeLineStatusPending, consts.SubTimeLineStatusProcessing}).
		Where(dao.SubscriptionTimeline.Columns().SubscriptionId, sub.SubscriptionId).
		OmitEmpty().Scan(&oldOne)
	if oldOne != nil {
		periodEnd := gtime.Now().Timestamp()
		if !config.GetConfigInstance().IsProd() {
			periodEnd = utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock)
		}
		nextStatus := consts.SubTimeLineStatusFinished
		if oldOne.Status == consts.SubTimeLineStatusPending {
			if endSubStatus == consts.SubStatusExpired {
				nextStatus = consts.SubTimeLineStatusExpired
			} else {
				nextStatus = consts.SubTimeLineStatusCancelled
			}
		}
		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(g.Map{
			dao.SubscriptionTimeline.Columns().Status:    nextStatus,
			dao.SubscriptionTimeline.Columns().PeriodEnd: periodEnd,
		}).Where(dao.SubscriptionTimeline.Columns().Id, oldOne.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, `FinishOldTimelineBySubEnd update old one failure %s`, err.Error())
		}
	}
}

func SubscriptionNewPendingTimeline(ctx context.Context, invoice *entity.Invoice) {
	g.Log().Infof(ctx, "SubscriptionNewPendingTimeline-%s", invoice.InvoiceId)
	utility.Assert(invoice != nil, "invoice is null ")
	utility.Assert(len(invoice.SubscriptionId) > 0, "not sub invoice")
	utility.Assert(invoice.PeriodStart > 0, "invalid invoice data")
	utility.Assert(invoice.PeriodEnd > 0, "invalid invoice data")
	sub := query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)
	utility.Assert(sub != nil, "sub not found")
	one := query.GetSubscriptionTimeLineByUniqueId(ctx, invoice.InvoiceId)
	if one == nil {
		//create pending one
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
			Status:          consts.SubTimeLineStatusPending,
			GatewayId:       sub.GatewayId,
			PeriodStart:     invoice.PeriodStart,
			PeriodEnd:       invoice.PeriodEnd,
			PeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
			PeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
			CreateTime:      gtime.Now().Timestamp(),
		}

		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			g.Log().Errorf(ctx, `SubscriptionNewPendingTimeline record insert failure %s`, err.Error())
		}
	}
}

func SubscriptionFirstPaidTimeline(ctx context.Context, invoice *entity.Invoice) {
	g.Log().Infof(ctx, "SubscriptionFirstPaidTimeline-%s", invoice.InvoiceId)
	utility.Assert(invoice != nil, "invoice is null ")
	utility.Assert(len(invoice.SubscriptionId) > 0, "not sub invoice")
	utility.Assert(invoice.Status == consts.InvoiceStatusPaid, "invoice not paid")
	utility.Assert(invoice.PeriodStart > 0, "invalid invoice data")
	utility.Assert(invoice.PeriodEnd > 0, "invalid invoice data")
	sub := query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)
	utility.Assert(sub != nil, "sub not found")
	one := query.GetSubscriptionTimeLineByUniqueId(ctx, invoice.InvoiceId)
	if one == nil {
		//create pending one
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
			Status:          consts.SubTimeLineStatusProcessing,
			GatewayId:       sub.GatewayId,
			PeriodStart:     invoice.PeriodStart,
			PeriodEnd:       invoice.PeriodEnd,
			PeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
			PeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
			CreateTime:      gtime.Now().Timestamp(),
		}

		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			g.Log().Errorf(ctx, `SubscriptionFirstPaidTimeline record insert failure %s`, err.Error())
		}
	} else {
		_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(g.Map{
			dao.SubscriptionTimeline.Columns().Status:    consts.SubTimeLineStatusProcessing,
			dao.SubscriptionTimeline.Columns().PeriodEnd: invoice.PeriodEnd,
		}).Where(dao.SubscriptionTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, `SubscriptionFirstPaidTimeline update old one failure %s`, err.Error())
		}
	}
}

func SubscriptionNewTimeline(ctx context.Context, invoice *entity.Invoice) {
	g.Log().Infof(ctx, "SubscriptionNewTimeline-%s", invoice.InvoiceId)
	utility.Assert(invoice != nil, "invoice is null ")
	utility.Assert(len(invoice.SubscriptionId) > 0, "not sub invoice")
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
			Where(dao.SubscriptionTimeline.Columns().Status, consts.SubTimeLineStatusProcessing).
			Where(dao.SubscriptionTimeline.Columns().SubscriptionId, sub.SubscriptionId).
			OmitEmpty().Scan(&oldOne)
		if oldOne != nil {
			periodEnd := oldOne.PeriodEnd
			if periodEnd > invoice.PeriodStart {
				periodEnd = invoice.PeriodStart
			}
			_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(g.Map{
				dao.SubscriptionTimeline.Columns().Status:    consts.SubTimeLineStatusFinished,
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
			Status:          consts.SubTimeLineStatusProcessing,
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
