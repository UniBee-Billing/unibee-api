// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantVatNumberVerifyHistory is the golang structure for table merchant_vat_number_verify_history.
type MerchantVatNumberVerifyHistory struct {
	Id              int64       `json:"id"              description:"Id"`                    // Id
	MerchantId      int64       `json:"merchantId"      description:"merchantId"`            // merchantId
	VatNumber       string      `json:"vatNumber"       description:"vat_number"`            // vat_number
	Valid           int64       `json:"valid"           description:"0-Invalid，1-Valid"`     // 0-Invalid，1-Valid
	ValidateGateway string      `json:"validateGateway" description:"validate_gateway"`      // validate_gateway
	CountryCode     string      `json:"countryCode"     description:"country_code"`          // country_code
	CompanyName     string      `json:"companyName"     description:"company_name"`          // company_name
	CompanyAddress  string      `json:"companyAddress"  description:"company_address"`       // company_address
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"create time"`           // create time
	GmtModify       *gtime.Time `json:"gmtModify"       description:"update time"`           // update time
	IsDeleted       int         `json:"isDeleted"       description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	ValidateMessage string      `json:"validateMessage" description:"validate_message"`      // validate_message
}
