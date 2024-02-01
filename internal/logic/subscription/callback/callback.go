package callback

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionPaymentCallback struct {
}

func (s SubscriptionPaymentCallback) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "PaymentNeedAuthorisedCallback sub not found:"+payment.PaymentId)
			user := query.GetUserAccountById(ctx, uint64(sub.UserId))
			plan := query.GetPlanById(ctx, sub.PlanId)
			utility.Assert(plan != nil, "PaymentNeedAuthorisedCallback plan not found:"+sub.SubscriptionId)
			if user != nil {
				merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
				if merchant != nil {
					err := email.SendTemplateEmail(ctx, merchant.Id, user.Email, email.TemplateSubscriptionNeedAuthorized, "", &email.TemplateVariable{
						UserName:            user.UserName,
						MerchantProductName: plan.ChannelProductName,
						MerchantCustomEmail: merchant.Email,
						MerchantName:        merchant.Name,
						PaymentAmount:       utility.ConvertCentToDollarStr(invoice.TotalAmount, invoice.Currency),
						Currency:            strings.ToUpper(invoice.Currency),
						PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd).Layout("2006-01-02"),
					})
					if err != nil {
						fmt.Printf("PaymentNeedAuthorisedCallback SendTemplateEmail err:%s", err.Error())
					}
				}
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		// better use redis mq to trace payment
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
				dao.Subscription.Columns().LatestInvoiceId: invoice.InvoiceId,
			}).Where(dao.Subscription.Columns().SubscriptionId, payment.SubscriptionId).OmitNil().Update()
			if err != nil {
				utility.AssertError(err, "PaymentCreateCallback")
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		// better use redis mq to trace payment
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpgrade := query.GetSubscriptionUpgradePendingUpdateByPendingUpdateId(ctx, payment.PaymentId)
			pendingSubDowngrade := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpgrade != nil && strings.Compare(payment.BillingReason, "SubscriptionUpgrade") == 0 {
				utility.Assert(strings.Compare(pendingSubUpgrade.SubscriptionId, payment.SubscriptionId) == 0, "payment sub_id not match pendingUpdate sub_id")
				utility.Assert(pendingSubUpgrade.Status == consts.PendingSubStatusCreate, "pendingUpdate has already finished or cancelled")
				// Upgrade
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpgrade.UpdateSubscriptionId)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_Upgrade")
				}
			} else if pendingSubDowngrade != nil && strings.Compare(payment.BillingReason, "SubscriptionDowngrade") == 0 {
				utility.Assert(strings.Compare(pendingSubUpgrade.ChannelUpdateId, payment.PaymentId) == 0, "paymentId not match pendingUpdate ChannelUpdateId")
				// Downgrade
				_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, pendingSubDowngrade.UpdateSubscriptionId)
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
			} else if strings.Compare(payment.BillingReason, "SubscriptionCreate") == 0 {
				// SubscriptionCycle
				err := handler.HandleSubscriptionCreatePaymentSuccess(ctx, sub, payment)
				if err != nil {
					utility.AssertError(err, "PaymentSuccessCallback_Finish_SubscriptionCreate")
				}
			} else {
				//todo mark
				utility.Assert(false, fmt.Sprintf("PaymentSuccessCallback_Finish Miss Match Payment:%s", payment.PaymentId))
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.UpdateSubscriptionId)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
			// billing cycle use cronjob check active status as contain other processing payment
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByChannelUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.UpdateSubscriptionId)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
			// billing cycle use cronjob check active status as contain other processing payment
		}
	}
}
