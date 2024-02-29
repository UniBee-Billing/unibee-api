// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CountryRate is the golang structure for table country_rate.
type CountryRate struct {
	Id                    uint64      `json:"id"                    description:"id"`                                                          // id
	MerchantId            uint64      `json:"merchantId"            description:""`                                                            //
	Gateway               string      `json:"gateway"               description:"vat gateway name, em vatsense"`                               // vat gateway name, em vatsense
	CountryCode           string      `json:"countryCode"           description:"country_code"`                                                // country_code
	CountryName           string      `json:"countryName"           description:"country_name"`                                                // country_name
	Latitude              string      `json:"latitude"              description:"latitude"`                                                    // latitude
	Longitude             string      `json:"longitude"             description:"longitude"`                                                   // longitude
	Vat                   int         `json:"vat"                   description:"vat contains，1-yes，2-no"`                                     // vat contains，1-yes，2-no
	Eu                    int         `json:"eu"                    description:"is eu member state, 1-yes, 2-no"`                             // is eu member state, 1-yes, 2-no
	StandardTaxPercentage int64       `json:"standardTaxPercentage" description:"Standard Tax Scale，1000 = 10%"`                               // Standard Tax Scale，1000 = 10%
	Other                 string      `json:"other"                 description:"other rates(json)"`                                           // other rates(json)
	StandardDescription   string      `json:"standardDescription"   description:"standard_description"`                                        // standard_description
	StandardTypes         string      `json:"standardTypes"         description:"standard_typs"`                                               // standard_typs
	Provinces             string      `json:"provinces"             description:"Whether the country has provinces with provincial sales tax"` // Whether the country has provinces with provincial sales tax
	Mamo                  string      `json:"mamo"                  description:"mamo"`                                                        // mamo
	GmtCreate             *gtime.Time `json:"gmtCreate"             description:"create time"`                                                 // create time
	GmtModify             *gtime.Time `json:"gmtModify"             description:"update time"`                                                 // update time
	IsDeleted             int         `json:"isDeleted"             description:"0-UnDeleted，1-Deleted"`                                       // 0-UnDeleted，1-Deleted
	CreateTime            int64       `json:"createTime"            description:"create utc time"`                                             // create utc time
}
