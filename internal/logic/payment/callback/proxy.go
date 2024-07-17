package callback

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"time"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	payment2 "unibee/internal/consumer/webhook/payment"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/callback/onetime"
	"unibee/internal/logic/subscription/callback"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type proxy struct {
	BizType int
}

func printChannelPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err)
	}
}

func (p proxy) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendRefundWebhookBackground(refund.RefundId, event.UNIBEE_WEBHOOK_EVENT_REFUND_CREATED)
		p.GetCallbackImpl().PaymentRefundCreateCallback(backgroundCtx, payment, refund)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentRefundCreateCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendRefundWebhookBackground(refund.RefundId, event.UNIBEE_WEBHOOK_EVENT_REFUND_SUCCESS)
		p.GetCallbackImpl().PaymentRefundSuccessCallback(backgroundCtx, payment, refund)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentRefundSuccessCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendRefundWebhookBackground(refund.RefundId, event.UNIBEE_WEBHOOK_EVENT_REFUND_FAILURE)
		p.GetCallbackImpl().PaymentRefundFailureCallback(backgroundCtx, payment, refund)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentRefundFailureCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendRefundWebhookBackground(refund.RefundId, event.UNIBEE_WEBHOOK_EVENT_REFUND_CANCELLED)
		p.GetCallbackImpl().PaymentRefundCancelCallback(backgroundCtx, payment, refund)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentRefundCancelCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendRefundWebhookBackground(refund.RefundId, event.UNIBEE_WEBHOOK_EVENT_REFUND_REVERSED)
		p.GetCallbackImpl().PaymentRefundReverseCallback(backgroundCtx, payment, refund)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentRefundReverseCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()
		p.GetCallbackImpl().PaymentNeedAuthorisedCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentNeedAuthorisedCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendPaymentWebhookBackground(payment.PaymentId, event.UNIBEE_WEBHOOK_EVENT_PAYMENT_CREATED)
		p.GetCallbackImpl().PaymentCreateCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentCreateCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	go func() {
		var err error
		backgroundCtx := context.Background()
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendPaymentWebhookBackground(payment.PaymentId, event.UNIBEE_WEBHOOK_EVENT_PAYMENT_SUCCESS)
		p.GetCallbackImpl().PaymentSuccessCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentSuccessCallback cost：%s \n", time.Now().Sub(startTime))
	}()

	return
}

func (p proxy) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendPaymentWebhookBackground(payment.PaymentId, event.UNIBEE_WEBHOOK_EVENT_PAYMENT_FAILURE)
		p.GetCallbackImpl().PaymentFailureCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentFailureCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	go func() {
		var err error
		backgroundCtx := context.Background()
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		startTime := time.Now()

		payment2.SendPaymentWebhookBackground(payment.PaymentId, event.UNIBEE_WEBHOOK_EVENT_PAYMENT_CANCEL)
		p.GetCallbackImpl().PaymentCancelCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentCancelCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) GetCallbackImpl() (one _interface.PaymentBizCallbackInterface) {
	utility.Assert(p.BizType >= 0, "bizType is not set")
	if p.BizType == consts.BizTypeOneTime {
		return &onetime.Onetime{}
	} else if p.BizType == consts.BizTypeSubscription {
		return &callback.SubscriptionPaymentCallback{}
	} else {
		return &Invalid{}
	}
}
