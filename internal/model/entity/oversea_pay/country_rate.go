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
	Gateway               string      `json:"gateway"               description:"vat channel name, em vatsense"`                               // vat channel name, em vatsense
	CountryCode           string      `json:"countryCode"           description:"country_code"`                                                // country_code
	CountryName           string      `json:"countryName"           description:"country_name"`                                                // country_name
	Latitude              string      `json:"latitude"              description:"latitude"`                                                    // latitude
	Longitude             string      `json:"longitude"             description:"longitude"`                                                   // longitude
	Vat                   int         `json:"vat"                   description:"是否有 vat，1-有，2-没有"`                                            // 是否有 vat，1-有，2-没有
	Eu                    int         `json:"eu"                    description:"是否欧盟成员国, 1-是，2-不是"`                                           // 是否欧盟成员国, 1-是，2-不是
	StandardTaxPercentage int64       `json:"standardTaxPercentage" description:"Standard Tax 万分比，1000 表示 10%"`                                // Standard Tax 万分比，1000 表示 10%
	Other                 string      `json:"other"                 description:"other rates(json格式)"`                                         // other rates(json格式)
	StandardDescription   string      `json:"standardDescription"   description:"standard_description"`                                        // standard_description
	StandardTypes         string      `json:"standardTypes"         description:"standard_typs 限定"`                                            // standard_typs 限定
	Provinces             string      `json:"provinces"             description:"Whether the country has provinces with provincial sales tax"` // Whether the country has provinces with provincial sales tax
	Mamo                  string      `json:"mamo"                  description:"备注"`                                                          // 备注
	GmtCreate             *gtime.Time `json:"gmtCreate"             description:"创建时间"`                                                        // 创建时间
	GmtModify             *gtime.Time `json:"gmtModify"             description:"更新时间"`                                                        // 更新时间
	IsDeleted             int         `json:"isDeleted"             description:"是否删除，0-未删除，1-删除"`                                             // 是否删除，0-未删除，1-删除
}
