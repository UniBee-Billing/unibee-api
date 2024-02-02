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
	Authorised         = TradeEventTypeEnum{2, "Authorised"} //
	Cancelled          = TradeEventTypeEnum{3, "Cancelled"}  //
	Expired            = TradeEventTypeEnum{4, "Expired"}    //
	ChargeBack         = TradeEventTypeEnum{5, "ChargeBack"}
	ChargeBackReversed = TradeEventTypeEnum{6, "ChargeBackReversed"}
	CaptureFailed      = TradeEventTypeEnum{7, "CaptureFailed"} //
	Error              = TradeEventTypeEnum{8, "Error"}
	Refunded           = TradeEventTypeEnum{9, "Refunded"}          //
	RefundFailed       = TradeEventTypeEnum{10, "RefundFailed"}     //
	RefundedReversed   = TradeEventTypeEnum{11, "RefundedReversed"} //

	Refused = TradeEventTypeEnum{12, "Refused"}

	SentForRefund = TradeEventTypeEnum{13, "SentForRefund"} //

	SentForSettle   = TradeEventTypeEnum{14, "SentForSettle"}
	Settled         = TradeEventTypeEnum{15, "Settled"}
	SettledReversed = TradeEventTypeEnum{16, "SettledReversed"}
)

func SaveEvent(ctx context.Context, overseaPayEvent entity.PaymentEvent) {
	_, err := dao.PaymentEvent.Ctx(ctx).Data(overseaPayEvent).OmitNil().Insert(overseaPayEvent)
	if err != nil {
		g.Log().Errorf(ctx, `record insert failure %s`, err)
	}
}
