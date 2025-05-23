package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type NewReq struct {
	g.Meta                `path:"/new" tags:"Plan" method:"post" summary:"Create Plan"`
	ExternalPlanId        string                               `json:"externalPlanId" dc:"ExternalPlanId"`
	PlanName              string                               `json:"planName" dc:"Plan Name"   v:"required" `
	InternalName          string                               `json:"internalName"              description:""`
	Amount                int64                                `json:"amount"   dc:"Plan CaptureAmount"   v:"required" `
	Currency              string                               `json:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit          string                               `json:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week"`
	IntervalCount         int                                  `json:"intervalCount"  dc:"Number Of IntervalUnit，em: day|month|year|week"`
	Description           string                               `json:"description"  dc:"Description"`
	Type                  int                                  `json:"type"  d:"1"  dc:"The type of plan, 1-main plan，2-addon plan, 3-onetime plan, default main plan" `
	ProductName           string                               `json:"productName" dc:"Default Copy PlanName"  `
	ProductDescription    string                               `json:"productDescription" dc:"Default Copy Description" `
	ImageUrl              string                               `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl               string                               `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds              []int64                              `json:"addonIds"  dc:"Plan Ids Of Recurring Addon Type" `
	OnetimeAddonIds       []int64                              `json:"onetimeAddonIds"  dc:"Plan Ids Of Onetime Addon Type" `
	MetricLimits          []*bean.PlanMetricLimitParam         `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge   []*bean.PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*bean.PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
	GasPayer              string                               `json:"gasPayer" dc:"who pay the gas for crypto payment, merchant|user"`
	Metadata              map[string]interface{}               `json:"metadata" dc:"Metadata，Map"`
	TrialAmount           int64                                `json:"trialAmount"                description:"price of trial period， not available for addon"` // price of trial period
	TrialDurationTime     int64                                `json:"trialDurationTime"         description:"duration of trial， not available for addon"`      // duration of trial
	TrialDemand           string                               `json:"trialDemand"               description:"demand of trial， not available for addon, example, paymentMethod, payment method will ask for subscription trial start"`
	CancelAtTrialEnd      int                                  `json:"cancelAtTrialEnd"          description:"whether cancel at subscription first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ProductId             int64                                `json:"productId"   dc:"Id of product which plan to linked" `
}
type NewRes struct {
	Plan *bean.Plan `json:"plan" dc:"Plan"`
}

type EditReq struct {
	g.Meta                `path:"/edit" tags:"Plan" method:"post" summary:"Edit Plan" dc:"Edit exist plan, amount|currency|intervalUnit|intervalCount is not editable when plan is active "`
	PlanId                uint64                                `json:"planId" dc:"Id of plan" v:"required"`
	ExternalPlanId        *string                               `json:"externalPlanId" dc:"ExternalPlanId"`
	PlanName              *string                               `json:"planName" dc:"Name of plan" `
	InternalName          *string                               `json:"internalName"              description:""`
	Amount                *int64                                `json:"amount"   dc:"CaptureAmount of plan, not editable when plan is active" `
	Currency              *string                               `json:"currency"   dc:"Currency of plan, not editable when plan is active"`
	IntervalUnit          *string                               `json:"intervalUnit" dc:"Interval unit of plan，em: day|month|year|week, not editable when plan is active"`
	IntervalCount         *int                                  `json:"intervalCount"  dc:"Number,intervalUnit of plan, not editable when plan is active" `
	Description           *string                               `json:"description"  dc:"Description of plan"`
	ProductName           *string                               `json:"productName" dc:"ProductName of plan, Default copy planName"  `
	ProductDescription    *string                               `json:"productDescription" dc:"ProductDescription of plan, Default copy description" `
	ImageUrl              *string                               `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl               *string                               `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds              []int64                               `json:"addonIds"  dc:"Plan Ids Of Recurring Addon Type" `
	OnetimeAddonIds       []int64                               `json:"onetimeAddonIds"  dc:"Plan Ids Of Onetime Addon Type" `
	MetricLimits          []*bean.PlanMetricLimitParam          `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge   *[]*bean.PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge *[]*bean.PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
	GasPayer              *string                               `json:"gasPayer" dc:"who pay the gas for crypto payment, merchant|user"`
	Metadata              *map[string]interface{}               `json:"metadata" dc:"Metadata，Map"`
	TrialAmount           *int64                                `json:"trialAmount"                description:"price of trial period， not available for addon"` // price of trial period
	TrialDurationTime     *int64                                `json:"trialDurationTime"         description:"duration of trial， not available for addon"`      // duration of trial
	TrialDemand           *string                               `json:"trialDemand"               description:"demand of trial, not available for addon, example, paymentMethod, payment method will ask for subscription trial start"`
	CancelAtTrialEnd      *int                                  `json:"cancelAtTrialEnd"          description:"whether cancel at subscription first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ProductId             *int64                                `json:"productId"   dc:"Id of product which plan to linked" `
}
type EditRes struct {
	Plan *bean.Plan `json:"plan" dc:"Plan"`
}

type AddonsBindingReq struct {
	g.Meta          `path:"/addons_binding" tags:"Plan" method:"post" summary:"Addon Binding"`
	PlanId          uint64  `json:"planId" dc:"PlanID" v:"required"`
	Action          int64   `json:"action" dc:"Action Type，0-override,1-add，2-delete" v:"required"`
	AddonIds        []int64 `json:"addonIds"  dc:"Plan Ids Of Recurring Addon Type"  v:"required" `
	OnetimeAddonIds []int64 `json:"onetimeAddonIds"  dc:"Plan Ids Of Onetime Addon Type"   v:"required" `
}
type AddonsBindingRes struct {
	Plan *bean.Plan `json:"plan" dc:"Plan"`
}

type ListReq struct {
	g.Meta        `path:"/list" tags:"Plan" method:"get,post" summary:"Get Plan List"`
	PlanIds       []int64 `json:"planIds"  dc:"filter id list of plan, default all" `
	ProductIds    []int64 `json:"productIds"  dc:"filter id list of product, default all" `
	Type          []int   `json:"type"  dc:"1-main plan，2-addon plan,3-onetime" `
	Status        []int   `json:"status" dc:"Filter, Default All，,Status，1-Editing，2-Active，3-InActive，4-SoftArchive, 5-HardArchive" `
	PublishStatus int     `json:"publishStatus" dc:"Filter, Default All，PublishStatus，1-UnPublished，2-Published" `
	Currency      string  `json:"currency" dc:"Filter Currency"  `
	SearchKey     string  `json:"searchKey" dc:"Search Key, plan name or description"  `
	SortField     string  `json:"sortField" dc:"Sort Field，plan_name|gmt_create|gmt_modify，Default gmt_create" `
	SortType      string  `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int     `json:"page"  dc:"Page, Start 0" `
	Count         int     `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Plans []*detail.PlanDetail `json:"plans" dc:"Plans"`
	Total int                  `json:"total" dc:"Total"`
}

type CopyReq struct {
	g.Meta `path:"/copy" tags:"Plan" method:"post" summary:"Copy Plan"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type CopyRes struct {
	Plan *bean.Plan `json:"plan" dc:"Plan"`
}

type ActivateReq struct {
	g.Meta `path:"/activate" tags:"Plan" method:"post" summary:"Activate Plan"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type ActivateRes struct {
}

type PublishReq struct {
	g.Meta `path:"/publish" tags:"Plan" method:"post" summary:"Publish Plan" dc:"Publish plan，a plan will display at user portal when its published"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type PublishRes struct {
}

type UnPublishReq struct {
	g.Meta `path:"/unpublished" tags:"Plan" method:"post" summary:"UnPublish Plan" `
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type UnPublishRes struct {
}

type DetailReq struct {
	g.Meta `path:"/detail" tags:"Plan" method:"get,post" summary:"Plan Detail"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type DetailRes struct {
	Plan *detail.PlanDetail `json:"plan" dc:"Plan Detail"`
}

type ArchiveReq struct {
	g.Meta      `path:"/archive" tags:"Plan" method:"post" summary:"Archive Plan"`
	PlanId      uint64 `json:"planId" dc:"PlanId" v:"required"`
	HardArchive bool   `json:"hardArchive" dc:"Hard Archive"`
}
type ArchiveRes struct {
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Plan" method:"post" summary:"Delete Plan"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type DeleteRes struct {
}
