package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Edit(ctx context.Context, req *product.EditReq) (res *product.EditRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
