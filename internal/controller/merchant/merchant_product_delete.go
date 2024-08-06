package merchant

import (
	"context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Delete(ctx context.Context, req *product.DeleteReq) (res *product.DeleteRes, err error) {

	err = product2.ProductDelete(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &product.DeleteRes{}, nil
}
