package service

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strings"
)

type PaymentListInternalReq struct {
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId, Default All " `
	SortField  string `p:"sortField" dc:"排序字段，merchant_id|user_id|gmt_create|status" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}

type PaymentListInternalRes struct {
	Payments []*entity.Payment `json:"payments" description:"payments明细"`
}

func PaymentList(ctx context.Context, req *PaymentListInternalReq) (res *PaymentListInternalRes, err error) {
	var mainList []*entity.Payment
	if req.Count <= 0 {
		req.Count = 10 //每页数量Default 10
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_modify desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("merchant_id|user_id|gmt_create|status", req.SortField), "sortField should one of merchant_id|user_id|gmt_create|status")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	err = dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().MerchantId, req.MerchantId).
		Where(dao.Payment.Columns().UserId, req.UserId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}

	return &PaymentListInternalRes{Payments: mainList}, nil
}
