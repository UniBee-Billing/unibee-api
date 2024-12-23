package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Delete(ctx context.Context, req *discount.DeleteReq) (res *discount.DeleteRes, err error) {
	err = discount2.DeleteMerchantDiscountCode(ctx, _interface.GetMerchantId(ctx), req.Id)
	if err != nil {
		return nil, err
	}
	return &discount.DeleteRes{}, nil
}
