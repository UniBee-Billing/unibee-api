package webhook

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"time"
)

type GatewayKeyEnum struct {
	Code int64
	Desc string
}

var (
	GatewayInvalid  = GatewayKeyEnum{-1, "无效支付"}
	GatewayGrab     = GatewayKeyEnum{0, "Grab支付"}
	GatewayKlarna   = GatewayKeyEnum{1, "Klarna支付"}
	GatewayEvonet   = GatewayKeyEnum{2, "Evonet支付"}
	GatewayPaypal   = GatewayKeyEnum{3, "Paypal支付"}
	GatewayStripe   = GatewayKeyEnum{4, "Stripe支付"}
	GatewayBlank    = GatewayKeyEnum{50, "0金额支付专用"}
	GatewayAutoTest = GatewayKeyEnum{500, "自动化测试支付专用"}
)

type GatewayWebhookProxy struct {
	PaymentChannel *entity.MerchantGateway
}

func (p GatewayWebhookProxy) getRemoteGateway() (one _interface.GatewayWebhookInterface) {
	utility.Assert(p.PaymentChannel != nil, "gateway is not set")
	if p.PaymentChannel.EnumKey == GatewayEvonet.Code {
		return &EvonetWebhook{}
	} else if p.PaymentChannel.EnumKey == GatewayPaypal.Code {
		return &PaypalWebhook{}
	} else if p.PaymentChannel.EnumKey == GatewayStripe.Code {
		return &StripeWebhook{}
	} else if p.PaymentChannel.EnumKey == GatewayBlank.Code {
		return &BlankWebhook{}
	} else if p.PaymentChannel.EnumKey == GatewayAutoTest.Code {
		return &AutoTestWebhook{}
	} else {
		return &InvalidWebhook{}
	}
}

func (p GatewayWebhookProxy) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
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
	err = p.getRemoteGateway().GatewayCheckAndSetupWebhook(ctx, gateway)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayCheckAndSetupWebhook cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return err
}

func (p GatewayWebhookProxy) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	startTime := time.Now()
	p.getRemoteGateway().GatewayWebhook(r, gateway)
	glog.Infof(r.Context(), "MeasureChannelFunction:GatewayWebhook cost：%s \n", time.Now().Sub(startTime))
}
func (p GatewayWebhookProxy) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *ro.GatewayRedirectInternalResp, err error) {
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayRedirect(r, gateway)
	glog.Infof(r.Context(), "MeasureChannelFunction:GatewayRedirect cost：%s \n", time.Now().Sub(startTime))
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
