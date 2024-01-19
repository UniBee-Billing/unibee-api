package sync

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/payment/handler"
	"go-oversea-pay/internal/query"
)

func PaymentBackgroundSync(channelId int64, channelPaymentId string) {
	if channelId <= 0 {
		return
	}
	if len(channelPaymentId) == 0 {
		return
	}
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				var err error
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "PaymentBackgroundSyncFromInvoice Background panic error:%s\n", err.Error())
				return
			}
		}()
		backgroundCtx := context.Background()
		payChannel := query.GetPayChannelById(backgroundCtx, channelId)
		details, err := gateway.GetPayChannelServiceProvider(backgroundCtx, channelId).DoRemoteChannelPaymentDetail(backgroundCtx, payChannel, channelPaymentId)
		if err == nil {
			err := handler.CreateOrUpdatePaymentByDetail(backgroundCtx, details, details.ChannelPaymentId)
			if err != nil {
				fmt.Printf("SubscriptionDetail Background Fetch error%s", err)
				return
			}
		}
	}()
}
