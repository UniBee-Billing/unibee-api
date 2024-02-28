package service

import (
	"context"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type SubscriptionTimeLineListInternalReq struct {
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId, Default All " `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}

type SubscriptionTimeLineListInternalRes struct {
	SubscriptionTimelines []*entity.SubscriptionTimeline `json:"subscriptionTimeline" dc:"SubscriptionTimeline"`
}

func SubscriptionTimeLineList(ctx context.Context, req *SubscriptionTimeLineListInternalReq) (res *SubscriptionTimeLineListInternalRes, err error) {
	var mainList []*entity.SubscriptionTimeline
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_modify desc"
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
		Where(dao.SubscriptionTimeline.Columns().Status, 2).
		Where(dao.SubscriptionTimeline.Columns().UserId, req.UserId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}

	return &SubscriptionTimeLineListInternalRes{SubscriptionTimelines: mainList}, nil
}
