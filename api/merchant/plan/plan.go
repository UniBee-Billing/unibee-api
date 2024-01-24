package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionPlanCreateReq struct {
	g.Meta             `path:"/subscription_plan_create" tags:"Merchant-Plan-Controller" method:"post" summary:"Create Plan"`
	MerchantId         int64   `p:"merchantId" dc:"MerchantId" v:"required"`
	PlanName           string  `p:"planName" dc:"Plan Name"   v:"required" `
	Amount             int64   `p:"amount"   dc:"Plan Amount"   v:"required" `
	Currency           string  `p:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit       string  `p:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week" v:"required" `
	IntervalCount      int     `p:"intervalCount"  d:"1" dc:"Default 1，Number Of IntervalUnit" `
	Type               int     `p:"type"  d:"1"  dc:"Default 1，,1-main plan，2-addon plan" `
	Description        string  `p:"description"  dc:"Description"`
	ProductName        string  `p:"productName" dc:"Default Copy PlanName"  `
	ProductDescription string  `p:"productDescription" dc:"Default Copy Description" `
	ImageUrl           string  `p:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl            string  `p:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds           []int64 `p:"addonIds"  dc:"Plan Ids Of Addon Type" `
}
type SubscriptionPlanCreateRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"Plan"`
}

type SubscriptionPlanEditReq struct {
	g.Meta             `path:"/subscription_plan_edit" tags:"Merchant-Plan-Controller" method:"post" summary:"Edit Plan"`
	PlanId             int64   `p:"planId" dc:"PlanId" v:"required"`
	PlanName           string  `p:"planName" dc:"Plan Name"   v:"required" `
	Amount             int64   `p:"amount"   dc:"Plan Amount"   v:"required" `
	Currency           string  `p:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit       string  `p:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week" v:"required" `
	IntervalCount      int     `p:"intervalCount"  d:"1" dc:"Default 1，Number Of IntervalUnit" `
	Description        string  `p:"description"  dc:"Description"`
	ProductName        string  `p:"productName" dc:"Default Copy PlanName"  `
	ProductDescription string  `p:"productDescription" dc:"Default Copy Description" `
	ImageUrl           string  `p:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl            string  `p:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds           []int64 `p:"addonIds"  dc:"Plan Ids Of Addon Type" `
}
type SubscriptionPlanEditRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"Plan"`
}

type SubscriptionPlanAddonsBindingReq struct {
	g.Meta   `path:"/subscription_plan_addons_binding" tags:"Merchant-Plan-Controller" method:"post" summary:"Plan Binding Addons"`
	PlanId   int64   `p:"planId" dc:"PlanID" v:"required"`
	Action   int64   `p:"action" dc:"Action Type，0-override,1-add，2-delete" v:"required"`
	AddonIds []int64 `p:"addonIds"  dc:"Plan Ids Of Addon Type"  v:"required" `
}
type SubscriptionPlanAddonsBindingRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"Plan"`
}

type SubscriptionPlanListReq struct {
	g.Meta        `path:"/subscription_plan_list" tags:"Merchant-Plan-Controller" method:"post" summary:"订阅计划列表"`
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Type          int    `p:"type"  d:"1"  dc:"Default 1，,1-main plan，2-addon plan" `
	Status        int    `p:"status" dc:"Filter, Default All，,Status，1-Editing，2-Active，3-InActive，4-Expired" `
	PublishStatus int    `p:"publishStatus" dc:"Filter, Default All，PublishStatus，1-UnPublished，2-Published" `
	Currency      string `p:"currency" dc:"Filter Currency"  `
	SortField     string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `p:"page"  dc:"Page, Start 0" `
	Count         int    `p:"count"  dc:"Count Of Per Page" `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.PlanDetailRo `p:"plans" dc:"Plans"`
}

type SubscriptionPlanChannelTransferAndActivateReq struct {
	g.Meta `path:"/subscription_plan_activate" tags:"Merchant-Plan-Controller" method:"post" summary:"Plan Sync To Gateway And Activate"`
	PlanId int64 `p:"planId" dc:"PlanId" v:"required"`
	//ChannelId int64 `p:"channelId"    v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionPlanChannelTransferAndActivateRes struct {
}

type SubscriptionPlanChannelActivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_activate" tags:"Merchant-Plan-Controller" method:"post" summary:"Plan Activate "  deprecated:"true" `
	PlanId    int64 `p:"planId" dc:"PlanId" v:"required"`
	ChannelId int64 `p:"channelId"    v:"required" `
}
type SubscriptionPlanChannelActivateRes struct {
}

type SubscriptionPlanChannelDeactivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_deactivate" tags:"Merchant-Plan-Controller" method:"post" summary:"Plan DeActivate" deprecated:"true" `
	PlanId    int64 `p:"planId" dc:"PlanId" v:"required"`
	ChannelId int64 `p:"channelId"    v:"required" `
}
type SubscriptionPlanChannelDeactivateRes struct {
}

type SubscriptionPlanPublishReq struct {
	g.Meta `path:"/subscription_plan_publish" tags:"Merchant-Plan-Controller" method:"post" summary:"Publish Plan，Will Be Visible To UserPortal" `
	PlanId int64 `p:"planId" dc:"PlanId" v:"required"`
}
type SubscriptionPlanPublishRes struct {
}

type SubscriptionPlanUnPublishReq struct {
	g.Meta `path:"/subscription_plan_unpublished" tags:"Merchant-Plan-Controller" method:"post" summary:"UnPublish Plan" `
	PlanId int64 `p:"planId" dc:"PlanId" v:"required"`
}
type SubscriptionPlanUnPublishRes struct {
}

type SubscriptionPlanDetailReq struct {
	g.Meta `path:"/subscription_plan_detail" tags:"Merchant-Plan-Controller" method:"post" summary:"Plan Detail"`
	PlanId int64 `p:"planId" dc:"PlanId" v:"required"`
}
type SubscriptionPlanDetailRes struct {
	Plan *ro.PlanDetailRo `p:"plan" dc:"Plan Detail"`
}

type SubscriptionPlanExpireReq struct {
	g.Meta    `path:"/subscription_plan_expire" tags:"Merchant-Plan-Controller" method:"post" summary:"Expired Plan"`
	PlanId    int64 `p:"planId" dc:"PlanId" v:"required"`
	EmailCode int64 `p:"emailCode" dc:"Code From Email" v:"required"`
}
type SubscriptionPlanExpireRes struct {
}
