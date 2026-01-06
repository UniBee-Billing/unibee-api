package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/preload"
	"unibee/internal/logic/subscription/config"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type CreditNoteListInternalReq struct {
	MerchantId      uint64   `json:"merchantId" dc:"MerchantId" v:"required"`
	SearchKey       string   `json:"searchKey" dc:"The search key of invoice" `
	Emails          []string `json:"emails" dc:"The email list of invoice user, split by commas or semicolons" `
	Status          []int    `json:"status" dc:"The status of invoice, 2-processing｜3-paid | 4-failed | 5-cancelled" `
	GatewayIds      []int64  `json:"gatewayIds" dc:"GatewayIds, Search Filter GatewayIds" `
	PlanIds         []int64  `json:"planIds" dc:"PlanIds, Search Filter PlanIds" `
	Currency        string   `json:"currency" dc:"The currency of invoice" `
	SortField       string   `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType        string   `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	Page            int      `json:"page"  dc:"Page, Start 0" `
	Count           int      `json:"count"  dc:"Count" dc:"Count By Page" `
	DeleteInclude   bool     `json:"deleteInclude" dc:"Is Delete Include" `
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd" `
	ReportTimeStart int64    `json:"reportTimeStart" dc:"ReportTimeStart" `
	ReportTimeEnd   int64    `json:"reportTimeEnd" dc:"ReportTimeEnd" `
	SkipTotal       bool
}

type CreditNoteListInternalRes struct {
	CreditNotes []*detail.CreditNoteDetail `json:"creditNotes" dc:"CreditNote Detail Object List"`
	Total       int                        `json:"total" dc:"Total"`
}

func CreditNoteList(ctx context.Context, req *CreditNoteListInternalReq) (res *CreditNoteListInternalRes, err error) {
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
		Where(dao.Invoice.Columns().Currency, strings.ToUpper(req.Currency)).
		WhereLT(dao.Invoice.Columns().TotalAmount, 0) // refund invoice
	if !config.GetMerchantSubscriptionConfig(ctx, req.MerchantId).ShowZeroInvoice {
		query = query.WhereNot(dao.Invoice.Columns().TotalAmount, 0)
	}
	//if len(req.SearchKey) > 0 {
	//	query = query.Where(query.Builder().WhereOrLike(dao.Invoice.Columns().InvoiceName, "%"+req.SearchKey+"%").WhereOrLike(dao.Invoice.Columns().ProductName, "%"+req.SearchKey+"%"))
	//}
	if req.GatewayIds != nil && len(req.GatewayIds) > 0 {
		query = query.WhereIn(dao.Invoice.Columns().GatewayId, req.GatewayIds)
	}
	if req.Status != nil && len(req.Status) > 0 {
		query = query.WhereIn(dao.Invoice.Columns().Status, req.Status)
	}
	if len(req.SearchKey) > 0 || len(req.Emails) > 0 {
		var userIdList = make([]uint64, 0)
		var list []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		if len(req.SearchKey) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().Email, "%"+req.SearchKey+"%")
		}
		if len(req.Emails) > 0 {
			userQuery = userQuery.WhereIn(dao.UserAccount.Columns().Email, req.Emails)
		}
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
		for _, user := range list {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) == 0 {
			return &CreditNoteListInternalRes{CreditNotes: make([]*detail.CreditNoteDetail, 0), Total: 0}, nil
		}
		query = query.WhereIn(dao.Invoice.Columns().UserId, userIdList)
	}
	if len(req.PlanIds) > 0 {
		linePlanQuery := query.Builder()
		for _, planId := range req.PlanIds {
			linePlanQuery = linePlanQuery.WhereOrLike(dao.Invoice.Columns().Lines, "%\"plan\":{\"id\":"+fmt.Sprintf("%d", planId)+"%")
		}
		query = query.Where(linePlanQuery)
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
	var resultList []*detail.CreditNoteDetail
	preload.CreditNoteListPreloadForContext(ctx, mainList)
	for _, invoice := range mainList {
		resultList = append(resultList, detail.ConvertInvoiceToCreditNoteDetail(ctx, invoice))
	}

	return &CreditNoteListInternalRes{CreditNotes: resultList, Total: total}, nil
}
