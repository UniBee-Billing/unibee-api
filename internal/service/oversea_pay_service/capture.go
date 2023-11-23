package oversea_pay_service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/paychannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func DoChannelCapture(ctx context.Context, overseaPay *entity.OverseaPay) (err error) {
	utility.Assert(overseaPay != nil, "entity not found")
	utility.Assert(overseaPay.PayStatus == consts.TO_BE_PAID, "payment not waiting for pay")
	utility.Assert(overseaPay.AuthorizeStatus != consts.WAITING_AUTHORIZED, "payment not authorised")
	utility.Assert(overseaPay.BuyerPayFee > 0, "capture value should > 0")
	utility.Assert(overseaPay.BuyerPayFee <= overseaPay.PaymentFee, "capture value should <= authorized value")

	// todo mark 事务实现 channel capture
	return g.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		//事务处理 channel capture
		result, err := transaction.Update("oversea_pay", g.Map{"authorize_status": consts.CAPTURE_REQUEST, "buyer_pay_fee": overseaPay.BuyerPayFee},
			g.Map{"id": overseaPay.Id, "pay_status": consts.TO_BE_PAID})
		if err != nil || result == nil {
			_ = transaction.Rollback()
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil || affected != 1 {
			_ = transaction.Rollback()
			return err
		}

		//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
		_, err = paychannel.GetPayChannelServiceProvider(int(overseaPay.ChannelId)).DoRemoteChannelCapture(ctx, overseaPay)
		if err != nil {
			return err
		}
		return nil
	})
}
