package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func PaymentGatewayCancel(ctx context.Context, payment *entity.Payment) (err error) {
	utility.Assert(payment != nil, "entity not found")
	utility.Assert(payment.Status == consts.PaymentCreated, "payment not waiting for pay")
	utility.Assert(payment.AuthorizeStatus < consts.CaptureRequest, "payment has capture request")

	return dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的 todo mark
		_, err = api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayCancel(ctx, payment)
		if err != nil {
			//_ = transaction.Rollback()
			return err
		}
		return nil
	})
}
