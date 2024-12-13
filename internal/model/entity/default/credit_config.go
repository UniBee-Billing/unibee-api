// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditConfig is the golang structure for table credit_config.
type CreditConfig struct {
	Id                    uint64      `json:"id"                    description:"Id"`                                                                                                                         // Id
	Type                  int         `json:"type"                  description:"type of credit account, 1-main account, 2-promo credit account"`                                                             // type of credit account, 1-main account, 2-promo credit account
	Currency              string      `json:"currency"              description:"currency"`                                                                                                                   // currency
	ExchangeRate          int64       `json:"exchangeRate"          description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	GmtCreate             *gtime.Time `json:"gmtCreate"             description:"create time"`                                                                                                                // create time
	GmtModify             *gtime.Time `json:"gmtModify"             description:"update time"`                                                                                                                // update time
	CreateTime            int64       `json:"createTime"            description:"create utc time"`                                                                                                            // create utc time
	MerchantId            uint64      `json:"merchantId"            description:"merchant id"`                                                                                                                // merchant id
	Recurring             int         `json:"recurring"             description:"apply to recurring, default no, 0-no,1-yes"`                                                                                 // apply to recurring, default no, 0-no,1-yes
	DiscountCodeExclusive int         `json:"discountCodeExclusive" description:"discount code exclusive when purchase, default no, 0-no, 1-yes"`                                                             // discount code exclusive when purchase, default no, 0-no, 1-yes
	Logo                  string      `json:"logo"                  description:"logo image base64, show when user purchase"`                                                                                 // logo image base64, show when user purchase
	Name                  string      `json:"name"                  description:"name"`                                                                                                                       // name
	Description           string      `json:"description"           description:"description"`                                                                                                                // description
	LogoUrl               string      `json:"logoUrl"               description:"logo url, show when user purchase"`                                                                                          // logo url, show when user purchase
	MetaData              string      `json:"metaData"              description:"meta_data(json)"`                                                                                                            // meta_data(json)
	RechargeEnable        int         `json:"rechargeEnable"        description:"0-no,1-yes"`                                                                                                                 // 0-no,1-yes
	PayoutEnable          int         `json:"payoutEnable"          description:"0-no,1-yes"`                                                                                                                 // 0-no,1-yes
	PreviewDefaultUsed    int         `json:"previewDefaultUsed"    description:"is default used when in purchase preview,0-no, 1-yes"`                                                                       // is default used when in purchase preview,0-no, 1-yes
}
