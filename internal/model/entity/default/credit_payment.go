// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditPayment is the golang structure for table credit_payment.
type CreditPayment struct {
	Id                      int64       `json:"id"                      description:"Id"`                                                                                                                         // Id
	UserId                  uint64      `json:"userId"                  description:"user_id"`                                                                                                                    // user_id
	CreditId                uint64      `json:"creditId"                description:"id of credit account"`                                                                                                       // id of credit account
	Currency                string      `json:"currency"                description:"currency"`                                                                                                                   // currency
	CreditPaymentId         string      `json:"creditPaymentId"         description:"credit payment id"`                                                                                                          // credit payment id
	ExternalCreditPaymentId string      `json:"externalCreditPaymentId" description:"external credit payment id"`                                                                                                 // external credit payment id
	TotalAmount             int64       `json:"totalAmount"             description:"total amount,cent"`                                                                                                          // total amount,cent
	PaidTime                int64       `json:"paidTime"                description:"paid time"`                                                                                                                  // paid time
	Name                    string      `json:"name"                    description:"recharge transaction title"`                                                                                                 // recharge transaction title
	Description             string      `json:"description"             description:"recharge transaction description"`                                                                                           // recharge transaction description
	GmtCreate               *gtime.Time `json:"gmtCreate"               description:"create time"`                                                                                                                // create time
	GmtModify               *gtime.Time `json:"gmtModify"               description:"update time"`                                                                                                                // update time
	CreateTime              int64       `json:"createTime"              description:"create utc time"`                                                                                                            // create utc time
	MerchantId              uint64      `json:"merchantId"              description:"merchant id"`                                                                                                                // merchant id
	InvoiceId               string      `json:"invoiceId"               description:"invoice_id"`                                                                                                                 // invoice_id
	TotalRefundAmount       int64       `json:"totalRefundAmount"       description:"total amount,cent"`                                                                                                          // total amount,cent
	ExchangeRate            int64       `json:"exchangeRate"            description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	PaidCurrencyAmount      int64       `json:"paidCurrencyAmount"      description:""`                                                                                                                           //
	AccountType             int         `json:"accountType"             description:"type of credit account, 1-main recharge account, 2-promo credit account"`                                                    // type of credit account, 1-main recharge account, 2-promo credit account
}
