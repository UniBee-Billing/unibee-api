package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	discount2 "unibee/internal/logic/discount"
	"unibee/internal/query"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) QuantityIncrement(ctx context.Context, req *discount.QuantityIncrementReq) (res *discount.QuantityIncrementRes, err error) {
	err = discount2.QuantityIncrement(ctx, _interface.GetMerchantId(ctx), req.Id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &discount.QuantityIncrementRes{DiscountCode: bean.SimplifyMerchantDiscountCode(query.GetDiscountById(ctx, req.Id))}, nil
}
