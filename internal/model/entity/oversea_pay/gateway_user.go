// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayUser is the golang structure for table gateway_user.
type GatewayUser struct {
	Id                          uint64      `json:"id"                          description:""`                               //
	GmtCreate                   *gtime.Time `json:"gmtCreate"                   description:"create time"`                    // create time
	GmtModify                   *gtime.Time `json:"gmtModify"                   description:"update time"`                    // update time
	UserId                      int64       `json:"userId"                      description:"userId"`                         // userId
	GatewayId                   int64       `json:"gatewayId"                   description:"gateway_id"`                     // gateway_id
	GatewayUserId               string      `json:"gatewayUserId"               description:"gateway_user_Id"`                // gateway_user_Id
	IsDeleted                   int         `json:"isDeleted"                   description:"0-UnDeleted，1-Deleted"`          // 0-UnDeleted，1-Deleted
	GatewayDefaultPaymentMethod string      `json:"gatewayDefaultPaymentMethod" description:"gateway_default_payment_method"` // gateway_default_payment_method
	CreateTime                  int64       `json:"createTime"                  description:"create utc time"`                // create utc time
}
