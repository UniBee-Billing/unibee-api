package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Edit(ctx context.Context, req *discount.EditReq) (res *discount.EditRes, err error) {
	err = discount2.EditMerchantDiscountCode(ctx, &discount2.CreateDiscountCodeInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		Code:               req.Code,
		Name:               req.Name,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		CycleLimit:         req.CycleLimit,
		//SubscriptionLimit:  req.SubscriptionLimit,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &discount.EditRes{}, nil
}
