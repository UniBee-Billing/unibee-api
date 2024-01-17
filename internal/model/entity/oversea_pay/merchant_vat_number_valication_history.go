// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantVatNumberValicationHistory is the golang structure for table merchant_vat_number_valication_history.
type MerchantVatNumberValicationHistory struct {
	Id              int64       `json:"id"              ` // ID
	MerchantId      int64       `json:"merchantId"      ` // merchantId
	VatNumber       string      `json:"vatNumber"       ` // vat_number
	Valid           int64       `json:"valid"           ` // 0-无效，1-有效
	ValidateChannel string      `json:"validateChannel" ` // validate_channel
	CountryCode     string      `json:"countryCode"     ` // country_code
	CompanyName     string      `json:"companyName"     ` // company_name
	CompanyAddress  string      `json:"companyAddress"  ` // company_address
	GmtCreate       *gtime.Time `json:"gmtCreate"       ` // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       ` // 修改时间
	IsDeleted       int         `json:"isDeleted"       ` // 是否删除，0-未删除，1-删除
	ValidateMessage string      `json:"validateMessage" ` // validate_message
}
