// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditAccount is the golang structure of table credit_account for DAO operations like Where/Data.
type CreditAccount struct {
	g.Meta         `orm:"table:credit_account, do:true"`
	Id             interface{} // Id
	UserId         interface{} // user_id
	Type           interface{} // type of credit account, 1-main recharge account, 2-promo credit account
	Currency       interface{} // currency
	Amount         interface{} // credit amount,cent
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	CreateTime     interface{} // create utc time
	MerchantId     interface{} // merchant id
	RechargeEnable interface{} // 0-no, 1-yes
	PayoutEnable   interface{} // 0-no, 1-yes
}
