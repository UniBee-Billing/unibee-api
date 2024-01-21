package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type TradeEventTypeEnum struct {
	Type int
	Desc string
}

var (
	Authorised         = TradeEventTypeEnum{2, "授权"}   //
	Cancelled          = TradeEventTypeEnum{3, "取消"}   //
	Expird             = TradeEventTypeEnum{4, "授权过期"} //
	ChargeBack         = TradeEventTypeEnum{5, "拒付"}
	ChargeBackReversed = TradeEventTypeEnum{6, "拒付资金退回"}
	CaptureFailed      = TradeEventTypeEnum{7, "请款失败"} //
	Error              = TradeEventTypeEnum{8, "异常错误"}
	Refunded           = TradeEventTypeEnum{9, "退款成功"}    //
	RefundFailed       = TradeEventTypeEnum{10, "退款失败"}   //
	RefundedReversed   = TradeEventTypeEnum{11, "退款资金退回"} //

	Refused = TradeEventTypeEnum{12, "授权失败"}

	SentForRefund = TradeEventTypeEnum{13, "退款请求成功"} //

	SentForSettle   = TradeEventTypeEnum{14, "请款请求成功"}
	Settled         = TradeEventTypeEnum{15, "扣款成功"}
	SettledReversed = TradeEventTypeEnum{16, "扣款资金退回"}
)

func SaveTimeLine(ctx context.Context, overseaPayEvent entity.PaymentEvent) {
	_, err := dao.PaymentEvent.Ctx(ctx).Data(overseaPayEvent).OmitNil().Insert(overseaPayEvent)
	if err != nil {
		g.Log().Errorf(ctx, `record insert failure %s`, err)
	}
}
