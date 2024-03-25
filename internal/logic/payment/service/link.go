package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type LinkCheckRes struct {
	Message string
	Link    string
	Payment *entity.Payment
}

func LinkCheck(ctx context.Context, paymentId string, time int64) *LinkCheckRes {
	var res = &LinkCheckRes{
		Message: "",
		Link:    "",
		Payment: nil,
	}
	one := query.GetPaymentByPaymentId(ctx, paymentId)
	if one == nil {
		g.Log().Errorf(ctx, "LinkEntry payment not found paymentId: %s", paymentId)
		res.Message = "Payment Not Found"
		return res
	}
	res.Payment = one
	if one.Status == consts.PaymentCancelled {
		res.Message = "Payment Cancelled"
	} else if one.Status == consts.PaymentFailed {
		res.Message = "Payment Failure"
	} else if one.Status == consts.PaymentSuccess {
		res.Message = "Payment Already Success"
	} else if one.ExpireTime != 0 && one.ExpireTime < time {
		res.Message = "Payment Expired"
	} else if len(one.GatewayLink) > 0 {
		res.Link = one.GatewayLink
	} else if strings.Contains(one.Link, "unibee.top") {
		res.Message = "Server Error"
	} else {
		res.Link = one.Link // old version
	}
	return res
}
