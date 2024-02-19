package service

import (
	"context"
	"strings"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/invoice/invoice_compute"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

type SubscriptionInvoiceListInternalReq struct {
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	FirstName     string `p:"firstName" dc:"FirstName" `
	LastName      string `p:"lastName" dc:"LastName" `
	Currency      string `p:"Currency" dc:"Currency" `
	Status        []int  `p:"status" dc:"Status" `
	AmountStart   int64  `p:"amountStart" dc:"AmountStart" `
	AmountEnd     int64  `p:"amountEnd" dc:"AmountEnd" `
	UserId        int    `p:"userId" dc:"Filter UserId Default All" `
	SendEmail     string `p:"sendEmail" dc:"Filter SendEmail Default All" `
	SortField     string `p:"sortField" dc:"Sort Field，invoice_id|gmt_create|period_end|total_amount" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"Is Delete Include" `
	Page          int    `p:"page"  dc:"Page, Start WIth 0" `
	Count         int    `p:"count"  dc:"Count Of Page"`
}

type SubscriptionInvoiceListInternalRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"Invoice Detail List"`
}

func SubscriptionInvoiceList(ctx context.Context, req *SubscriptionInvoiceListInternalReq) (res *SubscriptionInvoiceListInternalRes, err error) {
	var mainList []*entity.Invoice
	if req.Count <= 0 {
		req.Count = 20
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
	query := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
		Where(dao.Invoice.Columns().Currency, strings.ToUpper(req.Currency)).
		Where(dao.Invoice.Columns().SendEmail, req.SendEmail)
	if req.UserId > 0 {
		query = query.Where(dao.Invoice.Columns().UserId, req.UserId)
	}
	if req.Status != nil && len(req.Status) > 0 {
		query = query.WhereIn(dao.Invoice.Columns().Status, req.Status)
	}
	if req.AmountStart < req.AmountEnd {
		query = query.WhereGTE(dao.Invoice.Columns().TotalAmount, req.AmountStart)
		query = query.WhereLTE(dao.Invoice.Columns().TotalAmount, req.AmountEnd)
	}
	if len(req.FirstName) > 0 || len(req.LastName) > 0 {
		var userIdList []uint64
		var list []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx)
		if len(req.FirstName) > 0 {
			userQuery = userQuery.WhereOrLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
		}
		if len(req.LastName) > 0 {
			userQuery = userQuery.WhereOrLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
		}
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
		for _, user := range list {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) > 0 {
			query = query.WhereIn(dao.Invoice.Columns().UserId, userIdList)
		}
	}
	err = query.WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	var resultList []*ro.InvoiceDetailRo
	for _, invoice := range mainList {
		resultList = append(resultList, invoice_compute.ConvertInvoiceToRo(ctx, invoice))
	}

	return &SubscriptionInvoiceListInternalRes{Invoices: resultList}, nil
}

func SearchInvoice(ctx context.Context, searchKey string) (list []*entity.Invoice, err error) {
	//Will Search
	var mainList []*entity.Invoice
	var isDeletes = []int{0}
	var sortKey = "gmt_create desc"
	_ = dao.Invoice.Ctx(ctx).
		WhereOr(dao.Invoice.Columns().InvoiceId, searchKey).
		WhereOr(dao.Invoice.Columns().SendEmail, searchKey).
		WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(0, 10).
		OmitEmpty().Scan(&mainList)
	if len(mainList) < 10 {
		//模糊查询
		var likeList []*entity.Invoice
		_ = dao.Invoice.Ctx(ctx).
			WhereOrLike(dao.Invoice.Columns().InvoiceId, "%"+searchKey+"%").
			WhereOrLike(dao.Invoice.Columns().InvoiceName, "%"+searchKey+"%").
			WhereOrLike(dao.Invoice.Columns().SendEmail, "%"+searchKey+"%").
			WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
			Order(sortKey).
			Limit(0, 10).
			OmitEmpty().Scan(&likeList)
		if len(likeList) > 0 {
			for _, item := range likeList {
				mainList = append(mainList, item)
			}
		}
	}
	return mainList, err
}
