package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Activate(ctx context.Context, req *product.ActivateReq) (res *product.ActivateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
