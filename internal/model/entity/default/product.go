// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure for table product.
type Product struct {
	Id          uint64      `json:"id"          description:""`                                           //
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`                                // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update time"`                                // update time
	CompanyId   int64       `json:"companyId"   description:"company id"`                                 // company id
	MerchantId  uint64      `json:"merchantId"  description:"merchant id"`                                // merchant id
	ProductName string      `json:"productName" description:"ProductName"`                                // ProductName
	Description string      `json:"description" description:"description"`                                // description
	ImageUrl    string      `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     string      `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      int         `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"`                      // 0-UnDeleted，1-Deleted
	CreateTime  int64       `json:"createTime"  description:"create utc time"`                            // create utc time
	MetaData    string      `json:"metaData"    description:"meta_data(json)"`                            // meta_data(json)
}
