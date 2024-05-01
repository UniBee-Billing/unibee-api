package service

import (
	"context"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type InvoiceListInternalReq struct {
	MerchantId    uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	FirstName     string `json:"firstName" dc:"FirstName" `
	LastName      string `json:"lastName" dc:"LastName" `
	Currency      string `json:"Currency" dc:"Currency" `
	Status        []int  `json:"status" dc:"Status" `
	AmountStart   int64  `json:"amountStart" dc:"AmountStart" `
	AmountEnd     int64  `json:"amountEnd" dc:"AmountEnd" `
	UserId        uint64 `json:"userId" dc:"Filter UserId Default All" `
	SendEmail     string `json:"sendEmail" dc:"Filter SendEmail Default All" `
	SortField     string `json:"sortField" dc:"Sort Field，invoice_id|gmt_create|period_end|total_amount" `
	SortType      string `json:"sortType" dc:"Sort Type，asc|desc" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Is Delete Include" `
	Page          int    `json:"page"  dc:"Page, Start With 0" `
	Count         int    `json:"count"  dc:"Count Of Page"`
}

type InvoiceListInternalRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoice Detail List"`
}

func InvoiceList(ctx context.Context, req *InvoiceListInternalReq) (res *InvoiceListInternalRes, err error) {
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
		WhereNot(dao.Invoice.Columns().TotalAmount, 0).
		Where(dao.Invoice.Columns().Currency, strings.ToUpper(req.Currency))
	if len(req.SendEmail) > 0 {
		query = query.WhereLike(dao.Invoice.Columns().SendEmail, "%"+req.SendEmail+"%")
	}
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
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		if len(req.FirstName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
		}
		if len(req.LastName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
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
	var resultList []*detail.InvoiceDetail
	for _, invoice := range mainList {
		resultList = append(resultList, detail.ConvertInvoiceToDetail(ctx, invoice))
	}

	return &InvoiceListInternalRes{Invoices: resultList}, nil
}

func SearchInvoice(ctx context.Context, merchantId uint64, searchKey string) (list []*bean.InvoiceSimplify, err error) {
	//Will Search
	var mainList = make([]*bean.InvoiceSimplify, 0)
	var mainMap = make(map[uint64]*bean.InvoiceSimplify)
	var isDeletes = []int{0}
	var sortKey = "gmt_create desc"
	m := dao.Invoice.Ctx(ctx)
	_ = m.Where(m.Builder().WhereOr(dao.Invoice.Columns().InvoiceId, searchKey).
		WhereOr(dao.Invoice.Columns().SendEmail, searchKey)).
		WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Where(dao.Invoice.Columns().MerchantId, merchantId).
		Order(sortKey).
		Limit(0, 10).
		OmitEmpty().Scan(&mainList)
	for _, invoice := range mainList {
		mainMap[invoice.Id] = invoice
	}
	if len(mainList) < 10 {
		//like search
		var likeList []*entity.Invoice
		m = dao.Invoice.Ctx(ctx)
		_ = m.Where(m.Builder().WhereOrLike(dao.Invoice.Columns().InvoiceId, "%"+searchKey+"%").
			WhereOrLike(dao.Invoice.Columns().InvoiceName, "%"+searchKey+"%").
			WhereOrLike(dao.Invoice.Columns().SendEmail, "%"+searchKey+"%")).
			WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
			Where(dao.Invoice.Columns().MerchantId, merchantId).
			Order(sortKey).
			Limit(0, 10).
			OmitEmpty().Scan(&likeList)
		if len(likeList) > 0 {
			for _, invoice := range likeList {
				if mainMap[invoice.Id] == nil {
					mainList = append(mainList, bean.SimplifyInvoice(invoice))
					mainMap[invoice.Id] = bean.SimplifyInvoice(invoice)
				}
			}
		}
	}
	return mainList, err
}
