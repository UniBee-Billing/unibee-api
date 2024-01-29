package webhook

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/channel/ro"
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

type PayChannelWebhookProxy struct {
	PaymentChannel *entity.OverseaPayChannel
}

func (p PayChannelWebhookProxy) getRemoteChannel() (channelService _interface.RemotePaymentChannelWebhookInterface) {
	utility.Assert(p.PaymentChannel != nil, "channel is not set")
	if p.PaymentChannel.EnumKey == ChannelEvonet.Code {
		return &EvonetWebhook{}
	} else if p.PaymentChannel.EnumKey == ChannelPaypal.Code {
		return &PaypalWebhook{}
	} else if p.PaymentChannel.EnumKey == ChannelStripe.Code {
		return &StripeWebhook{}
	} else if p.PaymentChannel.EnumKey == ChannelBlank.Code {
		return &BlankWebhook{}
	} else if p.PaymentChannel.EnumKey == ChannelAutoTest.Code {
		return &AutoTestWebhook{}
	} else {
		return &InvalidWebhook{}
	}
}

func (p PayChannelWebhookProxy) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
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

func (p PayChannelWebhookProxy) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	startTime := time.Now()
	p.getRemoteChannel().DoRemoteChannelWebhook(r, payChannel)
	glog.Infof(r.Context(), "MeasureChannelFunction:DoRemoteChannelWebhook cost：%s \n", time.Now().Sub(startTime))
}
func (p PayChannelWebhookProxy) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	startTime := time.Now()
	res, err = p.getRemoteChannel().DoRemoteChannelRedirect(r, payChannel)
	glog.Infof(r.Context(), "MeasureChannelFunction:DoRemoteChannelRedirect cost：%s \n", time.Now().Sub(startTime))
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
