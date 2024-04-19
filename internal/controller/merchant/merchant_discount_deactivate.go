package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Deactivate(ctx context.Context, req *discount.DeactivateReq) (res *discount.DeactivateRes, err error) {
	err = discount2.DeactivateMerchantDiscountCode(ctx, _interface.GetMerchantId(ctx), req.Id)
	if err != nil {
		return nil, err
	}
	return &discount.DeactivateRes{}, nil
}
