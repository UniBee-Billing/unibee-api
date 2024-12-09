// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditPayment is the golang structure of table credit_payment for DAO operations like Where/Data.
type CreditPayment struct {
	g.Meta                  `orm:"table:credit_payment, do:true"`
	Id                      interface{} // Id
	UserId                  interface{} // user_id
	CreditId                interface{} // id of credit account
	Currency                interface{} // currency
	CreditPaymentId         interface{} // credit payment id
	ExternalCreditPaymentId interface{} // external credit payment id
	TotalAmount             interface{} // total amount,cent
	PaidTime                interface{} // paid time
	Name                    interface{} // recharge transaction title
	Description             interface{} // recharge transaction description
	GmtCreate               *gtime.Time // create time
	GmtModify               *gtime.Time // update time
	CreateTime              interface{} // create utc time
	MerchantId              interface{} // merchant id
	InvoiceId               interface{} // invoice_id
	TotalRefundAmount       interface{} // total amount,cent
	ExchangeRate            interface{} //
	PaidCurrencyAmount      interface{} //
	AccountType             interface{} // type of credit account, 1-main recharge account, 2-promo credit account
}
