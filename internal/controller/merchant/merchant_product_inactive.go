package merchant

import (
	"context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Inactive(ctx context.Context, req *product.InactiveReq) (res *product.InactiveRes, err error) {
	err = product2.ProductInactivate(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.InactiveRes{}, nil
}
