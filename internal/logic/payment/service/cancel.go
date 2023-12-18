package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func DoChannelCancel(ctx context.Context, overseaPay *entity.OverseaPay) (err error) {
	utility.Assert(overseaPay != nil, "entity not found")
	utility.Assert(overseaPay.PayStatus == consts.TO_BE_PAID, "payment not waiting for pay")
	utility.Assert(overseaPay.AuthorizeStatus < consts.CAPTURE_REQUEST, "payment has capture request")

	return dao.OverseaPay.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		////事务处理 outchannel capture
		//result, err := transaction.Update("oversea_pay", g.Map{"pay_status": consts.PAY_FAILED},
		//	g.Map{"id": overseaPay.Id, "pay_status": consts.TO_BE_PAID})
		//if err != nil || result == nil {
		//	//_ = transaction.Rollback()
		//	return err
		//}
		//affected, err := result.RowsAffected()
		//if err != nil || affected != 1 {
		//	//_ = transaction.Rollback()
		//	return err
		//}

		//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的 todo mark
		_, err = outchannel.GetPayChannelServiceProvider(ctx, overseaPay.ChannelId).DoRemoteChannelCancel(ctx, overseaPay)
		if err != nil {
			//_ = transaction.Rollback()
			return err
		}
		return nil
	})
}
