package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/product"
)

func (c *ControllerProduct) Copy(ctx context.Context, req *product.CopyReq) (res *product.CopyRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
