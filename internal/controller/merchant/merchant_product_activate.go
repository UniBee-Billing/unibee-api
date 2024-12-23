package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Activate(ctx context.Context, req *product.ActivateReq) (res *product.ActivateRes, err error) {
	err = product2.ProductActivate(ctx, _interface.GetMerchantId(ctx), req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.ActivateRes{}, nil
}
