package transaction

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type CreditTransactionListInternalReq struct {
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
	CreateTimeStart  int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd    int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal        bool
}

type CreditTransactionListInternalRes struct {
	CreditTransactions []*detail.CreditTransactionDetail `json:"creditTransactions" dc:"Credit Transaction List"`
	Total              int                               `json:"total" dc:"Total"`
}

func CreditTransactionList(ctx context.Context, req *CreditTransactionListInternalReq) (res *CreditTransactionListInternalRes, err error) {
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
			return &CreditTransactionListInternalRes{CreditTransactions: make([]*detail.CreditTransactionDetail, 0), Total: 0}, nil
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
	//err = query.ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, err
	}
	var resultList []*detail.CreditTransactionDetail
	for _, invoice := range mainList {
		resultList = append(resultList, detail.ConvertToCreditTransactionDetail(ctx, invoice))
	}

	return &CreditTransactionListInternalRes{CreditTransactions: resultList, Total: total}, nil
}
