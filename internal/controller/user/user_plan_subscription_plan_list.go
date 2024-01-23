package user

import (
	"context"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/plan"
)

func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error) {
	return &plan.SubscriptionPlanListRes{Plans: service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId:    req.MerchantId,
		Type:          req.Type,
		Status:        consts.PlanStatusActive,
		PublishStatus: 2,
		Currency:      req.Currency,
		SortField:     req.SortField,
		SortType:      req.SortType,
		Page:          req.Page,
		Count:         req.Count,
	})}, nil
}
