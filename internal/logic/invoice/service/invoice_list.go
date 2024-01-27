package service

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/invoice/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionInvoiceListInternalReq struct {
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId        int    `p:"userId" dc:"FilterUserId Default All" `
	SendEmail     int    `p:"sendEmail" dc:"Filter SendEmail Default All" `
	SortField     string `p:"sortField" dc:"Sort Field，invoice_id|gmt_create|period_end|total_amount" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"Is Delete Include" `
	Page          int    `p:"page"  dc:"Page, Start WIth 0" `
	Count         int    `p:"count"  dc:"Count Of Page"`
}

type SubscriptionInvoiceListInternalRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"invoices明细"`
}

func SubscriptionInvoiceList(ctx context.Context, req *SubscriptionInvoiceListInternalReq) (res *SubscriptionInvoiceListInternalRes, err error) {
	var mainList []*entity.Invoice
	if req.Count <= 0 {
		req.Count = 10 //每页数量Default 10
	}
	if req.Page < 0 {
		req.Page = 0
	}

	var isDeletes = []int{0}
	if req.DeleteInclude {
		isDeletes = append(isDeletes, 1)
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
	err = dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
		Where(dao.Invoice.Columns().UserId, req.UserId).
		Where(dao.Invoice.Columns().SendEmail, req.SendEmail).
		WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	var resultList []*ro.InvoiceDetailRo
	for _, invoice := range mainList {
		resultList = append(resultList, util.ConvertInvoiceToRo(invoice))
	}

	return &SubscriptionInvoiceListInternalRes{Invoices: resultList}, nil
}
