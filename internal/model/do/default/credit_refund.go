// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditRefund is the golang structure of table credit_refund for DAO operations like Where/Data.
type CreditRefund struct {
	g.Meta                 `orm:"table:credit_refund, do:true"`
	Id                     interface{} // Id
	UserId                 interface{} // user_id
	CreditId               interface{} // id of credit account
	Currency               interface{} // currency
	InvoiceId              interface{} // invoice_id
	CreditPaymentId        interface{} // credit refund id
	CreditRefundId         interface{} // credit refund id
	ExternalCreditRefundId interface{} // external credit refund id
	RefundAmount           interface{} // total refund amount,cent
	RefundTime             interface{} // refund time
	Name                   interface{} // recharge transaction title
	Description            interface{} // recharge transaction description
	GmtCreate              *gtime.Time // create time
	GmtModify              *gtime.Time // update time
	CreateTime             interface{} // create utc time
	MerchantId             interface{} // merchant id
	AccountType            interface{} // type of credit account, 1-main recharge account, 2-promo credit account
}
