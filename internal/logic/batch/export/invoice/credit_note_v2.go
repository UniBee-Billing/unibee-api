package invoice

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskCreditNoteV2Export struct {
}

func (t TaskCreditNoteV2Export) TaskName() string {
	return fmt.Sprintf("CreditNoteExport")
}

func (t TaskCreditNoteV2Export) Header() interface{} {
	return ExportCreditNoteEntity{}
}

func (t TaskCreditNoteV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.CreditNoteListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZone = zone
			}
		}
		if value, ok := payload["gatewayIds"].([]interface{}); ok {
			req.GatewayIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["searchKey"].(string); ok {
			req.SearchKey = value
		}
		if value, ok := payload["emails"].(string); ok {
			emails := make([]string, 0)

			// 1. Process directly passed emails parameter
			if len(value) > 0 {
				cleanedEmails := strings.ReplaceAll(value, ";", ",")
				emails = strings.Split(cleanedEmails, ",")
				// Clean each email, remove spaces
				for i, email := range emails {
					emails[i] = strings.TrimSpace(email)
				}
				// Filter empty strings
				var filteredEmails []string
				for _, email := range emails {
					if email != "" {
						filteredEmails = append(filteredEmails, email)
					}
				}
				emails = filteredEmails
			}
			req.Emails = emails
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeStart"].(float64); ok {
			req.ReportTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeEnd"].(float64); ok {
			req.ReportTimeEnd = int64(value) - timeZone
		}
	}
	req.SkipTotal = true
	result, err := service.CreditNoteList(ctx, req)
	if result != nil && result.CreditNotes != nil {
		// Preload all related data
		//preload := preload2.CreditNoteListPreload(ctx, result.CreditNotes)

		for _, one := range result.CreditNotes {
			// Populate data from preload
			//if one.Gateway == nil && one.GatewayId > 0 {
			//	if gateway, ok := preload.Gateways[one.GatewayId]; ok {
			//		one.Gateway = detail.ConvertGatewayDetail(ctx, gateway)
			//	}
			//}
			//if one.UserSnapshot == nil && one.UserId > 0 {
			//	if user, ok := preload.Users[one.UserId]; ok {
			//		one.UserSnapshot = bean.SimplifyUserAccount(user)
			//	}
			//}
			//if one.Payment == nil && len(one.PaymentId) > 0 {
			//	if payment, ok := preload.Payments[one.PaymentId]; ok {
			//		one.Payment = bean.SimplifyPayment(payment)
			//	}
			//}
			//if one.Refund == nil && len(one.RefundId) > 0 {
			//	if refund, ok := preload.Refunds[one.RefundId]; ok {
			//		one.Refund = bean.SimplifyRefund(refund)
			//	}
			//}
			//if one.Subscription == nil && len(one.SubscriptionId) > 0 {
			//	if subscription, ok := preload.Subscriptions[one.SubscriptionId]; ok {
			//		one.Subscription = bean.SimplifySubscription(ctx, subscription)
			//	}
			//}
			//if one.Discount == nil && len(one.DiscountCode) > 0 {
			//	if discount, ok := preload.Discounts[one.DiscountCode]; ok {
			//		one.Discount = bean.SimplifyMerchantDiscountCode(discount)
			//	}
			//}

			var creditNoteGateway = ""
			if one.Gateway != nil {
				creditNoteGateway = one.Gateway.GatewayName
			}
			if one.UserSnapshot == nil {
				one.UserSnapshot = &bean.UserAccount{}
			}
			if one.Refund == nil {
				one.Refund = &bean.Refund{}
			}
			var planIdStr string
			var planNameStr string
			if one.PlanSnapshot != nil && one.PlanSnapshot.Plan != nil {
				planIdStr = fmt.Sprintf("%v", one.PlanSnapshot.Plan.Id)
				planNameStr = one.PlanSnapshot.Plan.PlanName
			}
			mainList = append(mainList, &ExportCreditNoteEntity{
				CreditNoteId:        one.InvoiceId,
				UserId:              fmt.Sprintf("%v", one.UserId),
				Email:               one.UserSnapshot.Email,
				FirstName:           one.UserSnapshot.FirstName,
				LastName:            one.UserSnapshot.LastName,
				CreditNoteName:      one.InvoiceName,
				ProductName:         one.ProductName,
				Currency:            one.Currency,
				TotalAmount:         utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
				TaxAmount:           utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency),
				Status:              consts.InvoiceStatusToEnum(one.Status).Description(),
				Gateway:             creditNoteGateway,
				CreateTime:          gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				FinishTime:          gtime.NewFromTimeStamp(one.FinishTime + timeZone),
				RefundId:            one.RefundId,
				PaymentId:           one.PaymentId,
				PlanId:              planIdStr,
				PlanName:            planNameStr,
				SubscriptionId:      one.SubscriptionId,
				RefundReason:        one.Refund.RefundComment,
				PartialCreditAmount: utility.ConvertCentToDollarStr(one.PartialCreditPaidAmount, one.Currency),
			})
		}
	}
	return mainList, nil
}

//
//// creditNoteListInternalReq is a local copy of CreditNoteListInternalReq to avoid import cycles
//type creditNoteListInternalReq struct {
//	MerchantId      uint64   `json:"merchantId" dc:"merchantId"`
//	GatewayIds      []int64  `json:"gatewayIds" dc:"gateway ids"`
//	SearchKey       string   `json:"searchKey" dc:"search key"`
//	Emails          []string `json:"emails" dc:"emails"`
//	Currency        string   `json:"currency" dc:"currency"`
//	Status          []int    `json:"status" dc:"status"`
//	PlanIds         []int64  `json:"planIds" dc:"plan ids"`
//	SortField       string   `json:"sortField" dc:"sort field"`
//	SortType        string   `json:"sortType" dc:"sort type"`
//	Page            int      `json:"page" dc:"page"`
//	Count           int      `json:"count" dc:"count"`
//	CreateTimeStart int64    `json:"createTimeStart" dc:"create time start"`
//	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"create time end"`
//	ReportTimeStart int64    `json:"reportTimeStart" dc:"report time start"`
//	ReportTimeEnd   int64    `json:"reportTimeEnd" dc:"report time end"`
//	SkipTotal       bool
//}
//
//// creditNoteListInternalRes is a local copy of CreditNoteListInternalRes to avoid import cycles
//type creditNoteListInternalRes struct {
//	CreditNotes []*detail.CreditNoteDetail `json:"creditNotes" dc:"CreditNote Detail Object List"`
//	Total       int                        `json:"total" dc:"Total"`
//}
//
//// creditNoteList is a local copy of CreditNoteList to avoid import cycles
//func creditNoteList(ctx context.Context, req *creditNoteListInternalReq) *creditNoteListInternalRes {
//	var mainList []*entity.Invoice
//	var total = 0
//	if req.Count <= 0 {
//		req.Count = 20
//	}
//	if req.Page < 0 {
//		req.Page = 0
//	}
//
//	var isDeletes = []int{0}
//	utility.Assert(req.MerchantId > 0, "merchantId not found")
//	var sortKey = "gmt_create desc"
//	if len(req.SortField) > 0 {
//		utility.Assert(strings.Contains("invoice_id|gmt_create|gmt_modify|period_end|total_amount", req.SortField), "sortField should one of invoice_id|gmt_create|period_end|total_amount")
//		if len(req.SortType) > 0 {
//			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
//			sortKey = req.SortField + " " + req.SortType
//		} else {
//			sortKey = req.SortField + " desc"
//		}
//	}
//	query := dao.Invoice.Ctx(ctx).
//		Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
//		Where(dao.Invoice.Columns().Currency, strings.ToUpper(req.Currency)).
//		WhereLT(dao.Invoice.Columns().TotalAmount, 0) // refund invoice
//	if !config.GetMerchantSubscriptionConfig(ctx, req.MerchantId).ShowZeroInvoice {
//		query = query.WhereNot(dao.Invoice.Columns().TotalAmount, 0)
//	}
//	//if len(req.SearchKey) > 0 {
//	//	query = query.Where(query.Builder().WhereOrLike(dao.Invoice.Columns().InvoiceName, "%"+req.SearchKey+"%").WhereOrLike(dao.Invoice.Columns().ProductName, "%"+req.SearchKey+"%"))
//	//}
//	if req.GatewayIds != nil && len(req.GatewayIds) > 0 {
//		query = query.WhereIn(dao.Invoice.Columns().GatewayId, req.GatewayIds)
//	}
//	if req.Status != nil && len(req.Status) > 0 {
//		query = query.WhereIn(dao.Invoice.Columns().Status, req.Status)
//	}
//	if len(req.SearchKey) > 0 || len(req.Emails) > 0 {
//		var userIdList = make([]uint64, 0)
//		var list []*entity.UserAccount
//		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
//		if len(req.SearchKey) > 0 {
//			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().Email, "%"+req.SearchKey+"%")
//		}
//		if len(req.Emails) > 0 {
//			userQuery = userQuery.WhereIn(dao.UserAccount.Columns().Email, req.Emails)
//		}
//		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
//		for _, user := range list {
//			userIdList = append(userIdList, user.Id)
//		}
//		if len(userIdList) == 0 {
//			return &creditNoteListInternalRes{CreditNotes: make([]*detail.CreditNoteDetail, 0), Total: 0}
//		}
//		query = query.WhereIn(dao.Invoice.Columns().UserId, userIdList)
//	}
//	if len(req.PlanIds) > 0 {
//		linePlanQuery := query.Builder()
//		for _, planId := range req.PlanIds {
//			linePlanQuery = linePlanQuery.WhereOrLike(dao.Invoice.Columns().Lines, "%\"plan\":{\"id\":"+fmt.Sprintf("%d", planId)+"%")
//		}
//		query = query.Where(linePlanQuery)
//	}
//	if req.CreateTimeStart > 0 {
//		query = query.WhereGTE(dao.Invoice.Columns().CreateTime, req.CreateTimeStart)
//	}
//	if req.CreateTimeEnd > 0 {
//		query = query.WhereLTE(dao.Invoice.Columns().CreateTime, req.CreateTimeEnd)
//	}
//	if req.ReportTimeStart > 0 {
//		query = query.Where(query.Builder().WhereOrGTE(dao.Invoice.Columns().CreateTime, req.ReportTimeStart).
//			WhereOrGTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeStart)))
//	}
//	if req.ReportTimeEnd > 0 {
//		query = query.Where(query.Builder().WhereOrLTE(dao.Invoice.Columns().CreateTime, req.ReportTimeEnd).
//			WhereOrLTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeEnd)))
//	}
//	query = query.WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
//		Order(sortKey).
//		Limit(req.Page*req.Count, req.Count).
//		OmitEmpty()
//	if req.SkipTotal {
//		err := query.Scan(&mainList)
//		if err != nil {
//			return nil
//		}
//	} else {
//		err := query.ScanAndCount(&mainList, &total, true)
//		if err != nil {
//			return nil
//		}
//	}
//	var resultList []*detail.CreditNoteDetail
//	for _, invoice := range mainList {
//		resultList = append(resultList, detail.ConvertInvoiceToCreditNoteDetail(ctx, invoice))
//	}
//
//	return &creditNoteListInternalRes{CreditNotes: resultList, Total: total}
//}
