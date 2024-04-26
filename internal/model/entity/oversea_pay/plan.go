// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Plan is the golang structure for table plan.
type Plan struct {
	Id                        uint64      `json:"id"                        description:""`                                                                                                                //
	GmtCreate                 *gtime.Time `json:"gmtCreate"                 description:"create time"`                                                                                                     // create time
	GmtModify                 *gtime.Time `json:"gmtModify"                 description:"update time"`                                                                                                     // update time
	CompanyId                 int64       `json:"companyId"                 description:"company id"`                                                                                                      // company id
	MerchantId                uint64      `json:"merchantId"                description:"merchant id"`                                                                                                     // merchant id
	PlanName                  string      `json:"planName"                  description:"PlanName"`                                                                                                        // PlanName
	Amount                    int64       `json:"amount"                    description:"amount, cent, without tax"`                                                                                       // amount, cent, without tax
	Currency                  string      `json:"currency"                  description:"currency"`                                                                                                        // currency
	IntervalUnit              string      `json:"intervalUnit"              description:"period unit,day|month|year|week"`                                                                                 // period unit,day|month|year|week
	IntervalCount             int         `json:"intervalCount"             description:"period unit count"`                                                                                               // period unit count
	Description               string      `json:"description"               description:"description"`                                                                                                     // description
	ImageUrl                  string      `json:"imageUrl"                  description:"image_url"`                                                                                                       // image_url
	HomeUrl                   string      `json:"homeUrl"                   description:"home_url"`                                                                                                        // home_url
	GatewayProductName        string      `json:"gatewayProductName"        description:"gateway product name"`                                                                                            // gateway product name
	GatewayProductDescription string      `json:"gatewayProductDescription" description:"gateway product description"`                                                                                     // gateway product description
	TaxPercentage             int         `json:"taxPercentage"             description:"taxPercentage 1000 = 10%"`                                                                                        // taxPercentage 1000 = 10%
	TaxInclusive              int         `json:"taxInclusive"              description:"deperated"`                                                                                                       // deperated
	Type                      int         `json:"type"                      description:"type，1-main plan，2-recurring addon plan 3-onetime addon plan"`                                                    // type，1-main plan，2-recurring addon plan 3-onetime addon plan
	Status                    int         `json:"status"                    description:"status，1-editing，2-active，3-inactive，4-expired"`                                                                  // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 int         `json:"isDeleted"                 description:"0-UnDeleted，1-Deleted"`                                                                                           // 0-UnDeleted，1-Deleted
	BindingAddonIds           string      `json:"bindingAddonIds"           description:"binded recurring addon planIds，split with ,"`                                                                     // binded recurring addon planIds，split with ,
	BindingOnetimeAddonIds    string      `json:"bindingOnetimeAddonIds"    description:"binded onetime addon planIds，split with ,"`                                                                       // binded onetime addon planIds，split with ,
	PublishStatus             int         `json:"publishStatus"             description:"1-UnPublish,2-Publish, Use For Display Plan At UserPortal"`                                                       // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime                int64       `json:"createTime"                description:"create utc time"`                                                                                                 // create utc time
	ExtraMetricData           string      `json:"extraMetricData"           description:""`                                                                                                                //
	MetaData                  string      `json:"metaData"                  description:"meta_data(json)"`                                                                                                 // meta_data(json)
	GasPayer                  string      `json:"gasPayer"                  description:"who pay the gas, merchant|user"`                                                                                  // who pay the gas, merchant|user
	TrialAmount               int64       `json:"trialAmount"               description:"amount of trial, 0 for free"`                                                                                     // amount of trial, 0 for free
	TrialDurationTime         int64       `json:"trialDurationTime"         description:"duration of trial"`                                                                                               // duration of trial
	TrialDemand               string      `json:"trialDemand"               description:""`                                                                                                                //
	CancelAtTrialEnd          int         `json:"cancelAtTrialEnd"          description:"whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
}
