package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Delete(ctx context.Context, req *product.DeleteReq) (res *product.DeleteRes, err error) {

	err = product2.ProductDelete(ctx, _interface.GetMerchantId(ctx), req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.DeleteRes{}, nil
}
