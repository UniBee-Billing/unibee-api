package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/merchant/plan"
)

func (c *ControllerPlan) Expire(ctx context.Context, req *plan.ExpireReq) (res *plan.ExpireRes, err error) {

	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
