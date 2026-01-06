// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantCheckout is the golang structure for table merchant_checkout.
type MerchantCheckout struct {
	Id          int64       `json:"id"          description:"ID"`                    // ID
	MerchantId  uint64      `json:"merchantId"  description:"merchantId"`            // merchantId
	Name        string      `json:"name"        description:"name"`                  // name
	Description string      `json:"description" description:"description"`           // description
	Data        string      `json:"data"        description:"data(json)"`            // data(json)
	Staging     string      `json:"staging"     description:"staging_data(json)"`    // staging_data(json)
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update time"`           // update time
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime  int64       `json:"createTime"  description:"create utc time"`       // create utc time
}
