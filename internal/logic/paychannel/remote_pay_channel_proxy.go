package paychannel

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/paychannel/evonet"
	"go-oversea-pay/internal/logic/paychannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

type PayChannelProxy struct {
	channel int // todo mark 应该使用 enum key
}

func (p PayChannelProxy) getRemoteChannel() (channelService RemotePayChannelService) {
	utility.Assert(p.channel > 0, "channel is not set")
	//目前只有一个渠道 todo mark
	return &evonet.Evonet{}
}

func (p PayChannelProxy) DoRemoteChannelWebhook(r *ghttp.Request) {
	p.getRemoteChannel().DoRemoteChannelWebhook(r)
}

func (p PayChannelProxy) DoRemoteChannelRedirect(r *ghttp.Request) {
	p.getRemoteChannel().DoRemoteChannelRedirect(r)
}

func (p PayChannelProxy) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			fmt.Printf("exception panic error:%s\n", err)
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
			fmt.Printf("exception panic error:%s\n", err)
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
			fmt.Printf("exception panic error:%s\n", err)
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
			fmt.Printf("exception panic error:%s\n", err)
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
			fmt.Printf("exception panic error:%s\n", err)
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
			fmt.Printf("exception panic error:%s\n", err)
			return
		}
	}()

	return p.getRemoteChannel().DoRemoteChannelRefundStatusCheck(ctx, pay, refund)
}
