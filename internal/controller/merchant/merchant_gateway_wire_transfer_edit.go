package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferEdit(ctx context.Context, req *gateway.WireTransferEditReq) (res *gateway.WireTransferEditRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
