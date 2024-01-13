// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CountryRate is the golang structure of table country_rate for DAO operations like Where/Data.
type CountryRate struct {
	g.Meta                `orm:"table:country_rate, do:true"`
	Id                    interface{} // id
	VatName               interface{} // vat channel name, em vatsense
	CountryCode           interface{} // country_code
	CountryName           interface{} // country_name
	Latitude              interface{} // latitude
	Longitude             interface{} // longitude
	Vat                   interface{} // 是否有 vat，1-有，2-没有
	Eu                    interface{} // 是否欧盟成员国, 1-是，2-不是
	StandardTaxPercentage interface{} // Standard Tax百分比，10 表示 10%
	Other                 interface{} // other rates(json格式)
	StandardDescription   interface{} // standard_description
	StandardTypes         interface{} // standard_typs 限定
	Provinces             interface{} // Whether the country has provinces with provincial sales tax
	Mamo                  interface{} // 备注
	GmtCreate             *gtime.Time // 创建时间
	GmtModify             *gtime.Time // 更新时间
	IsDeleted             interface{} // 是否删除，0-未删除，1-删除
}
