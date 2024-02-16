package api

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
	"time"
)

type GatewayKeyEnum struct {
	Code int64
	Desc string
}

var (
	GatewayInvalid  = GatewayKeyEnum{-1, "Invalid Gateway"}
	GatewayGrab     = GatewayKeyEnum{0, "Grab"}
	GatewayKlarna   = GatewayKeyEnum{1, "Klarna"}
	GatewayEvonet   = GatewayKeyEnum{2, "Evonet"}
	GatewayPaypal   = GatewayKeyEnum{3, "Paypal"}
	GatewayStripe   = GatewayKeyEnum{4, "Stripe"}
	GatewayBlank    = GatewayKeyEnum{50, "0 Payment"}
	GatewayAutoTest = GatewayKeyEnum{500, "For Automatic Test Payment"}
)

type GatewayProxy struct {
	Gateway *entity.MerchantGateway
}

func (p GatewayProxy) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
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

	res, err = p.getRemoteGateway().GatewayUserPaymentMethodListQuery(ctx, gateway, userId)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
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

	res, err = p.getRemoteGateway().GatewayUserCreate(ctx, gateway, user)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) getRemoteGateway() (one _interface.GatewayInterface) {
	utility.Assert(p.Gateway != nil, "gateway is not set")
	if p.Gateway.EnumKey == GatewayEvonet.Code {
		return &Evonet{}
	} else if p.Gateway.EnumKey == GatewayPaypal.Code {
		return &Paypal{}
	} else if p.Gateway.EnumKey == GatewayStripe.Code {
		return &Stripe{}
	} else if p.Gateway.EnumKey == GatewayBlank.Code {
		return &Blank{}
	} else if p.Gateway.EnumKey == GatewayAutoTest.Code {
		return &AutoTest{}
	} else {
		return &Invalid{}
	}
}

func (p GatewayProxy) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
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

	res, err = p.getRemoteGateway().GatewayMerchantBalancesQuery(ctx, gateway)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayMerchantBalancesQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayUserDetailQuery(ctx, gateway, userId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserDetailQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayInvoiceCreateAndPay(ctx context.Context, gateway *entity.MerchantGateway, createInvoiceInternalReq *ro.GatewayCreateInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayInvoiceCreateAndPay(ctx, gateway, createInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayInvoiceCreateAndPay cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayInvoiceCancel(ctx context.Context, gateway *entity.MerchantGateway, cancelInvoiceInternalReq *ro.GatewayCancelInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayInvoiceCancel(ctx, gateway, cancelInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayInvoiceCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayInvoicePay(ctx context.Context, gateway *entity.MerchantGateway, payInvoiceInternalReq *ro.GatewayPayInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayInvoicePay(ctx, gateway, payInvoiceInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayInvoicePay cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayInvoiceDetails(ctx context.Context, gateway *entity.MerchantGateway, gatewayInvoiceId string) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayInvoiceDetails(ctx, gateway, gatewayInvoiceId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayInvoiceDetails cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionEndTrial(ctx, plan, gatewayPlan, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionNewTrialEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionNewTrialEnd(ctx, plan, gatewayPlan, subscription, newTrialEnd)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionNewTrialEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionPreviewInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionUpdateProrationPreview(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionUpdateProrationPreview cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionCreate(ctx context.Context, subscriptionRo *ro.GatewayCreateSubscriptionInternalReq) (res *ro.GatewayCreateSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionCreate(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.GatewayCancelSubscriptionInternalReq) (res *ro.GatewayCancelSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionCancel(ctx, subscriptionCancelInternalReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelAtPeriodEndSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionCancelAtPeriodEnd(ctx, plan, gatewayPlan, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionCancelAtPeriodEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx, plan, gatewayPlan, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionCancelLastCancelAtPeriodEnd cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionUpdate(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionUpdate(ctx, subscriptionRo)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionUpdate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewaySubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewaySubscriptionDetails(ctx, plan, gatewayPlan, subscription)
	glog.Infof(ctx, "MeasureChannelFunction:GatewaySubscriptionDetails cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
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
	err = p.getRemoteGateway().GatewayPlanActive(ctx, plan, gatewayPlan)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPlanActive cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p GatewayProxy) GatewayPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
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
	err = p.getRemoteGateway().GatewayPlanDeactivate(ctx, plan, gatewayPlan)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPlanDeactivate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p GatewayProxy) GatewayProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreateProductInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayProductCreate(ctx, plan, gatewayPlan)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayProductCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreatePlanInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayPlanCreateAndActivate(ctx, plan, gatewayPlan)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPlanCreateAndActivate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayPayment(ctx, createPayContext)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPayment cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayCapture(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayCapture cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayCancel(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayPayStatusCheck(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPayStatusCheck cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayPaymentList(ctx, gateway, listReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPaymentList cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayPaymentDetail(ctx, gateway, gatewayPaymentId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPaymentDetail cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayRefund(ctx, pay, refund)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefund cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayRefundStatusCheck(ctx, pay, refund)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundStatusCheck cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayRefundDetail(ctx, gateway, gatewayRefundId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundDetail cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
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
	res, err = p.getRemoteGateway().GatewayRefundList(ctx, gateway, gatewayPaymentId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundList cost：%s \n", time.Now().Sub(startTime))
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
