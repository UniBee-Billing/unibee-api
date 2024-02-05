package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"User-Subscription-Controller" method:"post" summary:"Subscription Detail"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionDetailRes struct {
	User                                *entity.UserAccount                 `json:"user" dc:"user"`
	Subscription                        *entity.Subscription                `json:"subscription" dc:"Subscription"`
	Plan                                *entity.SubscriptionPlan            `json:"plan" dc:"Plan"`
	Channel                             *ro.OutChannelRo                    `json:"channel" dc:"Channel"`
	Addons                              []*ro.SubscriptionPlanAddonRo       `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *ro.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionPayCheckReq struct {
	g.Meta         `path:"/subscription_pay_check" tags:"User-Subscription-Controller" method:"post" summary:"Subscription Pay Status Check"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionPayCheckRes struct {
	PayStatus    consts.SubscriptionStatusEnum `json:"payStatus" dc:"Pay Status，1-Pending，2-Paid，3-Suspend，4-Cancel, 5-Expired"`
	Subscription *entity.Subscription          `json:"subscription" dc:"Subscription"`
}

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"User-Subscription-Controller" method:"post" summary:"Query Subscription Support Gateway Channel"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}
type SubscriptionChannelsRes struct {
	Channels []*ro.OutChannelRo `json:"Channels"`
}

type SubscriptionCreatePreviewReq struct {
	g.Meta         `path:"/subscription_create_preview" tags:"User-Subscription-Controller" method:"post" summary:"User Create Subscription Preview"`
	PlanId         int64                              `p:"planId" dc:"PlanId" v:"required"`
	Quantity       int64                              `p:"quantity" dc:"Quantity，Default 1" `
	ChannelId      int64                              `p:"channelId" dc:"ChannelId"   v:"required" `
	UserId         int64                              `p:"userId" dc:"UserId" v:"required"`
	AddonParams    []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	VatCountryCode string                             `p:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                             `p:"vatNumber" dc:"VatNumber" `
}
type SubscriptionCreatePreviewRes struct {
	Plan              *entity.SubscriptionPlan           `json:"planId"`
	Quantity          int64                              `json:"quantity"`
	PayChannel        *ro.OutChannelRo                   `json:"payChannel"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount       int64                              `json:"totalAmount"                `
	Currency          string                             `json:"currency"              `
	Invoice           *ro.InvoiceDetailSimplify          `json:"invoice"`
	UserId            int64                              `json:"userId" `
	Email             string                             `json:"email" `
	VatCountryCode    string                             `json:"vatCountryCode"              `
	VatCountryName    string                             `json:"vatCountryName"              `
	TaxScale          int64                              `json:"taxScale"              `
	VatNumber         string                             `json:"vatNumber"              `
	VatNumberValidate *ro.ValidResult                    `json:"vatNumberValidate"              `
}

type SubscriptionCreateReq struct {
	g.Meta             `path:"/subscription_create_submit" tags:"User-Subscription-Controller" method:"post" summary:"User Create Subscription"`
	PlanId             int64                              `p:"planId" dc:"PlanId" v:"required"`
	Quantity           int64                              `p:"quantity" dc:"Quantity，Default 1" `
	ChannelId          int64                              `p:"channelId" dc:"ChannelId"   v:"required" `
	UserId             int64                              `p:"userId" dc:"UserId" v:"required"`
	AddonParams        []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                              `p:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                             `p:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl          string                             `p:"returnUrl"  dc:"RedirectUrl"  `
	VatCountryCode     string                             `p:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                             `p:"vatNumber" dc:"VatNumber" `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"Subscription"`
	Paid         bool                 `json:"paid"`
	Link         string               `json:"link"`
}

type SubscriptionUpdatePreviewReq struct {
	g.Meta              `path:"/subscription_update_preview" tags:"User-Subscription-Controller" method:"post" summary:"User Update Subscription Preview"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           int64                              `p:"newPlanId" dc:"NewPlanId" v:"required"`
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
	g.Meta              `path:"/subscription_update_submit" tags:"User-Subscription-Controller" method:"post" summary:"User Update Subscription"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           int64                              `p:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity            int64                              `p:"quantity" dc:"Quantity，Default 1" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount  int64                              `p:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency     string                             `p:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate       int64                              `p:"prorationDate" dc:"prorationDatem PaidDate Start Proration" v:"required" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
}
type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"SubscriptionPendingUpdate"`
	Paid                      bool                              `json:"paid" dc:"Paid，true|false"`
	Link                      string                            `json:"link" dc:"Pay Link"`
	Note                      string                            `json:"note" dc:"note"`
}

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"User-Subscription-Controller" method:"post" summary:"Subscription List (Return Latest Active One)"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int64 `p:"userId" dc:"UserId" v:"required|length:4,30" `
	//Status     int   `p:"status" dc:"Filter Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	//SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	//SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	//Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	//Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `p:"subscriptions" dc:"Subscription List"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"User-Subscription-Controller" method:"post" summary:"User Cancel Subscription Immediately (Should In Create Status)"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionUpdateCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_at_period_end" tags:"User-Subscription-Controller" method:"post" summary:"User Edit Subscription-Set Cancel Ad Period End"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionUpdateCancelAtPeriodEndRes struct {
}

type SubscriptionUpdateCancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_last_cancel_at_period_end" tags:"User-Subscription-Controller" method:"post" summary:"User Edit Subscription-Cancel Last CancelAtPeriod"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionUpdateCancelLastCancelAtPeriodEndRes struct {
}

type SubscriptionSuspendReq struct {
	g.Meta         `path:"/subscription_suspend" tags:"User-Subscription-Controller" method:"post" summary:"User Edit Subscription-Stop"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionSuspendRes struct {
}

type SubscriptionResumeReq struct {
	g.Meta         `path:"/subscription_resume" tags:"User-Subscription-Controller" method:"post" summary:"User Edit Subscription-Resume"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SubscriptionResumeRes struct {
}

type SubscriptionTimeLineListReq struct {
	g.Meta     `path:"/user_subscription_timeline_list" tags:"User-Subscription-Timeline-Controller" method:"post" summary:"User Subscription TimeLine List"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId, Default All " `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count" dc:"Count Of Page" `
}

type SubscriptionTimeLineListRes struct {
	SubscriptionTimeLines []*entity.SubscriptionTimeline `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
}
