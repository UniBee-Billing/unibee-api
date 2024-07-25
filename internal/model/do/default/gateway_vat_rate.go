// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayVatRate is the golang structure of table gateway_vat_rate for DAO operations like Where/Data.
type GatewayVatRate struct {
	g.Meta           `orm:"table:gateway_vat_rate, do:true"`
	Id               interface{} //
	GmtCreate        *gtime.Time // create time
	GmtModify        *gtime.Time // update time
	VatRateId        interface{} // vat_rate_id
	GatewayId        interface{} // gateway_id
	GatewayVatRateId interface{} // gateway_vat_rate_Id
	IsDeleted        interface{} // 0-UnDeleted，1-Deleted
	CreateTime       interface{} // create utc time
}
