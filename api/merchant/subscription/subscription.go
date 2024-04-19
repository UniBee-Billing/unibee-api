package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type ConfigReq struct {
	g.Meta `path:"/config" tags:"Subscription" method:"get" summary:"SubscriptionConfig"`
}
type ConfigRes struct {
	Config *bean.SubscriptionConfig `json:"config" dc:"Config"`
}

type ConfigUpdateReq struct {
	g.Meta                             `path:"/config/update" tags:"Subscription" method:"get" summary:"Update Merchant Subscription Config"`
	DowngradeEffectImmediately         *bool  `json:"downgradeEffectImmediately" dc:"DowngradeEffectImmediately, whether subscription downgrade should effect immediately or at period end, default at period end"`
	UpgradeProration                   *bool  `json:"upgradeProration" dc:"UpgradeProration, whether subscription update generation proration invoice or not, default yes"`
	IncompleteExpireTime               *int64 `json:"incompleteExpireTime" dc:"IncompleteExpireTime, em.. default 1day for plan of month type"`
	InvoiceEmail                       *bool  `json:"invoiceEmail" dc:"InvoiceEmail, whether to send invoice email to user, default yes"`
	TryAutomaticPaymentBeforePeriodEnd *int64 `json:"tryAutomaticPaymentBeforePeriodEnd" dc:"TryAutomaticPaymentBeforePeriodEnd, default 30 min"`
}

type ConfigUpdateRes struct {
	Config *bean.SubscriptionConfig `json:"config" dc:"Config"`
}

type DetailReq struct {
	g.Meta         `path:"/detail" tags:"Subscription" method:"get,post" summary:"SubscriptionDetail"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}

type DetailRes struct {
	User                                *bean.UserAccountSimplify               `json:"user" dc:"User"`
	Subscription                        *bean.SubscriptionSimplify              `json:"subscription" dc:"Subscription"`
	Plan                                *bean.PlanSimplify                      `json:"plan" dc:"Plan"`
	Gateway                             *bean.GatewaySimplify                   `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam                  `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	LatestInvoice                       *bean.InvoiceSimplify                   `json:"latestInvoice" dc:"LatestInvoice"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type ListReq struct {
	g.Meta    `path:"/list" tags:"Subscription" method:"get,post" summary:"SubscriptionList"`
	UserId    int64  `json:"userId"  dc:"UserId" `
	Status    []int  `json:"status" dc:"Filter, Default All，Status，0-Init | 1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page" dc:"Page, Start WIth 0" `
	Count     int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}
type ListRes struct {
	Subscriptions []*detail.SubscriptionDetail `json:"subscriptions" dc:"Subscriptions"`
}

type CancelReq struct {
	g.Meta         `path:"/cancel" tags:"Subscription" method:"post" summary:"CancelSubscriptionImmediately" dc:"Cancel subscription immediately, no proration invoice will generate"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	InvoiceNow     bool   `json:"invoiceNow" dc:"Default false"  deprecated:"true"`
	Prorate        bool   `json:"prorate" dc:"Prorate Generate Invoice，Default false"  deprecated:"true"`
}
type CancelRes struct {
}

type CancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_at_period_end" tags:"Subscription" method:"post" summary:"CancelSubscriptionAtPeriodEnd" dc:"Cancel subscription at period end, the subscription will not turn to 'cancelled' at once but will cancelled at period end time, no invoice will generate, the flag 'cancelAtPeriodEnd' of subscription will be enabled"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelAtPeriodEndRes struct {
}

type CancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_last_cancel_at_period_end" tags:"Subscription" method:"post" summary:"CancelLastCancelSubscriptionAtPeriodEnd" dc:"This action should be request before subscription's period end, If subscription's flag 'cancelAtPeriodEnd' is enabled, this action will resume it to disable, and subscription will continue cycle recurring seems no cancelAtPeriod action be setting"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type CancelLastCancelAtPeriodEndRes struct {
}

type SuspendReq struct {
	g.Meta         `path:"/suspend" tags:"Subscription" method:"post" summary:"Merchant Edit Subscription-Stop"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type SuspendRes struct {
}

type ResumeReq struct {
	g.Meta         `path:"/resume" tags:"Subscription" method:"post" summary:"Merchant Edit Subscription-Resume"  deprecated:"true"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}
type ResumeRes struct {
}

type ChangeGatewayReq struct {
	g.Meta          `path:"/change_gateway" tags:"Subscription" method:"post" summary:"ChangeSubscriptionGateway" `
	SubscriptionId  string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId" `
}
type ChangeGatewayRes struct {
}

type AddNewTrialStartReq struct {
	g.Meta             `path:"/add_new_trial_start" tags:"Subscription" method:"post" summary:"AppendSubscriptionTrialEnd"`
	SubscriptionId     string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	AppendTrialEndHour int64  `json:"appendTrialEndHour" dc:"add appendTrialEndHour For Free" v:"required"`
}
type AddNewTrialStartRes struct {
}

type RenewReq struct {
	g.Meta         `path:"/renew" tags:"Subscription" method:"post" summary:"RenewSubscription" dc:"renew an exist subscription "`
	UserId         uint64                      `json:"userId" dc:"UserId" v:"required"`
	SubscriptionId string                      `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	GatewayId      *uint64                     `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	TaxPercentage  *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	DiscountCode   string                      `json:"discountCode" dc:"DiscountCode, override subscription discount"`
	Discount       *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
}

type RenewRes struct {
	Subscription *bean.SubscriptionSimplify `json:"subscription" dc:"Subscription"`
	Paid         bool                       `json:"paid"`
	Link         string                     `json:"link"`
}

type CreatePreviewReq struct {
	g.Meta         `path:"/create_preview" tags:"Subscription" method:"post" summary:"CreateSubscriptionPreview"`
	PlanId         uint64                 `json:"planId" dc:"PlanId" v:"required"`
	UserId         uint64                 `json:"userId" dc:"UserId" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity" `
	GatewayId      uint64                 `json:"gatewayId" dc:"Id" v:"required" `
	AddonParams    []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage  *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	DiscountCode   string                 `json:"discountCode" dc:"DiscountCode"`
}

type CreatePreviewRes struct {
	Plan              *bean.PlanSimplify      `json:"plan"`
	Quantity          int64                   `json:"quantity"`
	Gateway           *bean.GatewaySimplify   `json:"gateway"`
	AddonParams       []*bean.PlanAddonParam  `json:"addonParams"`
	Addons            []*bean.PlanAddonDetail `json:"addons"`
	TotalAmount       int64                   `json:"totalAmount"                `
	DiscountAmount    int64                   `json:"discountAmount"`
	Currency          string                  `json:"currency"              `
	Invoice           *bean.InvoiceSimplify   `json:"invoice"`
	UserId            uint64                  `json:"userId" `
	Email             string                  `json:"email" `
	VatCountryCode    string                  `json:"vatCountryCode"              `
	VatCountryName    string                  `json:"vatCountryName"              `
	TaxPercentage     int64                   `json:"taxPercentage"              `
	VatNumber         string                  `json:"vatNumber"              `
	VatNumberValidate *bean.ValidResult       `json:"vatNumberValidate"   `
}

type CreateReq struct {
	g.Meta             `path:"/create_submit" tags:"Subscription" method:"post" summary:"CreateSubscription"`
	PlanId             uint64                      `json:"planId" dc:"PlanId" v:"required"`
	UserId             uint64                      `json:"userId" dc:"UserId" v:"required"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          uint64                      `json:"gatewayId" dc:"Id"   v:"required" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"            `
	ConfirmCurrency    string                      `json:"confirmCurrency"  dc:"Currency to verify if provide" `
	ReturnUrl          string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	VatCountryCode     string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	PaymentMethodId    string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata           map[string]string           `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
}

type CreateRes struct {
	Subscription *bean.SubscriptionSimplify `json:"subscription" dc:"Subscription"`
	Paid         bool                       `json:"paid"`
	Link         string                     `json:"link"`
}

type UpdatePreviewReq struct {
	g.Meta          `path:"/update_preview" tags:"Subscription" method:"post" summary:"UpdateSubscriptionPreview"`
	SubscriptionId  string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId       uint64                 `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity        int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId       uint64                 `json:"gatewayId" dc:"Id" `
	EffectImmediate int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
}
type UpdatePreviewRes struct {
	TotalAmount       int64                 `json:"totalAmount"                `
	DiscountAmount    int64                 `json:"discountAmount"`
	Currency          string                `json:"currency"              `
	Invoice           *bean.InvoiceSimplify `json:"invoice"`
	NextPeriodInvoice *bean.InvoiceSimplify `json:"nextPeriodInvoice"`
	ProrationDate     int64                 `json:"prorationDate"`
}

type UpdateReq struct {
	g.Meta             `path:"/update_submit" tags:"Subscription" method:"post" summary:"UpdateSubscription"`
	SubscriptionId     string                      `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId          uint64                      `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity           int64                       `json:"quantity" dc:"Quantity"  v:"required"`
	GatewayId          uint64                      `json:"gatewayId" dc:"Id" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	EffectImmediate    int                         `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"          `
	ConfirmCurrency    string                      `json:"confirmCurrency" dc:"Currency to verify if provide"   `
	ProrationDate      *int64                      `json:"prorationDate" dc:"The utc time to start Proration, default current time" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	Metadata           map[string]string           `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                      `json:"discountCode" dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
}

type UpdateRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
	Paid                      bool                                    `json:"paid"`
	Link                      string                                  `json:"link"`
}

type UserSubscriptionDetailReq struct {
	g.Meta `path:"/user_subscription_detail" tags:"Subscription" method:"get,post" summary:"SubscriptionDetail"`
	UserId uint64 `json:"userId" dc:"UserId" v:"required"`
}

type UserSubscriptionDetailRes struct {
	User                                *bean.UserAccountSimplify               `json:"user" dc:"user"`
	Subscription                        *bean.SubscriptionSimplify              `json:"subscription" dc:"Subscription"`
	Plan                                *bean.PlanSimplify                      `json:"plan" dc:"Plan"`
	Gateway                             *bean.GatewaySimplify                   `json:"gateway" dc:"Gateway"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"Processing Subscription Pending Update"`
}

type TimeLineListReq struct {
	g.Meta    `path:"/timeline_list" tags:"Subscription-Timeline" method:"get,post" summary:"SubscriptionTimeLineList"`
	UserId    uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start WIth 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*detail.SubscriptionTimeLineDetail `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
}

type PendingUpdateListReq struct {
	g.Meta         `path:"/pending_update_list" tags:"SubscriptionPendingUpdate" method:"get,post" summary:"SubscriptionPendingUpdateList"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page           int    `json:"page"  dc:"Page, Start WIth 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type PendingUpdateListRes struct {
	SubscriptionPendingUpdateDetails []*detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
}

type NewAdminNoteReq struct {
	g.Meta         `path:"/new_admin_note" tags:"Subscription-Note" method:"post" summary:"NewSubscriptionNote"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Note           string `json:"note" dc:"Note" v:"required"`
}

type NewAdminNoteRes struct {
}

type AdminNoteRo struct {
	Id             uint64 `json:"id"               description:"Id"`
	Note           string `json:"note"             description:"Note"`
	CreateTime     int64  `json:"createTime"       description:"CreateTime, UTC Time"`
	SubscriptionId string `json:"subscriptionId" description:"SubscriptionId"`
	UserName       string `json:"userName"   description:"UserName"`
	Mobile         string `json:"mobile"     description:"Mobile"`
	Email          string `json:"email"      description:"Email"`
	FirstName      string `json:"firstName"  description:"FirstName"`
	LastName       string `json:"lastName"   description:"LastName"`
}

type AdminNoteListReq struct {
	g.Meta         `path:"/admin_note_list" tags:"Subscription-Note" method:"get,post" summary:"SubscriptionNoteList"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Page           int    `json:"page"  dc:"Page, Start WIth 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type AdminNoteListRes struct {
	NoteLists []*AdminNoteRo `json:"noteLists"   description:""`
}

type OnetimeAddonNewReq struct {
	g.Meta             `path:"/new_onetime_addon_payment" tags:"Subscription" method:"post" summary:"NewSubscriptionOnetimeAddonPayment" dc:"Create payment for subscription onetime addon purchase"`
	SubscriptionId     string            `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached" v:"required"`
	AddonId            uint64            `json:"addonId" dc:"AddonId, id of one-time addon, the new payment will created base on the addon's amount'" v:"required"`
	Quantity           int64             `json:"quantity" dc:"Quantity, quantity of the new payment which one-time addon purchased"  v:"required"`
	ReturnUrl          string            `json:"returnUrl"  dc:"ReturnUrl, the addon's payment will redirect based on the returnUrl provided when it's back from gateway side"  `
	Metadata           map[string]string `json:"metadata" dc:"Metadata，custom data"`
	DiscountCode       string            `json:"discountCode" dc:"DiscountCode"`
	DiscountAmount     *int64            `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage *int64            `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount provide"`
	GatewayId          *uint64           `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	TaxPercentage      *int64            `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, use subscription's taxPercentage if not provide"`
}

type OnetimeAddonNewRes struct {
	SubscriptionOnetimeAddon *bean.SubscriptionOnetimeAddonSimplify `json:"subscriptionOnetimeAddon" dc:"SubscriptionOnetimeAddon, object of onetime-addon purchased"`
	Paid                     bool                                   `json:"paid" dc:"true|false,automatic payment is default behavior for one-time addon purchased, payment will create attach to the purchase, when payment is success, return false, otherwise false"`
	Link                     string                                 `json:"link" dc:"if automatic payment is false, Gateway Link will provided that manual payment needed"`
	Invoice                  *bean.InvoiceSimplify                  `json:"invoice" dc:"invoice of one-time payment"`
}

type OnetimeAddonListReq struct {
	g.Meta `path:"/onetime_addon_list" tags:"Subscription" method:"get" summary:"SubscriptionOnetimeAddonList"`
	UserId uint64 `json:"userId" dc:"UserId" v:"required"`
	Page   int    `json:"page"  dc:"Page, Start With 0" `
	Count  int    `json:"count" dc:"Count Of Page" `
}

type OnetimeAddonListRes struct {
	SubscriptionOnetimeAddons []*detail.SubscriptionOnetimeAddonDetail `json:"subscriptionOnetimeAddons" description:"SubscriptionOnetimeAddons" `
}
