package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) New(ctx context.Context, req *discount.NewReq) (res *discount.NewRes, err error) {
	err = discount2.NewMerchantDiscountCode(ctx, &discount2.CreateDiscountCodeInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		Code:               req.Code,
		Name:               req.Name,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		UserLimit:          req.UserLimit,
		SubscriptionLimit:  req.SubscriptionLimit,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &discount.NewRes{}, nil
}
