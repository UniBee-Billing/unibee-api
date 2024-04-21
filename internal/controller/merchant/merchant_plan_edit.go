package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	"unibee/utility"
)

func (c *ControllerPlan) Edit(ctx context.Context, req *plan.EditReq) (res *plan.EditRes, err error) {

	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	one, err := service.PlanEdit(ctx, &service.EditInternalReq{
		PlanId:             req.PlanId,
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
	})
	if err != nil {
		return nil, err
	}
	return &plan.EditRes{Plan: bean.SimplifyPlan(one)}, nil
}
