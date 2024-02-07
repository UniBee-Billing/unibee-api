// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantPayGatewayBinding is the golang structure for table merchant_pay_gateway_binding.
type MerchantPayGatewayBinding struct {
	Id         uint64      `json:"id"         description:""`                      //
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"`           // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"`           // update time
	MerchantId int64       `json:"merchantId" description:"merchant id"`           // merchant id
	GatewayId  string      `json:"gatewayId"  description:"gateway_id"`            // gateway_id
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime int64       `json:"createTime" description:"create utc time"`       // create utc time
}
