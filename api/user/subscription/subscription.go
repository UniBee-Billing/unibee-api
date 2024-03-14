package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
)

type DetailReq struct {
	g.Meta         `path:"/detail" tags:"User-Subscription" method:"get,post" summary:"Subscription Detail"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type DetailRes struct {
	User                                *bean.UserAccountSimplify               `json:"user" dc:"user"`
	Subscription                        *bean.SubscriptionSimplify              `json:"subscription" dc:"Subscription"`
	Plan                                *bean.PlanSimplify                      `json:"plan" dc:"Plan"`
	Gateway                             *bean.GatewaySimplify                   `json:"gateway" dc:"Gateway"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type PayCheckReq struct {
	g.Meta         `path:"/pay_check" tags:"User-Subscription" method:"get,post" summary:"Subscription Pay Status Check"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type PayCheckRes struct {
	PayStatus    consts.SubscriptionStatusEnum `json:"payStatus" dc:"Pay Status，1-Pending，2-Paid，3-Suspend，4-Cancel, 5-Expired"`
	Subscription *bean.SubscriptionSimplify    `json:"subscription" dc:"Subscription"`
}

type CreatePreviewReq struct {
	g.Meta         `path:"/create_preview" tags:"User-Subscription" method:"post" summary:"User Create Subscription Preview"`
	PlanId         uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity" `
	GatewayId      uint64                 `json:"gatewayId" dc:"Id" v:"required" `
	AddonParams    []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                 `json:"vatNumber" dc:"VatNumber" `
}
type CreatePreviewRes struct {
	Plan              *bean.PlanSimplify      `json:"plan"`
	Quantity          int64                   `json:"quantity"`
	Gateway           *bean.GatewaySimplify   `json:"gateway"`
	AddonParams       []*bean.PlanAddonParam  `json:"addonParams"`
	Addons            []*bean.PlanAddonDetail `json:"addons"`
	TotalAmount       int64                   `json:"totalAmount"                `
	Currency          string                  `json:"currency"              `
	Invoice           *bean.InvoiceSimplify   `json:"invoice"`
	UserId            int64                   `json:"userId" `
	Email             string                  `json:"email" `
	VatCountryCode    string                  `json:"vatCountryCode"              `
	VatCountryName    string                  `json:"vatCountryName"              `
	TaxScale          int64                   `json:"taxScale"              `
	VatNumber         string                  `json:"vatNumber"              `
	VatNumberValidate *bean.ValidResult       `json:"vatNumberValidate"              `
}

type CreateReq struct {
	g.Meta             `path:"/create_submit" tags:"User-Subscription" method:"post" summary:"User Create Subscription"`
	PlanId             uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity           int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          uint64                 `json:"gatewayId" dc:"Id"   v:"required" `
	AddonParams        []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                  `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                 `json:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl          string                 `json:"returnUrl"  dc:"RedirectUrl"  `
	VatCountryCode     string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                 `json:"vatNumber" dc:"VatNumber" `
	PaymentMethodId    string                 `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata           map[string]string      `json:"metadata" dc:"Metadata，Map"`
}

type CreateRes struct {
	Subscription *bean.SubscriptionSimplify `json:"subscription" dc:"Subscription"`
	Paid         bool                       `json:"paid"`
	Link         string                     `json:"link"`
}

type UpdatePreviewReq struct {
	g.Meta              `path:"/update_preview" tags:"User-Subscription" method:"post" summary:"User Update Subscription Preview"`
	SubscriptionId      string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           uint64                 `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity            int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId           uint64                 `json:"gatewayId" dc:"Id" `
	WithImmediateEffect int                    `json:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams         []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
}
type UpdatePreviewRes struct {
	TotalAmount       int64                 `json:"totalAmount"                `
	Currency          string                `json:"currency"              `
	Invoice           *bean.InvoiceSimplify `json:"invoice"`
	NextPeriodInvoice *bean.InvoiceSimplify `json:"nextPeriodInvoice"`
	ProrationDate     int64                 `json:"prorationDate"`
}

type UpdateReq struct {
	g.Meta              `path:"/update_submit" tags:"User-Subscription" method:"post" summary:"User Update Subscription"`
	SubscriptionId      string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId           uint64                 `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity            int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId           uint64                 `json:"gatewayId" dc:"Id" `
	AddonParams         []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount  int64                  `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency     string                 `json:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate       int64                  `json:"prorationDate" dc:"prorationDatem PaidDate Start Proration" v:"required" `
	WithImmediateEffect int                    `json:"withImmediateEffect" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	Metadata            map[string]string      `json:"metadata" dc:"Metadata，Map"`
}
type UpdateRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"SubscriptionPendingUpdate"`
	Paid                      bool                                    `json:"paid" dc:"Paid，true|false"`
	Link                      string                                  `json:"link" dc:"Pay Link"`
	Note                      string                                  `json:"note" dc:"note"`
}

type ListReq struct {
	g.Meta `path:"/list" tags:"User-Subscription" method:"get,post" summary:"Subscription List (Return Latest Active One - Later Deprecated) "`
}
type ListRes struct {
	Subscriptions []*detail.SubscriptionDetail `json:"subscriptions" dc:"Subscription List"`
}

type CancelReq struct {
	g.Meta         `path:"/cancel" tags:"User-Subscription" method:"post" summary:"User Cancel Subscription Immediately (Should In Create Status)"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelRes struct {
}

type CancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_at_period_end" tags:"User-Subscription" method:"post" summary:"User Edit Subscription-Set Cancel Ad Period End"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelAtPeriodEndRes struct {
}

type CancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_last_cancel_at_period_end" tags:"User-Subscription" method:"post" summary:"User Edit Subscription-Cancel Last CancelAtPeriod"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelLastCancelAtPeriodEndRes struct {
}

type SuspendReq struct {
	g.Meta         `path:"/suspend" tags:"User-Subscription" method:"post" summary:"User Edit Subscription-Stop"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SuspendRes struct {
}

type ResumeReq struct {
	g.Meta         `path:"/resume" tags:"User-Subscription" method:"post" summary:"User Edit Subscription-Resume"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type ResumeRes struct {
}

type TimeLineListReq struct {
	g.Meta    `path:"/timeline_list" tags:"User-Subscription-Timeline" method:"get,post" summary:"Subscription TimeLine List"`
	UserId    int    `json:"userId" dc:"Filter UserId, Default All " `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start WIth 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*detail.SubscriptionTimeLineDetail `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
}
