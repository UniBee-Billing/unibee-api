package bean

import (
	"strings"
	entity "unibee/internal/model/entity/default"
)

type VatCountryRate struct {
	Id                    uint64 `json:"id"`
	Gateway               string `json:"gateway"           `        // gateway
	CountryCode           string `json:"countryCode"           `    // country_code
	CountryName           string `json:"countryName"           `    // country_name
	VatSupport            bool   `json:"vatSupport"               ` // vat support true or false
	IsEU                  bool   `json:"isEU"                  `
	StandardTaxPercentage int64  `json:"standardTaxPercentage" `
	Mamo                  string `json:"mamo"                  description:"mamo"` // mamo
}

type ValidResult struct {
	Valid           bool   `json:"valid"           `
	VatNumber       string `json:"vatNumber"           `
	CountryCode     string `json:"countryCode"           `
	CompanyName     string `json:"companyName"           `
	CompanyAddress  string `json:"companyAddress"           `
	ValidateMessage string `json:"validateMessage"           `
}

type MerchantVatNumberVerifyHistory struct {
	Id              int64  `json:"id"              description:"Id"`         // Id
	MerchantId      uint64 `json:"merchantId"      description:"merchantId"` // merchantId
	VatNumber       string `json:"vatNumber"       description:"vat_number"` // vat_number
	Status          int64  `json:"status" dc:"status, 0-Invalid，1-Valid" `
	ValidateGateway string `json:"validateGateway" description:"validate_gateway"` // validate_gateway
	CountryCode     string `json:"countryCode"     description:"country_code"`     // country_code
	CompanyName     string `json:"companyName"     description:"company_name"`     // company_name
	CompanyAddress  string `json:"companyAddress"  description:"company_address"`  // company_address
	ValidateMessage string `json:"validateMessage" description:"validate_message"` // validate_message
	CreateTime      int64  `json:"createTime"      description:"create utc time"`  // create utc time
	ManualValidate  bool   `json:"manualValidate"  description:"manual_validate"`
}

func SimplifyMerchantVatNumberVerifyHistory(one *entity.MerchantVatNumberVerifyHistory) *MerchantVatNumberVerifyHistory {
	if one == nil {
		return nil
	}
	manualValidate := false
	if strings.HasPrefix(one.ValidateMessage, "[Manual]") {
		manualValidate = true
		one.ValidateMessage = strings.Replace(one.ValidateMessage, "[Manual]", "", 1)
	}
	return &MerchantVatNumberVerifyHistory{
		Id:              one.Id,
		MerchantId:      one.MerchantId,
		VatNumber:       one.VatNumber,
		Status:          one.Valid,
		ValidateGateway: one.ValidateGateway,
		CountryCode:     one.CountryCode,
		CompanyName:     one.CompanyName,
		CompanyAddress:  one.CompanyAddress,
		ValidateMessage: one.ValidateMessage,
		CreateTime:      one.CreateTime,
		ManualValidate:  manualValidate,
	}
}
