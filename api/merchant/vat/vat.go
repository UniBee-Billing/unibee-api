package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type SetupGatewayReq struct {
	g.Meta      `path:"/setup_gateway" tags:"Vat Gateway" method:"post" summary:"Vat Gateway Setup"`
	GatewayName string `json:"gatewayName" dc:"GatewayName, em. vatsense" v:"required"`
	Data        string `json:"data" dc:"Data" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}
type SetupGatewayRes struct {
	Data string `json:"data" dc:"Data" dc:"The hide star data"`
}

type InitDefaultGatewayReq struct {
	g.Meta `path:"/init_default_gateway" tags:"Vat Gateway" method:"post" summary:"Init Default Vat Gateway"`
}
type InitDefaultGatewayRes struct {
}

type CountryListReq struct {
	g.Meta `path:"/country_list" tags:"Vat Gateway" method:"get,post" summary:"Get Vat Country List"`
}
type CountryListRes struct {
	VatCountryList []*bean.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta    `path:"/vat_number_validate" tags:"Vat Gateway" method:"post" summary:"Vat Number Validation"`
	VatNumber string `json:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *bean.ValidResult `json:"vatNumberValidate"`
}

type NumberValidateHistoryReq struct {
	g.Meta          `path:"/vat_number_validate_history" tags:"Vat Gateway" method:"post" summary:"Vat Number Validation History"`
	SearchKey       string `json:"searchKey" dc:"Search Key, vatNumber, validateGateway, company, company address, message"  `
	VatNumber       string `json:"vatNumber" dc:"Filter Vat Number"`
	CountryCode     string `json:"countryCode" dc:"CountryCode"`
	ValidateGateway string `json:"validateGateway" dc:"Filter Validate Gateway, vatsense"`
	Status          []int  `json:"status" dc:"status, 0-Invalid，1-Valid" `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
}
type NumberValidateHistoryRes struct {
	NumberValidateHistoryList []*bean.MerchantVatNumberVerifyHistory `json:"numberValidateHistoryList" dc:"NumberValidateHistoryList"`
	Total                     int                                    `json:"total" dc:"Total"`
}

type NumberValidateHistoryActivateReq struct {
	g.Meta         `path:"/vat_number_validate_history_activate" tags:"Vat Gateway" method:"post" summary:"Vat Number Validation History Activate"`
	HistoryId      int64  `json:"historyId" dc:"History Id" `
	CountryCode    string `json:"countryCode" dc:"CountryCode"`
	CompanyName    string `json:"companyName"     description:"company_name"`    // company_name
	CompanyAddress string `json:"companyAddress"  description:"company_address"` // company_address
}
type NumberValidateHistoryActivateRes struct {
}

type NumberValidateHistoryDeactivateReq struct {
	g.Meta    `path:"/vat_number_validate_history_deactivate" tags:"Vat Gateway" method:"post" summary:"Vat Number Validation History Deactivate"`
	HistoryId int64 `json:"historyId" dc:"History Id" `
}
type NumberValidateHistoryDeactivateRes struct {
}
