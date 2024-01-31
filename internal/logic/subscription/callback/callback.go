package callback

import (
	"context"
	"fmt"
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
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpgrade := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			pendingSubDowngrade := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpgrade != nil && strings.Compare(payment.BillingReason, "SubscriptionUpgrade") == 0 {
				utility.Assert(strings.Compare(pendingSubUpgrade.SubscriptionId, payment.SubscriptionId) == 0, "payment sub_id not match pendingUpdate sub_id")
				// Upgrade
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpgrade)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_Upgrade")
				}
			} else if pendingSubDowngrade != nil && strings.Compare(payment.BillingReason, "SubscriptionDowngrade") == 0 {
				utility.Assert(strings.Compare(pendingSubUpgrade.ChannelUpdateId, payment.PaymentId) == 0, "paymentId not match pendingUpdate ChannelUpdateId")
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
			} else {
				//todo mark
				utility.Assert(false, fmt.Sprintf("PaymentSuccessCallback_Finish Miss Match Payment:%s", payment.PaymentId))
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				_, err := handler.PaymentFailureForPendingUpdate(ctx, pendingSubUpdate)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
			// billing cycle use cronjob check active status as contain other processing payment
		}
	}
}
