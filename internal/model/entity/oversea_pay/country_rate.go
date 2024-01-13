// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CountryRate is the golang structure for table country_rate.
type CountryRate struct {
	Id                    uint64      `json:"id"                    ` // id
	VatName               string      `json:"vatName"               ` // vat channel name, em vatsense
	CountryCode           string      `json:"countryCode"           ` // country_code
	CountryName           string      `json:"countryName"           ` // country_name
	Latitude              string      `json:"latitude"              ` // latitude
	Longitude             string      `json:"longitude"             ` // longitude
	Vat                   int         `json:"vat"                   ` // 是否有 vat，1-有，2-没有
	Eu                    int         `json:"eu"                    ` // 是否欧盟成员国, 1-是，2-不是
	StandardTaxPercentage int64       `json:"standardTaxPercentage" ` // Standard Tax百分比，10 表示 10%
	Other                 string      `json:"other"                 ` // other rates(json格式)
	StandardDescription   string      `json:"standardDescription"   ` // standard_description
	StandardTypes         string      `json:"standardTypes"         ` // standard_typs 限定
	Provinces             string      `json:"provinces"             ` // provinces 是否包含provinces Rate true or false
	Mamo                  string      `json:"mamo"                  ` // 备注
	GmtCreate             *gtime.Time `json:"gmtCreate"             ` // 创建时间
	GmtModify             *gtime.Time `json:"gmtModify"             ` // 更新时间
	IsDeleted             int         `json:"isDeleted"             ` // 是否删除，0-未删除，1-删除
}
