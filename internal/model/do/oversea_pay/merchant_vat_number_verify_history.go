// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantVatNumberVerifyHistory is the golang structure of table merchant_vat_number_verify_history for DAO operations like Where/Data.
type MerchantVatNumberVerifyHistory struct {
	g.Meta          `orm:"table:merchant_vat_number_verify_history, do:true"`
	Id              interface{} // ID
	MerchantId      interface{} // merchantId
	VatNumber       interface{} // vat_number
	Valid           interface{} // 0-无效，1-有效
	ValidateChannel interface{} // validate_channel
	CountryCode     interface{} // country_code
	CompanyName     interface{} // company_name
	CompanyAddress  interface{} // company_address
	GmtCreate       *gtime.Time // 创建时间
	GmtModify       *gtime.Time // 修改时间
	IsDeleted       interface{} // 0-UnDeleted，1-Deleted
	ValidateMessage interface{} // validate_message
}