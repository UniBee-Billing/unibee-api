// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure for table merchant_info.
type MerchantInfo struct {
	Id          int64       `json:"id"          description:"merchant_id"`           // merchant_id
	CompanyId   int64       `json:"companyId"   description:"company_id"`            // company_id
	UserId      int64       `json:"userId"      description:"create_user_id"`        // create_user_id
	Type        int         `json:"type"        description:"type"`                  // type
	CompanyName string      `json:"companyName" description:"company_name"`          // company_name
	Email       string      `json:"email"       description:"email"`                 // email
	BusinessNum string      `json:"businessNum" description:"business_num"`          // business_num
	Name        string      `json:"name"        description:"name"`                  // name
	Idcard      string      `json:"idcard"      description:"idcard"`                // idcard
	Location    string      `json:"location"    description:"location"`              // location
	Address     string      `json:"address"     description:"address"`               // address
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update_time"`           // update_time
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CompanyLogo string      `json:"companyLogo" description:"company_logo"`          // company_logo
	HomeUrl     string      `json:"homeUrl"     description:""`                      //
	FirstName   string      `json:"firstName"   description:"first_name"`            // first_name
	LastName    string      `json:"lastName"    description:"last_name"`             // last_name
	Phone       string      `json:"phone"       description:"phone"`                 // phone
	CreateAt    int64       `json:"createAt"    description:"create utc time"`       // create utc time
}
