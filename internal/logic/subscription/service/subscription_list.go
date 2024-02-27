package service

import (
	"context"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionListInternalReq struct {
	MerchantId uint64 `p:"merchantId" dc:"MerchantId"`
	UserId     int64  `p:"userId"  dc:"UserId" `
	Status     []int  `p:"status" dc:"Default All，,Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page" dc:"Page, Start WIth 0" `
	Count      int    `p:"count" dc:"Count Of Page" `
}

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*ro.SubscriptionDetailVo, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	{
		one.Data = ""
		one.ResponseData = ""
	}
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	if user != nil {
		user.Password = ""
	}
	return &ro.SubscriptionDetailVo{
		User:                                ro.SimplifyUserAccount(user),
		Subscription:                        ro.SimplifySubscription(one),
		Gateway:                             ConvertGatewayToRo(query.GetGatewayById(ctx, one.GatewayId)),
		Plan:                                ro.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:                              addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UnfinishedSubscriptionPendingUpdate: GetUnfinishedSubscriptionPendingUpdateDetailByUpdateSubscriptionId(ctx, one.PendingUpdateId),
	}, nil
}

func SubscriptionList(ctx context.Context, req *SubscriptionListInternalReq) (list []*ro.SubscriptionDetailVo) {
	var mainList []*entity.Subscription
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var sortKey = "gmt_modify desc"
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
	err := baseQuery.Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalPlanIds []uint64
	for _, sub := range mainList {
		totalPlanIds = append(totalPlanIds, sub.PlanId)
		var addonParams []*ro.SubscriptionPlanAddonParamRo
		if len(sub.AddonData) > 0 {
			err := utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
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
		user := query.GetUserAccountById(ctx, uint64(sub.UserId))
		if user != nil {
			user.Password = ""
		}
		list = append(list, &ro.SubscriptionDetailVo{
			User:         ro.SimplifyUserAccount(user),
			Subscription: ro.SimplifySubscription(sub),
			Gateway:      ConvertGatewayToRo(query.GetGatewayById(ctx, sub.GatewayId)),
			Plan:         nil,
			Addons:       nil,
			AddonParams:  addonParams,
		})
	}
	if len(totalPlanIds) > 0 {
		//查询所有 Plan
		var allPlanList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalPlanIds).OmitEmpty().Scan(&allPlanList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[uint64]*entity.SubscriptionPlan)
			for _, pair := range allPlanList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, subRo := range list {
				subRo.Plan = ro.SimplifyPlan(mapPlans[subRo.Subscription.PlanId])
				if len(subRo.AddonParams) > 0 {
					for _, param := range subRo.AddonParams {
						if mapPlans[param.AddonPlanId] != nil {
							subRo.Addons = append(subRo.Addons, &ro.PlanAddonVo{
								Quantity:  param.Quantity,
								AddonPlan: ro.SimplifyPlan(mapPlans[param.AddonPlanId]),
							})
						}
					}
				}
			}
		}
	}
	return list
}
