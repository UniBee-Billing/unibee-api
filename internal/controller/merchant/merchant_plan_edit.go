package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Edit(ctx context.Context, req *plan.EditReq) (res *plan.EditRes, err error) {

	one, err := plan2.PlanEdit(ctx, &plan2.EditInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		PlanId:             req.PlanId,
		ExternalPlanId:     req.ExternalPlanId,
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
		TrialAmount:        req.TrialAmount,
		TrialDurationTime:  req.TrialDurationTime,
		TrialDemand:        req.TrialDemand,
		CancelAtTrialEnd:   req.CancelAtTrialEnd,
	})
	if err != nil {
		return nil, err
	}
	return &plan.EditRes{Plan: bean.SimplifyPlan(one)}, nil
}
