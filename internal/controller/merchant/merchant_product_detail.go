package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Detail(ctx context.Context, req *product.DetailReq) (res *product.DetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
