package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) New(ctx context.Context, req *discount.NewReq) (res *discount.NewRes, err error) {
	one, err := discount2.NewMerchantDiscountCode(ctx, &discount2.CreateDiscountCodeInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		Code:               req.Code,
		Name:               req.Name,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		Type:               0,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		CycleLimit:         req.CycleLimit,
		//SubscriptionLimit:  req.SubscriptionLimit,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PlanIds:   req.PlanIds,
		Metadata:  req.Metadata,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &discount.NewRes{Discount: bean.SimplifyMerchantDiscountCode(one)}, nil
}
