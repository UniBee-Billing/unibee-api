// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditRecharge is the golang structure of table credit_recharge for DAO operations like Where/Data.
type CreditRecharge struct {
	g.Meta            `orm:"table:credit_recharge, do:true"`
	Id                interface{} // Id
	UserId            interface{} // user_id
	CreditId          interface{} // id of credit account
	RechargeId        interface{} // unique recharge id for credit account
	RechargeStatus    interface{} // recharge status, 10-in charging，20-recharge success，30-recharge failed
	Currency          interface{} // currency
	TotalAmount       interface{} // recharge total amount, cent
	PaymentAmount     interface{} // the payment amount for recharge
	Name              interface{} // recharge title
	Description       interface{} // recharge description
	PaidTime          interface{} // paid time
	GatewayId         interface{} // payment gateway id
	InvoiceId         interface{} // invoice_id
	PaymentId         interface{} // paymentId
	TotalRefundAmount interface{} // total refund amount,cent
	GmtCreate         *gtime.Time // create time
	GmtModify         *gtime.Time // update time
	CreateTime        interface{} // create utc time
	MerchantId        interface{} // merchant id
	AccountType       interface{} // type of credit account, 1-main recharge account, 2-promo credit account
}
