// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMember is the golang structure for table merchant_member.
type MerchantMember struct {
	Id         uint64      `json:"id"         description:"userId"`                // userId
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"`           // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"`           // update time
	MerchantId uint64      `json:"merchantId" description:"merchant id"`           // merchant id
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Password   string      `json:"password"   description:"password"`              // password
	UserName   string      `json:"userName"   description:"user name"`             // user name
	Mobile     string      `json:"mobile"     description:"mobile"`                // mobile
	Email      string      `json:"email"      description:"email"`                 // email
	FirstName  string      `json:"firstName"  description:"first name"`            // first name
	LastName   string      `json:"lastName"   description:"last name"`             // last name
	CreateTime int64       `json:"createTime" description:"create utc time"`       // create utc time
}
