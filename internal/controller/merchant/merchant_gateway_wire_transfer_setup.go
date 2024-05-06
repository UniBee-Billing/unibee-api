package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferSetup(ctx context.Context, req *gateway.WireTransferSetupReq) (res *gateway.WireTransferSetupRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
