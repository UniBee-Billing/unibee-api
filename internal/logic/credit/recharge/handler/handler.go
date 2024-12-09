package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/credit"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func HandleCreditRechargeSuccess(ctx context.Context, creditRechargeId string, invoice *entity.Invoice, payment *entity.Payment) (err error) {
	creditAccount := credit.QueryOrCreateCreditAccount(ctx, invoice.UserId, invoice.Currency, consts.CreditAccountTypeMain)
	utility.Assert(creditAccount != nil, "credit creditAccount failed")
	// update invoiceId to credit recharge
	one := query.GetCreditRechargeByRechargeId(ctx, creditRechargeId)
	if one.RechargeStatus == consts.CreditRechargeSuccess {
		return gerror.New("recharge already success")
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicCreditRechargeSuccess, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.CreditRecharge.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			_, err = transaction.Update(dao.CreditRecharge.Table(), g.Map{
				dao.CreditRecharge.Columns().InvoiceId:      invoice.InvoiceId,
				dao.CreditRecharge.Columns().PaymentId:      payment.PaymentId,
				dao.CreditRecharge.Columns().RechargeStatus: consts.CreditRechargeSuccess,
				dao.CreditRecharge.Columns().GmtModify:      gtime.Now(),
			}, g.Map{dao.CreditRecharge.Columns().RechargeId: creditRechargeId})
			if err != nil {
				g.Log().Errorf(ctx, "HandleCreditRechargeSuccess sendResult err=%s", err.Error())
				return err
			}
			creditAccount = query.GetCreditAccountById(ctx, creditAccount.Id)
			trans := &entity.CreditTransaction{
				UserId:             one.UserId,
				CreditId:           one.CreditId,
				Currency:           one.Currency,
				InvoiceId:          invoice.InvoiceId,
				TransactionId:      utility.CreateEventId(),
				TransactionType:    consts.CreditTransactionRechargeIncome,
				CreditAmountAfter:  creditAccount.Amount + one.TotalAmount,
				CreditAmountBefore: creditAccount.Amount,
				DeltaAmount:        one.TotalAmount,
				BizId:              one.RechargeId,
				Name:               one.Name,
				Description:        one.Description,
				CreateTime:         gtime.Now().Timestamp(),
				MerchantId:         one.MerchantId,
				AccountType:        creditAccount.Type,
			}
			_, err = dao.CreditTransaction.Ctx(ctx).Data(trans).OmitNil().Insert(trans)
			if err != nil {
				return err
			}
			// append the credit amount
			//update, err := dao.CreditAccount.Ctx(ctx).Data(g.Map{
			//	dao.CreditAccount.Columns().Amount:    fmt.Sprintf("%s+%d", dao.CreditAccount.Columns().Amount, one.TotalAmount),
			//	dao.CreditAccount.Columns().GmtModify: gtime.Now(),
			//}).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Update()
			update, err := dao.CreditAccount.Ctx(ctx).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Increment(dao.CreditAccount.Columns().Amount, one.TotalAmount)
			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("update credit amount failed, err: %v", err.Error()))
				return err
			}
			if update == nil {
				return gerror.New("update credit amount err")
			}
			if err != nil {
				return err
			}
			affected, err := update.RowsAffected()
			if err != nil {
				return err
			}
			if affected != 1 {
				return gerror.New("update credit amount err")
			}
			if err != nil {
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
		g.Log().Errorf(ctx, "HandleCreditRechargeSuccess sendResult err=%s", err.Error())
	} else {
		g.Log().Infof(ctx, "HandleCreditRechargeSuccess send success")
	}
	return err
}

func HandleCreditRechargeFailed(ctx context.Context, creditRechargeId string) (err error) {
	one := query.GetCreditRechargeByRechargeId(ctx, creditRechargeId)
	if one.RechargeStatus == consts.CreditRechargeSuccess {
		return gerror.New("recharge already success")
	}
	// update invoiceId to credit recharge
	_, err = dao.CreditRecharge.Ctx(ctx).Data(g.Map{
		dao.CreditRecharge.Columns().RechargeStatus: consts.CreditRechargeFailed,
		dao.CreditRecharge.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.CreditRecharge.Columns().RechargeId, creditRechargeId).OmitNil().Update()
	if err != nil {
		return err
	}
	return err
}
