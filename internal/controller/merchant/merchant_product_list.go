package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) List(ctx context.Context, req *product.ListReq) (res *product.ListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
