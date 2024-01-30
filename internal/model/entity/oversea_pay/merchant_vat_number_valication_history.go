// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantVatNumberValicationHistory is the golang structure for table merchant_vat_number_valication_history.
type MerchantVatNumberValicationHistory struct {
	Id              int64       `json:"id"              description:"ID"`                    // ID
	MerchantId      int64       `json:"merchantId"      description:"merchantId"`            // merchantId
	VatNumber       string      `json:"vatNumber"       description:"vat_number"`            // vat_number
	Valid           int64       `json:"valid"           description:"0-无效，1-有效"`             // 0-无效，1-有效
	ValidateChannel string      `json:"validateChannel" description:"validate_channel"`      // validate_channel
	CountryCode     string      `json:"countryCode"     description:"country_code"`          // country_code
	CompanyName     string      `json:"companyName"     description:"company_name"`          // company_name
	CompanyAddress  string      `json:"companyAddress"  description:"company_address"`       // company_address
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"创建时间"`                  // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       description:"修改时间"`                  // 修改时间
	IsDeleted       int         `json:"isDeleted"       description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	ValidateMessage string      `json:"validateMessage" description:"validate_message"`      // validate_message
}
