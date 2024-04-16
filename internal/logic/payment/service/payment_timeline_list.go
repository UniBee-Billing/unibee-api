package service

import (
	"context"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/utility"
)

type PaymentTimelineListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId     uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField  string `json:"sortField" dc:"Sort Field，merchant_id|gmt_create|gmt_modify|user_id" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `json:"page"  dc:"Page, Start WIth 0" `
	Count      int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type PaymentTimeLineListInternalRes struct {
	PaymentTimelines []*bean.PaymentTimelineSimplify `json:"paymentTimeline" dc:"paymentTimelines明细"`
}

func PaymentTimeLineList(ctx context.Context, req *PaymentTimelineListInternalReq) (res *PaymentTimeLineListInternalRes, err error) {
	var mainList []*bean.PaymentTimelineSimplify
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_modify desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("merchant_id|gmt_create|gmt_modify|user_id", req.SortField), "sortField should one of merchant_id|gmt_create|gmt_modify|user_id")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	err = dao.PaymentTimeline.Ctx(ctx).
		Where(dao.PaymentTimeline.Columns().MerchantId, req.MerchantId).
		Where(dao.PaymentTimeline.Columns().UserId, req.UserId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}

	return &PaymentTimeLineListInternalRes{PaymentTimelines: mainList}, nil
}
