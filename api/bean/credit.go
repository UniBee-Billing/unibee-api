package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
)

type CreditConfig struct {
	Id                    uint64                 `json:"id"                    description:"Id"`                                                                                                                         // Id
	Type                  int                    `json:"type"                  description:"type of credit account, 1-main account, 2-promo credit account"`                                                             // type of credit account, 1-main account, 2-promo credit account
	Currency              string                 `json:"currency"              description:"currency"`                                                                                                                   // currency
	ExchangeRate          int64                  `json:"exchangeRate"          description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	CreateTime            int64                  `json:"createTime"            description:"create utc time"`                                                                                                            // create utc time
	MerchantId            uint64                 `json:"merchantId"            description:"merchant id"`                                                                                                                // merchant id
	Recurring             int                    `json:"recurring"             description:"apply to recurring, default no, 0-no,1-yes"`                                                                                 // apply to reucrring, default yes, 0-yes,1-no
	DiscountCodeExclusive int                    `json:"discountCodeExclusive" description:"discount code exclusive when purchase, default no, 0-no, 1-yes"`                                                             // discount code exclusive when purchase, default yes, 0-yes, 1-no
	Logo                  string                 `json:"logo"                  description:"logo image base64, show when user purchase"`                                                                                 // logo image base64, show when user purchase
	Name                  string                 `json:"name"                  description:"name"`                                                                                                                       // name
	Description           string                 `json:"description"           description:"description"`                                                                                                                // description
	LogoUrl               string                 `json:"logoUrl"               description:"logo url, show when user purchase"`                                                                                          // logo url, show when user purchase
	MetaData              map[string]interface{} `json:"metaData"              description:"meta_data(json)"`
	RechargeEnable        int                    `json:"rechargeEnable"        description:"0-yes, 1-no"`
	PayoutEnable          int                    `json:"payoutEnable"          description:"0-yes, 1-no"`
	PreviewDefaultUsed    int                    `json:"previewDefaultUsed"    description:"is default used when in purchase preview, default no, 0-no, 1-yes"` // is default used when in purchase preview,0-yes, 1-no
}

func SimplifyCreditConfig(one *entity.CreditConfig) *CreditConfig {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyCreditConfig Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &CreditConfig{
		Id:                    one.Id,
		Type:                  one.Type,
		Currency:              one.Currency,
		ExchangeRate:          one.ExchangeRate,
		CreateTime:            one.CreateTime,
		MerchantId:            one.MerchantId,
		Recurring:             one.Recurring,
		DiscountCodeExclusive: one.DiscountCodeExclusive,
		Logo:                  one.Logo,
		Name:                  one.Name,
		Description:           one.Description,
		LogoUrl:               one.LogoUrl,
		MetaData:              metadata,
		RechargeEnable:        one.RechargeEnable,
		PayoutEnable:          one.PayoutEnable,
		PreviewDefaultUsed:    one.PreviewDefaultUsed,
	}
}

type CreditAccount struct {
	Id         uint64 `json:"id"         description:"Id"`                                                     // Id
	UserId     uint64 `json:"userId"     description:"user_id"`                                                // user_id
	Type       int    `json:"type"       description:"type of credit account, 1-main account, 2-gift account"` // type of credit account, 1-main account, 2-gift account
	Currency   string `json:"currency"   description:"currency"`                                               // currency
	Amount     int64  `json:"amount"     description:"credit amount,cent"`                                     // credit amount,cent
	CreateTime int64  `json:"createTime" description:"create utc time"`                                        // create utc time
}

func SimplifyCreditAccount(one *entity.CreditAccount) *CreditAccount {
	if one == nil {
		return nil
	}
	return &CreditAccount{
		Id:         one.Id,
		UserId:     one.UserId,
		Type:       one.Type,
		Currency:   one.Currency,
		Amount:     one.Amount,
		CreateTime: one.CreateTime,
	}
}

func SimplifyCreditAccountList(list []*entity.CreditAccount) []*CreditAccount {
	if list == nil || len(list) == 0 {
		return make([]*CreditAccount, 0)
	}
	result := make([]*CreditAccount, 0)
	for _, one := range list {
		result = append(result, SimplifyCreditAccount(one))
	}
	return result
}

type CreditRecharge struct {
	Id                int64  `json:"id"                description:"Id"`                                                                      // Id
	UserId            uint64 `json:"userId"            description:"user_id"`                                                                 // user_id
	CreditId          uint64 `json:"creditId"          description:"id of credit account"`                                                    // id of credit account
	RechargeId        string `json:"rechargeId"        description:"unique recharge id for credit account"`                                   // unique recharge id for credit account
	RechargeStatus    int    `json:"rechargeStatus"    description:"recharge status, 10-in charging，20-recharge success，30-recharge failed"`  // recharge status, 10-recharging，20-recharge success，30-recharge failed
	Currency          string `json:"currency"          description:"currency"`                                                                // currency
	TotalAmount       int64  `json:"totalAmount"       description:"recharge total amount, cent"`                                             // recharge total amount, cent
	PaymentAmount     string `json:"paymentAmount"     description:"the payment amount for recharge"`                                         // the payment amount for recharge
	Name              string `json:"name"              description:"recharge name"`                                                           // recharge name
	Description       string `json:"description"       description:"recharge description"`                                                    // recharge description
	PaidTime          int64  `json:"paidTime"          description:"paid time"`                                                               // paid time
	GatewayId         uint64 `json:"gatewayId"         description:"payment gateway id"`                                                      // payment gateway id
	InvoiceId         string `json:"invoiceId"         description:"invoice_id"`                                                              // invoice_id
	PaymentId         string `json:"paymentId"         description:"paymentId"`                                                               // paymentId
	TotalRefundAmount int64  `json:"totalRefundAmount" description:"total refund amount,cent"`                                                // total refund amount,cent
	CreateTime        int64  `json:"createTime"        description:"create utc time"`                                                         // create utc time
	AccountType       int    `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
}

func SimplifyCreditRecharge(one *entity.CreditRecharge) *CreditRecharge {
	if one == nil {
		return nil
	}
	return &CreditRecharge{
		Id:                one.Id,
		UserId:            one.UserId,
		CreditId:          one.CreditId,
		RechargeId:        one.RechargeId,
		RechargeStatus:    one.RechargeStatus,
		Currency:          one.Currency,
		TotalAmount:       one.TotalAmount,
		PaymentAmount:     one.PaymentId,
		Name:              one.Name,
		Description:       one.Description,
		PaidTime:          one.PaidTime,
		GatewayId:         one.GatewayId,
		InvoiceId:         one.InvoiceId,
		PaymentId:         one.PaymentId,
		TotalRefundAmount: one.TotalRefundAmount,
		CreateTime:        one.CreateTime,
		AccountType:       one.AccountType,
	}
}

type CreditTransaction struct {
	Id                 int64  `json:"id"                 description:"Id"`                                                                                                                                             // Id
	UserId             uint64 `json:"userId"             description:"user_id"`                                                                                                                                        // user_id
	CreditId           uint64 `json:"creditId"           description:"id of credit account"`                                                                                                                           // id of credit account
	Currency           string `json:"currency"           description:"currency"`                                                                                                                                       // currency
	TransactionId      string `json:"transactionId"      description:"unique id for timeline"`                                                                                                                         // unique id for timeline
	TransactionType    int    `json:"transactionType"    description:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out"` // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter  int64  `json:"creditAmountAfter"  description:"the credit amount after transaction,cent"`                                                                                                       // the credit amount after transaction,cent
	CreditAmountBefore int64  `json:"creditAmountBefore" description:"the credit amount before transaction,cent"`                                                                                                      // the credit amount before transaction,cent
	DeltaAmount        int64  `json:"deltaAmount"        description:"delta amount,cent"`                                                                                                                              // delta amount,cent
	BizId              string `json:"bizId"              description:"business id"`                                                                                                                                    // bisness id
	Name               string `json:"name"               description:"recharge transaction title"`                                                                                                                     // recharge transaction title
	Description        string `json:"description"        description:"recharge transaction description"`                                                                                                               // recharge transaction description 	// update time
	CreateTime         int64  `json:"createTime"         description:"create utc time"`                                                                                                                                // create utc time
	MerchantId         uint64 `json:"merchantId"         description:"merchant id"`                                                                                                                                    // merchant id
	InvoiceId          string `json:"invoiceId"         description:"invoice_id"`                                                                                                                                      // invoice_id
	AccountType        int    `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"`                                                                         // type of credit account, 1-main recharge account, 2-promo credit account
}

func SimplifyCreditTransaction(one *entity.CreditTransaction) *CreditTransaction {
	if one == nil {
		return nil
	}
	return &CreditTransaction{
		Id:                 one.Id,
		UserId:             one.UserId,
		CreditId:           one.CreditId,
		Currency:           one.Currency,
		TransactionId:      one.TransactionId,
		TransactionType:    one.TransactionType,
		CreditAmountAfter:  one.CreditAmountBefore,
		CreditAmountBefore: one.CreditAmountBefore,
		DeltaAmount:        one.DeltaAmount,
		BizId:              one.BizId,
		Name:               one.Name,
		Description:        one.Description,
		CreateTime:         one.CreateTime,
		MerchantId:         one.MerchantId,
		InvoiceId:          one.InvoiceId,
		AccountType:        one.AccountType,
	}
}

type CreditPayout struct {
	ExchangeRate   int64 `json:"exchangeRate"            description:"exchange rate, keep two decimal places，scale = 100, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"`
	CreditAmount   int64 `json:"creditAmount"      description:"credit amount, scale = 100"`
	CurrencyAmount int64 `json:"currencyAmount"      description:"currency amount,cent"`
}

type CreditPayment struct {
	Id                      int64  `json:"id"                      description:"Id"`                                                                // Id
	UserId                  uint64 `json:"userId"                  description:"user_id"`                                                           // user_id
	CreditId                uint64 `json:"creditId"                description:"id of credit account"`                                              // id of credit account
	Currency                string `json:"currency"                description:"currency"`                                                          // currency
	CreditPaymentId         string `json:"creditPaymentId"         description:"credit payment id"`                                                 // credit payment id
	ExternalCreditPaymentId string `json:"externalCreditPaymentId" description:"external credit payment id"`                                        // external credit payment id
	TotalAmount             int64  `json:"totalAmount"             description:"total amount,cent"`                                                 // total amount,cent
	PaidTime                int64  `json:"paidTime"                description:"paid time"`                                                         // paid time
	Name                    string `json:"name"                    description:"recharge transaction title"`                                        // recharge transaction title
	Description             string `json:"description"             description:"recharge transaction description"`                                  // recharge transaction description
	CreateTime              int64  `json:"createTime"              description:"create utc time"`                                                   // create utc time
	MerchantId              uint64 `json:"merchantId"              description:"merchant id"`                                                       // merchant id
	InvoiceId               string `json:"invoiceId"               description:"invoice_id"`                                                        // invoice_id
	TotalRefundAmount       int64  `json:"totalRefundAmount"       description:"total amount,cent"`                                                 // total amount,cent
	ExchangeRate            int64  `json:"exchangeRate"            description:""`                                                                  //
	PaidCurrencyAmount      int64  `json:"paidCurrencyAmount"      description:""`                                                                  //
	AccountType             int    `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
}

func SimplifyCreditPayment(one *entity.CreditPayment) *CreditPayment {
	if one == nil {
		return nil
	}
	return &CreditPayment{
		Id:                      one.Id,
		UserId:                  one.UserId,
		CreditId:                one.CreditId,
		Currency:                one.Currency,
		CreditPaymentId:         one.CreditPaymentId,
		ExternalCreditPaymentId: one.ExternalCreditPaymentId,
		TotalAmount:             one.TotalAmount,
		PaidTime:                one.PaidTime,
		Name:                    one.Name,
		Description:             one.Description,
		CreateTime:              one.CreateTime,
		MerchantId:              one.MerchantId,
		InvoiceId:               one.InvoiceId,
		TotalRefundAmount:       one.TotalRefundAmount,
		ExchangeRate:            one.ExchangeRate,
		PaidCurrencyAmount:      one.PaidCurrencyAmount,
		AccountType:             one.AccountType,
	}
}

type CreditRefund struct {
	Id                     int64  `json:"id"                     description:"Id"`                                                                 // Id
	UserId                 uint64 `json:"userId"                 description:"user_id"`                                                            // user_id
	CreditId               uint64 `json:"creditId"               description:"id of credit account"`                                               // id of credit account
	Currency               string `json:"currency"               description:"currency"`                                                           // currency
	InvoiceId              string `json:"invoiceId"              description:"invoice_id"`                                                         // invoice_id
	CreditPaymentId        string `json:"creditPaymentId"        description:"credit refund id"`                                                   // credit refund id
	CreditRefundId         string `json:"creditRefundId"         description:"credit refund id"`                                                   // credit refund id
	ExternalCreditRefundId string `json:"externalCreditRefundId" description:"external credit refund id"`                                          // external credit refund id
	TotalRefundAmount      int64  `json:"totalRefundAmount"      description:"total refund amount,cent"`                                           // total refund amount,cent
	RefundTime             int64  `json:"refundTime"             description:"refund time"`                                                        // refund time
	Name                   string `json:"name"                   description:"recharge transaction title"`                                         // recharge transaction title
	Description            string `json:"description"            description:"recharge transaction description"`                                   // recharge transaction description
	CreateTime             int64  `json:"createTime"             description:"create utc time"`                                                    // create utc time
	MerchantId             uint64 `json:"merchantId"             description:"merchant id"`                                                        // merchant id
	AccountType            int    `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
}

func SimplifyCreditRefund(one *entity.CreditRefund) *CreditRefund {
	if one == nil {
		return nil
	}
	return &CreditRefund{
		Id:                     one.Id,
		UserId:                 one.UserId,
		CreditId:               one.CreditId,
		Currency:               one.Currency,
		InvoiceId:              one.InvoiceId,
		CreditPaymentId:        one.CreditPaymentId,
		CreditRefundId:         one.CreditRefundId,
		ExternalCreditRefundId: one.ExternalCreditRefundId,
		TotalRefundAmount:      one.RefundAmount,
		RefundTime:             one.RefundTime,
		Name:                   one.Name,
		Description:            one.Description,
		CreateTime:             one.CreateTime,
		MerchantId:             one.MerchantId,
		AccountType:            one.AccountType,
	}
}