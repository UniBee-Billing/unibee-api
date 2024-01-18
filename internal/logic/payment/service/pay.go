package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	redismqcmd "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/event"
	"go-oversea-pay/internal/logic/payment/gateway"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

func DoChannelPay(ctx context.Context, createPayContext *ro.CreatePayContext) (channelInternalPayResult *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.Pay.BizType > 0, "pay bizType is nil")
	utility.Assert(createPayContext.PayChannel != nil, "pay gateway is nil")
	utility.Assert(createPayContext.Pay != nil, "pay is nil")
	utility.Assert(len(createPayContext.Pay.BizId) > 0, "支付单号为空")
	utility.Assert(createPayContext.Pay.ChannelId > 0, "pay channelId is nil")
	utility.Assert(createPayContext.Pay.PaymentFee > 0, "支付金额为空")
	utility.Assert(len(createPayContext.Pay.CountryCode) > 0, "contryCode为空")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(createPayContext.Pay.MerchantId > 0, "merchantId为空")
	utility.Assert(createPayContext.Pay.CompanyId > 0, "companyId为空")
	// 查询并处理所有待支付订单 todo mark

	createPayContext.Pay.Status = consts.TO_BE_PAID
	//createPayContext.Pay.AdditionalData = todo mark
	createPayContext.Pay.PaymentId = utility.CreatePaymentId()
	createPayContext.Pay.OpenApiId = createPayContext.OpenApiId
	//toSave.setServiceRate(iMerchantInfoService.getServiceDeductPoint(toSave.getMerchantId(),toSave.getChannelId()));//记录当下的服务费率
	redisKey := fmt.Sprintf("createPay-merchantId:%d-bizId:%s", createPayContext.Pay.MerchantId, createPayContext.Pay.BizId)
	isDuplicatedInvoke := false

	//- 并发调用严重，加上Redis排它锁, todo mark 使用数据库锁
	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		return nil, gerror.Newf(`too fast duplicate call %s`, createPayContext.Pay.BizId)
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCreated, createPayContext.Pay.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//事务处理 gateway refund
			//insert, err := transaction.Insert(dao.OverseaPay.Table(), createPayContext.Pay, 100)
			insert, err := dao.Payment.Ctx(ctx).Data(createPayContext.Pay).OmitEmpty().Insert(createPayContext.Pay)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			createPayContext.Pay.Id = id

			//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的todo mark
			channelInternalPayResult, err = gateway.GetPayChannelServiceProvider(ctx, createPayContext.Pay.ChannelId).DoRemoteChannelPayment(ctx, createPayContext)
			if err != nil {
				//_ = transaction.Rollback()
				return err
			}
			channelInternalPayResult.PayChannel = createPayContext.Pay.ChannelId
			channelInternalPayResult.PayOrderNo = createPayContext.Pay.PaymentId
			jsonData, err := gjson.Marshal(channelInternalPayResult)
			if err != nil {
				return err
			}
			createPayContext.Pay.PaymentData = string(jsonData)
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().PaymentData: createPayContext.Pay.PaymentData},
				g.Map{dao.Payment.Columns().Id: id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
		event.SaveTimeLine(ctx, entity.Timeline{
			BizType:   0,
			BizId:     createPayContext.Pay.PaymentId,
			Fee:       createPayContext.Pay.PaymentFee,
			EventType: event.SentForSettle.Type,
			Event:     event.SentForSettle.Desc,
			OpenApiId: createPayContext.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", createPayContext.Pay.PaymentId, "SentForSettle"),
		})
	}
	return channelInternalPayResult, nil
}
