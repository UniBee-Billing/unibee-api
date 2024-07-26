package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Delete(ctx context.Context, req *product.DeleteReq) (res *product.DeleteRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
