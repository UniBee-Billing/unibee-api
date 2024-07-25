package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
)

type UserCurrentSubscriptionDetailReq struct {
	g.Meta `path:"/current/detail" tags:"User-Subscription" method:"get" summary:"User Current Subscription Detail"`
}
type UserCurrentSubscriptionDetailRes struct {
	User                                *bean.UserAccount                       `json:"user" dc:"user"`
	Subscription                        *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Plan                                *bean.Plan                              `json:"plan" dc:"Plan"`
	Gateway                             *bean.Gateway                           `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam                  `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	LatestInvoice                       *bean.Invoice                           `json:"latestInvoice" dc:"LatestInvoice"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type DetailReq struct {
	g.Meta         `path:"/detail" tags:"User-Subscription" method:"get,post" summary:"Subscription Detail"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}

type DetailRes struct {
	User                                *bean.UserAccount                       `json:"user" dc:"user"`
	Subscription                        *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Plan                                *bean.Plan                              `json:"plan" dc:"Plan"`
	Gateway                             *bean.Gateway                           `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam                  `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	LatestInvoice                       *bean.Invoice                           `json:"latestInvoice" dc:"LatestInvoice"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type PayCheckReq struct {
	g.Meta         `path:"/pay_check" tags:"User-Subscription" method:"get,post" summary:"Subscription Pay Status Check"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type PayCheckRes struct {
	PayStatus    consts.SubscriptionStatusEnum `json:"payStatus" dc:"Pay Status，1-Pending，2-Paid，3-Suspend，4-Cancel, 5-Expired"`
	Subscription *bean.Subscription            `json:"subscription" dc:"Subscription"`
}

type CreatePreviewReq struct {
	g.Meta         `path:"/create_preview" tags:"User-Subscription" method:"post" summary:"User Create Subscription Preview"`
	PlanId         uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity" `
	GatewayId      *uint64                `json:"gatewayId" dc:"Id" `
	AddonParams    []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                 `json:"vatNumber" dc:"VatNumber" `
	DiscountCode   string                 `json:"discountCode"        dc:"DiscountCode"`
}
type CreatePreviewRes struct {
	Plan                     *bean.Plan                 `json:"plan"`
	TrialEnd                 int64                      `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	Quantity                 int64                      `json:"quantity"`
	Gateway                  *bean.Gateway              `json:"gateway"`
	AddonParams              []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                   []*bean.PlanAddonDetail    `json:"addons"`
	OriginAmount             int64                      `json:"originAmount"                `
	TotalAmount              int64                      `json:"totalAmount"                `
	DiscountAmount           int64                      `json:"discountAmount"`
	Currency                 string                     `json:"currency"              `
	Invoice                  *bean.Invoice              `json:"invoice"`
	UserId                   uint64                     `json:"userId" `
	Email                    string                     `json:"email" `
	VatCountryCode           string                     `json:"vatCountryCode"              `
	VatCountryName           string                     `json:"vatCountryName"              `
	TaxPercentage            int64                      `json:"taxPercentage"              `
	VatNumber                string                     `json:"vatNumber"              `
	VatNumberValidate        *bean.ValidResult          `json:"vatNumberValidate"              `
	Discount                 *bean.MerchantDiscountCode `json:"discount" `
	VatNumberValidateMessage string                     `json:"vatNumberValidateMessage" `
	DiscountMessage          string                     `json:"discountMessage" `
}

type CreateReq struct {
	g.Meta             `path:"/create_submit" tags:"User-Subscription" method:"post" summary:"User Create Subscription"`
	PlanId             uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity           int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          *uint64                `json:"gatewayId" dc:"Id" `
	AddonParams        []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                  `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                 `json:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl          string                 `json:"returnUrl"  dc:"RedirectUrl"  `
	VatCountryCode     string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                 `json:"vatNumber" dc:"VatNumber" `
	PaymentMethodId    string                 `json:"paymentMethodId" dc:"PaymentMethodId" `
	DiscountCode       string                 `json:"discountCode"        dc:"DiscountCode"`
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type CreateRes struct {
	Subscription *bean.Subscription `json:"subscription" dc:"Subscription"`
	Paid         bool               `json:"paid"`
	Link         string             `json:"link"`
}

type UpdatePreviewReq struct {
	g.Meta          `path:"/update_preview" tags:"User-Subscription" method:"post" summary:"User Update Subscription Preview"`
	SubscriptionId  string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId       uint64                 `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity        int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId       *uint64                `json:"gatewayId" dc:"Id" `
	EffectImmediate int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
}
type UpdatePreviewRes struct {
	OriginAmount      int64                      `json:"originAmount"                `
	TotalAmount       int64                      `json:"totalAmount"                `
	DiscountAmount    int64                      `json:"discountAmount"`
	Currency          string                     `json:"currency"              `
	Invoice           *bean.Invoice              `json:"invoice"`
	NextPeriodInvoice *bean.Invoice              `json:"nextPeriodInvoice"`
	ProrationDate     int64                      `json:"prorationDate"`
	Discount          *bean.MerchantDiscountCode `json:"discount" `
}

type UpdateReq struct {
	g.Meta             `path:"/update_submit" tags:"User-Subscription" method:"post" summary:"User Update Subscription"`
	SubscriptionId     string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId          uint64                 `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity           int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          *uint64                `json:"gatewayId" dc:"Id" `
	AddonParams        []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                  `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                 `json:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate      *int64                 `json:"prorationDate" dc:"The utc time to start Proration, default current time" `
	EffectImmediate    int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                 `json:"discountCode"        dc:"DiscountCode"`
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

type ChangeGatewayReq struct {
	g.Meta          `path:"/change_gateway" tags:"User-Subscription" method:"post" summary:"Change Subscription Gateway" `
	SubscriptionId  string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId" `
}
type ChangeGatewayRes struct {
}

type TimeLineListReq struct {
	g.Meta    `path:"/timeline_list" tags:"User-Subscription-Timeline" method:"get,post" summary:"Subscription TimeLine List"`
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start With 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*detail.SubscriptionTimeLineDetail `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
	Total                 int                                  `json:"total" dc:"Total"`
}

type OnetimeAddonNewReq struct {
	g.Meta         `path:"/new_onetime_addon_payment" tags:"User-Subscription" method:"post" summary:"New Subscription Onetime Addon Payment"`
	SubscriptionId string                 `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached" v:"required"`
	AddonId        uint64                 `json:"addonId" dc:"AddonId, id of one-time addon, the new payment will created base on the addon's amount'" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity, quantity of the new payment which one-time addon purchased"  v:"required"`
	ReturnUrl      string                 `json:"returnUrl"  dc:"ReturnUrl, the addon's payment will redirect based on the returnUrl provided when it's back from gateway side"  `
	Metadata       map[string]interface{} `json:"metadata" dc:"Metadata，custom data"`
	DiscountCode   string                 `json:"discountCode"        dc:"DiscountCode"`
	GatewayId      *uint64                `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
}

type OnetimeAddonNewRes struct {
	SubscriptionOnetimeAddon *bean.SubscriptionOnetimeAddon `json:"subscriptionOnetimeAddon" dc:"SubscriptionOnetimeAddon, object of onetime-addon purchased"`
	Paid                     bool                           `json:"paid" dc:"true|false,automatic payment is default behavior for one-time addon purchased, payment will create attach to the purchase, when payment is success, return false, otherwise false"`
	Link                     string                         `json:"link" dc:"if automatic payment is false, Gateway Link will provided that manual payment needed"`
	Invoice                  *bean.Invoice                  `json:"invoice" dc:"invoice of one-time payment"`
}

type OnetimeAddonListReq struct {
	g.Meta `path:"/onetime_addon_list" tags:"User-Subscription" method:"get" summary:"Subscription OnetimeAddon List"`
	Page   int `json:"page"  dc:"Page, Start With 0" `
	Count  int `json:"count" dc:"Count Of Page" `
}

type OnetimeAddonListRes struct {
	SubscriptionOnetimeAddons []*detail.SubscriptionOnetimeAddonDetail `json:"subscriptionOnetimeAddons" description:"SubscriptionOnetimeAddons" `
}

type MarkWireTransferPaidReq struct {
	g.Meta         `path:"/mark_wire_transfer_paid" tags:"User-Subscription" method:"post" summary:"MarkWireTransferInvoiceSuccess" dc:"Mark wire transfer subscription as paid, subscription will change to 8-Processed "`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}

type MarkWireTransferPaidRes struct {
}

type UserPendingCryptoSubscriptionDetailReq struct {
	g.Meta `path:"/user_pending_crypto_subscription_detail" tags:"User-Subscription" method:"get" summary:"UserPendingCryptoSubscriptionDetail"`
}

type UserPendingCryptoSubscriptionDetailRes struct {
	Subscription *detail.SubscriptionDetail `json:"subscription" dc:"Subscription"`
}
