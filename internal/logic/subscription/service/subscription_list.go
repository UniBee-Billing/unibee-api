package service

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/subscription/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

type SubscriptionListInternalReq struct {
	MerchantId int64 `p:"merchantId" dc:"MerchantId"`
	UserId     int64 `p:"userId"  dc:"UserId" `
	Status     int   `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	Page       int   `p:"page" d:"0"  dc:"分页页码,0开始" `
	Count      int   `p:"count" d:"20"  dc:"订阅计划货币" dc:"每页数量" `
}

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*subscription.SubscriptionDetailRes, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")

	return &subscription.SubscriptionDetailRes{
		Subscription: one,
		Plan:         query.GetPlanById(ctx, one.PlanId),
		Addons:       query.GetSubscriptionAddonsBySubscriptionId(ctx, one.SubscriptionId),
	}, nil
}

func SubscriptionList(ctx context.Context, req *SubscriptionListInternalReq) (list []*ro.SubscriptionDetailRo) {
	var mainList []*entity.Subscription
	if req.Count <= 0 {
		req.Count = 10 //每页数量默认 10
	}
	if req.Page < 0 {
		req.Page = 0
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, req.MerchantId).
		Where(dao.Subscription.Columns().UserId, req.UserId).
		Where(dao.Subscription.Columns().Status, req.Status).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalPlanIds []int64
	for _, sub := range mainList {
		totalPlanIds = append(totalPlanIds, sub.PlanId)
		var addonParams []*ro.SubscriptionPlanAddonParamRo
		if len(sub.AddonData) > 0 {
			err := utility.UnmarshalFromJsonString(sub.AddonData, addonParams)
			if err == nil {
				for _, s := range addonParams {
					totalPlanIds = append(totalPlanIds, s.AddonPlanId) // 添加到整数列表中
				}
			}
		}
		list = append(list, &ro.SubscriptionDetailRo{
			Subscription: sub,
			Plan:         nil,
			Addons:       nil,
			AddonParams:  addonParams,
		})
	}
	if len(totalPlanIds) > 0 {
		//查询所有 Plan
		var allPlanList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalPlanIds).Scan(&allPlanList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allPlanList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, subRo := range list {
				subRo.Plan = mapPlans[subRo.Subscription.PlanId]
				if len(subRo.AddonParams) > 0 {
					for _, param := range subRo.AddonParams {
						if mapPlans[param.AddonPlanId] != nil {
							subRo.Addons = append(subRo.Addons, &ro.SubscriptionPlanAddonRo{
								Quantity:  param.Quantity,
								AddonPlan: mapPlans[param.AddonPlanId],
							})
						}
					}
				}
			}
		}
	}
	return list
}
