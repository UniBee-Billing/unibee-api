// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditTransaction is the golang structure for table credit_transaction.
type CreditTransaction struct {
	Id                 int64       `json:"id"                 description:"Id"`                                                                                                                                             // Id
	UserId             uint64      `json:"userId"             description:"user_id"`                                                                                                                                        // user_id
	CreditId           uint64      `json:"creditId"           description:"id of credit account"`                                                                                                                           // id of credit account
	Currency           string      `json:"currency"           description:"currency"`                                                                                                                                       // currency
	TransactionId      string      `json:"transactionId"      description:"unique id for timeline"`                                                                                                                         // unique id for timeline
	TransactionType    int         `json:"transactionType"    description:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out"` // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter  int64       `json:"creditAmountAfter"  description:"the credit amount after transaction,cent"`                                                                                                       // the credit amount after transaction,cent
	CreditAmountBefore int64       `json:"creditAmountBefore" description:"the credit amount before transaction,cent"`                                                                                                      // the credit amount before transaction,cent
	DeltaAmount        int64       `json:"deltaAmount"        description:"delta amount,cent"`                                                                                                                              // delta amount,cent
	BizId              string      `json:"bizId"              description:"bisness id"`                                                                                                                                     // bisness id
	Name               string      `json:"name"               description:"recharge transaction title"`                                                                                                                     // recharge transaction title
	Description        string      `json:"description"        description:"recharge transaction description"`                                                                                                               // recharge transaction description
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`                                                                                                                                    // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`                                                                                                                                    // update time
	CreateTime         int64       `json:"createTime"         description:"create utc time"`                                                                                                                                // create utc time
	MerchantId         uint64      `json:"merchantId"         description:"merchant id"`                                                                                                                                    // merchant id
	InvoiceId          string      `json:"invoiceId"          description:"invoice_id"`                                                                                                                                     // invoice_id
	AccountType        int         `json:"accountType"        description:"type of credit account, 1-main recharge account, 2-promo credit account"`                                                                        // type of credit account, 1-main recharge account, 2-promo credit account
	AdminMemberId      uint64      `json:"adminMemberId"      description:"admin_member_id"`                                                                                                                                // admin_member_id
	ExchangeRate       int64       `json:"exchangeRate"       description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"`                     // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
}
