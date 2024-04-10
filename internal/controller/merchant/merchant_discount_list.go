package merchant

import (
	"context"
	"unibee/api/merchant/discount"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"
)

func (c *ControllerDiscount) List(ctx context.Context, req *discount.ListReq) (res *discount.ListRes, err error) {
	return &discount.ListRes{MerchantDiscountCodes: discount2.MerchantDiscountCodeList(ctx, _interface.GetMerchantId(ctx))}, nil
}
