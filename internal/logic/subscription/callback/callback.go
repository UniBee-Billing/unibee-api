package callback

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/logic/email"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionPaymentCallback struct {
}

func (s SubscriptionPaymentCallback) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "PaymentNeedAuthorisedCallback sub not found:"+payment.PaymentId)
			oneUser := query.GetUserAccountById(ctx, uint64(sub.UserId))
			plan := query.GetPlanById(ctx, sub.PlanId)
			utility.Assert(plan != nil, "PaymentNeedAuthorisedCallback plan not found:"+sub.SubscriptionId)
			if oneUser != nil {
				merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
				if merchant != nil {
					err := email.SendTemplateEmail(ctx, merchant.Id, oneUser.Email, oneUser.TimeZone, email.TemplateSubscriptionNeedAuthorized, "", &email.TemplateVariable{
						UserName:            oneUser.FirstName + " " + oneUser.LastName,
						MerchantProductName: plan.PlanName,
						MerchantCustomEmail: merchant.Email,
						MerchantName:        merchant.Name,
						PaymentAmount:       utility.ConvertCentToDollarStr(invoice.TotalAmount, invoice.Currency),
						Currency:            strings.ToUpper(invoice.Currency),
						PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
					})
					if err != nil {
						g.Log().Errorf(ctx, "PaymentNeedAuthorisedCallback SendTemplateEmail err:%s", err.Error())
					}
				}
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			user.UpdateUserDefaultSubscription(ctx, payment.UserId, payment.SubscriptionId)
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(invoice != nil, "payment of BizTypeSubscription invalid invoice")
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpgrade := query.GetSubscriptionUpgradePendingUpdateByGatewayUpdateId(ctx, payment.PaymentId)
			pendingSubDowngrade := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpgrade != nil && strings.Compare(payment.BillingReason, "SubscriptionUpgrade") == 0 {
				if strings.Compare(pendingSubUpgrade.SubscriptionId, payment.SubscriptionId) == 0 &&
					pendingSubUpgrade.Status == consts.PendingSubStatusCreate {
					// Upgrade
					_, err := handler.HandlePendingUpdatePaymentSuccess(ctx, sub, pendingSubUpgrade.UpdateSubscriptionId, invoice)
					if err != nil {
						g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_Upgrade error:%s", err.Error())
					}
				}
			} else if pendingSubDowngrade != nil && strings.Compare(payment.BillingReason, "SubscriptionDowngrade") == 0 {
				if strings.Compare(pendingSubUpgrade.SubscriptionId, payment.SubscriptionId) == 0 &&
					pendingSubUpgrade.Status == consts.PendingSubStatusCreate {
					// Downgrade
					_, err := handler.HandlePendingUpdatePaymentSuccess(ctx, sub, pendingSubDowngrade.UpdateSubscriptionId, invoice)
					if err != nil {
						g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_Downgrade error:%s", err.Error())
					} else {
						err = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, payment)
						if err != nil {
							g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_Downgrade error:%s", err.Error())
						}
					}
				}
			} else if strings.Compare(payment.BillingReason, "SubscriptionCycle") == 0 && sub.Amount == payment.TotalAmount && strings.Compare(sub.LatestInvoiceId, invoice.InvoiceId) == 0 {
				// SubscriptionCycle
				err := handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, payment)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_SubscriptionCycle error:%s", err.Error())
				}
			} else if strings.Compare(payment.BillingReason, "SubscriptionCreate") == 0 {
				// SubscriptionCycle
				err := handler.HandleSubscriptionFirstPaymentSuccess(ctx, sub, payment)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_SubscriptionCreate error:%s", err.Error())
				}
			} else {
				utility.Assert(false, fmt.Sprintf("PaymentSuccessCallback_Finish Miss Match Subscription Action:%s", payment.PaymentId))
			}
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicSubscriptionPaymentSuccess.Topic,
				Tag:   redismq2.TopicSubscriptionPaymentSuccess.Tag,
				Body:  payment.SubscriptionId,
			})
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByGatewayUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.UpdateSubscriptionId)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByGatewayUpdateId(ctx, payment.PaymentId)
			if pendingSubUpdate != nil {
				_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.UpdateSubscriptionId)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
		}
	}
}
