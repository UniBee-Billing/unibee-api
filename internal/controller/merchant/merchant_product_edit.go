package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	product2 "unibee/internal/logic/product"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Edit(ctx context.Context, req *product.EditReq) (res *product.EditRes, err error) {
	one, err := product2.ProductEdit(ctx, &product2.EditInternalReq{
		ProductId:   req.ProductId,
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
	return &product.EditRes{Product: bean.SimplifyProduct(one)}, nil
}
