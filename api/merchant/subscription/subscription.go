package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type DetailReq struct {
	g.Meta         `path:"/detail" tags:"Merchant-Subscription" method:"get,post" summary:"Subscription Detail"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type DetailRes struct {
	User                                *ro.UserAccountSimplify               `json:"user" dc:"User"`
	Subscription                        *ro.SubscriptionSimplify              `json:"subscription" dc:"Subscription"`
	Plan                                *ro.PlanSimplify                      `json:"plan" dc:"Plan"`
	Gateway                             *ro.GatewaySimplify                   `json:"gateway" dc:"Gateway"`
	Addons                              []*ro.PlanAddonVo                     `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *ro.SubscriptionPendingUpdateDetailVo `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type ListReq struct {
	g.Meta    `path:"/list" tags:"Merchant-Subscription" method:"get,post" summary:"Subscription List"`
	UserId    int64  `json:"userId"  dc:"UserId" `
	Status    []int  `json:"status" dc:"Filter, Default All，Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page" dc:"Page, Start WIth 0" `
	Count     int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}
type ListRes struct {
	Subscriptions []*ro.SubscriptionDetailVo `json:"subscriptions" dc:"Subscriptions"`
}

type CancelReq struct {
	g.Meta         `path:"/cancel" tags:"Merchant-Subscription" method:"post" summary:"Merchant Cancel Subscription Immediately (Will Not Generate Proration Invoice)"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	InvoiceNow     bool   `json:"invoiceNow" dc:"Default false"  deprecated:"true"`
	Prorate        bool   `json:"prorate" dc:"Prorate Generate Invoice，Default false"  deprecated:"true"`
}
type CancelRes struct {
}

type CancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_at_period_end" tags:"Merchant-Subscription" method:"post" summary:"Merchant Edit Subscription-Set Cancel Ad Period End"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelAtPeriodEndRes struct {
}

type CancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_last_cancel_at_period_end" tags:"Merchant-Subscription" method:"post" summary:"Merchant Edit Subscription-Cancel Last CancelAtPeriod"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelLastCancelAtPeriodEndRes struct {
}

type SuspendReq struct {
	g.Meta         `path:"/suspend" tags:"Merchant-Subscription" method:"post" summary:"Merchant Edit Subscription-Stop"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SuspendRes struct {
}

type ResumeReq struct {
	g.Meta         `path:"/resume" tags:"Merchant-Subscription" method:"post" summary:"Merchant Edit Subscription-Resume"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type ResumeRes struct {
}

type AddNewTrialStartReq struct {
	g.Meta             `path:"/add_new_trial_start" tags:"Merchant-Subscription" method:"post" summary:"Merchant Edit Subscription-add appendTrialEndHour For Free"`
	SubscriptionId     string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	AppendTrialEndHour int64  `json:"appendTrialEndHour" dc:"add appendTrialEndHour For Free" v:"required"`
}
type AddNewTrialStartRes struct {
}

type UpdatePreviewReq struct {
	g.Meta              `path:"/update_preview" tags:"Merchant-Subscription" method:"post" summary:"Merchant Update Subscription Preview"`
	SubscriptionId      string                             `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           uint64                             `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity            int64                              `json:"quantity" dc:"Quantity，Default 1" `
	WithImmediateEffect int                                `json:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `json:"addonParams" dc:"addonParams" `
}
type UpdatePreviewRes struct {
	TotalAmount       int64                     `json:"totalAmount"                `
	Currency          string                    `json:"currency"              `
	Invoice           *ro.InvoiceDetailSimplify `json:"invoice"`
	NextPeriodInvoice *ro.InvoiceDetailSimplify `json:"nextPeriodInvoice"`
	ProrationDate     int64                     `json:"prorationDate"`
}

type UpdateReq struct {
	g.Meta              `path:"/update_submit" tags:"Merchant-Subscription" method:"post" summary:"Merchant Update Subscription Submit"`
	SubscriptionId      string                             `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           uint64                             `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity            int64                              `json:"quantity" dc:"Quantity，Default 1" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `json:"addonParams" dc:"addonParams" `
	WithImmediateEffect int                                `json:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	ConfirmTotalAmount  int64                              `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency     string                             `json:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate       int64                              `json:"prorationDate" dc:"prorationDate date to start Proration，Get From Preview" v:"required" `
}

type UpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
	Paid                      bool                              `json:"paid"`
	Link                      string                            `json:"link"`
}

type UserSubscriptionDetailReq struct {
	g.Meta `path:"/user_subscription_detail" tags:"Merchant-Subscription" method:"get,post" summary:"Subscription Detail"`
	UserId int64 `json:"userId" dc:"UserId" v:"required"`
}

type UserSubscriptionDetailRes struct {
	User                                *ro.UserAccountSimplify               `json:"user" dc:"user"`
	Subscription                        *ro.SubscriptionSimplify              `json:"subscription" dc:"Subscription"`
	Plan                                *ro.PlanSimplify                      `json:"plan" dc:"Plan"`
	Gateway                             *ro.GatewaySimplify                   `json:"gateway" dc:"Gateway"`
	Addons                              []*ro.PlanAddonVo                     `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *ro.SubscriptionPendingUpdateDetailVo `json:"unfinishedSubscriptionPendingUpdate" dc:"Processing Subscription Pending Update"`
}

type TimeLineListReq struct {
	g.Meta    `path:"/timeline_list" tags:"Merchant-Subscription-Timeline" method:"get,post" summary:"Merchant Subscription TimeLine List"`
	UserId    int    `json:"userId" dc:"Filter UserId, Default All " `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start WIth 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*ro.SubscriptionTimeLineDetailVo `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
}

type PendingUpdateListReq struct {
	g.Meta         `path:"/pending_update_list" tags:"Merchant-SubscriptionPendingUpdate" method:"get,post" summary:"Merchant SubscriptionPendingUpdate List"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page           int    `json:"page"  dc:"Page, Start WIth 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type PendingUpdateListRes struct {
	SubscriptionPendingUpdateDetails []*ro.SubscriptionPendingUpdateDetailVo `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
}

type NewAdminNoteReq struct {
	g.Meta         `path:"/new_admin_note" tags:"Merchant-Subscription-Note" method:"post" summary:"Merchant New Subscription Note"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Note           string `json:"note" dc:"note" v:"required"`
}

type NewAdminNoteRes struct {
}

type AdminNoteRo struct {
	Id             uint64      `json:"id"               description:"id"`           // id
	Note           string      `json:"note"             description:"note"`         // note
	GmtCreate      *gtime.Time `json:"gmtCreate"  description:"创建时间"`               // 创建时间
	GmtModify      *gtime.Time `json:"gmtModify"  description:"修改时间"`               // 修改时间
	SubscriptionId string      `json:"subscriptionId" description:"SubscriptionId"` // 用户ID
	UserName       string      `json:"userName"   description:"用户名"`                // 用户名
	Mobile         string      `json:"mobile"     description:"手机号"`                // 手机号
	Email          string      `json:"email"      description:"邮箱"`                 // 邮箱
	FirstName      string      `json:"firstName"  description:""`                   //
	LastName       string      `json:"lastName"   description:""`                   //
}

type AdminNoteListReq struct {
	g.Meta         `path:"/admin_note_list" tags:"Merchant-Subscription-Note" method:"get,post" summary:"Merchant Subscription Note List"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Page           int    `json:"page"  dc:"Page, Start WIth 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type AdminNoteListRes struct {
	NoteLists []*AdminNoteRo `json:"noteLists"   description:""`
}
