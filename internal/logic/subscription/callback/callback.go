package callback

import (
	"context"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionPaymentCallback struct {
}

func (s SubscriptionPaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment) {

	if consts.ProrationUsingUniBeeCompute {
		// better use redis mq to trace payment
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			pendingSubUpgrade := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			pendingSubDowngrade := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			utility.Assert(strings.Compare(pendingSubUpgrade.SubscriptionId, payment.SubscriptionId) == 0, "payment sub_id not match pendingUpdate sub_id")
			if pendingSubUpgrade != nil && strings.Compare(payment.BillingReason, "SubscriptionUpgrade") == 0 {
				// Upgrade
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpgrade)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_Upgrade")
				}
			} else if pendingSubDowngrade != nil && strings.Compare(payment.BillingReason, "SubscriptionDowngrade") == 0 {
				// Downgrade
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubDowngrade)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_Downgrade")
				}

				err = handler.FinishNextBillingCycleForSubscription(ctx, sub, payment)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_Downgrade")
				}
			} else if strings.Compare(payment.BillingReason, "SubscriptionCycle") == 0 {
				// SubscriptionCycle
				err := handler.FinishNextBillingCycleForSubscription(ctx, sub, payment)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_SubscriptionCycle")
				}
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, Payment *entity.Payment) {

}
