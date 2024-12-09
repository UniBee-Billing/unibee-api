// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditTransaction is the golang structure of table credit_transaction for DAO operations like Where/Data.
type CreditTransaction struct {
	g.Meta             `orm:"table:credit_transaction, do:true"`
	Id                 interface{} // Id
	UserId             interface{} // user_id
	CreditId           interface{} // id of credit account
	Currency           interface{} // currency
	TransactionId      interface{} // unique id for timeline
	TransactionType    interface{} // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter  interface{} // the credit amount after transaction,cent
	CreditAmountBefore interface{} // the credit amount before transaction,cent
	DeltaAmount        interface{} // delta amount,cent
	BizId              interface{} // bisness id
	Name               interface{} // recharge transaction title
	Description        interface{} // recharge transaction description
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	CreateTime         interface{} // create utc time
	MerchantId         interface{} // merchant id
	InvoiceId          interface{} // invoice_id
	AccountType        interface{} // type of credit account, 1-main recharge account, 2-promo credit account
	AdminMemberId      interface{} // admin_member_id
}
