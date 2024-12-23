package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface/context"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) New(ctx context.Context, req *plan.NewReq) (res *plan.NewRes, err error) {

	one, err := plan2.PlanCreate(ctx, &plan2.PlanInternalReq{
		ExternalPlanId:     req.ExternalPlanId,
		Type:               req.Type,
		PlanName:           req.PlanName,
		Amount:             req.Amount,
		Currency:           req.Currency,
		IntervalUnit:       req.IntervalUnit,
		IntervalCount:      req.IntervalCount,
		Description:        req.Description,
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		ImageUrl:           req.ImageUrl,
		HomeUrl:            req.HomeUrl,
		AddonIds:           req.AddonIds,
		OnetimeAddonIds:    req.OnetimeAddonIds,
		MetricLimits:       req.MetricLimits,
		GasPayer:           req.GasPayer,
		Metadata:           req.Metadata,
		MerchantId:         _interface.GetMerchantId(ctx),
		TrialDemand:        req.TrialDemand,
		TrialAmount:        req.TrialAmount,
		TrialDurationTime:  req.TrialDurationTime,
		CancelAtTrialEnd:   req.CancelAtTrialEnd,
		ProductId:          req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &plan.NewRes{Plan: bean.SimplifyPlan(one)}, nil
}
