package refund

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
	"unibee/internal/logic/credit/account"
	"unibee/internal/logic/credit/credit_query"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type CreditRefundInternalReq struct {
	UserId                 uint64 `json:"userId"`
	MerchantId             uint64 `json:"merchantId"`
	CreditPaymentId        string `json:"creditPaymentId"`
	ExternalCreditRefundId string `json:"externalCreditRefundId"`
	InvoiceId              string `json:"invoiceId"`
	RefundAmount           int64  `json:"refundAmount"`
	Currency               string `json:"currency"`
	Name                   string `json:"name"             description:"credit refund name"`
	Description            string `json:"description"       description:"credit refund description"`
}

type CreditRefundInternalRes struct {
	User          *entity.UserAccount   `json:"user"`
	CreditAccount *entity.CreditAccount `json:"creditAccount"`
	CreditRefund  *entity.CreditRefund  `json:"creditRefund"`
}

func NewCreditRefund(ctx context.Context, req *CreditRefundInternalReq) (res *CreditRefundInternalRes, err error) {
	utility.Assert(req != nil, "invalid request")
	utility.Assert(len(req.CreditPaymentId) > 0, "invalid creditPaymentId")
	utility.Assert(req.UserId > 0, "invalid userId")
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.ExternalCreditRefundId) > 0, "invalid externalCreditRefundId")
	utility.Assert(len(req.Currency) > 0, "invalid currency")
	utility.Assert(len(req.InvoiceId) > 0, "invalid invoiceId")
	utility.Assert(req.RefundAmount > 0, "invalid totalRefundAmount")
	user := credit_query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.MerchantId == req.MerchantId, "invalid merchantId")
	payment := credit_query.GetCreditPaymentByCreditPaymentId(ctx, req.CreditPaymentId)
	utility.Assert(payment != nil, "credit payment not found")
	if payment.TotalAmount-payment.TotalRefundAmount < req.RefundAmount {
		return nil, gerror.New("no enough amount can refund")
	}
	creditAccount := account.QueryOrCreateCreditAccount(ctx, req.UserId, req.Currency, payment.AccountType)
	utility.Assert(creditAccount != nil, "credit creditAccount failed")
	// check exist externalCreditRefundId
	one := credit_query.GetCreditRefundByExternalCreditRefundId(ctx, req.MerchantId, req.ExternalCreditRefundId)
	utility.Assert(one == nil, "credit payment exist with same externalCreditRefundId")
	{
		//name and description
		if len(req.Name) == 0 {
			req.Name = "Credit Refund"
		}
		if len(req.Description) == 0 {
			req.Description = "Credit Refund"
		}
	}

	creditRefundId := utility.CreateCreditRefundId()
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicCreditRefundSuccess, creditRefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.CreditRecharge.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {

			one = &entity.CreditRefund{
				UserId:                 req.UserId,
				CreditId:               creditAccount.Id,
				Currency:               req.Currency,
				InvoiceId:              req.InvoiceId,
				CreditRefundId:         creditRefundId,
				ExternalCreditRefundId: req.ExternalCreditRefundId,
				RefundAmount:           req.RefundAmount,
				RefundTime:             gtime.Now().Timestamp(),
				Name:                   req.Name,
				Description:            req.Description,
				CreateTime:             gtime.Now().Timestamp(),
				MerchantId:             req.MerchantId,
				AccountType:            creditAccount.Type,
			}
			result, err := dao.CreditRefund.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				return err
			}
			id, _ := result.LastInsertId()
			one.Id = id

			payment = credit_query.GetCreditPaymentByCreditPaymentId(ctx, req.CreditPaymentId)
			if payment == nil {
				return gerror.New("credit payment not found")
			}
			if payment.TotalAmount-payment.TotalRefundAmount < req.RefundAmount {
				return gerror.New("no enough amount can refund")
			}

			trans := &entity.CreditTransaction{
				UserId:             one.UserId,
				CreditId:           one.CreditId,
				Currency:           one.Currency,
				InvoiceId:          req.InvoiceId,
				TransactionId:      utility.CreateEventId(),
				TransactionType:    consts.CreditTransactionRefundIncome,
				CreditAmountAfter:  creditAccount.Amount + one.RefundAmount,
				CreditAmountBefore: creditAccount.Amount,
				DeltaAmount:        one.RefundAmount,
				BizId:              one.CreditRefundId,
				Name:               one.Name,
				Description:        one.Description,
				CreateTime:         gtime.Now().Timestamp(),
				MerchantId:         one.MerchantId,
				ExchangeRate:       payment.ExchangeRate,
				AccountType:        creditAccount.Type,
			}
			_, err = dao.CreditTransaction.Ctx(ctx).Data(trans).OmitNil().Insert(trans)
			if err != nil {
				return err
			}
			// append the payment total refund amount
			//update, err := dao.CreditPayment.Ctx(ctx).Data(g.Map{
			//	dao.CreditPayment.Columns().TotalRefundAmount: fmt.Sprintf("%s+%d", dao.CreditPayment.Columns().TotalRefundAmount, one.RefundAmount),
			//	dao.CreditPayment.Columns().GmtModify:         gtime.Now(),
			//}).Where(dao.CreditPayment.Columns().CreditPaymentId, req.CreditPaymentId).Update()
			update, err := dao.CreditPayment.Ctx(ctx).Where(dao.CreditPayment.Columns().CreditPaymentId, req.CreditPaymentId).Increment(dao.CreditPayment.Columns().TotalRefundAmount, req.RefundAmount)
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

			// append the credit amount
			//update, err = dao.CreditAccount.Ctx(ctx).Data(g.Map{
			//	dao.CreditAccount.Columns().Amount:    fmt.Sprintf("%s+%d", dao.CreditAccount.Columns().Amount, one.RefundAmount),
			//	dao.CreditAccount.Columns().GmtModify: gtime.Now(),
			//}).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Update()
			update, err = dao.CreditAccount.Ctx(ctx).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Increment(dao.CreditAccount.Columns().Amount, req.RefundAmount)
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
			affected, err = update.RowsAffected()
			if err != nil {
				return err
			}
			if affected != 1 {
				return gerror.New("update credit amount err")
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
		g.Log().Errorf(ctx, "NewCreditRefund sendResult err=%s", err.Error())
		return nil, err
	} else {
		// send message
		g.Log().Infof(ctx, "NewCreditRefund send success")
		return &CreditRefundInternalRes{
			User:          user,
			CreditAccount: creditAccount,
			CreditRefund:  one,
		}, nil
	}
}
