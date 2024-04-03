package detail

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/query"
	"unibee/utility"
)

func GetRefundDetail(ctx context.Context, merchantId uint64, refundId string) *detail.RefundDetail {
	one := query.GetRefundByRefundId(ctx, refundId)
	utility.Assert(merchantId == one.MerchantId, "merchant not match")
	if one != nil {
		return &detail.RefundDetail{
			User:    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
			Payment: bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, one.PaymentId)),
			Refund:  bean.SimplifyRefund(one),
		}
	}
	return nil
}
