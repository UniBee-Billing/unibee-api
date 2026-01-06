package discount

import (
	"context"
	"fmt"
	"strconv"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskMultiUserDiscountV2Export struct {
}

func (t TaskMultiUserDiscountV2Export) TaskName() string {
	return "MultiUserDiscountExport"
}

func (t TaskMultiUserDiscountV2Export) Header() interface{} {
	return ExportUserDiscountEntity{}
}

func (t TaskMultiUserDiscountV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	if merchant == nil {
		return mainList, nil
	}
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	var ids []uint64
	if _, ok := payload["exportAll"].(interface{}); ok {
		ids = query.GetAllMerchantDiscountIds(ctx, merchant.Id)
	} else if value, ok := payload["ids"].([]interface{}); ok {
		ids = export.JsonArrayTypeConvertUint64(ctx, value)
	}
	if len(ids) <= 0 {
		return mainList, nil
	}
	req := &userDiscountListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZoneStr = value
				timeZone = zone
			}
		}
		if value, ok := payload["userIds"].([]interface{}); ok {
			req.UserIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
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
	}
	req.SkipTotal = true

	// Collect all user discount data first
	var allUserDiscounts []*userDiscountDetail
	for _, id := range ids {
		req.Id = id
		result := merchantUserDiscountCodeList(ctx, req)
		if result != nil {
			allUserDiscounts = append(allUserDiscounts, result...)
		}
	}

	if len(allUserDiscounts) > 0 {
		// Preload all related data
		preload := preloadMultiUserDiscountData(ctx, allUserDiscounts)

		for _, one := range allUserDiscounts {
			// Populate user and plan data from preload
			if one.User == nil && one.UserId > 0 {
				if user, ok := preload.Users[one.UserId]; ok {
					one.User = user
				}
			}
			if one.Plan == nil && len(one.PlanId) > 0 {
				if planId, err := strconv.ParseUint(one.PlanId, 10, 64); err == nil {
					if plan, ok := preload.Plans[planId]; ok {
						one.Plan = plan
					}
				}
			}
			var firstName = ""
			var lastName = ""
			var email = ""
			if one.User != nil {
				firstName = one.User.FirstName
				lastName = one.User.LastName
				email = one.User.Email
			} else {
				one.User = &bean.UserAccount{}
			}
			if one.Plan == nil {
				one.Plan = &bean.Plan{}
			}
			recurring := "No"
			if one.Recurring == 1 {
				recurring = "Yes"
			}
			statusStr := "Finished"
			if one.Status == 2 {
				statusStr = "Rollback"
			}
			mainList = append(mainList, &ExportMultiUserDiscountEntity{
				Id:             fmt.Sprintf("%v", one.Id),
				UserId:         fmt.Sprintf("%v", one.User.Id),
				ExternalUserId: fmt.Sprintf("%v", one.User.ExternalUserId),
				MerchantName:   merchant.Name,
				FirstName:      firstName,
				LastName:       lastName,
				Email:          email,
				PlanId:         fmt.Sprintf("%v", one.Plan.Id),
				ExternalPlanId: fmt.Sprintf("%v", one.Plan.ExternalPlanId),
				PlanName:       one.Plan.PlanName,
				Code:           one.Code,
				Status:         statusStr,
				SubscriptionId: one.SubscriptionId,
				TransactionId:  one.PaymentId,
				InvoiceId:      one.InvoiceId,
				CreateTime:     gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				ApplyAmount:    utility.ConvertCentToDollarStr(one.ApplyAmount, one.Currency),
				Currency:       one.Currency,
				Recurring:      recurring,
				TimeZone:       timeZoneStr,
			})
		}
	}
	return mainList, nil
}

// PreloadData holds all preloaded data for multi user discounts
type PreloadData struct {
	Users map[uint64]*bean.UserAccount
	Plans map[uint64]*bean.Plan
}

// preloadMultiUserDiscountData preloads all related data for multi user discounts in bulk
func preloadMultiUserDiscountData(ctx context.Context, userDiscounts []*userDiscountDetail) *PreloadData {
	preload := &PreloadData{
		Users: make(map[uint64]*bean.UserAccount),
		Plans: make(map[uint64]*bean.Plan),
	}

	if len(userDiscounts) == 0 {
		return preload
	}

	// Collect unique user IDs and plan IDs
	userIds := make(map[uint64]bool)
	planIds := make(map[uint64]bool)

	for _, userDiscount := range userDiscounts {
		if userDiscount.UserId > 0 {
			userIds[userDiscount.UserId] = true
		}
		if len(userDiscount.PlanId) > 0 {
			if planId, err := strconv.ParseUint(userDiscount.PlanId, 10, 64); err == nil {
				planIds[planId] = true
			}
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

	// Bulk query plans
	if len(planIds) > 0 {
		planIdList := make([]uint64, 0, len(planIds))
		for id := range planIds {
			planIdList = append(planIdList, id)
		}

		var plans []*entity.Plan
		err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, planIdList).Scan(&plans)
		if err == nil {
			for _, plan := range plans {
				preload.Plans[plan.Id] = bean.SimplifyPlan(plan)
			}
		}
	}

	return preload
}

// userDiscountListInternalReq is a local copy of UserDiscountListInternalReq to avoid import cycles
type userDiscountListInternalReq struct {
	MerchantId      uint64   `json:"merchantId" dc:"merchantId"`
	Id              uint64   `json:"id" dc:"discount code id"`
	UserIds         []uint64 `json:"userIds" dc:"user ids"`
	PlanIds         []uint64 `json:"planIds" dc:"plan ids"`
	Email           string   `json:"email" dc:"email"`
	SortField       string   `json:"sortField" dc:"sort field"`
	SortType        string   `json:"sortType" dc:"sort type"`
	Page            int      `json:"page" dc:"page"`
	Count           int      `json:"count" dc:"count"`
	CreateTimeStart int64    `json:"createTimeStart" dc:"create time start"`
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"create time end"`
	SkipTotal       bool
}

// userDiscountDetail is a local copy of UserDiscountDetail to avoid import cycles
type userDiscountDetail struct {
	Id             int64             `json:"id"`
	MerchantId     uint64            `json:"merchantId"`
	UserId         uint64            `json:"userId"`
	PlanId         string            `json:"planId"`
	Code           string            `json:"code"`
	Status         int               `json:"status"`
	SubscriptionId string            `json:"subscriptionId"`
	PaymentId      string            `json:"paymentId"`
	InvoiceId      string            `json:"invoiceId"`
	CreateTime     int64             `json:"createTime"`
	ApplyAmount    int64             `json:"applyAmount"`
	Currency       string            `json:"currency"`
	Recurring      int               `json:"recurring"`
	User           *bean.UserAccount `json:"user"`
	Plan           *bean.Plan        `json:"plan"`
}

// merchantUserDiscountCodeList is a local copy of MerchantUserDiscountCodeList to avoid import cycles
func merchantUserDiscountCodeList(ctx context.Context, req *userDiscountListInternalReq) []*userDiscountDetail {
	var mainList []*entity.MerchantUserDiscountCode
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
		if len(req.SortType) > 0 {
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	query := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId)

	if req.Id > 0 {
		query = query.Where(dao.MerchantUserDiscountCode.Columns().Id, req.Id)
	}
	if len(req.UserIds) > 0 {
		query = query.WhereIn(dao.MerchantUserDiscountCode.Columns().UserId, req.UserIds)
	}
	if len(req.PlanIds) > 0 {
		query = query.WhereIn(dao.MerchantUserDiscountCode.Columns().PlanId, req.PlanIds)
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
		if len(userIdList) > 0 {
			query = query.WhereIn(dao.MerchantUserDiscountCode.Columns().UserId, userIdList)
		}
	}
	if req.CreateTimeStart > 0 {
		query = query.WhereGTE(dao.MerchantUserDiscountCode.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		query = query.WhereLTE(dao.MerchantUserDiscountCode.Columns().CreateTime, req.CreateTimeEnd)
	}
	query = query.
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	if req.SkipTotal {
		err := query.Scan(&mainList)
		if err != nil {
			return nil
		}
	} else {
		err := query.ScanAndCount(&mainList, &total, true)
		if err != nil {
			return nil
		}
	}

	var resultList []*userDiscountDetail
	for _, item := range mainList {
		resultList = append(resultList, &userDiscountDetail{
			Id:             item.Id,
			MerchantId:     item.MerchantId,
			UserId:         item.UserId,
			PlanId:         item.PlanId,
			Code:           item.Code,
			Status:         item.Status,
			SubscriptionId: item.SubscriptionId,
			PaymentId:      item.PaymentId,
			InvoiceId:      item.InvoiceId,
			CreateTime:     item.CreateTime,
			ApplyAmount:    item.ApplyAmount,
			Currency:       item.Currency,
			Recurring:      item.Recurring,
			User:           nil, // Will be populated by preload
			Plan:           nil, // Will be populated by preload
		})
	}

	return resultList
}
