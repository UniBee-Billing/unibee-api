package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) New(ctx context.Context, req *plan.NewReq) (res *plan.NewRes, err error) {

	one, err := service.PlanCreate(ctx, &service.PlanInternalReq{
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
	})
	if err != nil {
		return nil, err
	}
	return &plan.NewRes{Plan: bean.SimplifyPlan(one)}, nil
}
