package merchant

import (
	"context"
	"unibee/api/bean"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Copy(ctx context.Context, req *product.CopyReq) (res *product.CopyRes, err error) {
	one, err := product2.ProductCopy(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.CopyRes{Product: bean.SimplifyProduct(one)}, nil
}
