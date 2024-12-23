package credit

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type CreditAccountListInternalReq struct {
	MerchantId      uint64 `json:"merchantId"  description:"merchantId"`
	UserId          uint64 `json:"userId"  description:"filter id of user"`
	Email           string `json:"email"  description:"filter email of user"`
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type CreditAccountListInternalRes struct {
	CreditAccounts []*detail.CreditAccountDetail `json:"creditAccounts" dc:"Credit Account List"`
	Total          int                           `json:"total" dc:"Total"`
}

func CreditAccountList(ctx context.Context, req *CreditAccountListInternalReq) (res *CreditAccountListInternalRes, err error) {
	var mainList []*entity.CreditAccount
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
	q := dao.CreditAccount.Ctx(ctx).
		Where(dao.CreditAccount.Columns().MerchantId, req.MerchantId)
	if req.UserId > 0 {
		q = q.Where(dao.CreditAccount.Columns().UserId, req.UserId)
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
			return &CreditAccountListInternalRes{CreditAccounts: make([]*detail.CreditAccountDetail, 0), Total: 0}, nil
		}
		q = q.WhereIn(dao.Invoice.Columns().UserId, userIdList)

	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.Invoice.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.Invoice.Columns().CreateTime, req.CreateTimeEnd)
	}
	q = q.
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	err = q.ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, err
	}
	var resultList []*detail.CreditAccountDetail
	for _, invoice := range mainList {
		resultList = append(resultList, detail.ConvertToCreditAccountDetail(ctx, invoice))
	}

	return &CreditAccountListInternalRes{CreditAccounts: resultList, Total: total}, nil
}
