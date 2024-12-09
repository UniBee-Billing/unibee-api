package payment

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/credit"
	"unibee/internal/logic/credit/config"
	"unibee/internal/logic/credit/refund"
	currency2 "unibee/internal/logic/currency"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func CheckCreditUserPayout(ctx context.Context, merchantId uint64, userId uint64, creditType int, currency string, currencyAmount int64) (*bean.CreditAccount, *bean.CreditPayout, error) {
	if merchantId < 0 || userId < 0 {
		return nil, nil, gerror.New("invalid merchantId or userId")
	}
	if merchantId <= 0 {
		return nil, nil, gerror.New("invalid merchantId")
	}
	if creditType != 1 && creditType != 2 {
		return nil, nil, gerror.New("invalid creditType")
	}
	if !currency2.IsCurrencySupport(currency) {
		return nil, nil, gerror.New("invalid currency")
	}
	one := query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return nil, nil, gerror.New("credit config need setup")
	}
	if one.PayoutEnable > 0 {
		return nil, nil, gerror.New("credit account payout disable")
	}
	account := query.GetCreditAccountByUserId(ctx, userId, creditType, currency)
	if account == nil {
		return nil, nil, gerror.New("credit account not found")
	}
	if account.Amount == 0 {
		return nil, nil, gerror.New("credit account amount is zero")
	}
	creditPaymentAmount, _ := config.ConvertCurrencyAmountToCreditAmount(ctx, merchantId, creditType, currency, currencyAmount)
	if creditPaymentAmount > account.Amount {
		creditPaymentAmount = account.Amount
		currencyAmount, _ = config.ConvertCurrencyAmountToCreditAmount(ctx, merchantId, creditType, currency, creditPaymentAmount)
	}
	return bean.SimplifyCreditAccount(account), &bean.CreditPayout{
		ExchangeRate:   one.ExchangeRate,
		CreditAmount:   creditPaymentAmount,
		CurrencyAmount: currencyAmount,
	}, nil
}

type CreditPaymentInternalReq struct {
	UserId                  uint64 `json:"userId"`
	MerchantId              uint64 `json:"merchantId"`
	ExternalCreditPaymentId string `json:"externalCreditPaymentId"`
	InvoiceId               string `json:"invoiceId"`
	CurrencyAmount          int64  `json:"currencyAmount"`
	Currency                string `json:"currency"`
	CreditType              int    `json:"creditType"`
	Name                    string `json:"name"             description:"credit payment name"`
	Description             string `json:"description"       description:"credit payment  description"`
}

type CreditPaymentInternalRes struct {
	User          *entity.UserAccount   `json:"user"`
	CreditAccount *entity.CreditAccount `json:"creditAccount"`
	CreditPayment *entity.CreditPayment `json:"creditPayment"`
}

func NewCreditPayment(ctx context.Context, req *CreditPaymentInternalReq) (res *CreditPaymentInternalRes, err error) {
	utility.Assert(req != nil, "invalid request")
	utility.Assert(req.UserId > 0, "invalid userId")
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.ExternalCreditPaymentId) > 0, "invalid externalCreditPaymentId")
	utility.Assert(len(req.Currency) > 0, "invalid currency")
	utility.Assert(len(req.InvoiceId) > 0, "invalid invoiceId")
	utility.Assert(req.CurrencyAmount > 0, "invalid currencyAmount")
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.MerchantId == req.MerchantId, "invalid merchantId")
	utility.Assert(req.CreditType == consts.CreditAccountTypeMain || req.CreditType == consts.CreditAccountTypePromo, "invalid creditType")
	creditAccount := credit.QueryOrCreateCreditAccount(ctx, req.UserId, req.Currency, req.CreditType)
	utility.Assert(creditAccount != nil, "credit creditAccount failed")
	utility.Assert(creditAccount.Type == req.CreditType, "invalid credit account type, should be main account")
	utility.AssertError(config.CheckCreditConfigPayout(ctx, req.MerchantId, creditAccount.Type, req.Currency), "Credit Config Error")
	creditPaymentAmount, exchangeRate := config.ConvertCurrencyAmountToCreditAmount(ctx, req.MerchantId, creditAccount.Type, req.Currency, req.CurrencyAmount)
	if creditAccount.Amount < creditPaymentAmount {
		return nil, gerror.New("credit amount is not enough")
	}
	// check exist externalCreditPaymentId
	one := query.GetCreditPaymentByExternalCreditPaymentId(ctx, req.MerchantId, req.ExternalCreditPaymentId)
	utility.Assert(one == nil, "credit payment exist with same externalCreditPaymentId")
	{
		//name and description
		if len(req.Name) == 0 {
			req.Name = "Credit Payment"
		}
		if len(req.Description) == 0 {
			req.Description = "Credit Payment"
		}
	}

	creditPaymentId := utility.CreateCreditPaymentId()
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicCreditPaymentSuccess, creditPaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.CreditRecharge.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			one = &entity.CreditPayment{
				UserId:                  req.UserId,
				CreditId:                creditAccount.Id,
				Currency:                req.Currency,
				CreditPaymentId:         creditPaymentId,
				ExternalCreditPaymentId: req.ExternalCreditPaymentId,
				TotalAmount:             creditPaymentAmount,
				PaidTime:                gtime.Now().Timestamp(),
				Name:                    req.Name,
				Description:             req.Description,
				CreateTime:              gtime.Now().Timestamp(),
				MerchantId:              req.MerchantId,
				InvoiceId:               req.InvoiceId,
				ExchangeRate:            exchangeRate,
				PaidCurrencyAmount:      req.CurrencyAmount,
				AccountType:             creditAccount.Type,
			}
			result, err := dao.CreditPayment.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				return err
			}
			id, _ := result.LastInsertId()
			one.Id = id

			creditAccount = query.GetCreditAccountById(ctx, creditAccount.Id)
			if creditAccount.Amount < creditPaymentAmount {
				return gerror.New("credit amount is less than total amount")
			}
			trans := &entity.CreditTransaction{
				UserId:             one.UserId,
				CreditId:           one.CreditId,
				Currency:           one.Currency,
				InvoiceId:          req.InvoiceId,
				TransactionId:      utility.CreateEventId(),
				TransactionType:    consts.CreditTransactionPayout,
				CreditAmountAfter:  creditAccount.Amount - one.TotalAmount,
				CreditAmountBefore: creditAccount.Amount,
				DeltaAmount:        one.TotalAmount,
				BizId:              one.CreditPaymentId,
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
			//	dao.CreditAccount.Columns().Amount:    fmt.Sprintf("%s-%d", dao.CreditAccount.Columns().Amount, one.TotalAmount),
			//	dao.CreditAccount.Columns().GmtModify: gtime.Now(),
			//}).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Update()
			update, err := dao.CreditAccount.Ctx(ctx).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Decrement(dao.CreditAccount.Columns().Amount, one.TotalAmount)
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
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction, nil
		} else {
			return redismq.RollbackTransaction, err
		}
	})
	if err != nil {
		g.Log().Errorf(ctx, "NewCreditPayment sendResult err=%s", err.Error())
		return nil, err
	} else {
		// send message
		g.Log().Infof(ctx, "NewCreditPayment send success")
		return &CreditPaymentInternalRes{
			User:          user,
			CreditAccount: creditAccount,
			CreditPayment: one,
		}, nil
	}
}

func printChannelPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "CallbackException panic error:%s", err)
	}
}

func RollbackCreditPayment(ctx context.Context, merchantId uint64, externalCreditPaymentId string) error {
	one := query.GetCreditPaymentByExternalCreditPaymentId(ctx, merchantId, externalCreditPaymentId)
	if one == nil {
		return gerror.New("credit payment not found")
	}
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				printChannelPanic(backgroundCtx, err)
				return
			}
		}()
		_, err = refund.NewCreditRefund(ctx, &refund.CreditRefundInternalReq{
			UserId:                 one.UserId,
			MerchantId:             one.MerchantId,
			CreditPaymentId:        one.CreditPaymentId,
			ExternalCreditRefundId: externalCreditPaymentId,
			InvoiceId:              one.InvoiceId,
			RefundAmount:           one.TotalAmount,
			Currency:               one.Currency,
			Name:                   one.Name,
			Description:            one.Description,
		})
		if err != nil {
			g.Log().Error(context.Background(), "RollbackCreditPayment NewCreditRefund externalCreditPaymentId:%s error:%s", externalCreditPaymentId, err.Error())
		}
	}()
	return nil
}
