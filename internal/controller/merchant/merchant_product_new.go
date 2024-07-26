package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) New(ctx context.Context, req *product.NewReq) (res *product.NewRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
