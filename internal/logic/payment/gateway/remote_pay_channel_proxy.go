package gateway

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/gateway/out"
	"go-oversea-pay/internal/logic/payment/gateway/out/evonet"
	"go-oversea-pay/internal/logic/payment/gateway/out/paypal"
	"go-oversea-pay/internal/logic/payment/gateway/out/stripe"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

type PayChannelKeyEnum struct {
	Code int64
	Desc string
}

var (
	Invalid  = PayChannelKeyEnum{-1, "无效支付"}
	Grab     = PayChannelKeyEnum{0, "Grab支付"}
	Klarna   = PayChannelKeyEnum{1, "Klarna支付"}
	Evonet   = PayChannelKeyEnum{2, "Evonet支付"}
	Paypal   = PayChannelKeyEnum{3, "Paypal支付"}
	Stripe   = PayChannelKeyEnum{4, "Stripe支付"}
	Blank    = PayChannelKeyEnum{50, "0金额支付专用"}
	AutoTest = PayChannelKeyEnum{500, "自动化测试支付专用"}
)

type PayChannelProxy struct {
	channel *entity.OverseaPayChannel
}

func (p PayChannelProxy) DoRemoteChannelCustomerBalanceQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, customerId string) (res *ro.ChannelCustomerBalanceQueryInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelCustomerBalanceQuery(ctx, payChannel, customerId)
}

func (p PayChannelProxy) DoRemoteChannelInvoiceCreate(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelInvoiceCreate(ctx, payChannel, createInvoiceInternalReq)
}

func (p PayChannelProxy) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelInvoicePay(ctx, payChannel, payInvoiceInternalReq)
}

func (p PayChannelProxy) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelInvoiceDetails(ctx, payChannel, channelInvoiceId)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionNewTrailEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrailEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionNewTrailEnd(ctx, plan, planChannel, subscription, newTrailEnd)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionUpdateProrationPreview(ctx, subscriptionRo)
}

func (p PayChannelProxy) getRemoteChannel() (channelService RemotePayChannelInterface) {
	utility.Assert(p.channel != nil, "channel is not set")
	if p.channel.EnumKey == Evonet.Code {
		return &evonet.Evonet{}
	} else if p.channel.EnumKey == Paypal.Code {
		return &paypal.Paypal{}
	} else if p.channel.EnumKey == Stripe.Code {
		return &stripe.Stripe{}
	} else if p.channel.EnumKey == Blank.Code {
		return &out.Blank{}
	} else if p.channel.EnumKey == AutoTest.Code {
		return &out.AutoTest{}
	} else {
		return &out.Invalid{}
	}
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionCreate(ctx, subscriptionRo)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionCancel(ctx, subscriptionCancelInternalReq)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx, plan, planChannel, subscription)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx, plan, planChannel, subscription)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionUpdate(ctx, subscriptionRo)
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, subscription)
}

func (p PayChannelProxy) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelCheckAndSetupWebhook(ctx, payChannel)
}

func (p PayChannelProxy) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelPlanActive(ctx, plan, planChannel)
}

func (p PayChannelProxy) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelPlanDeactivate(ctx, plan, planChannel)
}

func (p PayChannelProxy) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelProductCreate(ctx, plan, planChannel)
}

func (p PayChannelProxy) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
}

func (p PayChannelProxy) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	p.getRemoteChannel().DoRemoteChannelWebhook(r, payChannel)
}

func (p PayChannelProxy) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	return p.getRemoteChannel().DoRemoteChannelRedirect(r, payChannel)
}

func (p PayChannelProxy) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelPayment(ctx, createPayContext)
}

func (p PayChannelProxy) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelCapture(ctx, pay)
}

func (p PayChannelProxy) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelCancel(ctx, pay)
}

func (p PayChannelProxy) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelPayStatusCheck(ctx, pay)
}

func (p PayChannelProxy) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()
	return p.getRemoteChannel().DoRemoteChannelRefund(ctx, pay, refund)
}

func (p PayChannelProxy) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "ChannelException error:%s", err.Error())
			return
		}
	}()

	return p.getRemoteChannel().DoRemoteChannelRefundStatusCheck(ctx, pay, refund)
}
