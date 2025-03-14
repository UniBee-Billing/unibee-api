package service

import (
	"context"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type PaymentItemListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId     uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField  string `json:"sortField" dc:"Sort Field，merchant_id|gmt_create|gmt_modify|user_id" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `json:"page"  dc:"Page, Start With 0" `
	Count      int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type PaymentItemListInternalRes struct {
	PaymentItems []*bean.PaymentItem `json:"paymentItem" dc:"paymentItems"`
	Total        int                 `json:"total" dc:"Total"`
}

func OneTimePaymentItemList(ctx context.Context, req *PaymentItemListInternalReq) (res *PaymentItemListInternalRes, err error) {
	var mainList []*entity.PaymentItem
	var total int
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("merchant_id|gmt_create|gmt_modify|user_id", req.SortField), "sortField should one of merchant_id|gmt_create|gmt_modify|user_id")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	err = dao.PaymentItem.Ctx(ctx).
		Where(dao.PaymentItem.Columns().MerchantId, req.MerchantId).
		Where(dao.PaymentItem.Columns().UserId, req.UserId).
		Where(dao.PaymentItem.Columns().BizType, consts.BizTypeOneTime).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, err
	}

	var resultList = make([]*bean.PaymentItem, 0)
	for _, one := range mainList {
		resultList = append(resultList, bean.SimplifyPaymentItemTimeline(one))
	}

	return &PaymentItemListInternalRes{PaymentItems: resultList, Total: total}, nil
}
