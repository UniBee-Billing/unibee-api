package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Activate(ctx context.Context, req *discount.ActivateReq) (res *discount.ActivateRes, err error) {
	err = discount2.ActivateMerchantDiscountCode(ctx, _interface.GetMerchantId(ctx), req.Code)
	if err != nil {
		return nil, err
	}
	return &discount.ActivateRes{}, nil
}