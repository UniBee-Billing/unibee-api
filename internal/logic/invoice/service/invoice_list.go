package service

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type InvoiceListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	FirstName       string `json:"firstName" dc:"FirstName" `
	LastName        string `json:"lastName" dc:"LastName" `
	Currency        string `json:"Currency" dc:"Currency" `
	Status          []int  `json:"status" dc:"Status" `
	AmountStart     *int64 `json:"amountStart" dc:"AmountStart" `
	AmountEnd       *int64 `json:"amountEnd" dc:"AmountEnd" `
	UserId          uint64 `json:"userId" dc:"Filter UserId Default All" `
	SendEmail       string `json:"sendEmail" dc:"Filter SendEmail Default All" `
	SortField       string `json:"sortField" dc:"Sort Field，invoice_id|gmt_create|period_end|total_amount" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc" `
	DeleteInclude   bool   `json:"deleteInclude" dc:"Is Delete Include" `
	Type            *int   `json:"type"  dc:"invoice Type, 0-payment, 1-refund" `
	Page            int    `json:"page"  dc:"Page, Start With 0" `
	Count           int    `json:"count"  dc:"Count Of Page"`
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	ReportTimeStart int64  `json:"reportTimeStart" dc:"ReportTimeStart" `
	ReportTimeEnd   int64  `json:"reportTimeEnd" dc:"ReportTimeEnd" `
	SkipTotal       bool
}

type InvoiceListInternalRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoice Detail List"`
	Total    int                     `json:"total" dc:"Total"`
}

func InvoiceList(ctx context.Context, req *InvoiceListInternalReq) (res *InvoiceListInternalRes, err error) {
	var mainList []*entity.Invoice
	var total = 0
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
	var sortKey = "gmt_create desc"
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
	if req.AmountStart != nil && req.AmountEnd != nil {
		utility.Assert(*req.AmountStart <= *req.AmountEnd, "amountStart should lower then amountEnd")
	}
	if req.AmountStart != nil {
		query = query.WhereGTE(dao.Invoice.Columns().TotalAmount, &req.AmountStart)
	}
	if req.AmountEnd != nil {
		query = query.WhereLTE(dao.Invoice.Columns().TotalAmount, &req.AmountEnd)
	}
	if req.Type != nil {
		utility.Assert(*req.Type == 0 || *req.Type == 1, "type should be 0 or 1")
		if *req.Type == 0 {
			query = query.WhereNull(dao.Invoice.Columns().RefundId)
		} else if *req.Type == 1 {
			query = query.WhereNotNull(dao.Invoice.Columns().RefundId)
		}
	}
	if len(req.FirstName) > 0 || len(req.LastName) > 0 {
		var userIdList = make([]uint64, 0)
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
		if len(userIdList) == 0 {
			return &InvoiceListInternalRes{Invoices: make([]*detail.InvoiceDetail, 0), Total: 0}, nil
		}
		query = query.WhereIn(dao.Invoice.Columns().UserId, userIdList)

	}
	if req.CreateTimeStart > 0 {
		query = query.WhereGTE(dao.Invoice.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		query = query.WhereLTE(dao.Invoice.Columns().CreateTime, req.CreateTimeEnd)
	}
	if req.ReportTimeStart > 0 {
		query = query.Where(query.Builder().WhereOrGTE(dao.Invoice.Columns().CreateTime, req.ReportTimeStart).
			WhereOrGTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeStart)))
	}
	if req.ReportTimeEnd > 0 {
		query = query.Where(query.Builder().WhereOrLTE(dao.Invoice.Columns().CreateTime, req.ReportTimeEnd).
			WhereOrLTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeEnd)))
	}
	query = query.WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	if req.SkipTotal {
		err = query.Scan(&mainList)
	} else {
		err = query.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return nil, err
	}
	var resultList []*detail.InvoiceDetail
	for _, invoice := range mainList {
		resultList = append(resultList, detail.ConvertInvoiceToDetail(ctx, invoice))
	}

	return &InvoiceListInternalRes{Invoices: resultList, Total: total}, nil
}

func SearchInvoice(ctx context.Context, merchantId uint64, searchKey string) (list []*bean.Invoice, err error) {
	//Will Search
	var mainList = make([]*bean.Invoice, 0)
	var mainMap = make(map[uint64]*bean.Invoice)
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
