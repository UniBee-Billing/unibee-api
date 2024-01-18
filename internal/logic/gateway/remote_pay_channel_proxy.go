package gateway

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	_interface "go-oversea-pay/internal/interface"
	out2 "go-oversea-pay/internal/logic/gateway/out"
	"go-oversea-pay/internal/logic/gateway/out/evonet"
	"go-oversea-pay/internal/logic/gateway/out/paypal"
	"go-oversea-pay/internal/logic/gateway/out/stripe"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"time"
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

//func channelFunctionAop(ctx context.Context, fn func(), res interface{}, err error) {
//	startTime := time.Now()
//
//	// 调用目标函数
//	res, err = fn()
//
//	endTime := time.Now()
//	fmt.Printf("执行完成，耗时：%v\n", endTime.Sub(startTime))
//}

func (p PayChannelProxy) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()

	res, err = p.getRemoteChannel().DoRemoteChannelMerchantBalancesQuery(ctx, payChannel)

	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelMerchantBalancesQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelUserBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, customerId string) (res *ro.ChannelUserBalanceQueryInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelUserBalancesQuery(ctx, payChannel, customerId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelUserBalancesQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelInvoiceCreateAndPay(ctx, payChannel, createInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelInvoiceCreateAndPay cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.OverseaPayChannel, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelInvoiceCancel(ctx, payChannel, cancelInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelInvoiceCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelInvoicePay(ctx, payChannel, payInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelInvoicePay cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelInvoiceDetails(ctx, payChannel, channelInvoiceId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelInvoiceDetails cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionNewTrialEnd(ctx, plan, planChannel, subscription, newTrialEnd)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionNewTrialEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionUpdateProrationPreview(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionUpdateProrationPreview cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
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
		return &out2.Blank{}
	} else if p.channel.EnumKey == AutoTest.Code {
		return &out2.AutoTest{}
	} else {
		return &out2.Invalid{}
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
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionCreate(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionCancel(ctx, subscriptionCancelInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx, plan, planChannel, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionCancelAtPeriodEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx, plan, planChannel, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionUpdate(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionUpdate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionDetails cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	err = p.getRemoteChannel().DoRemoteChannelCheckAndSetupWebhook(ctx, payChannel)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelCheckAndSetupWebhook cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p PayChannelProxy) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	err = p.getRemoteChannel().DoRemoteChannelPlanActive(ctx, plan, planChannel)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPlanActive cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p PayChannelProxy) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	err = p.getRemoteChannel().DoRemoteChannelPlanDeactivate(ctx, plan, planChannel)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPlanDeactivate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p PayChannelProxy) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelProductCreate(ctx, plan, planChannel)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelProductCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPlanCreateAndActivate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	startTime := time.Now()
	p.getRemoteChannel().DoRemoteChannelWebhook(r, payChannel)
	glog.Infof(r.Context(), "MeasureChannelFunction:DoRemoteChannelWebhook cost：%s \n", time.Now().Sub(startTime))
}
func (p PayChannelProxy) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelRedirect(r, payChannel)
	glog.Infof(r.Context(), "MeasureChannelFunction:DoRemoteChannelRedirect cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelPayment(ctx, createPayContext)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPayment cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelCapture(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelCapture cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelCancel(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.OutPayRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelPayStatusCheck(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPayStatusCheck cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelRefund(ctx, pay, refund)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelRefund cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func printChannelPanic(ctx context.Context, err error) {
	g.Log().Errorf(ctx, "ChannelException panic requestId:%s error:%s", _interface.BizCtx().Get(ctx).RequestId, err.Error())
}

func (p PayChannelProxy) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelRefundStatusCheck(ctx, pay, refund)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelRefundStatusCheck cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}
