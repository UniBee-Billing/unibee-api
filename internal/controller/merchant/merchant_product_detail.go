package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Detail(ctx context.Context, req *product.DetailReq) (res *product.DetailRes, err error) {
	one := query.GetProductById(ctx, req.ProductId)
	if one == nil {
		return nil, gerror.New("product not found")
	}
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "No Permission")
	return &product.DetailRes{Product: bean.SimplifyProduct(one)}, nil
}