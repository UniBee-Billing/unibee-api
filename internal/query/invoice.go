package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
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

func GetSubLatestPaidInvoice(ctx context.Context, subId string) (one *entity.Invoice) {
	if len(subId) == 0 {
		return nil
	}
	err := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().SubscriptionId, subId).
		Where(dao.Invoice.Columns().BizType, consts.BizTypeSubscription).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusPaid).
		OrderDesc(dao.Invoice.Columns().CreateTime).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return one
}
