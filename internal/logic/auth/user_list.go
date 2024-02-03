package auth

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type UserListInternalReq struct {
	MerchantId         int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId             int    `p:"userId" dc:"Filter UserId, Default All" `
	Email              int    `p:"email" dc:"Filter Email, Default All" `
	UserName           int    `p:"userName" dc:"Filter UserName, Default All" `
	SubscriptionName   int    `p:"subscriptionName" dc:"Filter SubscriptionName, Default All" `
	SubscriptionStatus int    `p:"subscriptionStatus" dc:"Filter SubscriptionStatus, Default All" `
	PaymentMethod      int    `p:"paymentMethod" dc:"Filter ChannelDefaultPaymentMethod, Default All" `
	BillingType        int    `p:"billingType" dc:"Filter BillingType, Default All" `
	DeleteInclude      bool   `p:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField          string `p:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType           string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page               int    `p:"page"  dc:"Page,Start 0" `
	Count              int    `p:"count" dc:"Count Of Page" `
}

type UserListInternalRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
}

func UserAccountList(ctx context.Context, req *UserListInternalReq) (res *UserListInternalRes, err error) {
	var mainList []*entity.UserAccount
	if req.Count <= 0 {
		req.Count = 10 //每页数量Default 10
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
	err = dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Id, req.UserId).
		Where(dao.UserAccount.Columns().Email, req.Email).
		Where(dao.UserAccount.Columns().UserName, req.UserName).
		Where(dao.UserAccount.Columns().SubscriptionName, req.SubscriptionName).
		Where(dao.UserAccount.Columns().SubscriptionStatus, req.SubscriptionStatus).
		Where(dao.UserAccount.Columns().PaymentMethod, req.PaymentMethod).
		Where(dao.UserAccount.Columns().BillingType, req.BillingType).
		WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	for _, user := range mainList {
		user.Password = ""
	}
	return &UserListInternalRes{UserAccounts: mainList}, nil
}

func SearchUser(ctx context.Context, searchKey string) (list []*entity.UserAccount, err error) {
	//Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId
	var mainList []*entity.UserAccount
	var isDeletes = []int{0}
	var sortKey = "gmt_create desc"
	_ = dao.UserAccount.Ctx(ctx).
		WhereOr(dao.UserAccount.Columns().Id, searchKey).
		WhereOr(dao.UserAccount.Columns().SubscriptionId, searchKey).
		WhereOr(dao.UserAccount.Columns().VATNumber, searchKey).
		WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(0, 10).
		OmitEmpty().Scan(&mainList)
	if len(mainList) < 10 {
		//继续查 InvoiceId 和 PaymentId
		invoice := query.GetInvoiceByInvoiceId(ctx, searchKey)
		if invoice != nil && invoice.UserId > 0 {
			user := query.GetUserAccountById(ctx, uint64(invoice.UserId))
			if user != nil {
				mainList = append(mainList, user)
			}
		}
		payment := query.GetPaymentByPaymentId(ctx, searchKey)
		if payment != nil && payment.UserId > 0 {
			user := query.GetUserAccountById(ctx, uint64(payment.UserId))
			if user != nil {
				mainList = append(mainList, user)
			}
		}
	}
	if len(mainList) < 10 {
		//模糊查询
		var likeList []*entity.UserAccount
		_ = dao.UserAccount.Ctx(ctx).
			WhereOrLike(dao.UserAccount.Columns().Email, "%"+searchKey+"%").
			WhereOrLike(dao.UserAccount.Columns().UserName, "%"+searchKey+"%").
			WhereOrLike(dao.UserAccount.Columns().CompanyName, "%"+searchKey+"%").
			WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes).
			Order(sortKey).
			Limit(0, 10).
			OmitEmpty().Scan(&likeList)
		if len(likeList) > 0 {
			for _, item := range likeList {
				mainList = append(mainList, item)
			}
		}
	}
	for _, user := range mainList {
		user.Password = ""
	}
	return mainList, err
}
