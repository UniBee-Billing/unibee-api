package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/merchant/invoice"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Send-Email:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	if one != nil && len(one.SendPdf) == 0 && req.SendPdf != nil && len(*req.SendPdf) > 0 {
		utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "no permission")
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendPdf:   *req.SendPdf,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		utility.AssertError(err, "Update invoice pdf error")
	}
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SendEmailRes{}, nil
}
