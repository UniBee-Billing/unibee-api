package detail

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type CreditAccountDetail struct {
	Id         uint64            `json:"id"         description:"Id"` // Id
	User       *bean.UserAccount `json:"user"`
	Type       int               `json:"type"       description:"type of credit account, 1-main account, 2-gift account"` // type of credit account, 1-main account, 2-gift account
	Currency   string            `json:"currency"   description:"currency"`                                               // currency
	Amount     int64             `json:"amount"     description:"credit amount, in cent if type is main"`                 // credit amount,cent
	CreateTime int64             `json:"createTime" description:"create utc time"`                                        // create utc time
}

func ConvertToCreditAccountDetail(ctx context.Context, one *entity.CreditAccount) *CreditAccountDetail {
	if one == nil {
		return nil
	}
	return &CreditAccountDetail{
		Id:         one.Id,
		User:       bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		Type:       one.Type,
		Currency:   one.Currency,
		Amount:     one.Amount,
		CreateTime: one.CreateTime,
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
	Id                 int64                `json:"id"                 description:"Id"` // Id
	User               *bean.UserAccount    `json:"user"`
	CreditAccount      *bean.CreditAccount  `json:"creditAccount"`
	Currency           string               `json:"currency"           description:"currency"`                                                                                                                                       // currency
	TransactionId      string               `json:"transactionId"      description:"unique id for timeline"`                                                                                                                         // unique id for timeline
	TransactionType    int                  `json:"transactionType"    description:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out"` // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter  int64                `json:"creditAmountAfter"  description:"the credit amount after transaction,cent"`                                                                                                       // the credit amount after transaction,cent
	CreditAmountBefore int64                `json:"creditAmountBefore" description:"the credit amount before transaction,cent"`                                                                                                      // the credit amount before transaction,cent
	DeltaAmount        int64                `json:"deltaAmount"        description:"delta amount,cent"`                                                                                                                              // delta amount,cent
	BizId              string               `json:"bizId"              description:"bisness id"`                                                                                                                                     // bisness id
	Name               string               `json:"name"               description:"recharge transaction title"`                                                                                                                     // recharge transaction title
	Description        string               `json:"description"        description:"recharge transaction description"`                                                                                                               // recharge transaction description 	// update time
	CreateTime         int64                `json:"createTime"         description:"create utc time"`                                                                                                                                // create utc time
	MerchantId         uint64               `json:"merchantId"         description:"merchant id"`                                                                                                                                    // merchant id
	InvoiceId          string               `json:"invoiceId"         description:"invoice_id"`                                                                                                                                      // invoice_id
	AccountType        int                  `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"`                                                                         // type of credit account, 1-main recharge account, 2-promo credit account
	AdminMember        *bean.MerchantMember `json:"adminMember"       description:"admin member"`
}

func ConvertToCreditTransactionDetail(ctx context.Context, one *entity.CreditTransaction) *CreditTransactionDetail {
	if one == nil {
		return nil
	}
	return &CreditTransactionDetail{
		Id:                 one.Id,
		User:               bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		CreditAccount:      bean.SimplifyCreditAccount(query.GetCreditAccountById(ctx, one.CreditId)),
		Currency:           one.Currency,
		TransactionId:      one.TransactionId,
		TransactionType:    one.TransactionType,
		CreditAmountAfter:  one.CreditAmountAfter,
		CreditAmountBefore: one.CreditAmountBefore,
		DeltaAmount:        one.DeltaAmount,
		BizId:              one.BizId,
		Name:               one.Name,
		Description:        one.Description,
		CreateTime:         one.CreateTime,
		MerchantId:         one.MerchantId,
		InvoiceId:          one.InvoiceId,
		AccountType:        one.AccountType,
		AdminMember:        bean.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, one.AdminMemberId)),
	}
}
