package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) ProcessInvoiceForPay(ctx context.Context, req *invoice.ProcessInvoiceForPayReq) (res *invoice.ProcessInvoiceForPayRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
