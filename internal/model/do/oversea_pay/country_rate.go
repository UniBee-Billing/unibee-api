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
	Gateway               interface{} // vat gateway name, em vatsense
	CountryCode           interface{} // country_code
	CountryName           interface{} // country_name
	Latitude              interface{} // latitude
	Longitude             interface{} // longitude
	Vat                   interface{} // vat contains，1-yes，2-no
	Eu                    interface{} // is eu member state, 1-yes, 2-no
	StandardTaxPercentage interface{} // Standard Tax Scale，1000 = 10%
	Other                 interface{} // other rates(json)
	StandardDescription   interface{} // standard_description
	StandardTypes         interface{} // standard_typs
	Provinces             interface{} // Whether the country has provinces with provincial sales tax
	Mamo                  interface{} // mamo
	GmtCreate             *gtime.Time // create time
	GmtModify             *gtime.Time // update time
	IsDeleted             interface{} // 0-UnDeleted，1-Deleted
	CreateAt              interface{} // create utc time
}
