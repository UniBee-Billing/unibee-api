package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetInvoiceByInvoiceId(ctx context.Context, invoiceId string) (one *entity.Invoice) {
	if len(invoiceId) == 0 {
		return nil
	}
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{InvoiceId: invoiceId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetInvoiceByPaymentId(ctx context.Context, paymentId string) (one *entity.Invoice) {
	if len(paymentId) == 0 {
		return nil
	}
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{PaymentId: paymentId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetInvoiceByRefundId(ctx context.Context, refundId string) (one *entity.Invoice) {
	if len(refundId) == 0 {
		return nil
	}
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{RefundId: refundId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
