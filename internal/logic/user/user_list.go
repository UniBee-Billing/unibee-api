package user

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type UserListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId          int64  `json:"userId" dc:"Filter UserId, Default All" `
	Email           string `json:"email" dc:"Search Email" `
	FirstName       string `json:"firstName" dc:"Search FirstName" `
	LastName        string `json:"lastName" dc:"Search LastName" `
	SubscriptionId  string `json:"subscriptionId" dc:"Search Filter SubscriptionId" `
	SubStatus       []int  `json:"subStatus" dc:"Filter, Default All，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing | 9-Failed" `
	Status          []int  `json:"status" dc:"Status, 0-Active｜2-Frozen" `
	PlanIds         []int  `json:"planIds" dc:"PlanIds, Search Filter PlanIds" `
	DeleteInclude   bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField       string `json:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

type UserListInternalRes struct {
	UserAccounts []*detail.UserAccountDetail `json:"userAccounts" description:"UserAccounts" `
	Total        int                         `json:"total" dc:"Total"`
}

func UserList(ctx context.Context, req *UserListInternalReq) (res *UserListInternalRes, err error) {
	var mainList []*entity.UserAccount
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
		utility.Assert(strings.Contains("user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type", req.SortField), "sortField should one of user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().MerchantId, req.MerchantId).
		WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes)
	if req.UserId > 0 {
		q = q.Where(dao.UserAccount.Columns().Id, req.UserId)
	}
	if len(req.SubscriptionId) > 0 {
		q = q.Where(dao.UserAccount.Columns().SubscriptionId, req.SubscriptionId)
	}
	if len(req.Email) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().Email, "%"+req.Email+"%")
	}
	if len(req.FirstName) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
	}
	if len(req.LastName) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
	}
	if req.SubStatus != nil && len(req.SubStatus) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().SubscriptionStatus, req.SubStatus)
	}
	if req.Status != nil && len(req.Status) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().Status, req.Status)
	}
	if req.PlanIds != nil && len(req.PlanIds) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().PlanId, req.PlanIds)
	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeEnd)
	}
	q = q.Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	if req.SkipTotal {
		err = q.Scan(&mainList)
	} else {
		err = q.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return nil, err
	}
	var resultList = make([]*detail.UserAccountDetail, 0)
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertUserAccountToDetail(ctx, one))
	}
	return &UserListInternalRes{UserAccounts: resultList, Total: total}, nil
}

func SearchUser(ctx context.Context, merchantId uint64, searchKey string) (list []*detail.UserAccountDetail, err error) {
	//Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId
	var mainList = make([]*entity.UserAccount, 0)
	var mainMap = make(map[uint64]*entity.UserAccount)
	var isDeletes = []int{0}
	var sortKey = "gmt_create desc"
	m := dao.UserAccount.Ctx(ctx)
	_ = m.
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Where(
			m.Builder().WhereOr(dao.UserAccount.Columns().Id, searchKey).
				WhereOr(dao.UserAccount.Columns().SubscriptionId, searchKey).
				WhereOr(dao.UserAccount.Columns().VATNumber, searchKey)).
		WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(0, 10).
		OmitEmpty().Scan(&mainList)
	for _, user := range mainList {
		mainMap[user.Id] = user
	}
	if len(mainList) < 10 {
		//keep go on InvoiceId and PaymentId
		invoice := query.GetInvoiceByInvoiceId(ctx, searchKey)
		if invoice != nil && invoice.UserId > 0 && invoice.MerchantId == merchantId {
			user := query.GetUserAccountById(ctx, invoice.UserId)
			if user != nil && mainMap[user.Id] == nil {
				mainList = append(mainList, user)
				mainMap[user.Id] = user
			}
		}
		payment := query.GetPaymentByPaymentId(ctx, searchKey)
		if payment != nil && payment.UserId > 0 && payment.MerchantId == merchantId {
			user := query.GetUserAccountById(ctx, payment.UserId)
			if user != nil && mainMap[user.Id] == nil {
				mainList = append(mainList, user)
				mainMap[user.Id] = user
			}
		}
	}
	if len(mainList) < 10 {
		//like search
		var likeList []*entity.UserAccount
		m := dao.UserAccount.Ctx(ctx)
		_ = m.
			Where(dao.UserAccount.Columns().MerchantId, merchantId).
			Where(m.Builder().WhereOrLike(dao.UserAccount.Columns().Email, "%"+searchKey+"%").
				WhereOrLike(dao.UserAccount.Columns().FirstName, "%"+searchKey+"%").
				WhereOrLike(dao.UserAccount.Columns().LastName, "%"+searchKey+"%").
				WhereOrLike(dao.UserAccount.Columns().CompanyName, "%"+searchKey+"%")).
			WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes).
			Order(sortKey).
			Limit(0, 10).
			OmitEmpty().Scan(&likeList)
		if len(likeList) > 0 {
			for _, user := range likeList {
				if mainMap[user.Id] == nil {
					mainList = append(mainList, user)
					mainMap[user.Id] = user
				}
			}
		}
	}
	var resultList = make([]*detail.UserAccountDetail, 0)
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertUserAccountToDetail(ctx, one))
	}
	return resultList, err
}
