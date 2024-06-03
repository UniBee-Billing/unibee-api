package subscription

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/payment"
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

type UserPendingCryptoSubscriptionDetailReq struct {
	g.Meta         `path:"/user_pending_crypto_subscription_detail" tags:"Subscription" method:"get,post" summary:"UserPendingCryptoSubscriptionDetail"`
	UserId         uint64 `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
}

type UserPendingCryptoSubscriptionDetailRes struct {
	Subscription *detail.SubscriptionDetail `json:"subscription" dc:"Subscription"`
}

type ListReq struct {
	g.Meta    `path:"/list" tags:"Subscription" method:"get,post" summary:"SubscriptionList"`
	UserId    int64  `json:"userId"  dc:"UserId" `
	Status    []int  `json:"status" dc:"Filter, Default All，Status，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing" `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page" dc:"Page, Start With 0" `
	Count     int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}
type ListRes struct {
	Subscriptions []*detail.SubscriptionDetail `json:"subscriptions" dc:"Subscriptions"`
	Total         int                          `json:"total" dc:"Total"`
}

type CancelReq struct {
	g.Meta         `path:"/cancel" tags:"Subscription" method:"post" summary:"CancelSubscriptionImmediately" dc:"Cancel subscription immediately, no proration invoice will generate"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
	InvoiceNow     bool   `json:"invoiceNow" dc:"Default false"  deprecated:"true"`
	Prorate        bool   `json:"prorate" dc:"Prorate Generate Invoice，Default false"  deprecated:"true"`
}
type CancelRes struct {
}

type CancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_at_period_end" tags:"Subscription" method:"post" summary:"CancelSubscriptionAtPeriodEnd" dc:"Cancel subscription at period end, the subscription will not turn to 'cancelled' at once but will cancelled at period end time, no invoice will generate, the flag 'cancelAtPeriodEnd' of subscription will be enabled"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
}
type CancelAtPeriodEndRes struct {
}

type CancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_last_cancel_at_period_end" tags:"Subscription" method:"post" summary:"CancelLastCancelSubscriptionAtPeriodEnd" dc:"This action should be request before subscription's period end, If subscription's flag 'cancelAtPeriodEnd' is enabled, this action will resume it to disable, and subscription will continue cycle recurring seems no cancelAtPeriod action be setting"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
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
	SubscriptionId string                      `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached, either SubscriptionId or UserId needed, The only one active subscription or latest subscription will renew if userId provide instead of subscriptionId"`
	UserId         uint64                      `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription or latest cancel|expire subscription will renew if userId provide instead of subscriptionId"`
	GatewayId      *uint64                     `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	TaxPercentage  *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	DiscountCode   string                      `json:"discountCode" dc:"DiscountCode, override subscription discount"`
	Discount       *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment  bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl      string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	ProductData    *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
}

type RenewRes struct {
	Subscription *bean.SubscriptionSimplify `json:"subscription" dc:"Subscription"`
	Paid         bool                       `json:"paid"`
	Link         string                     `json:"link"`
}

type CreatePreviewReq struct {
	g.Meta         `path:"/create_preview" tags:"Subscription" method:"post" summary:"CreateSubscriptionPreview"`
	ExternalUserId string                 `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	Email          string                 `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId         uint64                 `json:"userId" dc:"UserId"`
	PlanId         uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity" `
	GatewayId      *uint64                `json:"gatewayId" dc:"GatewayId" `
	AddonParams    []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage  *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	DiscountCode   string                 `json:"discountCode" dc:"DiscountCode"`
	TrialEnd       int64                  `json:"trialEnd" dc:"trial_end, utc time"` // trial_end, utc time
}

type CreatePreviewRes struct {
	Plan                     *bean.PlanSimplify                 `json:"plan"`
	TrialEnd                 int64                              `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	Quantity                 int64                              `json:"quantity"`
	Gateway                  *bean.GatewaySimplify              `json:"gateway"`
	AddonParams              []*bean.PlanAddonParam             `json:"addonParams"`
	Addons                   []*bean.PlanAddonDetail            `json:"addons"`
	OriginAmount             int64                              `json:"originAmount"                `
	TotalAmount              int64                              `json:"totalAmount"                `
	DiscountAmount           int64                              `json:"discountAmount"`
	Currency                 string                             `json:"currency"              `
	Invoice                  *bean.InvoiceSimplify              `json:"invoice"`
	UserId                   uint64                             `json:"userId" `
	Email                    string                             `json:"email" `
	VatCountryCode           string                             `json:"vatCountryCode"              `
	VatCountryName           string                             `json:"vatCountryName"              `
	TaxPercentage            int64                              `json:"taxPercentage"              `
	VatNumber                string                             `json:"vatNumber"              `
	VatNumberValidate        *bean.ValidResult                  `json:"vatNumberValidate"   `
	Discount                 *bean.MerchantDiscountCodeSimplify `json:"discount" `
	VatNumberValidateMessage string                             `json:"vatNumberValidateMessage" `
	DiscountMessage          string                             `json:"discountMessage" `
}

type CreateReq struct {
	g.Meta             `path:"/create_submit" tags:"Subscription" method:"post" summary:"CreateSubscription"`
	PlanId             uint64                      `json:"planId" dc:"PlanId" v:"required"`
	ExternalUserId     string                      `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	Email              string                      `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId             uint64                      `json:"userId" dc:"UserId"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          *uint64                     `json:"gatewayId" dc:"GatewayId" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"            `
	ConfirmCurrency    string                      `json:"confirmCurrency"  dc:"Currency to verify if provide" `
	ReturnUrl          string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	VatCountryCode     string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	PaymentMethodId    string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata           map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	TrialEnd           int64                       `json:"trialEnd"                    dc:"trial_end, utc time"` // trial_end, utc time
	StartIncomplete    bool                        `json:"startIncomplete"        dc:"StartIncomplete, use now pay later, subscription will generate invoice and start with incomplete status if set"`
	ProductData        *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
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
	OriginAmount      int64                              `json:"originAmount"                `
	TotalAmount       int64                              `json:"totalAmount"                `
	DiscountAmount    int64                              `json:"discountAmount"`
	Currency          string                             `json:"currency"              `
	Invoice           *bean.InvoiceSimplify              `json:"invoice"`
	NextPeriodInvoice *bean.InvoiceSimplify              `json:"nextPeriodInvoice"`
	ProrationDate     int64                              `json:"prorationDate"`
	Discount          *bean.MerchantDiscountCodeSimplify `json:"discount" `
}

type UpdateReq struct {
	g.Meta             `path:"/update_submit" tags:"Subscription" method:"post" summary:"UpdateSubscription"`
	SubscriptionId     string                      `json:"subscriptionId" dc:"SubscriptionId, either SubscriptionId or UserId needed, The only one active subscription of userId will update"`
	UserId             uint64                      `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will update if userId provide instead of subscriptionId"`
	NewPlanId          uint64                      `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity           int64                       `json:"quantity" dc:"Quantity"  v:"required"`
	GatewayId          *uint64                     `json:"gatewayId" dc:"Id of gateway" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	EffectImmediate    int                         `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"          `
	ConfirmCurrency    string                      `json:"confirmCurrency" dc:"Currency to verify if provide"   `
	ProrationDate      *int64                      `json:"prorationDate" dc:"The utc time to start Proration, default current time" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	Metadata           map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                      `json:"discountCode" dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment      bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl          string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	ProductData        *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
}

type UpdateRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
	Paid                      bool                                    `json:"paid"`
	Link                      string                                  `json:"link"`
}

type UserSubscriptionDetailReq struct {
	g.Meta         `path:"/user_subscription_detail" tags:"Subscription" method:"get,post" summary:"UserSubscriptionDetail"`
	UserId         uint64 `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
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
	Page      int    `json:"page"  dc:"Page, Start With 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*detail.SubscriptionTimeLineDetail `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
	Total                 int                                  `json:"total" dc:"Total"`
}

type PendingUpdateListReq struct {
	g.Meta         `path:"/pending_update_list" tags:"SubscriptionPendingUpdate" method:"get,post" summary:"SubscriptionPendingUpdateList"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type PendingUpdateListRes struct {
	SubscriptionPendingUpdateDetails []*detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
	Total                            int                                       `json:"total" dc:"Total"`
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

type ActiveTemporarilyReq struct {
	g.Meta         `path:"/active_temporarily" tags:"Subscription" method:"post" summary:"SubscriptionActiveTemporarily" dc:"Subscription active temporarily, status will transmit from pending to incomplete"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	ExpireTime     int64  `json:"expireTime"  dc:"ExpireTime, the expire utc time if not paid"  v:"required"`
}

type ActiveTemporarilyRes struct {
}

type AdminNoteListReq struct {
	g.Meta         `path:"/admin_note_list" tags:"Subscription-Note" method:"get,post" summary:"SubscriptionNoteList"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type AdminNoteListRes struct {
	NoteLists []*AdminNoteRo `json:"noteLists"   description:""`
}

type OnetimeAddonNewReq struct {
	g.Meta             `path:"/new_onetime_addon_payment" tags:"Subscription" method:"post" summary:"NewSubscriptionOnetimeAddonPayment" dc:"Create payment for subscription onetime addon purchase"`
	SubscriptionId     string                 `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached, either SubscriptionId or UserId needed, The only one active subscription of userId will attach the addon"`
	UserId             uint64                 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will update if userId provide instead of subscriptionId"`
	AddonId            uint64                 `json:"addonId" dc:"AddonId, id of one-time addon, the new payment will created base on the addon's amount'" v:"required"`
	Quantity           int64                  `json:"quantity" dc:"Quantity, quantity of the new payment which one-time addon purchased"  v:"required"`
	ReturnUrl          string                 `json:"returnUrl"  dc:"ReturnUrl, the addon's payment will redirect based on the returnUrl provided when it's back from gateway side"  `
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，custom data"`
	DiscountCode       string                 `json:"discountCode" dc:"DiscountCode"`
	DiscountAmount     *int64                 `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage *int64                 `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount provide"`
	GatewayId          *uint64                `json:"gatewayId" dc:"GatewayId, use user's gateway if not provide"`
	TaxPercentage      *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, use subscription's taxPercentage if not provide"`
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

type NewPaymentReq struct {
	g.Meta            `path:"/payment/new" tags:"Subscription" method:"post" summary:"NewSubscriptionPayment"`
	ExternalPaymentId string                 `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment"`
	ExternalUserId    string                 `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	Email             string                 `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId            uint64                 `json:"userId" dc:"UserId, either ExternalUserId&Email or UserId needed"`
	Currency          string                 `json:"currency" dc:"Currency, either Currency&TotalAmount or PlanId needed" `
	TotalAmount       int64                  `json:"totalAmount" dc:"Total PaymentAmount, Cent, either TotalAmount&Currency or PlanId needed"`
	PlanId            uint64                 `json:"planId" dc:"PlanId, either TotalAmount&Currency or PlanId needed"`
	GatewayId         uint64                 `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl       string                 `json:"redirectUrl" dc:"Redirect Url"`
	CountryCode       string                 `json:"countryCode" dc:"CountryCode"`
	Name              string                 `json:"name" dc:"Name"`
	Description       string                 `json:"description" dc:"Description"`
	Items             []*payment.Item        `json:"items" dc:"Items"`
	Metadata          map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	GasPayer          string                 `json:"gasPayer" dc:"who pay the gas, merchant|user"`
}

type NewPaymentRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"The unique id of payment"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"The external unique id of payment"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}
