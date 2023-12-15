package handler

import (
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type HandlePayReq struct {
	MerchantOrderNo string
	ChannelPayId    string
	ChannelTradeNo  string
	PayFee          int64
	PayStatusEnum   consts.PayStatusEnum
	PaidTime        *gtime.Time
	ReceiptFee      int64
	Reason          string
	PaymentMethod   string
}

func HandlePayAuthorized(pay *entity.OverseaPay) (err error) {
	return nil
}

func HandlePayFailure(req *HandlePayReq) (err error) {
	return nil
}

func HandlePaySuccess(req *HandlePayReq) (err error) {
	return nil
}
