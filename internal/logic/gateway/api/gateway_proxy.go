package api

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"time"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type GatewayKeyEnum struct {
	Name string
	Desc string
}

var (
	GatewayInvalid  = GatewayKeyEnum{"invalid", "Invalid Gateway"}
	GatewayPaypal   = GatewayKeyEnum{"paypal", "Paypal"}
	GatewayStripe   = GatewayKeyEnum{"stripe", "Stripe"}
	GatewayBlank    = GatewayKeyEnum{"0", "0 Payment Gateway"}
	GatewayAutoTest = GatewayKeyEnum{"autotest", "Auto Test"}
)

type GatewayProxy struct {
	Gateway     *entity.MerchantGateway
	GatewayName string
}

func (p GatewayProxy) getRemoteGateway() (one _interface.GatewayInterface) {
	utility.Assert(len(p.GatewayName) > 0, "gateway is not set")
	if p.GatewayName == GatewayPaypal.Name {
		return &Paypal{}
	} else if p.GatewayName == GatewayStripe.Name {
		return &Stripe{}
	} else if p.GatewayName == GatewayBlank.Name {
		return &Blank{}
	} else if p.GatewayName == GatewayAutoTest.Name {
		return &AutoTest{}
	} else {
		return &Invalid{}
	}
}

func (p GatewayProxy) GatewayTest(ctx context.Context, key string, secret string) (err error) {
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
	err = p.getRemoteGateway().GatewayTest(ctx, key, secret)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayTest cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p GatewayProxy) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserAttachPaymentMethodInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayUserAttachPaymentMethodQuery(ctx, gateway, userId, gatewayPaymentMethod)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserDeAttachPaymentMethodInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayUserDeAttachPaymentMethodQuery(ctx, gateway, userId, gatewayPaymentMethod)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
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

func (p GatewayProxy) GatewayNewPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
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
	res, err = p.getRemoteGateway().GatewayNewPayment(ctx, createPayContext)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayNewPayment cost：%s \n", time.Now().Sub(startTime))
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
