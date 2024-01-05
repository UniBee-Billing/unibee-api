package merchant

import (
	"context"
	"fmt"
	v1 "go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {
	fmt.Println("context: ", ctx)
	plans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId: req.MerchantId,
		Type:       req.Type,
		Status:     req.Status,
		Currency:   req.Currency,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &v1.SubscriptionPlanListRes{Plans: plans}, nil
}
