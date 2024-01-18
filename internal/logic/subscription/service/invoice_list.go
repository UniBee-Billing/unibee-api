package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionInvoiceListInternalReq struct {
	g.Meta        `path:"/subscription_invoice_list" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice列表"`
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId        int    `p:"userId" dc:"UserId 不填查询所有" `
	SendEmail     int    `p:"sendEmail" dc:"SendEmail 不填查询所有" `
	SortField     string `p:"sortField" dc:"排序字段，invoice_id|gmt_create|period_end|total_amount" `
	SortType      string `p:"sortType" dc:"排序类型，asc|desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"是否包含删除，查看已删除发票需要超级管理员权限" `
	Page          int    `p:"page"  dc:"分页页码,0开始" `
	Count         int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}

type SubscriptionInvoiceListInternalRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"invoices明细"`
}

func SubscriptionInvoiceList(ctx context.Context, req *SubscriptionInvoiceListInternalReq) (res *SubscriptionInvoiceListInternalRes, err error) {
	var mainList []*entity.Invoice
	if req.Count <= 0 {
		req.Count = 10 //每页数量默认 10
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
		resultList = append(resultList, ConvertInvoiceToRo(invoice))
	}

	return &SubscriptionInvoiceListInternalRes{Invoices: resultList}, nil
}

func ConvertInvoiceToRo(invoice *entity.Invoice) *ro.InvoiceDetailRo {
	var lines []*ro.ChannelDetailInvoiceItem
	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	return &ro.InvoiceDetailRo{
		Id:                             invoice.Id,
		MerchantId:                     invoice.MerchantId,
		SubscriptionId:                 invoice.SubscriptionId,
		InvoiceId:                      invoice.InvoiceId,
		GmtCreate:                      invoice.GmtCreate,
		TotalAmount:                    invoice.TotalAmount,
		TaxAmount:                      invoice.TaxAmount,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		Currency:                       invoice.Currency,
		Lines:                          lines,
		ChannelId:                      invoice.ChannelId,
		Status:                         invoice.Status,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      invoice.SendEmail,
		SendPdf:                        invoice.SendPdf,
		UserId:                         invoice.UserId,
		Data:                           invoice.Data,
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           invoice.Link,
		ChannelStatus:                  invoice.ChannelStatus,
		ChannelInvoiceId:               invoice.ChannelInvoiceId,
		ChannelInvoicePdf:              invoice.ChannelInvoicePdf,
		TaxPercentage:                  invoice.TaxPercentage,
		SendNote:                       invoice.SendNote,
		SendTerms:                      invoice.SendTerms,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
	}
}
