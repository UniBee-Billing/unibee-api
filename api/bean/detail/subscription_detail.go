package detail

import (
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
)

type SubscriptionDetail struct {
	User                                *bean.UserAccountSimplify        `json:"user" dc:"user"`
	Subscription                        *bean.SubscriptionSimplify       `json:"subscription" dc:"Subscription"`
	Plan                                *bean.PlanSimplify               `json:"plan" dc:"Plan"`
	Gateway                             *bean.GatewaySimplify            `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam           `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail          `json:"addons" dc:"Addon"`
	UnfinishedSubscriptionPendingUpdate *SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionTimeLineDetail struct {
	MerchantId      uint64                  `json:"merchantId"      description:"merchant id"`             // merchant id
	UserId          int64                   `json:"userId"          description:"userId"`                  // userId
	SubscriptionId  string                  `json:"subscriptionId"  description:"subscription id"`         // subscription id
	PeriodStart     int64                   `json:"periodStart"     description:"period_start"`            // period_start
	PeriodEnd       int64                   `json:"periodEnd"       description:"period_end"`              // period_end
	PeriodStartTime *gtime.Time             `json:"periodStartTime" description:"period start (datetime)"` // period start (datetime)
	PeriodEndTime   *gtime.Time             `json:"periodEndTime"   description:"period end (datatime)"`   // period end (datatime)
	InvoiceId       string                  `json:"invoiceId"       description:"invoice id"`              // invoice id
	UniqueId        string                  `json:"uniqueId"        description:"unique id"`               // unique id
	Currency        string                  `json:"currency"        description:"currency"`                // currency
	PlanId          uint64                  `json:"planId"          description:"PlanId"`                  // PlanId
	Plan            *bean.PlanSimplify      `json:"plan" description:"Plan"`
	Quantity        int64                   `json:"quantity"        description:"quantity"` // quantity
	Addons          []*bean.PlanAddonDetail `json:"addons" description:"Addon"`
	GatewayId       uint64                  `json:"gatewayId"       description:"gateway_id"`      // gateway_id
	CreateTime      int64                   `json:"createTime"      description:"create utc time"` // create utc time
}

type SubscriptionPendingUpdateDetail struct {
	MerchantId           uint64                       `json:"merchantId"           description:"MerchantId"`
	SubscriptionId       string                       `json:"subscriptionId"       description:"SubscriptionId"`
	UpdateSubscriptionId string                       `json:"updateSubscriptionId" description:"UpdateSubscriptionId"`
	GmtCreate            *gtime.Time                  `json:"gmtCreate"            description:"GmtCreate"`
	Amount               int64                        `json:"amount"               description:"CaptureAmount, Cent"`
	Status               int                          `json:"status"               description:"Status，0-Init | 1-Create｜2-Finished｜3-Cancelled"`
	UpdateAmount         int64                        `json:"updateAmount"         description:"UpdateAmount, Cents"`
	ProrationAmount      int64                        `json:"prorationAmount"      description:"ProrationAmount,Cents"`
	Currency             string                       `json:"currency"             description:"Currency"`
	UpdateCurrency       string                       `json:"updateCurrency"       description:"UpdateCurrency"`
	PlanId               uint64                       `json:"planId"               description:"PlanId"`
	UpdatePlanId         uint64                       `json:"updatePlanId"         description:"UpdatePlanId"`
	Quantity             int64                        `json:"quantity"             description:"quantity"`
	UpdateQuantity       int64                        `json:"updateQuantity"       description:"UpdateQuantity"`
	AddonData            string                       `json:"addonData"            description:"plan addon json data"`
	UpdateAddonData      string                       `json:"updateAddonData"     description:"UpdateAddonData"`
	GatewayId            uint64                       `json:"gatewayId"            description:"Id"`
	UserId               int64                        `json:"userId"               description:"UserId"`
	GmtModify            *gtime.Time                  `json:"gmtModify"            description:"GmtModify"`
	Paid                 int                          `json:"paid"                 description:"Paid"`
	Link                 string                       `json:"link"                 description:"Link"`
	MerchantMember       *bean.MerchantMemberSimplify `json:"merchantMember"       description:"Merchant Member"`
	EffectImmediate      int                          `json:"effectImmediate"      description:"EffectImmediate"`
	EffectTime           int64                        `json:"effectTime"           description:"effect_immediate=0, EffectTime unit_time"`
	Note                 string                       `json:"note"            description:"Update Note"`
	Plan                 *bean.PlanSimplify           `json:"plan" dc:"Plan"`
	Addons               []*bean.PlanAddonDetail      `json:"addons" dc:"Addons"`
	UpdatePlan           *bean.PlanSimplify           `json:"updatePlan" dc:"UpdatePlan"`
	UpdateAddons         []*bean.PlanAddonDetail      `json:"updateAddons" dc:"UpdateAddons"`
	Metadata             map[string]string            `json:"metadata" description:""`
}