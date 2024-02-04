package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	v1 "go-oversea-pay/api/open/payment"
	redismqcmd "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/payment/event"
	"go-oversea-pay/internal/logic/payment/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
	"strings"
)

func DoChannelRefund(ctx context.Context, bizType int, req *v1.RefundsReq, openApiId int64) (refund *entity.Refund, err error) {
	utility.Assert(len(req.PaymentId) > 0, "invalid paymentId")
	utility.Assert(len(req.MerchantRefundId) > 0, "invalid merchantRefundId")
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.TotalAmount > 0, "TotalAmount fee error")
	utility.Assert(strings.Compare(payment.Currency, req.Amount.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(payment.Status == consts.PAY_SUCCESS, "payment not success")

	payChannel := query.GetPayChannelById(ctx, payment.ChannelId)
	utility.Assert(payChannel != nil, "channel not found")

	utility.Assert(req.Amount.Amount > 0, "refund value should > 0")
	utility.Assert(req.Amount.Amount <= payment.TotalAmount, "refund value should <= TotalAmount value")

	redisKey := fmt.Sprintf("createRefund-paymentId:%s-bizId:%s", payment.PaymentId, req.MerchantRefundId)
	isDuplicatedInvoke := false

	//- 退款并发调用严重，加上Redis排它锁, todo mark use db lock
	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		utility.Assert(true, "Submit Too Fast")
	}

	var (
		one *entity.Refund
	)
	err = dao.Refund.Ctx(ctx).Where(entity.Refund{
		PaymentId: req.PaymentId,
		BizId:     req.MerchantRefundId,
		BizType:   bizType,
	}).OmitEmpty().Scan(&one)
	utility.Assert(err == nil && one == nil, "Duplicate Submit")

	one = &entity.Refund{
		CompanyId:     payment.CompanyId,
		MerchantId:    payment.MerchantId,
		BizId:         req.MerchantRefundId,
		BizType:       bizType,
		PaymentId:     payment.PaymentId,
		RefundId:      utility.CreateRefundId(),
		RefundAmount:  req.Amount.Amount,
		Status:        consts.REFUND_ING,
		ChannelId:     payment.ChannelId,
		AppId:         payment.AppId,
		Currency:      payment.Currency,
		CountryCode:   payment.CountryCode,
		RefundComment: req.Reason,
		OpenApiId:     openApiId,
		UserId:        payment.UserId,
		//AdditionalData: req.
		//RefundComment: payBizTypeEnum.getDesc() +"退款",

	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundCreated, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//事务处理 channel refund
			//insert, err := transaction.Insert(dao.OverseaRefund.Table(), overseaRefund, 100) //todo mark ignore nil field
			one.UniqueId = one.RefundId
			insert, err := dao.Refund.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			one.Id = id

			//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
			channelResult, err := out.GetPayChannelServiceProvider(ctx, payment.ChannelId).DoRemoteChannelRefund(ctx, payment, one)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}

			one.ChannelRefundId = channelResult.ChannelRefundId
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().ChannelRefundId: channelResult.ChannelRefundId},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.REFUND_ING})
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
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.SentForRefund.Type,
			Event:     event.SentForRefund.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "SentForRefund", one.RefundId),
		})
		err = handler.CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
	}
	return one, nil
}
