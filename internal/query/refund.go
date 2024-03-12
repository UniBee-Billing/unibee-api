package query

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/payment"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func GetRefundByRefundId(ctx context.Context, refundId string) (one *entity.Refund) {
	if len(refundId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns().RefundId, refundId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetRefundByGatewayRefundId(ctx context.Context, gatewayRefundId string) (one *entity.Refund) {
	if len(gatewayRefundId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns().GatewayRefundId, gatewayRefundId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetRefundDetail(ctx context.Context, merchantId uint64, refundId string) *payment.RefundDetail {
	one := GetRefundByRefundId(ctx, refundId)
	utility.Assert(merchantId == one.MerchantId, "merchant not match")
	if one != nil {
		return &payment.RefundDetail{
			User:    bean.SimplifyUserAccount(GetUserAccountById(ctx, uint64(one.UserId))),
			Payment: bean.SimplifyPayment(GetPaymentByPaymentId(ctx, one.PaymentId)),
			Refund:  bean.SimplifyRefund(one),
		}
	}
	return nil
}
