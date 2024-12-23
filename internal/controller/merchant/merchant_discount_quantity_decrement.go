package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	discount2 "unibee/internal/logic/discount"
	"unibee/internal/query"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) QuantityDecrement(ctx context.Context, req *discount.QuantityDecrementReq) (res *discount.QuantityDecrementRes, err error) {
	err = discount2.QuantityDecrement(ctx, _interface.GetMerchantId(ctx), req.Id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &discount.QuantityDecrementRes{DiscountCode: bean.SimplifyMerchantDiscountCode(query.GetDiscountById(ctx, req.Id))}, nil
}
