package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Inactive(ctx context.Context, req *product.InactiveReq) (res *product.InactiveRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
