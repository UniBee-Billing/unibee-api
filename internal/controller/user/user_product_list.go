package user

import (
	"context"
	_interface "unibee/internal/interface"
	product2 "unibee/internal/logic/product"

	"unibee/api/user/product"
)

func (c *ControllerProduct) List(ctx context.Context, req *product.ListReq) (res *product.ListRes, err error) {
	list, total := product2.ProductList(ctx, &product2.ListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		Status:     []int{1, 2},
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &product.ListRes{
		Products: list,
		Total:    total,
	}, nil
}
