// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmailHistory is the golang structure for table email_history.
type EmailHistory struct {
	Id         int64       `json:"id"         description:""`            //
	MerchantId int64       `json:"merchantId" description:""`            //
	Email      string      `json:"email"      description:""`            //
	Title      string      `json:"title"      description:""`            //
	Content    string      `json:"content"    description:""`            //
	AttachFile string      `json:"attachFile" description:""`            //
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"` // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"` // update time
	Response   string      `json:"response"   description:""`            //
}