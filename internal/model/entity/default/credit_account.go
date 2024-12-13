// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditAccount is the golang structure for table credit_account.
type CreditAccount struct {
	Id             uint64      `json:"id"             description:"Id"`                                                                      // Id
	UserId         uint64      `json:"userId"         description:"user_id"`                                                                 // user_id
	Type           int         `json:"type"           description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
	Currency       string      `json:"currency"       description:"currency"`                                                                // currency
	Amount         int64       `json:"amount"         description:"credit amount,cent"`                                                      // credit amount,cent
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`                                                             // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`                                                             // update time
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                                                         // create utc time
	MerchantId     uint64      `json:"merchantId"     description:"merchant id"`                                                             // merchant id
	RechargeEnable int         `json:"rechargeEnable" description:"0-no, 1-yes"`                                                             // 0-no, 1-yes
	PayoutEnable   int         `json:"payoutEnable"   description:"0-no, 1-yes"`                                                             // 0-no, 1-yes
}
