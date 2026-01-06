package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/checkout"
)

func (c *ControllerCheckout) Detail(ctx context.Context, req *checkout.DetailReq) (res *checkout.DetailRes, err error) {
	one := query.GetMerchantCheckoutById(ctx, _interface.GetMerchantId(ctx), uint64(req.CheckoutId))
	utility.Assert(one != nil, "Checkout not found, please setup first")
	return &checkout.DetailRes{MerchantCheckout: bean.SimplifyMerchantCheckout(ctx, one)}, nil
}
