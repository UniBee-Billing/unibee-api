package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"Merchant-Subscription-Controller" method:"post" summary:"Subscription Detail"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionDetailRes struct {
	User                                *entity.UserAccount                 `json:"user" dc:"User"`
	Subscription                        *entity.Subscription                `json:"subscription" dc:"Subscription"`
	Plan                                *entity.SubscriptionPlan            `json:"plan" dc:"Plan"`
	Gateway                             *ro.OutGatewayRo                    `json:"gateway" dc:"Gateway"`
	Addons                              []*ro.SubscriptionPlanAddonRo       `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *ro.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"Merchant-Subscription-Controller" method:"post" summary:"Subscription List"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int64  `p:"userId"  dc:"UserId" `
	Status     []int  `p:"status" dc:"Filter, Default All，Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page" dc:"Page, Start WIth 0" `
	Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `json:"subscriptions" dc:"Subscriptions"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Cancel Subscription Immediately (Will Not Generate Proration Invoice)"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	InvoiceNow     bool   `p:"invoiceNow" dc:"Default false"  deprecated:"true"`
	Prorate        bool   `p:"prorate" dc:"Prorate Generate Invoice，Default false"  deprecated:"true"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionUpdateCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_at_period_end" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Edit Subscription-Set Cancel Ad Period End"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionUpdateCancelAtPeriodEndRes struct {
}

type SubscriptionUpdateCancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_last_cancel_at_period_end" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Edit Subscription-Cancel Last CancelAtPeriod"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionUpdateCancelLastCancelAtPeriodEndRes struct {
}

type SubscriptionSuspendReq struct {
	g.Meta         `path:"/subscription_suspend" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Edit Subscription-Stop"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionSuspendRes struct {
}

type SubscriptionResumeReq struct {
	g.Meta         `path:"/subscription_resume" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Edit Subscription-Resume"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionResumeRes struct {
}

type SubscriptionAddNewTrialStartReq struct {
	g.Meta             `path:"/subscription_add_new_trial_start" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Edit Subscription-add appendTrialEndHour For Free"`
	SubscriptionId     string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	AppendTrialEndHour int64  `p:"appendTrialEndHour" dc:"add appendTrialEndHour For Free" v:"required"`
}
type SubscriptionAddNewTrialStartRes struct {
}

type SubscriptionUpdatePreviewReq struct {
	g.Meta              `path:"/subscription_update_preview" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Update Subscription Preview"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           int64                              `p:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity            int64                              `p:"quantity" dc:"Quantity，Default 1" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
}
type SubscriptionUpdatePreviewRes struct {
	TotalAmount       int64                     `json:"totalAmount"                `
	Currency          string                    `json:"currency"              `
	Invoice           *ro.InvoiceDetailSimplify `json:"invoice"`
	NextPeriodInvoice *ro.InvoiceDetailSimplify `json:"nextPeriodInvoice"`
	ProrationDate     int64                     `json:"prorationDate"`
}

type SubscriptionUpdateReq struct {
	g.Meta              `path:"/subscription_update_submit" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant Update Subscription Submit"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           int64                              `p:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity            int64                              `p:"quantity" dc:"Quantity，Default 1" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	ConfirmTotalAmount  int64                              `p:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency     string                             `p:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate       int64                              `p:"prorationDate" dc:"prorationDate date to start Proration，Get From Preview" v:"required" `
}

type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
	Paid                      bool                              `json:"paid"`
	Link                      string                            `json:"link"`
}

type UserSubscriptionDetailReq struct {
	g.Meta     `path:"/user_subscription_detail" tags:"Merchant-Subscription-Controller" method:"post" summary:"Subscription Detail"`
	UserId     int64  `p:"userId" dc:"UserId" v:"required"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
}

type UserSubscriptionDetailRes struct {
	User                                *entity.UserAccount                 `json:"user" dc:"user"`
	Subscription                        *entity.Subscription                `json:"subscription" dc:"Subscription"`
	Plan                                *entity.SubscriptionPlan            `json:"plan" dc:"Plan"`
	Gateway                             *ro.OutGatewayRo                    `json:"gateway" dc:"Gateway"`
	Addons                              []*ro.SubscriptionPlanAddonRo       `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *ro.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"Processing Subscription Pending Update"`
}

type SubscriptionTimeLineListReq struct {
	g.Meta     `path:"/subscription_timeline_list" tags:"Merchant-Subscription-Timeline-Controller" method:"post" summary:"Merchant Subscription TimeLine List"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId, Default All " `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count" dc:"Count Of Page" `
}

type SubscriptionTimeLineListRes struct {
	SubscriptionTimeLines []*entity.SubscriptionTimeline `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
}

type SubscriptionMerchantPendingUpdateListReq struct {
	g.Meta         `path:"/subscription_merchant_pending_update_list" tags:"Merchant-SubscriptionPendingUpdate-Controller" method:"post" summary:"Merchant-SubscriptionPendingUpdate List"`
	MerchantId     int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	SortField      string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType       string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page           int    `p:"page"  dc:"Page, Start WIth 0" `
	Count          int    `p:"count" dc:"Count Of Page" `
}

type SubscriptionMerchantPendingUpdateListRes struct {
	SubscriptionPendingUpdateDetails []*ro.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
}

type SubscriptionNewAdminNoteReq struct {
	g.Meta         `path:"/subscription_new_admin_note" tags:"Merchant-Subscription-Note-Controller" method:"post" summary:"Merchant New Subscription Note"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	MerchantUserId int64  `p:"merchantUserId" dc:"MerchantUserId" v:"required"`
	Note           string `p:"note" dc:"note" v:"required"`
}

type SubscriptionNewAdminNoteRes struct {
}

type SubscriptionAdminNoteRo struct {
	GmtCreate      *gtime.Time `json:"gmtCreate"  description:"创建时间"`               // 创建时间
	GmtModify      *gtime.Time `json:"gmtModify"  description:"修改时间"`               // 修改时间
	SubscriptionId string      `json:"subscriptionId" description:"SubscriptionId"` // 用户ID
	UserName       string      `json:"userName"   description:"用户名"`                // 用户名
	Mobile         string      `json:"mobile"     description:"手机号"`                // 手机号
	Email          string      `json:"email"      description:"邮箱"`                 // 邮箱
	FirstName      string      `json:"firstName"  description:""`                   //
	LastName       string      `json:"lastName"   description:""`                   //
}

type SubscriptionAdminNoteListReq struct {
	g.Meta         `path:"/subscription_admin_note_list" tags:"Merchant-Subscription-Note-Controller" method:"post" summary:"Merchant Subscription Note List"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Page           int    `p:"page"  dc:"Page, Start WIth 0" `
	Count          int    `p:"count" dc:"Count Of Page" `
}

type SubscriptionAdminNoteListRes struct {
	NoteLists []*SubscriptionAdminNoteRo `json:"noteLists"   description:""`
}
