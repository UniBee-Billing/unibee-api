package api

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	_interface "go-oversea-pay/internal/interface"
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
	ChannelInvalid  = PayChannelKeyEnum{-1, "无效支付"}
	ChannelGrab     = PayChannelKeyEnum{0, "Grab支付"}
	ChannelKlarna   = PayChannelKeyEnum{1, "Klarna支付"}
	ChannelEvonet   = PayChannelKeyEnum{2, "Evonet支付"}
	ChannelPaypal   = PayChannelKeyEnum{3, "Paypal支付"}
	ChannelStripe   = PayChannelKeyEnum{4, "Stripe支付"}
	ChannelBlank    = PayChannelKeyEnum{50, "0金额支付专用"}
	ChannelAutoTest = PayChannelKeyEnum{500, "自动化测试支付专用"}
)

type PayChannelProxy struct {
	PaymentChannel *entity.MerchantGateway
}

func (p PayChannelProxy) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
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

	res, err = p.getRemoteChannel().DoRemoteChannelUserPaymentMethodListQuery(ctx, payChannel, userId)

	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantGateway, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
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

	res, err = p.getRemoteChannel().DoRemoteChannelUserCreate(ctx, payChannel, user)

	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelUserCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) getRemoteChannel() (channelService _interface.RemotePayChannelInterface) {
	utility.Assert(p.PaymentChannel != nil, "channel is not set")
	if p.PaymentChannel.EnumKey == ChannelEvonet.Code {
		return &Evonet{}
	} else if p.PaymentChannel.EnumKey == ChannelPaypal.Code {
		return &Paypal{}
	} else if p.PaymentChannel.EnumKey == ChannelStripe.Code {
		return &Stripe{}
	} else if p.PaymentChannel.EnumKey == ChannelBlank.Code {
		return &Blank{}
	} else if p.PaymentChannel.EnumKey == ChannelAutoTest.Code {
		return &AutoTest{}
	} else {
		return &Invalid{}
	}
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

func (p PayChannelProxy) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantGateway) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelUserDetailQuery(ctx, payChannel, userId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelUserDetailQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantGateway, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantGateway, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantGateway, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantGateway, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelSubscriptionEndTrial(ctx, plan, planChannel, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelSubscriptionNewTrialEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
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

func (p PayChannelProxy) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
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

func (p PayChannelProxy) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
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

func (p PayChannelProxy) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
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

func (p PayChannelProxy) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantGateway, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelPaymentList(ctx, payChannel, listReq)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPaymentList cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelPaymentDetail(ctx, payChannel, channelPaymentId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelPaymentDetail cost：%s \n", time.Now().Sub(startTime))
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

func (p PayChannelProxy) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelRefundDetail(ctx, payChannel, channelRefundId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelRefundDetail cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p PayChannelProxy) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteChannel().DoRemoteChannelRefundList(ctx, payChannel, channelPaymentId)
	glog.Infof(ctx, "MeasureChannelFunction:DoRemoteChannelRefundList cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func printChannelPanic(ctx context.Context, err error) {
	var requestId = "init"
	if _interface.BizCtx().Get(ctx) != nil {
		requestId = _interface.BizCtx().Get(ctx).RequestId
	}
	g.Log().Errorf(ctx, "ChannelException panic requestId:%s error:%s", requestId, err.Error())
}
