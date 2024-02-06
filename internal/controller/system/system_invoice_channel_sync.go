package system

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/system/invoice"
)

func (c *ControllerInvoice) ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
