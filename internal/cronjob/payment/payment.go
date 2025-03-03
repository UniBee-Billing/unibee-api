package payment

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func TaskForCancelExpiredPayment(ctx context.Context) {
	var list []*entity.Payment
	err := dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().Status, consts.PaymentCreated).
		WhereNotNull(dao.Payment.Columns().GatewayPaymentIntentId).
		WhereLT(dao.Payment.Columns().CreateTime, gtime.Now().Timestamp()-(7*86400)).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForCancelExpiredPayment error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForCancelExpiredPayment-%s", one.PaymentId)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Debugf(ctx, "TaskForCancelExpiredPayment GetLock 60s, key:%s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Errorf(ctx, "TaskForCancelExpiredPayment ReleaseLock, key:%s", key)
			}()
			if len(one.InvoiceId) > 0 {
				in := query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)
				if in != nil && (in.PaymentId != one.PaymentId || in.Status == consts.InvoiceStatusCancelled || in.Status == consts.InvoiceStatusFailed) {
					err = service.PaymentGatewayCancel(ctx, one)
					if err != nil {
						g.Log().Errorf(ctx, "TaskForCancelExpiredPayment PaymentGatewayCancel, error:%s", err.Error())
					} else {
						g.Log().Errorf(ctx, "TaskForCancelExpiredPayment PaymentGatewayCancel, success paymentId:%s", one.PaymentId)
					}
				}
			}
		} else {
			g.Log().Errorf(ctx, "TaskForCancelExpiredPayment GetLock Failure, key:%s", key)
			return
		}
	}
}
