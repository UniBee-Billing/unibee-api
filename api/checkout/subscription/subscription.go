package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type CreatePreviewReq struct {
	g.Meta                 `path:"/create_preview" tags:"Checkout" method:"post" summary:"CreateSubscriptionPreview"`
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
	TrialEnd               int64                  `json:"trialEnd" dc:"trial_end, utc time"`
	ApplyPromoCredit       *bool                  `json:"applyPromoCredit" `
	ApplyPromoCreditAmount *int64                 `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type CreatePreviewRes struct {
	Plan                           *bean.Plan                 `json:"plan"`
	TrialEnd                       int64                      `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	Quantity                       int64                      `json:"quantity"`
	Gateway                        *detail.Gateway            `json:"gateway"`
	AddonParams                    []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                         []*bean.PlanAddonDetail    `json:"addons"`
	OriginAmount                   int64                      `json:"originAmount"                `
	TotalAmount                    int64                      `json:"totalAmount"                `
	DiscountAmount                 int64                      `json:"discountAmount"`
	Currency                       string                     `json:"currency"              `
	Invoice                        *bean.Invoice              `json:"invoice"`
	UserId                         uint64                     `json:"userId" `
	Email                          string                     `json:"email" `
	VatCountryCode                 string                     `json:"vatCountryCode"              `
	VatCountryName                 string                     `json:"vatCountryName"              `
	TaxPercentage                  int64                      `json:"taxPercentage"              `
	SubscriptionAmountExcludingTax int64                      `json:"subscriptionAmountExcludingTax"`
	TaxAmount                      int64                      `json:"taxAmount"`
	VatNumber                      string                     `json:"vatNumber"              `
	VatNumberValidate              *bean.ValidResult          `json:"vatNumberValidate"   `
	Discount                       *bean.MerchantDiscountCode `json:"discount" `
	VatNumberValidateMessage       string                     `json:"vatNumberValidateMessage" `
	DiscountMessage                string                     `json:"discountMessage" `
	OtherPendingCryptoSubscription *detail.SubscriptionDetail `json:"otherPendingCryptoSubscription" `
	OtherActiveSubscriptionId      string                     `json:"otherActiveSubscriptionId" description:"other active or incomplete subscription id "`
	ApplyPromoCredit               bool                       `json:"applyPromoCredit" `
	SignIn                         *bean.CheckoutSignIn       `json:"signIn" dc:"Info of sign in"`
}

type CreateReq struct {
	g.Meta                 `path:"/create_submit" tags:"Checkout" method:"post" summary:"CreateSubscription"`
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
