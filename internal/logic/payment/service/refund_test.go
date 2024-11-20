package service

import (
	"context"
	"testing"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/query"
	"unibee/utility"
)

func TestGenerateInvoicePdf(t *testing.T) {
	ctx := context.Background()
	one := query.GetPaymentByPaymentId(ctx, "pay20241109rNB7Ycsr34Q0pxp")
	utility.Assert(one != nil, "one not found")
	refund := query.GetRefundByRefundId(ctx, "ref20241109UXZ4yebNyVhXv7q")
	utility.Assert(refund != nil, "refund not found")
	invoice_compute.CreateInvoiceSimplifyForRefund(ctx, one, refund)
	//one.RefundId = "refundId"
	//one.SendNote = "iv20240202ERExKnb6OhMfyyY"
	//var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	//err := createInvoicePdf(detail.ConvertInvoiceToDetail(ctx, one), query.GetMerchantById(ctx, one.MerchantId), query.GetUserAccountById(ctx, one.UserId), query.GetGatewayById(ctx, one.GatewayId), savePath)
	//utility.AssertError(err, "Pdf Generator Error")
	//err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	//if err != nil {
	//	return
	//}
	//err = os.Remove("pdf_test.pdf")
	//if err != nil {
	//	return
	//}
}
