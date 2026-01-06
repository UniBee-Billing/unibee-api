package credit

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskCreditTransactionV2Export struct {
}

func (t TaskCreditTransactionV2Export) TaskName() string {
	return "CreditTransactionExport"
}

func (t TaskCreditTransactionV2Export) Header() interface{} {
	return ExportCreditTransactionEntity{}
}

func (t TaskCreditTransactionV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &creditTransactionListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	timeZone := 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(float64); ok {
			timeZone = int(value)
			if timeZone > 0 {
				timeZoneStr = fmt.Sprintf("UTC+%d", timeZone)
			} else if timeZone < 0 {
				timeZoneStr = fmt.Sprintf("UTC%d", timeZone)
			}
		}
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = uint64(value)
		}
		if value, ok := payload["accountType"].(float64); ok {
			req.AccountType = int(value)
		}
		if value, ok := payload["transactionTypes"].([]interface{}); ok {
			req.TransactionTypes = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
	}
	req.SkipTotal = true
	result, _ := creditTransactionList(ctx, req)
	if result != nil && result.CreditTransactions != nil {
		for _, one := range result.CreditTransactions {
			if one.User == nil {
				one.User = &bean.UserAccount{}
			}

			by := ""
			if one.AdminMember != nil {
				by = one.AdminMember.Email
			}

			mainList = append(mainList, &ExportCreditTransactionEntity{
				Id:              fmt.Sprintf("%v", one.Id),
				ChangedAmount:   utility.ConvertCreditAmountToDollarStr(one.DeltaAmount, one.Currency, one.AccountType),
				Email:           one.User.Email,
				TransactionType: consts.CreditTransactionTypeToEnum(one.TransactionType).ExportDescription(one.DeltaAmount),
				Currency:        one.Currency,
				InvoiceId:       one.InvoiceId,
				By:              by,
				CreateTime:      gtime.NewFromTimeStamp(one.CreateTime + int64(timeZone*3600)),
				Name:            one.Name,
				TimeZone:        timeZoneStr,
			})
		}
	}
	return mainList, nil
}

// PreloadData holds all preloaded data for credit transactions
type PreloadData struct {
	Users        map[uint64]*bean.UserAccount
	AdminMembers map[uint64]*bean.MerchantMember
}

// preloadCreditTransactionData preloads all related data for credit transactions in bulk
func preloadCreditTransactionData(ctx context.Context, transactions []*entity.CreditTransaction) *PreloadData {
	preload := &PreloadData{
		Users:        make(map[uint64]*bean.UserAccount),
		AdminMembers: make(map[uint64]*bean.MerchantMember),
	}

	if len(transactions) == 0 {
		return preload
	}

	// Collect unique user IDs and admin member IDs
	userIds := make(map[uint64]bool)
	adminMemberIds := make(map[uint64]bool)

	for _, transaction := range transactions {
		if transaction.UserId > 0 {
			userIds[transaction.UserId] = true
		}
		if transaction.AdminMemberId > 0 {
			adminMemberIds[transaction.AdminMemberId] = true
		}
	}

	// Bulk query users
	if len(userIds) > 0 {
		userIdList := make([]uint64, 0, len(userIds))
		for id := range userIds {
			userIdList = append(userIdList, id)
		}

		var users []*entity.UserAccount
		err := dao.UserAccount.Ctx(ctx).WhereIn(dao.UserAccount.Columns().Id, userIdList).Scan(&users)
		if err == nil {
			for _, user := range users {
				preload.Users[user.Id] = bean.SimplifyUserAccount(user)
			}
		}
	}

	// Bulk query admin members
	if len(adminMemberIds) > 0 {
		adminMemberIdList := make([]uint64, 0, len(adminMemberIds))
		for id := range adminMemberIds {
			adminMemberIdList = append(adminMemberIdList, id)
		}

		var members []*entity.MerchantMember
		err := dao.MerchantMember.Ctx(ctx).WhereIn(dao.MerchantMember.Columns().Id, adminMemberIdList).Scan(&members)
		if err == nil {
			for _, member := range members {
				preload.AdminMembers[member.Id] = bean.SimplifyMerchantMember(member)
			}
		}
	}

	return preload
}

// creditTransactionListInternalReq is a local copy of CreditTransactionListInternalReq to avoid import cycles
type creditTransactionListInternalReq struct {
	MerchantId       uint64 `json:"merchantId"  description:"merchantId"`
	UserId           uint64 `json:"userId"  description:"filter id of user"`
	AccountType      int    `json:"accountType"  description:"filter type of account"`
	Currency         string `json:"currency"  description:"filter currency of account"`
	Email            string `json:"email"  description:"filter email of user"`
	SortField        string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType         string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	TransactionTypes []int  `json:"transactionTypes" dc:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out" `
	Page             int    `json:"page"  dc:"Page, Start 0" `
	Count            int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart  int64  `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd    int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal        bool
}

// creditTransactionListInternalRes is a local copy of CreditTransactionListInternalRes to avoid import cycles
type creditTransactionListInternalRes struct {
	CreditTransactions []*creditTransactionDetail `json:"creditTransactions" dc:"Credit Transaction List"`
	Total              int                        `json:"total" dc:"Total"`
}

// creditTransactionDetail is a local copy of CreditTransactionDetail to avoid import cycles
type creditTransactionDetail struct {
	Id                  int64                `json:"id"                 description:"Id"`
	User                *bean.UserAccount    `json:"user"`
	CreditAccount       *bean.CreditAccount  `json:"creditAccount"`
	Currency            string               `json:"currency"           description:"currency"`
	TransactionId       string               `json:"transactionId"      description:"unique id for timeline"`
	TransactionType     int                  `json:"transactionType"    description:"transaction type"`
	CreditAmountAfter   int64                `json:"creditAmountAfter"  description:"the credit amount after transaction,cent"`
	CreditAmountBefore  int64                `json:"creditAmountBefore" description:"the credit amount before transaction,cent"`
	DeltaAmount         int64                `json:"deltaAmount"        description:"delta amount,cent"`
	DeltaCurrencyAmount int64                `json:"deltaCurrencyAmount"     description:"delta currency amount, in cent"`
	ExchangeRate        int64                `json:"exchangeRate"          description:"ExchangeRate for transaction"`
	BizId               string               `json:"bizId"              description:"business id"`
	Name                string               `json:"name"               description:"recharge transaction title"`
	Description         string               `json:"description"        description:"recharge transaction description"`
	CreateTime          int64                `json:"createTime"         description:"create utc time"`
	MerchantId          uint64               `json:"merchantId"         description:"merchant id"`
	InvoiceId           string               `json:"invoiceId"         description:"invoice_id"`
	AccountType         int                  `json:"accountType"       description:"type of credit account"`
	AdminMember         *bean.MerchantMember `json:"adminMember"       description:"admin member"`
	By                  string               `json:"by"  dc:"" `
}

// creditTransactionList is a local copy of CreditTransactionList to avoid import cycles
func creditTransactionList(ctx context.Context, req *creditTransactionListInternalReq) (res *creditTransactionListInternalRes, err error) {
	var mainList []*entity.CreditTransaction
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("id|gmt_create|gmt_modify|amount", req.SortField), "sortField should one of id|gmt_create|gmt_modify|amount")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	query := dao.CreditTransaction.Ctx(ctx).
		Where(dao.CreditTransaction.Columns().MerchantId, req.MerchantId)

	if req.UserId > 0 {
		query = query.Where(dao.CreditTransaction.Columns().UserId, req.UserId)
	}
	if len(req.Email) > 0 {
		var userIdList = make([]uint64, 0)
		var list []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		userQuery = userQuery.WhereLike(dao.UserAccount.Columns().Email, "%"+req.Email+"%")
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
		for _, user := range list {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) == 0 {
			return &creditTransactionListInternalRes{CreditTransactions: make([]*creditTransactionDetail, 0), Total: 0}, nil
		}
		query = query.WhereIn(dao.CreditTransaction.Columns().UserId, userIdList)
	}
	if req.AccountType > 0 {
		query = query.Where(dao.CreditTransaction.Columns().AccountType, req.AccountType)
	}
	if len(req.Currency) > 0 {
		query = query.Where(dao.CreditTransaction.Columns().Currency, req.Currency)
	}
	if req.CreateTimeStart > 0 {
		query = query.WhereGTE(dao.CreditTransaction.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		query = query.WhereLTE(dao.CreditTransaction.Columns().CreateTime, req.CreateTimeEnd)
	}
	if len(req.TransactionTypes) > 0 {
		query = query.WhereIn(dao.CreditTransaction.Columns().TransactionType, req.TransactionTypes)
	}
	query = query.
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

	// Preload all related data
	preload := preloadCreditTransactionData(ctx, mainList)

	var resultList []*creditTransactionDetail
	for _, transaction := range mainList {
		// Convert to detail using preloaded data
		detail := convertToCreditTransactionDetailWithPreload(ctx, transaction, preload)
		resultList = append(resultList, detail)
	}

	return &creditTransactionListInternalRes{CreditTransactions: resultList, Total: total}, nil
}

// convertToCreditTransactionDetailWithPreload converts CreditTransaction to CreditTransactionDetail using preloaded data
func convertToCreditTransactionDetailWithPreload(ctx context.Context, one *entity.CreditTransaction, preload *PreloadData) *creditTransactionDetail {
	if one == nil {
		return nil
	}

	// Get user from preloaded data
	var user *bean.UserAccount
	if preload != nil {
		if u, ok := preload.Users[one.UserId]; ok {
			user = u
		}
	}
	if user == nil {
		user = bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId))
	}

	// Get admin member from preloaded data
	var adminMember *bean.MerchantMember
	if preload != nil {
		if m, ok := preload.AdminMembers[one.AdminMemberId]; ok {
			adminMember = m
		}
	}
	if adminMember == nil {
		adminMember = bean.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, one.AdminMemberId))
	}

	// Calculate by field
	by := "-"
	if one.AdminMemberId > 0 && adminMember != nil {
		by = adminMember.Email
	} else if one.UserId > 0 && user != nil {
		by = user.Email
	}

	return &creditTransactionDetail{
		Id:                  one.Id,
		User:                user,
		CreditAccount:       nil, // Not used in export
		Currency:            one.Currency,
		TransactionId:       one.TransactionId,
		TransactionType:     one.TransactionType,
		CreditAmountAfter:   one.CreditAmountAfter,
		CreditAmountBefore:  one.CreditAmountBefore,
		DeltaAmount:         one.DeltaAmount,
		DeltaCurrencyAmount: 0, // Not used in export
		ExchangeRate:        0, // Not used in export
		BizId:               one.BizId,
		Name:                one.Name,
		Description:         one.Description,
		CreateTime:          one.CreateTime,
		MerchantId:          one.MerchantId,
		InvoiceId:           one.InvoiceId,
		AccountType:         one.AccountType,
		AdminMember:         adminMember,
		By:                  by,
	}
}
