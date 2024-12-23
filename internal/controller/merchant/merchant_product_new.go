package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) New(ctx context.Context, req *product.NewReq) (res *product.NewRes, err error) {
	one, err := product2.ProductNew(ctx, &product2.NewInternalReq{
		MerchantId:  _interface.GetMerchantId(ctx),
		ProductName: req.ProductName,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		HomeUrl:     req.HomeUrl,
		Status:      req.Status,
		Metadata:    req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	return &product.NewRes{Product: bean.SimplifyProduct(one)}, nil
}
