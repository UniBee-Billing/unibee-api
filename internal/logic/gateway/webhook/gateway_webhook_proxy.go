package webhook

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"time"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

var GatewayWebhookNameMapping = map[string]_interface.GatewayWebhookInterface{
	"stripe":          &StripeWebhook{},
	"changelly":       &ChangellyWebhook{},
	"paypal":          &PaypalWebhook{},
	"invalid":         &InvalidWebhook{},
	"autotest_crypto": &AutoTestCryptoWebhook{},
	"autotest":        &AutoTestWebhook{},
}

type GatewayWebhookProxy struct {
	Gateway     *entity.MerchantGateway
	GatewayName string
}

func (p GatewayWebhookProxy) getRemoteGateway() (one _interface.GatewayWebhookInterface) {
	utility.Assert(len(p.GatewayName) > 0, "gateway is not set")
	one = GatewayWebhookNameMapping[p.GatewayName]
	utility.Assert(one != nil, "gateway not support:"+p.GatewayName)
	return
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
func (p GatewayWebhookProxy) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(r.Context(), err)
			return
		}
	}()
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
	if _interface.Context().Get(ctx) != nil {
		requestId = _interface.Context().Get(ctx).RequestId
	}
	g.Log().Errorf(ctx, "ChannelException panic requestId:%s error:%s", requestId, err.Error())
}
