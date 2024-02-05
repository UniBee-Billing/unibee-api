package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/api"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func DoChannelCapture(ctx context.Context, payment *entity.Payment) (err error) {
	utility.Assert(payment != nil, "entity not found")
	utility.Assert(payment.Status == consts.TO_BE_PAID, "payment not waiting for pay")
	utility.Assert(payment.AuthorizeStatus != consts.WAITING_AUTHORIZED, "payment not authorised")
	utility.Assert(payment.PaymentAmount > 0, "capture value should > 0")
	utility.Assert(payment.PaymentAmount <= payment.TotalAmount, "capture value should <= authorized value")

	return dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		//事务处理 channel capture
		result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.CAPTURE_REQUEST, dao.Payment.Columns().PaymentAmount: payment.PaymentAmount},
			g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
		if err != nil || result == nil {
			//_ = transaction.Rollback()
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil || affected != 1 {
			//_ = transaction.Rollback()
			return err
		}

		//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
		_, err = api.GetPayChannelServiceProvider(ctx, payment.GatewayId).DoRemoteChannelCapture(ctx, payment)
		if err != nil {
			//_ = transaction.Rollback()
			return err
		}
		return nil
	})
}
