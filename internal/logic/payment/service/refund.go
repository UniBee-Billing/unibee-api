package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "go-oversea-pay/api/out/v1"
	redismqcmd "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/event"
	"go-oversea-pay/internal/logic/payment/outchannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
	"strings"
)

func DoChannelRefund(ctx context.Context, bizType int, req *v1.RefundsReq, openApiId int64) (refund *entity.OverseaRefund, err error) {
	overseaPay := query.GetOverseaPayByMerchantOrderNo(ctx, req.PaymentsPspReference)
	utility.Assert(overseaPay != nil, "payment not found")
	utility.Assert(overseaPay.PaymentFee > 0, "payment fee error")
	utility.Assert(strings.Compare(overseaPay.Currency, req.Amount.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(overseaPay.PayStatus == consts.PAY_SUCCESS, "payment not success")

	channel := query.GetPayChannelById(ctx, overseaPay.ChannelId)
	utility.Assert(channel != nil, "支付渠道异常 outchannel not found")

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
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundCreated, overseaRefund.OutRefundNo), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.OverseaRefund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//事务处理 outchannel refund
			//insert, err := transaction.Insert(dao.OverseaRefund.Table(), overseaRefund, 100) //todo mark 需要忽略空字段
			insert, err := dao.OverseaRefund.Ctx(ctx).Data(overseaRefund).OmitEmpty().Insert(overseaRefund)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			overseaRefund.Id = id

			//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
			channelResult, err := outchannel.GetPayChannelServiceProvider(ctx, overseaPay.ChannelId).DoRemoteChannelRefund(ctx, overseaPay, overseaRefund)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}

			overseaRefund.ChannelRefundNo = channelResult.ChannelRefundNo
			result, err := transaction.Update(dao.OverseaRefund.Table(), g.Map{dao.OverseaRefund.Columns().ChannelRefundNo: channelResult.ChannelRefundNo},
				g.Map{dao.OverseaRefund.Columns().Id: overseaRefund.Id, dao.OverseaRefund.Columns().RefundStatus: consts.REFUND_ING})
			if err != nil || result == nil {
				//_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				//_ = transaction.Rollback()
				return err
			}
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction, nil
		} else {
			return redismq.RollbackTransaction, err
		}
	})

	if err != nil {
		return nil, err
	} else {
		//交易事件记录
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     overseaPay.Id,
			Fee:       overseaPay.PaymentFee,
			EventType: event.SentForRefund.Type,
			Event:     event.SentForRefund.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", overseaPay.MerchantOrderNo, "SentForRefund", overseaRefund.OutRefundNo),
		})
	}
	return overseaRefund, nil
}
