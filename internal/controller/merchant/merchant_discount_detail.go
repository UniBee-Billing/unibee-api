package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Detail(ctx context.Context, req *discount.DetailReq) (res *discount.DetailRes, err error) {
	one := query.GetDiscountByCode(ctx, _interface.GetMerchantId(ctx), req.Code)
	utility.Assert(one != nil, "code not found")
	return &discount.DetailRes{Discount: bean.SimplifyMerchantDiscountCode(one)}, nil
}
