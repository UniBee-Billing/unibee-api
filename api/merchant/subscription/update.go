package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PendingUpdateListReq struct {
	g.Meta         `path:"/pending_update_list" tags:"Subscription Update" method:"get,post" summary:"Get Subscription Pending Update List"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type PendingUpdateListRes struct {
	SubscriptionPendingUpdateDetails []*detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"Subscription Pending Update Details"`
	Total                            int                                       `json:"total" dc:"Total"`
}

type PendingUpdateDetailReq struct {
	g.Meta                      `path:"/pending_update_detail" tags:"Subscription Update" method:"get" summary:"Subscription Pending Update Detail"`
	SubscriptionPendingUpdateId string `json:"subscriptionPendingUpdateId" dc:"SubscriptionPendingUpdateId" v:"required"`
}

type PendingUpdateDetailRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"SubscriptionPendingUpdate" dc:"subscription pending update"`
}

type RenewReq struct {
	g.Meta                 `path:"/renew" tags:"Subscription Update" method:"post" summary:"Renew Subscription" dc:"renew an exist subscription "`
	SubscriptionId         string                      `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached, either SubscriptionId or UserId needed, The only one active subscription or latest subscription will renew if userId provide instead of subscriptionId"`
	UserId                 uint64                      `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription or latest cancel|expire subscription will renew if userId provide instead of subscriptionId"`
	ProductId              int64                       `json:"productId" dc:"Id of product" dc:"default product will use if not specified"`
	GatewayId              *uint64                     `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	GatewayPaymentType     string                      `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	TaxPercentage          *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	DiscountCode           string                      `json:"discountCode" dc:"DiscountCode, override subscription discount"`
	Discount               *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment          bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl              string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl              string                      `json:"cancelUrl" dc:"CancelUrl"`
	ProductData            *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
	Metadata               map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	ApplyPromoCredit       *bool                       `json:"applyPromoCredit" dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                      `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type RenewRes struct {
	Subscription *bean.Subscription `json:"subscription" dc:"Subscription"`
	Paid         bool               `json:"paid"`
	Link         string             `json:"link"`
}

type UpdatePreviewReq struct {
	g.Meta                 `path:"/update_preview" tags:"Subscription Update" method:"post" summary:"Update Subscription Preview"`
	SubscriptionId         string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId              uint64                 `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity               int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId              uint64                 `json:"gatewayId" dc:"Id" `
	EffectImmediate        int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams            []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode           string                 `json:"discountCode"        dc:"DiscountCode"`
	ApplyPromoCredit       *bool                  `json:"applyPromoCredit" dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                 `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
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
	DiscountMessage   string                     `json:"discountMessage" `
	ApplyPromoCredit  bool                       `json:"applyPromoCredit" dc:"apply promo credit or not"`
}

type UpdateReq struct {
	g.Meta                 `path:"/update_submit" tags:"Subscription Update" method:"post" summary:"Update Subscription"`
	SubscriptionId         string                      `json:"subscriptionId" dc:"SubscriptionId, either SubscriptionId or UserId needed, The only one active subscription of userId will update"`
	UserId                 uint64                      `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will update if userId provide instead of subscriptionId"`
	NewPlanId              uint64                      `json:"newPlanId" dc:"New PlanId" v:"required"`
	Quantity               int64                       `json:"quantity" dc:"Quantity"  v:"required"`
	GatewayId              *uint64                     `json:"gatewayId" dc:"Id of gateway" `
	GatewayPaymentType     string                      `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	AddonParams            []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	EffectImmediate        int                         `json:"effectImmediate" dc:"Force Effect Immediate，1-Immediate，2-Next Period, this api will check upgrade|downgrade automatically" `
	ConfirmTotalAmount     int64                       `json:"confirmTotalAmount"  dc:"TotalAmount to verify if provide"          `
	ConfirmCurrency        string                      `json:"confirmCurrency" dc:"Currency to verify if provide"   `
	ProrationDate          *int64                      `json:"prorationDate" dc:"The utc time to start Proration, default current time" `
	TaxPercentage          *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	Metadata               map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	DiscountCode           string                      `json:"discountCode" dc:"DiscountCode"`
	Discount               *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment          bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl              string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl              string                      `json:"cancelUrl" dc:"CancelUrl"`
	ProductData            *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
	ApplyPromoCredit       bool                        `json:"applyPromoCredit" dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                      `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type UpdateRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
	Paid                      bool                                    `json:"paid"`
	Link                      string                                  `json:"link"`
	Note                      string                                  `json:"note" dc:"note"`
}
