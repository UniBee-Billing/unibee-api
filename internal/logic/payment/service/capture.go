package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func PaymentGatewayCapture(ctx context.Context, payment *entity.Payment) (err error) {
	utility.Assert(payment != nil, "entity not found")
	utility.Assert(payment.Status == consts.PaymentCreated, "payment not waiting for pay")
	utility.Assert(payment.AuthorizeStatus == consts.Authorized, "payment not authorised")
	_, err = api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayCapture(ctx, payment)
	if err != nil {
		g.Log().Errorf(ctx, "PaymentGatewayCapture paymentId:%s, error:%s", payment.PaymentId, err.Error())
		return err
	}
	return dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		_, err = transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.CaptureRequest},
			g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.PaymentCreated})
		if err != nil {
			return err
		}
		return nil
	})
}
