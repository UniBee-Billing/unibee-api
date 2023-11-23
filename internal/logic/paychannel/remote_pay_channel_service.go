package paychannel

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type RemotePayChannelService interface {
	DoRemoteChannelPayment(ctx context.Context, createPayContext interface{}) (res interface{}, err error)
	DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error)
	DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error)
	DoRemoteChannelStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error)
	DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res interface{}, err error)
}
