package service

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type PaymentTimelineListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId     uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField  string `json:"sortField" dc:"Sort Field，merchant_id|gmt_create|gmt_modify|user_id" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `json:"page"  dc:"Page, Start With 0" `
	Count      int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type PaymentTimeLineListInternalRes struct {
	PaymentTimelines []*detail.PaymentTimelineDetail `json:"paymentTimeline" dc:"paymentTimelines"`
}

func PaymentTimeLineList(ctx context.Context, req *PaymentTimelineListInternalReq) (res *PaymentTimeLineListInternalRes, err error) {
	var mainList []*entity.PaymentTimeline
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

	var resultList = make([]*detail.PaymentTimelineDetail, 0)
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertPaymentTimeline(ctx, one))
	}

	return &PaymentTimeLineListInternalRes{PaymentTimelines: resultList}, nil
}
