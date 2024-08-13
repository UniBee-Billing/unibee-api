package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionListInternalReq struct {
	MerchantId      uint64   `json:"merchantId" dc:"MerchantId"`
	UserId          int64    `json:"userId"  dc:"UserId" `
	Status          []int    `json:"status" dc:"Default All，,Status，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	Currency        string   `json:"currency" dc:"The currency of subscription" `
	PlanIds         []uint64 `json:"planIds" dc:"The filter ids of plan" `
	ProductIds      []int64  `json:"productIds" dc:"The filter ids of product" `
	AmountStart     *int64   `json:"amountStart" dc:"The filter start amount of subscription" `
	AmountEnd       *int64   `json:"amountEnd" dc:"The filter end amount of subscription" `
	SortField       string   `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string   `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int      `json:"page" dc:"Page, Start With 0" `
	Count           int      `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*detail.SubscriptionDetail, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	{
		one.Data = ""
		one.ResponseData = ""
	}
	user := query.GetUserAccountById(ctx, one.UserId)
	var addonParams []*bean.PlanAddonParam
	if len(one.AddonData) > 0 {
		err := utility.UnmarshalFromJsonString(one.AddonData, &addonParams)
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionDetail parse addon param:%s", err.Error())
		}
	}
	latestInvoiceOne := bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId))
	if latestInvoiceOne != nil {
		latestInvoiceOne.Discount = bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, one.MerchantId, latestInvoiceOne.DiscountCode))
	}
	return &detail.SubscriptionDetail{
		User:                                bean.SimplifyUserAccount(user),
		Subscription:                        bean.SimplifySubscription(ctx, one),
		Gateway:                             bean.SimplifyGateway(query.GetGatewayById(ctx, one.GatewayId)),
		Plan:                                bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		AddonParams:                         addonParams,
		Addons:                              addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		LatestInvoice:                       latestInvoiceOne,
		Discount:                            bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, one.MerchantId, one.DiscountCode)),
		UnfinishedSubscriptionPendingUpdate: GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId(ctx, one.PendingUpdateId),
	}, nil
}

func SubscriptionList(ctx context.Context, req *SubscriptionListInternalReq) (list []*detail.SubscriptionDetail, total int) {
	var mainList []*entity.Subscription
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
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
		Where(dao.Subscription.Columns().MerchantId, req.MerchantId).
		Where(dao.Subscription.Columns().UserId, req.UserId)
	if req.Status != nil && len(req.Status) > 0 {
		baseQuery = baseQuery.WhereIn(dao.Subscription.Columns().Status, req.Status)
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
		utility.Assert(*req.AmountStart <= *req.AmountEnd, "amountStart should lower then amountEnd")
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
		return nil, 0
	}
	var totalPlanIds []uint64
	for _, sub := range mainList {
		totalPlanIds = append(totalPlanIds, sub.PlanId)
		var addonParams []*bean.PlanAddonParam
		if len(sub.AddonData) > 0 {
			err = utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
			if err == nil {
				for _, s := range addonParams {
					totalPlanIds = append(totalPlanIds, s.AddonPlanId) // 添加到整数列表中
				}
			}
		}
		{
			sub.Data = ""
			sub.ResponseData = ""
		}
		user := query.GetUserAccountById(ctx, sub.UserId)
		if user != nil {
			user.Password = ""
		}
		latestInvoiceOne := bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId))
		if latestInvoiceOne != nil {
			latestInvoiceOne.Discount = bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, sub.MerchantId, latestInvoiceOne.DiscountCode))
		}
		list = append(list, &detail.SubscriptionDetail{
			User:          bean.SimplifyUserAccount(user),
			Subscription:  bean.SimplifySubscription(ctx, sub),
			Gateway:       bean.SimplifyGateway(query.GetGatewayById(ctx, sub.GatewayId)),
			Plan:          nil,
			Addons:        nil,
			LatestInvoice: latestInvoiceOne,
			Discount:      bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, sub.MerchantId, sub.DiscountCode)),
			AddonParams:   addonParams,
		})
	}
	if len(totalPlanIds) > 0 {
		var allPlanList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalPlanIds).OmitEmpty().Scan(&allPlanList)
		if err == nil {
			mapPlans := make(map[uint64]*entity.Plan)
			for _, pair := range allPlanList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, subRo := range list {
				subRo.Plan = bean.SimplifyPlan(mapPlans[subRo.Subscription.PlanId])
				if len(subRo.AddonParams) > 0 {
					for _, param := range subRo.AddonParams {
						if mapPlans[param.AddonPlanId] != nil {
							subRo.Addons = append(subRo.Addons, &bean.PlanAddonDetail{
								Quantity:  param.Quantity,
								AddonPlan: bean.SimplifyPlan(mapPlans[param.AddonPlanId]),
							})
						}
					}
				}
			}
		}
	}
	return list, total
}

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
