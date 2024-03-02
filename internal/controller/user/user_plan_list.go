package user

import (
	"context"
	"unibee/api/user/plan"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPlan) List(ctx context.Context, req *plan.ListReq) (res *plan.ListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
	}

	publishPlans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId:    _interface.GetMerchantId(ctx),
		Status:        []int{consts.PlanStatusActive},
		PublishStatus: consts.PlanPublishStatusPublished,
		Currency:      req.Currency,
		Page:          0,
		Count:         10,
	})
	sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), _interface.GetMerchantId(ctx))
	if sub != nil {
		subPlan := query.GetPlanById(ctx, sub.PlanId)
		if subPlan != nil && subPlan.PublishStatus != consts.PlanPublishStatusPublished {
			userPlan, _ := service.SubscriptionPlanDetail(ctx, subPlan.Id)
			if userPlan != nil {
				publishPlans = append(publishPlans, userPlan.Plan)
			}
		}
	}
	return &plan.ListRes{Plans: publishPlans}, nil
}
