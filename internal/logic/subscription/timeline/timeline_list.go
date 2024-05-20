package timeline

import (
	"context"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionTimeLineListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId     uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField  string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `json:"page"  dc:"Page, Start With 0" `
	Count      int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type SubscriptionTimeLineListInternalRes struct {
	SubscriptionTimelines []*detail.SubscriptionTimeLineDetail
	Total                 int `json:"total" dc:"Total"`
}

func SubscriptionTimeLineList(ctx context.Context, req *SubscriptionTimeLineListInternalReq) (res *SubscriptionTimeLineListInternalRes, err error) {
	var mainList []*entity.SubscriptionTimeline
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("gmt_create|gmt_modify", req.SortField), "sortField should one of gmt_create|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	err = dao.SubscriptionTimeline.Ctx(ctx).
		Where(dao.SubscriptionTimeline.Columns().MerchantId, req.MerchantId).
		//Where(dao.SubscriptionTimeline.Columns().Status, consts.SubTimeLineStatusFinished).
		Where(dao.SubscriptionTimeline.Columns().UserId, req.UserId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, err
	}
	var timelines []*detail.SubscriptionTimeLineDetail
	for _, one := range mainList {
		timelines = append(timelines, &detail.SubscriptionTimeLineDetail{
			MerchantId:      one.MerchantId,
			UserId:          one.UserId,
			SubscriptionId:  one.SubscriptionId,
			PeriodStart:     one.PeriodStart,
			PeriodEnd:       one.PeriodEnd,
			PeriodStartTime: one.PeriodStartTime,
			PeriodEndTime:   one.PeriodEndTime,
			InvoiceId:       one.InvoiceId,
			UniqueId:        one.UniqueId,
			Currency:        one.Currency,
			PlanId:          one.PlanId,
			Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Quantity:        one.Quantity,
			Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			GatewayId:       one.GatewayId,
			CreateTime:      one.CreateTime,
		})
	}

	return &SubscriptionTimeLineListInternalRes{SubscriptionTimelines: timelines, Total: total}, nil
}
