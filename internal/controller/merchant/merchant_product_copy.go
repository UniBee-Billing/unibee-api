package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Copy(ctx context.Context, req *product.CopyReq) (res *product.CopyRes, err error) {
	one, err := product2.ProductCopy(ctx, _interface.GetMerchantId(ctx), req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.CopyRes{Product: bean.SimplifyProduct(one)}, nil
}
