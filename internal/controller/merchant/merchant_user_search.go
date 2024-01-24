package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/user"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
