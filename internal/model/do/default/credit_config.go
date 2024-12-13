// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditConfig is the golang structure of table credit_config for DAO operations like Where/Data.
type CreditConfig struct {
	g.Meta                `orm:"table:credit_config, do:true"`
	Id                    interface{} // Id
	Type                  interface{} // type of credit account, 1-main account, 2-promo credit account
	Currency              interface{} // currency
	ExchangeRate          interface{} // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	GmtCreate             *gtime.Time // create time
	GmtModify             *gtime.Time // update time
	CreateTime            interface{} // create utc time
	MerchantId            interface{} // merchant id
	Recurring             interface{} // apply to recurring, default no, 0-no,1-yes
	DiscountCodeExclusive interface{} // discount code exclusive when purchase, default no, 0-no, 1-yes
	Logo                  interface{} // logo image base64, show when user purchase
	Name                  interface{} // name
	Description           interface{} // description
	LogoUrl               interface{} // logo url, show when user purchase
	MetaData              interface{} // meta_data(json)
	RechargeEnable        interface{} // 0-no,1-yes
	PayoutEnable          interface{} // 0-no,1-yes
	PreviewDefaultUsed    interface{} // is default used when in purchase preview,0-no, 1-yes
}
