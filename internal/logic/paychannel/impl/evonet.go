package impl

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Evonet struct{}

func (e Evonet) DoRemoteChannelPayment(ctx context.Context, createPayContext interface{}) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}
