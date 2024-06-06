package event

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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

func printChannelPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err)
	}
}

func SaveEvent(ctx context.Context, overseaPayEvent entity.PaymentEvent) {
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		overseaPayEvent.CreateTime = gtime.Now().Timestamp()
		_, err = dao.PaymentEvent.Ctx(backgroundCtx).Data(overseaPayEvent).OmitNil().Insert(overseaPayEvent)
		if err != nil {
			g.Log().Errorf(backgroundCtx, `SaveEvent record insert failure %s`, err)
		}
	}()
}
