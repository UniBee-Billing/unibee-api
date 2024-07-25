// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantOperationLog is the golang structure for table merchant_operation_log.
type MerchantOperationLog struct {
	Id                 uint64      `json:"id"                 description:"id"`                                                   // id
	CompanyId          int64       `json:"companyId"          description:"company id"`                                           // company id
	MerchantId         uint64      `json:"merchantId"         description:"merchant Id"`                                          // merchant Id
	MemberId           uint64      `json:"memberId"           description:"member_id"`                                            // member_id
	UserId             uint64      `json:"userId"             description:"user_id"`                                              // user_id
	OptAccount         string      `json:"optAccount"         description:"admin account"`                                        // admin account
	OptAccountType     int         `json:"optAccountType"     description:"opt_account_type, 0-Member|1-User|2-OpenApi|3-System"` // opt_account_type, 0-Member|1-User|2-OpenApi|3-System
	OptAccountId       string      `json:"optAccountId"       description:"operation account id"`                                 // operation account id
	ClientType         int         `json:"clientType"         description:"client type"`                                          // client type
	BizType            int         `json:"bizType"            description:"biz_type"`                                             // biz_type
	OptTarget          string      `json:"optTarget"          description:"operation target"`                                     // operation target
	OptContent         string      `json:"optContent"         description:"operation content"`                                    // operation content
	CreateTime         int64       `json:"createTime"         description:"operation create utc time"`                            // operation create utc time
	IsDelete           int         `json:"isDelete"           description:"0-UnDeleted，1-Deleted"`                                // 0-UnDeleted，1-Deleted
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`                                          // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`                                          // update time
	QueryportRequestId string      `json:"queryportRequestId" description:"queryport id"`                                         // queryport id
	ServerType         int         `json:"serverType"         description:"server type"`                                          // server type
	ServerTypeDesc     string      `json:"serverTypeDesc"     description:"server type description"`                              // server type description
	SubscriptionId     string      `json:"subscriptionId"     description:"subscription_id"`                                      // subscription_id
	InvoiceId          string      `json:"invoiceId"          description:"invoice id"`                                           // invoice id
	PlanId             uint64      `json:"planId"             description:"plan id"`                                              // plan id
	DiscountCode       string      `json:"discountCode"       description:"discount_code"`                                        // discount_code
}
