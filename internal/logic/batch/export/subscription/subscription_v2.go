package subscription

import (
	"context"
	"fmt"
	"strings"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/gateway/util"
	preload2 "unibee/internal/logic/preload"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	dao "unibee/internal/dao/default"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskSubscriptionV2Export struct {
}

func (t TaskSubscriptionV2Export) TaskName() string {
	return "SubscriptionExport"
}

func (t TaskSubscriptionV2Export) Header() interface{} {
	return ExportSubscriptionEntity{}
}

func (t TaskSubscriptionV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.SubscriptionListInternalReq{
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
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = int64(value)
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
		if value, ok := payload["productIds"].([]interface{}); ok {
			req.ProductIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["amountStart"].(float64); ok {
			req.AmountStart = unibee.Int64(int64(value))
		}
		if value, ok := payload["amountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
		}
	}
	req.SkipTotal = true

	// Get subscriptions directly from database to avoid N+1 queries in SubscriptionList
	subscriptions := subscriptionList(ctx, req)
	if len(subscriptions) > 0 {
		// Preload all related data to avoid N+1 queries
		preload := preload2.SubscriptionListPreload(ctx, subscriptions)

		for _, one := range subscriptions {
			var subGateway = ""
			var stripeUserId = ""
			var stripePaymentMethod = ""
			var paypalVaultId = ""

			// Use preloaded data instead of individual queries
			if preload.Gateways[one.GatewayId] != nil {
				gateway := preload.Gateways[one.GatewayId]
				subGateway = gateway.GatewayName
				if gateway.GatewayType == consts.GatewayTypeCard {
					gatewayUser := util.GetGatewayUser(ctx, one.UserId, gateway.Id)
					if gatewayUser != nil {
						stripeUserId = gatewayUser.GatewayUserId
						stripePaymentMethod = one.GatewayDefaultPaymentMethod
					}
				} else if gateway.GatewayType == consts.GatewayTypePaypal {
					paypalVaultId = one.GatewayDefaultPaymentMethod
				}
			}

			var canAtPeriodEnd = "No"
			if one.CancelAtPeriodEnd == 1 {
				canAtPeriodEnd = "Yes"
			}

			var firstName = ""
			var lastName = ""
			var email = ""
			user := &entity.UserAccount{}
			if preload.Users[one.UserId] != nil {
				user = preload.Users[one.UserId]
				firstName = user.FirstName
				lastName = user.LastName
				email = user.Email
			}

			plan := &entity.Plan{}
			if preload.Plans[one.PlanId] != nil {
				plan = preload.Plans[one.PlanId]
			}

			var productName = ""
			if plan.ProductId > 0 {
				if preload.Products[uint64(plan.ProductId)] != nil {
					productName = preload.Products[uint64(plan.ProductId)].ProductName
				}
			}

			mainList = append(mainList, &ExportSubscriptionEntity{
				SubscriptionId:         one.SubscriptionId,
				ExternalSubscriptionId: one.ExternalSubscriptionId,
				UserId:                 fmt.Sprintf("%v", user.Id),
				ExternalUserId:         fmt.Sprintf("%v", user.ExternalUserId),
				FirstName:              firstName,
				LastName:               lastName,
				Email:                  email,
				MerchantName:           merchant.Name,
				Amount:                 utility.ConvertCentToDollarStr(one.Amount, one.Currency),
				Currency:               one.Currency,
				ProductId:              fmt.Sprintf("%v", plan.ProductId),
				ProductName:            productName,
				PlanId:                 fmt.Sprintf("%v", plan.Id),
				ExternalPlanId:         fmt.Sprintf("%v", plan.ExternalPlanId),
				PlanName:               plan.PlanName,
				PlanInternalName:       plan.InternalName,
				PlanIntervalUnit:       plan.IntervalUnit,
				PlanIntervalCount:      fmt.Sprintf("%d", plan.IntervalCount),
				Quantity:               fmt.Sprintf("%v", one.Quantity),
				Gateway:                subGateway,
				Status:                 consts.SubStatusToEnum(one.Status).Description(),
				CancelAtPeriodEnd:      canAtPeriodEnd,
				CurrentPeriodStart:     gtime.NewFromTimeStamp(one.CurrentPeriodStart + timeZone),
				CurrentPeriodEnd:       gtime.NewFromTimeStamp(one.CurrentPeriodEnd + timeZone),
				BillingCycleAnchor:     gtime.NewFromTimeStamp(one.BillingCycleAnchor + timeZone),
				DunningTime:            gtime.NewFromTimeStamp(one.DunningTime + timeZone),
				TrialEnd:               gtime.NewFromTimeStamp(one.TrialEnd + timeZone),
				FirstPaidTime:          gtime.NewFromTimeStamp(one.FirstPaidTime + timeZone),
				CancelReason:           one.CancelReason,
				CountryCode:            one.CountryCode,
				TaxPercentage:          utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
				CreateTime:             gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				StripeUserId:           stripeUserId,
				StripePaymentMethod:    stripePaymentMethod,
				PaypalVaultId:          paypalVaultId,
				TimeZone:               timeZoneStr,
			})
		}
	}
	return mainList, nil
}

// subscriptionList directly queries subscription data from database to avoid N+1 queries
func subscriptionList(ctx context.Context, req *service.SubscriptionListInternalReq) []*entity.Subscription {
	var mainList = make([]*entity.Subscription, 0)
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
		utility.Assert(strings.Contains("gmt_create|gmt_modify", req.SortField), "sortField should one of gmt_create|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}

	baseQuery := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, req.MerchantId)

	if req.Status != nil && len(req.Status) > 0 {
		baseQuery = baseQuery.WhereIn(dao.Subscription.Columns().Status, req.Status)
	}

	if len(req.Email) > 0 {
		var userIdList = make([]uint64, 0)
		var userList []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		if req.UserId > 0 {
			userQuery = userQuery.Where(dao.UserAccount.Columns().Id, req.UserId)
		}
		userQuery = userQuery.WhereLike(dao.UserAccount.Columns().Email, "%"+req.Email+"%")
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&userList)
		for _, user := range userList {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) == 0 {
			return mainList
		}
		baseQuery = baseQuery.WhereIn(dao.Subscription.Columns().UserId, userIdList)
	} else if req.UserId > 0 {
		baseQuery = baseQuery.Where(dao.Subscription.Columns().UserId, req.UserId)
	}

	if req.ProductIds != nil && len(req.ProductIds) > 0 {
		if req.PlanIds == nil {
			req.PlanIds = make([]uint64, 0)
		}
		var plans []*entity.Plan
		planQuery := dao.Plan.Ctx(ctx)
		if isInt64InArray(req.ProductIds, 0) {
			planQuery = planQuery.Where(planQuery.Builder().WhereOrIn(dao.Plan.Columns().ProductId, req.ProductIds).WhereOrNull(dao.Plan.Columns().ProductId))
		} else {
			planQuery = planQuery.WhereIn(dao.Plan.Columns().ProductId, req.ProductIds)
		}
		_ = planQuery.Where(dao.Plan.Columns().IsDeleted, 0).Scan(&plans)
		for _, plan := range plans {
			req.PlanIds = append(req.PlanIds, plan.Id)
		}
	}

	if req.PlanIds != nil && len(req.PlanIds) > 0 {
		baseQuery = baseQuery.WhereIn(dao.Subscription.Columns().PlanId, req.PlanIds)
	}

	if req.CreateTimeStart > 0 {
		baseQuery = baseQuery.WhereGTE(dao.Subscription.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		baseQuery = baseQuery.WhereLTE(dao.Subscription.Columns().CreateTime, req.CreateTimeEnd)
	}
	if req.AmountStart != nil && req.AmountEnd != nil {
		utility.Assert(*req.AmountStart <= *req.AmountEnd, "amountStart should lower than amountEnd")
	}
	if req.AmountStart != nil {
		baseQuery = baseQuery.WhereGTE(dao.Subscription.Columns().Amount, &req.AmountStart)
	}
	if req.AmountEnd != nil {
		baseQuery = baseQuery.WhereLTE(dao.Subscription.Columns().Amount, &req.AmountEnd)
	}
	if len(req.Currency) > 0 {
		baseQuery = baseQuery.Where(dao.Subscription.Columns().Currency, strings.ToUpper(req.Currency))
	}

	var err error
	baseQuery = baseQuery.Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty()

	if req.SkipTotal {
		err = baseQuery.Scan(&mainList)
	} else {
		err = baseQuery.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return mainList
	}

	return mainList
}

// isInt64InArray checks if a target value exists in an array
func isInt64InArray(arr []int64, target int64) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}
