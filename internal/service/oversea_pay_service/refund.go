package oversea_pay_service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "go-oversea-pay/api/out/v1"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/paychannel"
	"go-oversea-pay/internal/logic/paychannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strings"
)

func DoChannelRefund(ctx context.Context, bizType int, req *v1.RefundsReq, openApiId int64) (refund *entity.OverseaRefund, err error) {
	var (
		overseaPay *entity.OverseaPay
	)
	err = dao.OverseaPay.Ctx(ctx).Where(entity.OverseaPay{MerchantOrderNo: req.PaymentsPspReference}).OmitEmpty().Scan(&overseaPay)
	if err != nil {
		return nil, err
	}

	utility.Assert(overseaPay != nil, "payment not found")
	utility.Assert(overseaPay.PaymentFee > 0, "payment fee error")
	utility.Assert(strings.Compare(overseaPay.Currency, req.Amount.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(overseaPay.PayStatus == consts.PAY_SUCCESS, "payment not success")

	channel := util.GetOverseaPayChannel(ctx, uint64(overseaPay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")

	utility.Assert(req.Amount.Value > 0, "refund value should > 0")
	utility.Assert(req.Amount.Value <= overseaPay.PaymentFee, "refund value should <= PaymentFee value")

	redisKey := fmt.Sprintf("createRefund-payMerchantOrderNo:%s-bizId:%s", overseaPay.MerchantOrderNo, req.Reference)
	isDuplicatedInvoke := false

	//- 退款并发调用严重，加上Redis排它锁, todo mark 使用数据库锁
	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		return nil, gerror.Newf(`too fast duplicate call %s`, req.Reference)
	}

	var (
		overseaRefund *entity.OverseaRefund
	)
	err = dao.OverseaRefund.Ctx(ctx).Where(entity.OverseaRefund{
		OutTradeNo: req.PaymentsPspReference,
		BizId:      req.Reference,
		BizType:    bizType,
	}).OmitEmpty().Scan(&overseaRefund)
	utility.Assert(err == nil && overseaRefund == nil, "duplicate reference call")

	overseaRefund = &entity.OverseaRefund{
		Id:           utility.GenerateNextInt(),
		CompanyId:    overseaPay.CompanyId,
		MerchantId:   overseaPay.MerchantId,
		BizId:        req.Reference,
		BizType:      bizType,
		OutTradeNo:   overseaPay.MerchantOrderNo,
		OutRefundNo:  utility.CreateOutRefundNo(),
		RefundFee:    req.Amount.Value,
		RefundStatus: consts.REFUND_ING,
		ChannelId:    overseaPay.ChannelId,
		AppId:        overseaPay.AppId,
		Currency:     overseaPay.Currency,
		CountryCode:  overseaPay.CountryCode,
		OpenApiId:    openApiId,
		//AdditionalData: req.
		//RefundComment: payBizTypeEnum.getDesc() +"退款",

	}

	err = dao.OverseaRefund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		//事务处理 channel refund
		insert, err := transaction.Insert("oversea_refund", overseaRefund, 100)
		if err != nil {
			_ = transaction.Rollback()
			return err
		}
		id, err := insert.LastInsertId()
		if err != nil {
			_ = transaction.Rollback()
			return err
		}
		overseaRefund.Id = id

		//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
		channelResult, err := paychannel.GetPayChannelServiceProvider(int(overseaPay.ChannelId)).DoRemoteChannelRefund(ctx, overseaPay, overseaRefund)
		if err != nil {
			_ = transaction.Rollback()
			return err
		}

		result, err := transaction.Update("oversea_refund", g.Map{"channel_refund_no": channelResult.ChannelRefundNo},
			g.Map{"id": overseaRefund.Id, "refund_status": consts.REFUND_ING})
		if err != nil || result == nil {
			_ = transaction.Rollback()
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil || affected != 1 {
			_ = transaction.Rollback()
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return overseaRefund, nil
}
