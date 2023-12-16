package handler

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
)

type HandleRefundReq struct {
	ChargeRefundNo   string
	ChannelRefundNo  string
	RefundFee        int64
	RefundStatusEnum consts.RefundStatusEnum
	RefundTime       *gtime.Time
	Reason           string
}

func HandleRefundFailure(ctx context.Context, req *HandleRefundReq) (err error) {
	return nil
}

func HandleRefundSuccess(ctx context.Context, req *HandleRefundReq) (err error) {
	return nil
}
