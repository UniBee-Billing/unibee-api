// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchExportTemplate is the golang structure of table merchant_batch_export_template for DAO operations like Where/Data.
type MerchantBatchExportTemplate struct {
	g.Meta        `orm:"table:merchant_batch_export_template, do:true"`
	Id            interface{} // id
	MerchantId    interface{} // merchant_id
	MemberId      interface{} // member_id
	Name          interface{} // name
	Task          interface{} // task
	Format        interface{} // format
	Payload       interface{} // payload(json)
	ExportColumns interface{} // export_columns(json)
	IsDeleted     interface{} // 0-UnDeleted，1-Deleted
	GmtCreate     *gtime.Time // gmt_create
	GmtModify     *gtime.Time // update time
	CreateTime    interface{} // create utc time
}
