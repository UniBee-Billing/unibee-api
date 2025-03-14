package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type DetailReq struct {
	g.Meta         `path:"/detail" tags:"Subscription" method:"get,post" summary:"Subscription Detail"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}

type DetailRes struct {
	User                                *bean.UserAccount                       `json:"user" dc:"User"`
	Subscription                        *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Plan                                *bean.Plan                              `json:"plan" dc:"Plan"`
	Gateway                             *detail.Gateway                         `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam                  `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	LatestInvoice                       *bean.Invoice                           `json:"latestInvoice" dc:"LatestInvoice"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type UserPendingCryptoSubscriptionDetailReq struct {
	g.Meta         `path:"/user_pending_crypto_subscription_detail" tags:"Subscription" method:"get,post" summary:"User Pending Crypto Subscription Detail"`
	UserId         uint64 `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type UserPendingCryptoSubscriptionDetailRes struct {
	Subscription *detail.SubscriptionDetail `json:"subscription" dc:"Subscription"`
}

type ListReq struct {
	g.Meta          `path:"/list" tags:"Subscription" method:"get,post" summary:"Get Subscription List"`
	UserId          int64    `json:"userId"  dc:"UserId" `
	Status          []int    `json:"status" dc:"Filter, Default All，Status，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing | 9-Failed" `
	Currency        string   `json:"currency" dc:"The currency of subscription" `
	PlanIds         []uint64 `json:"planIds" dc:"The filter ids of plan" `
	ProductIds      []int64  `json:"productIds" dc:"The filter ids of product" `
	AmountStart     *int64   `json:"amountStart" dc:"The filter start amount of subscription" `
	AmountEnd       *int64   `json:"amountEnd" dc:"The filter end amount of subscription" `
	SortField       string   `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string   `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int      `json:"page" dc:"Page, Start With 0" `
	Count           int      `json:"count"  dc:"Count" dc:"Count Of Page" `
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd" `
}
type ListRes struct {
	Subscriptions []*detail.SubscriptionDetail `json:"subscriptions" dc:"Subscriptions"`
	Total         int                          `json:"total" dc:"Total"`
}

type CancelReq struct {
	g.Meta         `path:"/cancel" tags:"Subscription" method:"post" summary:"Cancel Subscription Immediately" dc:"Cancel subscription immediately, no proration invoice will generate"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
	ProductId      int64  `json:"productId" dc:"default product will use if productId not specified and subscriptionId is blank"`
	Reason         string `json:"reason" dc:"Reason"`
	InvoiceNow     bool   `json:"invoiceNow" dc:"Default false"  deprecated:"true"`
	Prorate        bool   `json:"prorate" dc:"Prorate Generate Invoice，Default false"  deprecated:"true"`
}
type CancelRes struct {
}

type CancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_at_period_end" tags:"Subscription" method:"post" summary:"Cancel Subscription At Period End" dc:"Cancel subscription at period end, the subscription will not turn to 'cancelled' at once but will cancelled at period end time, no invoice will generate, the flag 'cancelAtPeriodEnd' of subscription will be enabled"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}
type CancelAtPeriodEndRes struct {
}

type CancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/cancel_last_cancel_at_period_end" tags:"Subscription" method:"post" summary:"Cancel Last Cancel Subscription At Period End" dc:"This action should be request before subscription's period end, If subscription's flag 'cancelAtPeriodEnd' is enabled, this action will resume it to disable, and subscription will continue cycle recurring seems no cancelAtPeriod action be setting"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId, id of subscription, either SubscriptionId or UserId needed, The only one active subscription of userId will effect"`
	UserId         uint64 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will effect if userId provide instead of subscriptionId"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}
type CancelLastCancelAtPeriodEndRes struct {
}

type ChangeGatewayReq struct {
	g.Meta          `path:"/change_gateway" tags:"Subscription" method:"post" summary:"Change Subscription Gateway" `
	SubscriptionId  string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId" `
}
type ChangeGatewayRes struct {
}

type AddNewTrialStartReq struct {
	g.Meta             `path:"/add_new_trial_start" tags:"Subscription" method:"post" summary:"Append Subscription TrialEnd"`
	SubscriptionId     string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	AppendTrialEndHour int64  `json:"appendTrialEndHour" dc:"add appendTrialEndHour For Free" v:"required"`
}
type AddNewTrialStartRes struct {
}

type CreatePreviewReq struct {
	g.Meta                 `path:"/create_preview" tags:"Subscription" method:"post" summary:"Create Subscription Preview"`
	PlanId                 uint64                 `json:"planId" dc:"PlanId" v:"required"`
	Email                  string                 `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId                 uint64                 `json:"userId" dc:"UserId"`
	ExternalUserId         string                 `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	User                   *bean.NewUser          `json:"user" dc:"User Object"`
	Quantity               int64                  `json:"quantity" dc:"Quantity" `
	GatewayId              *uint64                `json:"gatewayId" dc:"GatewayId" `
	GatewayPaymentType     string                 `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	AddonParams            []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode         string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber              string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage          *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	DiscountCode           string                 `json:"discountCode" dc:"DiscountCode"`
	TrialEnd               int64                  `json:"trialEnd" dc:"trial_end, utc time"` // trial_end, utc time
	ApplyPromoCredit       *bool                  `json:"applyPromoCredit"  dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                 `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type CreatePreviewRes struct {
	Plan                           *bean.Plan                 `json:"plan"`
	TrialEnd                       int64                      `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	Quantity                       int64                      `json:"quantity"`
	Gateway                        *detail.Gateway            `json:"gateway"`
	AddonParams                    []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                         []*bean.PlanAddonDetail    `json:"addons"`
	SubscriptionAmountExcludingTax int64                      `json:"subscriptionAmountExcludingTax"                `
	TaxAmount                      int64                      `json:"taxAmount"                `
	DiscountAmount                 int64                      `json:"discountAmount"`
	TotalAmount                    int64                      `json:"totalAmount"                `
	OriginAmount                   int64                      `json:"originAmount"                `
	Currency                       string                     `json:"currency"              `
	Invoice                        *bean.Invoice              `json:"invoice"`
	UserId                         uint64                     `json:"userId" `
	Email                          string                     `json:"email" `
	VatCountryCode                 string                     `json:"vatCountryCode"              `
	VatCountryName                 string                     `json:"vatCountryName"              `
	TaxPercentage                  int64                      `json:"taxPercentage"              `
	VatNumber                      string                     `json:"vatNumber"              `
	VatNumberValidate              *bean.ValidResult          `json:"vatNumberValidate"   `
	Discount                       *bean.MerchantDiscountCode `json:"discount" `
	VatNumberValidateMessage       string                     `json:"vatNumberValidateMessage" `
	DiscountMessage                string                     `json:"discountMessage" `
	OtherPendingCryptoSubscription *detail.SubscriptionDetail `json:"otherPendingCryptoSubscription" `
	OtherActiveSubscriptionId      string                     `json:"otherActiveSubscriptionId" description:"other active or incomplete subscription id "`
	ApplyPromoCredit               bool                       `json:"applyPromoCredit"  dc:"apply promo credit or not"`
}

type CreateReq struct {
	g.Meta                 `path:"/create_submit" tags:"Subscription" method:"post" summary:"Create Subscription"`
	PlanId                 uint64                      `json:"planId" dc:"PlanId" v:"required"`
	UserId                 uint64                      `json:"userId" dc:"UserId"`
	Email                  string                      `json:"email" dc:"Email, one of ExternalUserId&Email, UserId or User needed"`
	ExternalUserId         string                      `json:"externalUserId" dc:"ExternalUserId, unique, one of ExternalUserId&Email, UserId or User needed"`
	User                   *bean.NewUser               `json:"user" dc:"User Object"`
	Quantity               int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId              *uint64                     `json:"gatewayId" dc:"GatewayId" `
	GatewayPaymentType     string                      `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	AddonParams            []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount     int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"            `
	ConfirmCurrency        string                      `json:"confirmCurrency"  dc:"Currency to verify if provide" `
	ReturnUrl              string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl              string                      `json:"cancelUrl" dc:"CancelUrl"`
	VatCountryCode         string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber              string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage          *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	PaymentMethodId        string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata               map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	DiscountCode           string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount               *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	TrialEnd               int64                       `json:"trialEnd"                    dc:"trial_end, utc time"` // trial_end, utc time
	StartIncomplete        bool                        `json:"startIncomplete"        dc:"StartIncomplete, use now pay later, subscription will generate invoice and start with incomplete status if set"`
	ProductData            *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
	ApplyPromoCredit       bool                        `json:"applyPromoCredit" dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                      `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type CreateRes struct {
	Subscription                   *bean.Subscription         `json:"subscription" dc:"Subscription"`
	User                           *bean.UserAccount          `json:"user" dc:"user"`
	Paid                           bool                       `json:"paid"`
	Link                           string                     `json:"link"`
	Token                          string                     `json:"token" dc:"token"`
	OtherPendingCryptoSubscription *detail.SubscriptionDetail `json:"otherPendingCryptoSubscription" `
}

type UserSubscriptionDetailReq struct {
	g.Meta         `path:"/user_subscription_detail" tags:"Subscription" method:"get,post" summary:"User Subscription Detail"`
	UserId         uint64 `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type UserSubscriptionDetailRes struct {
	User                                *bean.UserAccount                       `json:"user" dc:"user"`
	Subscription                        *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Plan                                *bean.Plan                              `json:"plan" dc:"Plan"`
	Gateway                             *detail.Gateway                         `json:"gateway" dc:"Gateway"`
	Addons                              []*bean.PlanAddonDetail                 `json:"addons" dc:"Plan Addon"`
	LatestInvoice                       *bean.Invoice                           `json:"latestInvoice" dc:"LatestInvoice"`
	UnfinishedSubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"Processing Subscription Pending Update"`
}

type ActiveTemporarilyReq struct {
	g.Meta         `path:"/active_temporarily" tags:"Subscription Update" method:"post" summary:"Subscription Active Temporarily" dc:"Subscription active temporarily, status will transmit from pending to incomplete"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	ExpireTime     int64  `json:"expireTime"  dc:"ExpireTime, the expire utc time if not paid"  v:"required"`
}

type ActiveTemporarilyRes struct {
}
