package callback

import (
	"context"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

type SubscriptionPaymentCallback struct {
}

func (s SubscriptionPaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment) {

	if consts.ProrationUsingUniBeeCompute {
		// better use redis mq to trace payment
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			pendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
				sub := query.GetSubscriptionBySubscriptionId(ctx, pendingSubUpdate.SubscriptionId)
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpdate)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback")
				}
			}
			// billing cycle payment
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, Payment *entity.Payment) {

}
