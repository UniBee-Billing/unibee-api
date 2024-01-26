package user

import (
	"context"
	"go-oversea-pay/api/user/plan"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/plan/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

// SubscriptionPlanList todo mark 修改成 User Protal Plan List， Only Return Publish Plans And User Sub Plan
func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error) {
	// service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
	}

	publishPlans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId:    req.MerchantId,
		Type:          req.Type,
		Status:        consts.PlanStatusActive,
		PublishStatus: consts.PlanPublishStatusPublished,
		Currency:      req.Currency,
		//SortField:     req.SortField,
		//SortType:      req.SortType,
		Page:  0,
		Count: 10,
	})
	sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), req.MerchantId)
	if sub != nil {
		subPlan := query.GetPlanById(ctx, sub.PlanId)
		if subPlan != nil && subPlan.PublishStatus != consts.PlanPublishStatusPublished {
			userPlan, _ := service.SubscriptionPlanDetail(ctx, int64(subPlan.Id))
			if userPlan != nil {
				publishPlans = append(publishPlans, userPlan.Plan)
			}
		}
	}
	return &plan.SubscriptionPlanListRes{Plans: publishPlans}, nil
}
