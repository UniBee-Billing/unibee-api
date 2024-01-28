package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetInvoiceByInvoiceId(ctx context.Context, invoiceId string) (one *entity.Invoice) {
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{InvoiceId: invoiceId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetInvoiceByChannelInvoiceId(ctx context.Context, channelInvoiceId string) (one *entity.Invoice) {
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{ChannelInvoiceId: channelInvoiceId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetInvoiceByPaymentId(ctx context.Context, paymentId string) (one *entity.Invoice) {
	err := dao.Invoice.Ctx(ctx).Where(entity.Invoice{PaymentId: paymentId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
