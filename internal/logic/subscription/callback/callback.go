package callback

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/email"
	"unibee/internal/logic/subscription/handler"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPaymentCallback struct {
}

func (s SubscriptionPaymentCallback) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
}

func (s SubscriptionPaymentCallback) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me
	panic("implement me")
}

func (s SubscriptionPaymentCallback) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	if payment.TotalAmount <= payment.RefundAmount {
		err := discount.UserDiscountRollbackFromPayment(ctx, payment.PaymentId)
		if err != nil {
			fmt.Printf("UserDiscountRollbackFromPayment error:%s", err.Error())
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me
	panic("implement me")
}

func (s SubscriptionPaymentCallback) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me
	panic("implement me")
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
				merchant := query.GetMerchantById(ctx, sub.MerchantId)
				if merchant != nil {
					err := email.SendTemplateEmail(ctx, merchant.Id, oneUser.Email, oneUser.TimeZone, oneUser.Language, email.TemplateSubscriptionNeedAuthorized, "", &email.TemplateVariable{
						UserName:              oneUser.FirstName + " " + oneUser.LastName,
						MerchantProductName:   plan.PlanName,
						MerchantCustomerEmail: merchant.Email,
						MerchantName:          query.GetMerchantCountryConfigName(ctx, payment.MerchantId, oneUser.CountryCode),
						PaymentAmount:         utility.ConvertCentToDollarStr(invoice.TotalAmount, invoice.Currency),
						Currency:              strings.ToUpper(invoice.Currency),
						PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
					})
					if err != nil {
						g.Log().Errorf(ctx, "SendTemplateEmail PaymentNeedAuthorisedCallback err:%s", err.Error())
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
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(invoice != nil, "payment of BizTypeSubscription invalid invoice")
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "sub not found")
			gateway := query.GetGatewayById(ctx, payment.GatewayId)
			utility.Assert(gateway != nil, "gateway not found")
			_ = handler.UpdateSubscriptionDefaultPaymentMethod(ctx, sub.SubscriptionId, payment.GatewayPaymentMethod)
			pendingUpdate := query.GetSubscriptionPendingUpdateByInvoiceId(ctx, invoice.InvoiceId)
			if pendingUpdate != nil {
				// PendingUpdate
				_, err := handler.HandlePendingUpdatePaymentSuccess(ctx, sub, pendingUpdate.PendingUpdateId, invoice)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_Update error:%s", err.Error())
				}
			} else if strings.Compare(payment.BillingReason, "SubscriptionCreate") == 0 {
				// SubscriptionCreate
				err := handler.HandleSubscriptionFirstInvoicePaid(ctx, sub, invoice)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_SubscriptionCreate error:%s", err.Error())
				}
			} else if strings.Compare(sub.LatestInvoiceId, invoice.InvoiceId) == 0 ||
				(gateway.GatewayType == consts.GatewayTypeCrypto && strings.Compare(payment.BillingReason, "SubscriptionRenew") == 0) ||
				(gateway.GatewayType == consts.GatewayTypeCrypto && strings.Compare(payment.BillingReason, "SubscriptionCycle") == 0) {
				// SubscriptionCycle or SubscriptionRenew
				err := handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, invoice)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentSuccessCallback_Finish_SubscriptionCycle error:%s", err.Error())
				}
			} else {
				g.Log().Infof(ctx, "PaymentSuccessCallback_Finish Miss Match Subscription Action:%s", payment.PaymentId)
			}
			if invoice != nil && len(invoice.CreateFrom) > 0 && invoice.CreateFrom == "AutoRenew" &&
				utility.TryLock(ctx, fmt.Sprintf("PaymentSuccessCallback_%s", invoice.InvoiceId), 60) {
				_, _ = redismq.Send(&redismq.Message{
					Topic:      redismq2.TopicSubscriptionAutoRenewSuccess.Topic,
					Tag:        redismq2.TopicSubscriptionAutoRenewSuccess.Tag,
					Body:       payment.SubscriptionId,
					CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
				})
			}
		}
	}
}

func (s SubscriptionPaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByInvoiceId(ctx, invoice.InvoiceId)
			if pendingSubUpdate != nil {
				_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.PendingUpdateId)
				if err != nil {
					utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
				}
			}
		}
	}
	err := discount.UserDiscountRollbackFromPayment(ctx, payment.PaymentId)
	if err != nil {
		fmt.Printf("UserDiscountRollbackFromPayment error:%s", err.Error())
	}
}

func (s SubscriptionPaymentCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	if consts.ProrationUsingUniBeeCompute {
		if payment.BizType == consts.BizTypeSubscription {
			utility.Assert(len(payment.SubscriptionId) > 0, "payment sub biz_type contain no sub_id")
			sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
			utility.Assert(sub != nil, "payment sub not found")
			if invoice != nil {
				pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByInvoiceId(ctx, invoice.InvoiceId)
				if pendingSubUpdate != nil {
					_, err := handler.HandlePendingUpdatePaymentFailure(ctx, pendingSubUpdate.PendingUpdateId)
					if err != nil {
						utility.AssertError(err, "PaymentFailureCallback_PaymentFailureForPendingUpdate")
					}
				}
			}
		}
	}
	err := discount.UserDiscountRollbackFromPayment(ctx, payment.PaymentId)
	if err != nil {
		fmt.Printf("UserDiscountRollbackFromPayment error:%s", err.Error())
	}
}
