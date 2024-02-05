package system

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/api/system/invoice"
)

func (c *ControllerInvoice) BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
