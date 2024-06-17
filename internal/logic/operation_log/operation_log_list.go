package operation_log

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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
	query := dao.MerchantOperationLog.Ctx(ctx).
		Where(dao.MerchantOperationLog.Columns().MerchantId, req.MerchantId).
		WhereGT(dao.MerchantOperationLog.Columns().MemberId, 0).
		Where(dao.MerchantOperationLog.Columns().IsDelete, 0).
		Limit(req.Page*req.Count, req.Count)
	if len(req.SubscriptionId) > 0 {
		query = query.Where(dao.MerchantOperationLog.Columns().SubscriptionId, req.SubscriptionId)
	}
	if len(req.InvoiceId) > 0 {
		query = query.Where(dao.MerchantOperationLog.Columns().InvoiceId, req.InvoiceId)
	}
	if req.PlanId > 0 {
		query = query.Where(dao.MerchantOperationLog.Columns().PlanId, req.PlanId)
	}
	if len(req.DiscountCode) > 0 {
		query = query.Where(dao.MerchantOperationLog.Columns().DiscountCode, req.DiscountCode)
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
		query = query.WhereIn(dao.MerchantOperationLog.Columns().UserId, userIdList)
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
		query = query.WhereIn(dao.MerchantOperationLog.Columns().MerchantId, memberIdList)
	}

	err := query.ScanAndCount(&mainList, &total, true)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantOperationLogList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertOperationLogToDetail(ctx, one))
	}
	return resultList, total
}
