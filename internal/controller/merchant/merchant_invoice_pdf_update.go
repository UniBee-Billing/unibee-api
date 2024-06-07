package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/merchant/invoice"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerInvoice) PdfUpdate(ctx context.Context, req *invoice.PdfUpdateReq) (res *invoice.PdfUpdateRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "invalid invoiceId")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("createInvoicePdf Unmarshal Metadata error:%s", err.Error())
		}
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	if req.IssueAddress != nil {
		metadata["IssueAddress"] = *req.IssueAddress
	}
	if req.IssueCompanyName != nil {
		metadata["IssueCompanyName"] = *req.IssueCompanyName
	}
	if req.IssueRegNumber != nil {
		metadata["IssueRegNumber"] = *req.IssueRegNumber
	}
	if req.IssueVatNumber != nil {
		metadata["IssueVatNumber"] = *req.IssueVatNumber
	}
	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().MetaData:  utility.MarshalToJsonString(metadata),
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "update Invoice error")
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail)
	return &invoice.PdfUpdateRes{}, nil
}
