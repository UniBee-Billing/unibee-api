package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) Detail(ctx context.Context, req *discount.DetailReq) (res *discount.DetailRes, err error) {
	one := query.GetDiscountById(ctx, req.Id)
	utility.Assert(one != nil, "code not found")
	return &discount.DetailRes{Discount: detail.ConvertMerchantDiscountCodeDetail(ctx, one)}, nil
}
