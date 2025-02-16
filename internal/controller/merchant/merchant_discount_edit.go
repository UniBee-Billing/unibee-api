package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Edit(ctx context.Context, req *discount.EditReq) (res *discount.EditRes, err error) {
	one, err := discount2.EditMerchantDiscountCode(ctx, &discount2.CreateDiscountCodeInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		Id:                 req.Id,
		Name:               req.Name,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		CycleLimit:         req.CycleLimit,
		//SubscriptionLimit:  req.SubscriptionLimit,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		Metadata:          req.Metadata,
		PlanApplyType:     req.PlanApplyType,
		PlanIds:           req.PlanIds,
		Quantity:          req.Quantity,
		Advance:           req.Advance,
		UserLimit:         req.UserLimit,
		UserScope:         req.UserScope,
		UpgradeLongerOnly: req.UpgradeLongerOnly,
		UpgradeOnly:       req.UpgradeOnly,
	})
	if err != nil {
		return nil, err
	}
	return &discount.EditRes{Discount: bean.SimplifyMerchantDiscountCode(one)}, nil
}
