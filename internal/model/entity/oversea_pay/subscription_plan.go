// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlan is the golang structure for table subscription_plan.
type SubscriptionPlan struct {
	Id                        uint64      `json:"id"                        description:""`                                               //
	GmtCreate                 *gtime.Time `json:"gmtCreate"                 description:"create time"`                                    // create time
	GmtModify                 *gtime.Time `json:"gmtModify"                 description:"update time"`                                    // update time
	CompanyId                 int64       `json:"companyId"                 description:"company id"`                                     // company id
	MerchantId                int64       `json:"merchantId"                description:"merchant id"`                                    // merchant id
	PlanName                  string      `json:"planName"                  description:"PlanName"`                                       // PlanName
	Amount                    int64       `json:"amount"                    description:"amount, cent, without tax"`                      // amount, cent, without tax
	Currency                  string      `json:"currency"                  description:"currency"`                                       // currency
	IntervalUnit              string      `json:"intervalUnit"              description:"period unit,day|month|year|week"`                // period unit,day|month|year|week
	IntervalCount             int         `json:"intervalCount"             description:"period unit count"`                              // period unit count
	Description               string      `json:"description"               description:"description"`                                    // description
	ImageUrl                  string      `json:"imageUrl"                  description:"image_url"`                                      // image_url
	HomeUrl                   string      `json:"homeUrl"                   description:"home_url"`                                       // home_url
	GatewayProductName        string      `json:"gatewayProductName"        description:"gateway product name"`                           // gateway product name
	GatewayProductDescription string      `json:"gatewayProductDescription" description:"gateway product description"`                    // gateway product description
	TaxScale                  int         `json:"taxScale"                  description:"tax scale 1000 = 10%"`                           // tax scale 1000 = 10%
	TaxInclusive              int         `json:"taxInclusive"              description:"deperated"`                                      // deperated
	Type                      int         `json:"type"                      description:"type，1-main plan，2-addon plan"`                  // type，1-main plan，2-addon plan
	Status                    int         `json:"status"                    description:"status，1-editing，2-active，3-inactive，4-expired"` // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 int         `json:"isDeleted"                 description:"0-UnDeleted，1-Deleted"`                          // 0-UnDeleted，1-Deleted
	BindingAddonIds           string      `json:"bindingAddonIds"           description:"binded addon planIds，split with ,"`              // binded addon planIds，split with ,
	PublishStatus             int         `json:"publishStatus"             description:"1-UnPublish,2-Publish,用于控制是否在 UserPortal 端展示"`   // 1-UnPublish,2-Publish,用于控制是否在 UserPortal 端展示
}
