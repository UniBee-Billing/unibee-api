package operation_log

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type OperationLogListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	MemberFirstName string `json:"memberFirstName" dc:"Filter Member's FirstName Default All" `
	MemberLastName  string `json:"memberLastName" dc:"Filter Member's LastName, Default All" `
	MemberEmail     string `json:"memberEmail" dc:"Filter Member's Email, Default All" `
	//MemberId        uint64 `json:"memberId" dc:"Filter MemberId Default All" `
	//UserId          uint64 `json:"userId" dc:"Filter UserId Default All" `
	FirstName      string `json:"firstName" dc:"FirstName" `
	LastName       string `json:"lastName" dc:"LastName" `
	Email          string `json:"email" dc:"Email" `
	SubscriptionId string `json:"subscriptionId"     dc:"subscription_id"` // subscription_id
	InvoiceId      string `json:"invoiceId"          dc:"invoice id"`      // invoice id
	PlanId         uint64 `json:"planId"             dc:"plan id"`         // plan id
	DiscountCode   string `json:"discountCode"       dc:"discount_code"`   // discount_code
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count"  dc:"Count Of Page"`
}

func MerchantOperationLogList(ctx context.Context, req *OperationLogListInternalReq) ([]*detail.MerchantOperationLogDetail, int) {
	utility.Assert(req.MerchantId > 0, "Invalid MerchantId")
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var resultList = make([]*detail.MerchantOperationLogDetail, 0)
	var mainList = make([]*entity.MerchantOperationLog, 0)
	q := dao.MerchantOperationLog.Ctx(ctx).
		Where(dao.MerchantOperationLog.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantOperationLog.Columns().IsDelete, 0).
		Limit(req.Page*req.Count, req.Count).
		OrderDesc(dao.MerchantOperationLog.Columns().CreateTime)
	if len(req.SubscriptionId) > 0 {
		q = q.Where(dao.MerchantOperationLog.Columns().SubscriptionId, req.SubscriptionId)
	}
	if len(req.InvoiceId) > 0 {
		q = q.Where(dao.MerchantOperationLog.Columns().InvoiceId, req.InvoiceId)
	}
	if req.PlanId > 0 {
		q = q.Where(dao.MerchantOperationLog.Columns().PlanId, req.PlanId)
	}
	if len(req.DiscountCode) > 0 {
		q = q.Where(dao.MerchantOperationLog.Columns().DiscountCode, req.DiscountCode)
	}
	if len(req.FirstName) > 0 || len(req.LastName) > 0 || len(req.Email) > 0 {
		var userIdList = make([]uint64, 0)
		var list []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		if len(req.FirstName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
		}
		if len(req.LastName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
		}
		if len(req.Email) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().Email, "%"+req.Email+"%")
		}
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
		for _, user := range list {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) == 0 {
			return make([]*detail.MerchantOperationLogDetail, 0), 0
		}
		q = q.WhereIn(dao.MerchantOperationLog.Columns().UserId, userIdList)
	}

	if len(req.MemberLastName) > 0 || len(req.MemberFirstName) > 0 || len(req.MemberEmail) > 0 {
		var memberIdList = make([]uint64, 0)
		var list []*entity.MerchantMember
		memberQuery := dao.MerchantMember.Ctx(ctx).Where(dao.MerchantMember.Columns().MerchantId, req.MerchantId)
		if len(req.MemberFirstName) > 0 {
			memberQuery = memberQuery.WhereLike(dao.MerchantMember.Columns().FirstName, "%"+req.MemberFirstName+"%")
		}
		if len(req.MemberLastName) > 0 {
			memberQuery = memberQuery.WhereLike(dao.MerchantMember.Columns().LastName, "%"+req.MemberLastName+"%")
		}
		if len(req.MemberEmail) > 0 {
			memberQuery = memberQuery.WhereLike(dao.MerchantMember.Columns().Email, "%"+req.MemberEmail+"%")
		}
		_ = memberQuery.Where(dao.MerchantMember.Columns().IsDeleted, 0).Scan(&list)
		for _, one := range list {
			memberIdList = append(memberIdList, one.Id)
		}
		if len(memberIdList) == 0 {
			return make([]*detail.MerchantOperationLogDetail, 0), 0
		}
		q = q.WhereIn(dao.MerchantOperationLog.Columns().MemberId, memberIdList)
	}

	err := q.ScanAndCount(&mainList, &total, true)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantOperationLogList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, convertOperationLogToDetail(ctx, one))
	}
	return resultList, total
}

func convertOperationLogToDetail(ctx context.Context, one *entity.MerchantOperationLog) *detail.MerchantOperationLogDetail {
	if one == nil {
		return nil
	}
	var optAccount = ""
	if one.MemberId > 0 {
		member := query.GetMerchantMemberById(ctx, one.MemberId)
		if member != nil {
			optAccount = fmt.Sprintf("Member (%s)", member.Email)
		}
	} else if one.OptAccountType == 1 && len(one.OptAccountId) > 0 {
		id, err := strconv.ParseInt(one.OptAccountId, 10, 64)
		if err != nil {
			return nil
		}
		user := query.GetUserAccountById(ctx, uint64(id))
		if user != nil {
			optAccount = fmt.Sprintf("%s %s (%s)", user.FirstName, user.LastName, user.Email)
		}
	} else if one.OptAccountType == 2 {
		one.OptAccountId = utility.HideStar(one.OptAccountId)
		optAccount = fmt.Sprintf("OpenApi (%s)", one.OptAccountId)
	} else {
		optAccount = one.OptAccount
	}

	return &detail.MerchantOperationLogDetail{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		MemberId:       one.MemberId,
		OptAccount:     optAccount,
		OptAccountId:   one.OptAccountId,
		OptAccountType: one.OptAccountType,
		OptTarget:      one.OptTarget,
		OptContent:     one.OptContent,
		CreateTime:     one.CreateTime,
		SubscriptionId: one.SubscriptionId,
		UserId:         one.UserId,
		InvoiceId:      one.InvoiceId,
		PlanId:         one.PlanId,
		DiscountCode:   one.DiscountCode,
		Member:         detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, one.MemberId)),
		//User:           bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		//Subscription:   bean.SimplifySubscription(query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)),
	}
}
