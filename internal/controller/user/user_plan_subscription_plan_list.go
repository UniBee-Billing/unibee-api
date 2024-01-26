package user

import (
	"context"
	"go-oversea-pay/api/user/plan"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/plan/service"
	service2 "go-oversea-pay/internal/logic/subscription/service"
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
	subs := service2.SubscriptionList(ctx, &service2.SubscriptionListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     int64(_interface.BizCtx().Get(ctx).User.Id),
		Status:     consts.SubStatusActive,
		SortField:  "gmt_create",
		SortType:   "desc",
		Page:       0,
		Count:      1,
	})
	if len(subs) > 0 && subs[0].Plan.PublishStatus != consts.PlanPublishStatusPublished {
		userPlan, _ := service.SubscriptionPlanDetail(ctx, subs[0].Subscription.PlanId)
		if userPlan != nil {
			publishPlans = append(publishPlans, userPlan.Plan)
		}
	}
	return &plan.SubscriptionPlanListRes{Plans: publishPlans}, nil
}
