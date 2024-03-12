package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewReq struct {
	g.Meta             `path:"/new" tags:"Plan" method:"post" summary:"Create Plan"`
	PlanName           string                                  `json:"planName" dc:"Plan Name"   v:"required" `
	Amount             int64                                   `json:"amount"   dc:"Plan CaptureAmount"   v:"required" `
	Currency           string                                  `json:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit       string                                  `json:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week" v:"required" `
	IntervalCount      int                                     `json:"intervalCount"  dc:"Number Of IntervalUnit，em: day|month|year|week" v:"required" `
	Description        string                                  `json:"description"  dc:"Description"`
	Type               int                                     `json:"type"  d:"1"  dc:"Default 1，,1-main plan，2-addon plan" `
	ProductName        string                                  `json:"productName" dc:"Default Copy PlanName"  `
	ProductDescription string                                  `json:"productDescription" dc:"Default Copy Description" `
	ImageUrl           string                                  `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl            string                                  `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds           []int64                                 `json:"addonIds"  dc:"Plan Ids Of Addon Type" `
	MetricLimits       []*bean.BulkMetricLimitPlanBindingParam `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	GasPayer           string                                  `json:"gasPayer" dc:"who pay the gas, merchant|user"`
	Metadata           map[string]string                       `json:"metadata" dc:"Metadata，Map"`
}
type NewRes struct {
	Plan *bean.PlanSimplify `json:"plan" dc:"Plan"`
}

type EditReq struct {
	g.Meta             `path:"/edit" tags:"Plan" method:"post" summary:"Edit Plan"`
	PlanId             uint64                                  `json:"planId" dc:"PlanId" v:"required"`
	PlanName           string                                  `json:"planName" dc:"Plan Name"   v:"required" `
	Amount             int64                                   `json:"amount"   dc:"Plan CaptureAmount"   v:"required" `
	Currency           string                                  `json:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit       string                                  `json:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week" v:"required" `
	IntervalCount      int                                     `json:"intervalCount"  dc:"Number Of IntervalUnit" `
	Description        string                                  `json:"description"  dc:"Description"`
	ProductName        string                                  `json:"productName" dc:"Default Copy PlanName"  `
	ProductDescription string                                  `json:"productDescription" dc:"Default Copy Description" `
	ImageUrl           string                                  `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl            string                                  `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds           []int64                                 `json:"addonIds"  dc:"Plan Ids Of Addon Type" `
	MetricLimits       []*bean.BulkMetricLimitPlanBindingParam `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	GasPayer           string                                  `json:"gasPayer" dc:"who pay the gas, merchant|user"`
	Metadata           map[string]string                       `json:"metadata" dc:"Metadata，Map"`
}
type EditRes struct {
	Plan *bean.PlanSimplify `json:"plan" dc:"Plan"`
}

type AddonsBindingReq struct {
	g.Meta   `path:"/addons_binding" tags:"Plan" method:"post" summary:"Plan Binding Addons"`
	PlanId   uint64  `json:"planId" dc:"PlanID" v:"required"`
	Action   int64   `json:"action" dc:"Action Type，0-override,1-add，2-delete" v:"required"`
	AddonIds []int64 `json:"addonIds"  dc:"Plan Ids Of Addon Type"  v:"required" `
}
type AddonsBindingRes struct {
	Plan *bean.PlanSimplify `json:"plan" dc:"Plan"`
}

type ListReq struct {
	g.Meta        `path:"/list" tags:"Plan" method:"get,post" summary:"Plan List"`
	Type          []int  `json:"type"  dc:"1-main plan，2-addon plan" `
	Status        []int  `json:"status" dc:"Filter, Default All，,Status，1-Editing，2-Active，3-InActive，4-Expired" `
	PublishStatus int    `json:"publishStatus" dc:"Filter, Default All，PublishStatus，1-UnPublished，2-Published" `
	Currency      string `json:"currency" dc:"Filter Currency"  `
	SortField     string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `json:"page"  dc:"Page, Start 0" `
	Count         int    `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Plans []*PlanDetail `json:"plans" dc:"Plans"`
}

type ActivateReq struct {
	g.Meta `path:"/activate" tags:"Plan" method:"post" summary:"Plan Sync To Gateway And Activate"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type ActivateRes struct {
}

type PublishReq struct {
	g.Meta `path:"/publish" tags:"Plan" method:"post" summary:"Publish Plan，Will Be Visible To UserPortal" `
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
	Plan *PlanDetail `json:"plan" dc:"Plan Detail"`
}

type PlanDetail struct {
	Plan             *bean.PlanSimplify              `json:"plan" dc:"Plan"`
	MetricPlanLimits []*bean.MerchantMetricPlanLimit `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	Addons           []*bean.PlanSimplify            `json:"addons" dc:"Addons"`
	AddonIds         []int64                         `json:"addonIds" dc:"AddonIds"`
}

type ExpireReq struct {
	g.Meta    `path:"/expire" tags:"Plan" method:"post" summary:"Expire A Plan"`
	PlanId    uint64 `json:"planId" dc:"PlanId" v:"required"`
	EmailCode int64  `json:"emailCode" dc:"Code From Email" v:"required"`
}
type ExpireRes struct {
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Plan" method:"post" summary:"Delete A Plan Before Activate"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type DeleteRes struct {
}
