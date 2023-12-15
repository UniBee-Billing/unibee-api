package handler

import (
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

func HandleRefundFailure(req *HandleRefundReq) (err error) {
	return nil
}

func HandleRefundSuccess(req *HandleRefundReq) (err error) {
	return nil
}
