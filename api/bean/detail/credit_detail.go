package detail

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type CreditAccountDetail struct {
	Id                   uint64            `json:"id"         description:"Id"` // Id
	User                 *bean.UserAccount `json:"user"`
	Type                 int               `json:"type"       description:"type of credit account, 1-main account, 2-gift account"`                                                                                // type of credit account, 1-main account, 2-gift account
	Currency             string            `json:"currency"   description:"currency"`                                                                                                                              // currency
	Amount               int64             `json:"amount"     description:"credit amount, in cent if type is main"`                                                                                                // credit amount,cent
	CurrencyAmount       int64             `json:"currencyAmount"     description:"currency amount, in cent"`                                                                                                      // currency amount,cent
	ExchangeRate         int64             `json:"exchangeRate"          description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	CreateTime           int64             `json:"createTime" description:"create utc time"`                                                                                                                       // create utc time
	TotalIncrementAmount int64             `json:"totalIncrementAmount"     description:"the total increment amount"`
	TotalDecrementAmount int64             `json:"totalDecrementAmount"     description:"the total decrement amount"`
}

func ConvertToCreditAccountDetail(ctx context.Context, one *entity.CreditAccount) *CreditAccountDetail {
	if one == nil {
		return nil
	}
	currencyAmount, exchangeRate := bean.ConvertCreditAmountToCurrency(ctx, one.MerchantId, one.Type, one.Currency, one.Amount)
	return &CreditAccountDetail{
		Id:                   one.Id,
		User:                 bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		Type:                 one.Type,
		Currency:             one.Currency,
		Amount:               one.Amount,
		CurrencyAmount:       currencyAmount,
		ExchangeRate:         exchangeRate,
		CreateTime:           one.CreateTime,
		TotalDecrementAmount: int64(bean.GetCreditAccountTotalDecrementAmount(ctx, one.Id)),
		TotalIncrementAmount: int64(bean.GetCreditAccountTotalIncrementAmount(ctx, one.Id)),
	}
}

type CreditRechargeDetail struct {
	User           *bean.UserAccount    `json:"user"`
	Merchant       *bean.Merchant       `json:"merchant"`
	Gateway        *bean.Gateway        `json:"gateway"`
	CreditAccount  *bean.CreditAccount  `json:"creditAccount"`
	CreditRecharge *bean.CreditRecharge `json:"creditRecharge"`
	Invoice        *bean.Invoice        `json:"invoice"`
	Payment        *bean.Payment        `json:"payment"`
	Link           string               `json:"link"`
}

type CreditTransactionDetail struct {
	Id                  int64                `json:"id"                 description:"Id"` // Id
	User                *bean.UserAccount    `json:"user"`
	CreditAccount       *bean.CreditAccount  `json:"creditAccount"`
	Currency            string               `json:"currency"           description:"currency"`                                                                                                                                                    // currency
	TransactionId       string               `json:"transactionId"      description:"unique id for timeline"`                                                                                                                                      // unique id for timeline
	TransactionType     int                  `json:"transactionType"    description:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out"`              // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter   int64                `json:"creditAmountAfter"  description:"the credit amount after transaction,cent"`                                                                                                                    // the credit amount after transaction,cent
	CreditAmountBefore  int64                `json:"creditAmountBefore" description:"the credit amount before transaction,cent"`                                                                                                                   // the credit amount before transaction,cent
	DeltaAmount         int64                `json:"deltaAmount"        description:"delta amount,cent"`                                                                                                                                           // delta amount,cent
	DeltaCurrencyAmount int64                `json:"deltaCurrencyAmount"     description:"delta currency amount, in cent"`                                                                                                                         // currency amount,cent
	ExchangeRate        int64                `json:"exchangeRate"          description:"ExchangeRate for transaction, keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	BizId               string               `json:"bizId"              description:"business id"`                                                                                                                                                 // bisness id
	Name                string               `json:"name"               description:"recharge transaction title"`                                                                                                                                  // recharge transaction title
	Description         string               `json:"description"        description:"recharge transaction description"`                                                                                                                            // recharge transaction description 	// update time
	CreateTime          int64                `json:"createTime"         description:"create utc time"`                                                                                                                                             // create utc time
	MerchantId          uint64               `json:"merchantId"         description:"merchant id"`                                                                                                                                                 // merchant id
	InvoiceId           string               `json:"invoiceId"         description:"invoice_id"`                                                                                                                                                   // invoice_id
	AccountType         int                  `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"`                                                                                      // type of credit account, 1-main recharge account, 2-promo credit account
	AdminMember         *bean.MerchantMember `json:"adminMember"       description:"admin member"`
	By                  string               `json:"by"  dc:"" `
}

func ConvertToCreditTransactionDetail(ctx context.Context, one *entity.CreditTransaction) *CreditTransactionDetail {
	if one == nil {
		return nil
	}
	creditAccount := bean.SimplifyCreditAccount(ctx, query.GetCreditAccountById(ctx, one.CreditId))
	deltaCurrencyAmount, exchangeRate := bean.ConvertTransactionCreditAmountToCurrency(ctx, one.MerchantId, one.AccountType, one.Currency, one.DeltaAmount, one.ExchangeRate)
	by := "-"
	if one.AdminMemberId > 0 {
		member := bean.GetMerchantMemberById(ctx, one.AdminMemberId)
		if member != nil {
			by = member.Email
		}
	} else if one.UserId > 0 {
		user := bean.GetUserAccountById(ctx, one.UserId)
		if user != nil {
			by = user.Email
		}
	}
	return &CreditTransactionDetail{
		Id:                  one.Id,
		User:                bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		CreditAccount:       creditAccount,
		Currency:            one.Currency,
		TransactionId:       one.TransactionId,
		TransactionType:     one.TransactionType,
		CreditAmountAfter:   one.CreditAmountAfter,
		CreditAmountBefore:  one.CreditAmountBefore,
		DeltaAmount:         one.DeltaAmount,
		DeltaCurrencyAmount: deltaCurrencyAmount,
		ExchangeRate:        exchangeRate,
		BizId:               one.BizId,
		Name:                one.Name,
		Description:         one.Description,
		CreateTime:          one.CreateTime,
		MerchantId:          one.MerchantId,
		InvoiceId:           one.InvoiceId,
		AccountType:         one.AccountType,
		AdminMember:         bean.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, one.AdminMemberId)),
		By:                  by,
	}
}
