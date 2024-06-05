package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) PdfUpdate(ctx context.Context, req *invoice.PdfUpdateReq) (res *invoice.PdfUpdateRes, err error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	if one != nil && len(one.SendPdf) == 0 && req.SendPdf != nil && len(*req.SendPdf) > 0 {
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendPdf:   *req.SendPdf,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		utility.AssertError(err, "Update invoice pdf error")
	}
	return &invoice.PdfUpdateRes{}, nil
}
