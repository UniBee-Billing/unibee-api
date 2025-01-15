package detail

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	detail2 "unibee/internal/logic/invoice/detail"
	"unibee/internal/query"
	"unibee/utility"
)

func GetPaymentDetail(ctx context.Context, merchantId uint64, paymentId string) *detail.PaymentDetail {
	one := query.GetPaymentByPaymentId(ctx, paymentId)
	utility.Assert(one != nil, "payment not found")
	utility.Assert(merchantId == one.MerchantId, "merchant not match")
	return &detail.PaymentDetail{
		User:    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		Payment: bean.SimplifyPayment(one),
		Gateway: detail.ConvertGatewayDetail(ctx, query.GetGatewayById(ctx, one.GatewayId)),
		Invoice: detail2.InvoiceDetail(ctx, one.InvoiceId),
	}
}
