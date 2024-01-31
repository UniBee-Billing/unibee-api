package callback

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/order/calllback"
	"go-oversea-pay/internal/logic/subscription/callback"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"time"
)

type proxy struct {
	BizType int
}

func printChannelPanic(ctx context.Context, err error) {
	g.Log().Errorf(ctx, "CallbackException panic error:%s", err.Error())
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

		p.GetCallbackImpl().PaymentCreateCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentFailureCallback cost：%s \n", time.Now().Sub(startTime))
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

		p.GetCallbackImpl().PaymentSuccessCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentFailureCallback cost：%s \n", time.Now().Sub(startTime))
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

		p.GetCallbackImpl().PaymentCancelCallback(backgroundCtx, payment, invoice)

		glog.Infof(backgroundCtx, "MeasurePaymentCallbackFunction:PaymentFailureCallback cost：%s \n", time.Now().Sub(startTime))
	}()
}

func (p proxy) GetCallbackImpl() (channelService _interface.PaymentBizCallbackInterface) {
	utility.Assert(p.BizType >= 0, "bizType is not set")
	if p.BizType == consts.BIZ_TYPE_ONE_TIME {
		return &calllback.MerchantOneTimePaymentCallback{}
	} else if p.BizType == consts.BIZ_TYPE_SUBSCRIPTION {
		return &callback.SubscriptionPaymentCallback{}
	} else {
		return &Invalid{}
	}
}
