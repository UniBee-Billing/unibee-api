package discount

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/credit/payment"
	"unibee/internal/logic/discount"
	"unibee/internal/query"
)

// InvoiceRollbackAllDiscountsFromPayment payment total refund, partial refund is not involved
func InvoiceRollbackAllDiscountsFromPayment(ctx context.Context, invoiceId string, paymentId string) error {
	if len(paymentId) == 0 {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invalid paymentId:%s", paymentId)
		return gerror.New("invalid paymentId")
	}
	invoice := query.GetInvoiceByPaymentId(ctx, paymentId)
	if invoice == nil {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invoice not found, paymentId:%s", paymentId)
		return gerror.New("invoice not found")
	}
	if len(invoiceId) == 0 {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invalid invoiceId:%s", invoiceId)
		return gerror.New("invalid invoiceId")
	}

	return InvoiceRollbackAllDiscountsFromInvoice(ctx, invoiceId)
}

// InvoiceRollbackAllDiscountsFromInvoice invoice create failed|cancel|failed, partial refund is not involved
func InvoiceRollbackAllDiscountsFromInvoice(ctx context.Context, invoiceId string) error {
	if len(invoiceId) == 0 {
		g.Log().Error(ctx, "InvoiceRollbackAllDiscountsFromInvoice invalid invoiceId:%s", invoiceId)
		return gerror.New("invalid invoiceId")
	}
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		return gerror.New("invoice not found")
	}

	err := discount.UserDiscountRollbackFromInvoice(context.Background(), invoiceId)
	if err != nil {
		g.Log().Error(context.Background(), "InvoiceRollbackAllDiscountsFromInvoice UserDiscountRollbackFromInvoice invoiceId:%s error:%s", invoiceId, err.Error())
	} else {
		g.Log().Info(ctx, "UserDiscountRollbackFromInvoice success invoiceId:%s", invoiceId)
	}
	err = payment.RollbackCreditPayment(context.Background(), one.MerchantId, one.InvoiceId)
	if err != nil {
		g.Log().Error(context.Background(), "InvoiceRollbackAllDiscountsFromInvoice RollbackCreditPayment invoiceId:%s error:%s", invoiceId, err.Error())
	} else {
		g.Log().Info(ctx, "RollbackCreditPayment success invoiceId:%s", invoiceId)
	}

	return nil
}
