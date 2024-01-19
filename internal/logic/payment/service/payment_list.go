package service

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strings"
)

type PaymentListInternalReq struct {
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int    `p:"userId" dc:"UserId 不填查询所有" `
	SortField  string `p:"sortField" dc:"排序字段，payment_id|gmt_create|status|payment_fee" `
	SortType   string `p:"sortType" dc:"排序类型，asc|desc" `
	Page       int    `p:"page"  dc:"分页页码,0开始" `
	Count      int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}

type PaymentListInternalRes struct {
	Payments []*entity.Payment `p:"payments" dc:"payments明细"`
}

func PaymentList(ctx context.Context, req *PaymentListInternalReq) (res *PaymentListInternalRes, err error) {
	var mainList []*entity.Payment
	if req.Count <= 0 {
		req.Count = 10 //每页数量默认 10
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_modify desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("invoice_id|gmt_create|gmt_modify|period_end|total_amount", req.SortField), "sortField should one of invoice_id|gmt_create|period_end|total_amount")
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
