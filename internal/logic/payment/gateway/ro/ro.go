package ro

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/vat_gateway"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type CreatePayContext struct {
	OpenApiId                int64                     `json:"openApiId"`
	AppId                    string                    `json:"appId"`
	Desc                     string                    `json:"desc"`
	Pay                      *entity.OverseaPay        `json:"pay"`
	PayChannel               *entity.OverseaPayChannel `json:"payChannel"`
	PaymentBrandAddition     *gjson.Json               `json:"paymentBrandAddition"`
	TerminalIp               string                    `json:"terminalIp"`
	UserId                   string                    `json:"userId"`
	ShopperEmail             string                    `json:"shopperEmail"`
	ShopperLocale            string                    `json:"shopperLocale"`
	Mobile                   string                    `json:"mobile"`
	MediaInfo                *gjson.Json               `json:"mediaInfo"`
	Items                    []*v1.OutLineItem         `json:"items"`
	BillingDetails           *v1.OutPayAddress         `json:"billingDetails"`
	ShippingDetails          *v1.OutPayAddress         `json:"shippingDetails"`
	ShopperName              *v1.OutShopperName        `json:"shopperName"`
	ShopperInteraction       string                    `json:"shopperInteraction"`
	RecurringProcessingModel string                    `json:"recurringProcessingModel"`
	StorePaymentMethod       bool                      `json:"storePaymentMethod"`
	TokenId                  string                    `json:"tokenId"`
	DeviceFingerprint        string                    `json:"deviceFingerprint"`
	MerchantOrderReference   string                    `json:"merchantOrderReference"`
	DateOfBirth              *gtime.Time               `json:"dateOfBirth"`
	Platform                 string                    `json:"platform"`
	DeviceType               string                    `json:"deviceType"`
}

type CreatePayInternalResp struct {
	AlipayOrderNo  string      `json:"alipayOrderNo"`
	PayOrderNo     string      `json:"payOrderNo"`
	AlreadyPaid    bool        `json:"alreadyPaid"`
	OrderString    string      `json:"orderString"`
	Message        string      `json:"message"`
	TppOrderNo     string      `json:"tppOrderNo"`
	TppPayId       string      `json:"tppPayId"`
	PayChannel     int64       `json:"payChannel"`
	PayChannelType string      `json:"payChannelType"`
	Action         *gjson.Json `json:"action"`
	AdditionalData *gjson.Json `json:"additionalData"`
}

// OutPayCaptureRo is the golang structure for table oversea_pay.
type OutPayCaptureRo struct {
	MerchantId   string          `json:"merchantId"         `      // 商户ID
	PspReference string          `json:"pspReference"            ` // 业务类型。1-订单
	Reference    string          `json:"reference"              `  // 业务id-即商户订单号
	Amount       *v1.PayAmountVo `json:"amount"`
	Status       string          `json:"status"`
}

// OutPayCancelRo is the golang structure for table oversea_pay.
type OutPayCancelRo struct {
	MerchantId   string `json:"merchantId"         `      // 商户ID
	PspReference string `json:"pspReference"            ` // 业务类型。1-订单
	Reference    string `json:"reference"              `  // 业务id-即商户订单号
	Status       string `json:"status"`
}

// OutPayRefundRo is the golang structure for table oversea_pay.
type OutPayRefundRo struct {
	MerchantId      string      `json:"merchantId"         `          // 商户ID
	ChannelRefundNo string      `json:"channelRefundNo"            `  // 业务类型。1-订单
	ChargeRefundNo  string      `json:"chargeRefundNo"              ` // 业务id-即商户订单号
	RefundStatus    int         `json:"refundStatus"`
	Reason          string      `json:"reason"              `    // 业务id-即商户订单号
	RefundFee       int64       `json:"refundFee"              ` // 业务id-即商户订单号
	RefundTime      *gtime.Time `json:"refundTime" `             // 创建时间
}

// OutPayRo is the golang structure for table oversea_pay.
type OutPayRo struct {
	MerchantId      string      `json:"merchantId"         `        // 商户ID
	MerchantOrderNo string      `json:"merchantOrderNo"         `   // 商户ID
	ChannelTradeNo  string      `json:"ChannelTradeNo"            ` // 业务类型。1-订单
	ChannelPayId    string      `json:"channelPayId"              ` // 业务id-即商户订单号
	PayStatus       int         `json:"payStatus"`
	Reason          string      `json:"reason"              ` // 业务id-即商户订单号
	PayFee          int64       `json:"PayFee"              ` // 业务id-即商户订单号
	PayTime         *gtime.Time `json:"PayTime" `             // 创建时间
}

type OutChannelRo struct {
	ChannelId   uint64 `json:"channelId"`
	ChannelName string `json:"channelName"`
}

type ChannelCreateProductInternalResp struct {
	ChannelProductId     string `json:"channelProductId"`
	ChannelProductStatus string `json:"channelProductStatus"`
}

type ChannelCreatePlanInternalResp struct {
	ChannelPlanId     string                                   `json:"channelPlanId"`
	ChannelPlanStatus string                                   `json:"channelPlanStatus"`
	Data              string                                   `json:"data"`
	Status            consts.SubscriptionPlanChannelStatusEnum `json:"status"`
}

type ChannelCreateSubscriptionInternalResp struct {
	ChannelUserId             string                                   `json:"channelUserId"`
	ChannelSubscriptionId     string                                   `json:"channelSubscriptionId"`
	ChannelSubscriptionStatus string                                   `json:"channelSubscriptionStatus"`
	Data                      string                                   `json:"data"`
	Link                      string                                   `json:"link"`
	Status                    consts.SubscriptionPlanChannelStatusEnum `json:"status"`
	Paid                      bool                                     `json:"paid"`
}

type ChannelCreateSubscriptionInternalReq struct {
	Plan           *entity.SubscriptionPlan        `json:"plan"`
	AddonPlans     []*SubscriptionPlanAddonRo      `json:"addonPlans"`
	PlanChannel    *entity.SubscriptionPlanChannel `json:"planChannel"`
	Subscription   *entity.Subscription            `json:"subscription"`
	VatCountryRate *vat_gateway.VatCountryRate     `json:"vatCountryRate"`
}

type ChannelUpdateSubscriptionInternalReq struct {
	Plan            *entity.SubscriptionPlan        `json:"plan"`
	Quantity        int64                           `json:"quantity" dc:"数量" `
	OldPlan         *entity.SubscriptionPlan        `json:"oldPlan"`
	AddonPlans      []*SubscriptionPlanAddonRo      `json:"addonPlans"`
	PlanChannel     *entity.SubscriptionPlanChannel `json:"planChannel"`
	OldPlanChannel  *entity.SubscriptionPlanChannel `json:"oldPlanChannel"`
	Subscription    *entity.Subscription            `json:"subscription"`
	ProrationDate   int64                           `json:"prorationDate"`
	EffectImmediate bool                            `json:"EffectImmediate"`
}

type ChannelCancelSubscriptionInternalReq struct {
	Plan         *entity.SubscriptionPlan        `json:"plan"`
	PlanChannel  *entity.SubscriptionPlanChannel `json:"planChannel"`
	Subscription *entity.Subscription            `json:"subscription"`
	InvoiceNow   bool                            `json:"invoiceNow"`
	Prorate      bool                            `json:"prorate"`
}

type ChannelCancelSubscriptionInternalResp struct {
}

type ChannelCancelAtPeriodEndSubscriptionInternalResp struct {
}

type ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp struct {
}

type ChannelUpdateSubscriptionPreviewInternalResp struct {
	Data          string                            `json:"data"`
	TotalAmount   int64                             `json:"totalAmount"`
	Currency      string                            `json:"currency"`
	Invoice       *ChannelDetailInvoiceInternalResp `json:"invoice"`
	ProrationDate int64                             `json:"prorationDate"`
}

type ChannelUpdateSubscriptionInternalResp struct {
	ChannelSubscriptionId     string                                   `json:"channelSubscriptionId"`
	ChannelSubscriptionStatus string                                   `json:"channelSubscriptionStatus"`
	ChannelInvoiceId          string                                   `json:"channelInvoiceId"`
	Data                      string                                   `json:"data"`
	LatestInvoiceLink         string                                   `json:"latestInvoiceLink"`
	Status                    consts.SubscriptionPlanChannelStatusEnum `json:"status"`
	Paid                      bool                                     `json:"paid"`
}

type ChannelDetailSubscriptionInternalResp struct {
	Status                 consts.SubscriptionStatusEnum `json:"status"`
	ChannelSubscriptionId  string                        `json:"channelSubscriptionId"`
	ChannelStatus          string                        `json:"channelStatus"                  ` // 货币
	Data                   string                        `json:"data"`
	ChannelLatestInvoiceId string                        `json:"channelLatestInvoiceId"`
	CancelAtPeriodEnd      bool                          `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd       int64                         `json:"currentPeriodEnd"`
	CurrentPeriodStart     int64                         `json:"currentPeriodStart"`
	TrailEnd               int64                         `json:"trailEnd"`
}

type ChannelBalance struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type ChannelUserBalanceQueryInternalResp struct {
	Balance              *ChannelBalance   `json:"balance"`
	CashBalance          []*ChannelBalance `json:"cashBalance"`
	InvoiceCreditBalance []*ChannelBalance `json:"invoiceCreditBalance"`
	Email                string            `json:"email"`
	Description          string            `json:"description"`
}

type ChannelMerchantBalanceQueryInternalResp struct {
	AvailableBalance       []*ChannelBalance `json:"available"`
	ConnectReservedBalance []*ChannelBalance `json:"connectReserved"`
	PendingBalance         []*ChannelBalance `json:"pending"`
}

type ChannelWebhookSubscriptionInternalResp struct {
}

type ChannelRedirectInternalResp struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ReturnUrl string `json:"returnUrl"`
	QueryPath string `json:"queryPath"`
}

type ChannelCreateInvoiceInternalReq struct {
	Invoice      *entity.Invoice   `json:"invoice"`
	InvoiceLines []*NewInvoiceItem `json:"invoiceLines"`
	PayMethod    int               `json:"payMethod"` // 1-自动支付， 2-发送邮件支付
	DaysUtilDue  int               `json:"daysUtilDue"`
}

type NewInvoiceItem struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type ChannelPayInvoiceInternalReq struct {
	ChannelInvoiceId string `json:"channelInvoiceId"`
}

type ChannelCancelInvoiceInternalReq struct {
	ChannelInvoiceId string `json:"channelInvoiceId"`
}

type ChannelDetailInvoiceInternalResp struct {
	ChannelSubscriptionId          string                      `json:"channelSubscriptionId"           ` // 货币
	TotalAmount                    int64                       `json:"totalAmount"        `              // 金额,单位：分
	TotalAmountExcludingTax        int64                       `json:"totalAmountExcludingTax"        `  // 金额,单位：分
	TaxAmount                      int64                       `json:"taxAmount"          `              // Tax金额,单位：分
	SubscriptionAmount             int64                       `json:"subscriptionAmount" `              // Sub金额,单位：分
	SubscriptionAmountExcludingTax int64                       `json:"subscriptionAmountExcludingTax" `  // Sub金额,单位：分
	Currency                       string                      `json:"currency"           `              // 货币
	Lines                          []*ChannelDetailInvoiceItem `json:"lines"              `              // lines json data
	ChannelId                      int64                       `json:"channelId"          `              // 支付渠道Id
	Status                         consts.InvoiceStatusEnum    `json:"status"             `              // 订阅单状态，0-Init | 1-Pending ｜2-Processing｜3-paid | 4-failed | 5-cancelled
	ChannelUserId                  string                      `json:"channelUserId"             `       // channelUserId
	Link                           string                      `json:"link"               `              //
	ChannelStatus                  string                      `json:"channelStatus"      `              // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelInvoiceId               string                      `json:"channelInvoiceId"   `              // 关联渠道发票 Id
	ChannelInvoicePdf              string                      `json:"ChannelInvoicePdf"   `             // 关联渠道发票 Pdf
	PeriodEnd                      int64                       `json:"periodEnd"`
	PeriodStart                    int64                       `json:"periodStart"`
}

type ChannelDetailInvoiceRo struct {
	TotalAmount                    int64                       `json:"totalAmount"`
	TotalAmountExcludingTax        int64                       `json:"totalAmountExcludingTax"`
	Currency                       string                      `json:"currency"`
	TaxAmount                      int64                       `json:"taxAmount"`
	SubscriptionAmount             int64                       `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                       `json:"subscriptionAmountExcludingTax"`
	Lines                          []*ChannelDetailInvoiceItem `json:"lines"`
}

type ChannelDetailInvoiceItem struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	Tax                    int64  `json:"tax"`
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Proration              bool   `json:"proration"`
	Quantity               int64  `json:"quantity"`
	PeriodEnd              int64  `json:"periodEnd"`
	PeriodStart            int64  `json:"periodStart"`
}

type InvoiceDetailRo struct {
	Id                             uint64                      `json:"id"                             ` //
	MerchantId                     int64                       `json:"merchantId"                     ` // 商户Id
	SubscriptionId                 string                      `json:"subscriptionId"                 ` // 订阅id（内部编号）
	InvoiceId                      string                      `json:"invoiceId"                      ` // 发票ID（内部编号）
	GmtCreate                      *gtime.Time                 `json:"gmtCreate"                      ` // 创建时间
	TotalAmount                    int64                       `json:"totalAmount"                    ` // 金额,单位：分
	TaxAmount                      int64                       `json:"taxAmount"                      ` // Tax金额,单位：分
	SubscriptionAmount             int64                       `json:"subscriptionAmount"             ` // Sub金额,单位：分
	Currency                       string                      `json:"currency"                       ` // 货币
	Lines                          []*ChannelDetailInvoiceItem `json:"lines"                          ` // lines json data
	ChannelId                      int64                       `json:"channelId"                      ` // 支付渠道Id
	Status                         int                         `json:"status"                         ` // 订阅单状态，0-Init | 1-draft｜2-open｜3-paid | 4-uncollectible | 5-void
	SendStatus                     int                         `json:"sendStatus"                     ` // 订阅单发送状态，0-No | 1- YES
	SendEmail                      string                      `json:"sendEmail"                      ` // send_email
	SendPdf                        string                      `json:"sendPdf"                        ` // send_pdf
	UserId                         int64                       `json:"userId"                         ` // userId
	Data                           string                      `json:"data"                           ` // 渠道额外参数，JSON格式
	GmtModify                      *gtime.Time                 `json:"gmtModify"                      ` // 修改时间
	IsDeleted                      int                         `json:"isDeleted"                      ` //
	Link                           string                      `json:"link"                           ` //
	ChannelStatus                  string                      `json:"channelStatus"                  ` // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelInvoiceId               string                      `json:"channelInvoiceId"               ` // 关联渠道发票 Id
	ChannelInvoicePdf              string                      `json:"channelInvoicePdf"              ` // 关联渠道发票 pdf
	TaxPercentage                  int64                       `json:"taxPercentage"                  ` // Tax税率，万分位，1000 表示 10%
	SendNote                       string                      `json:"sendNote"                       ` // send_note
	SendTerms                      string                      `json:"sendTerms"                      ` // send_terms
	TotalAmountExcludingTax        int64                       `json:"totalAmountExcludingTax"        ` // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax int64                       `json:"subscriptionAmountExcludingTax" ` // Sub金额(不含税）,单位：分
	PeriodStart                    int64                       `json:"periodStart"                    ` // period_start
	PeriodEnd                      int64                       `json:"periodEnd"                      ` // period_end
}

type PlanDetailRo struct {
	Plan     *entity.SubscriptionPlan   `p:"plan" json:"plan" dc:"订阅计划"`
	Channels []*OutChannelRo            `p:"channels" json:"channels" dc:"订阅计划 Channel 开通明细"`
	Addons   []*entity.SubscriptionPlan `p:"addons" json:"addons" dc:"订阅计划 Addons 明细"`
	AddonIds []int64                    `p:"addonIds" json:"addonIds" dc:"订阅计划 Addon Ids"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int64 `p:"quantity" json:"quantity" dc:"数量，默认 1" `
	AddonPlanId int64 `p:"addonPlanId" json:"addonPlanId" dc:"订阅计划Addon ID"`
}

type SubscriptionPlanAddonRo struct {
	Quantity         int64                           `p:"quantity"  json:"quantity" dc:"数量" `
	AddonPlan        *entity.SubscriptionPlan        `p:"addonPlan"  json:"addonPlan" dc:"addonPlan" `
	AddonPlanChannel *entity.SubscriptionPlanChannel `p:"addonPlanChannel"   json:"addonPlanChannel" dc:"addonPlanChannel" `
}

type SubscriptionDetailRo struct {
	User         *entity.UserAccount             `json:"user" dc:"user"`
	Subscription *entity.Subscription            `p:"subscription" json:"subscription" dc:"订阅"`
	Plan         *entity.SubscriptionPlan        `p:"plan" json:"plan" dc:"订阅计划"`
	Channel      *OutChannelRo                   `p:"channel" json:"channel" dc:"订阅渠道"`
	AddonParams  []*SubscriptionPlanAddonParamRo `p:"addonParams" json:"addonParams" dc:"订阅Addon参数"`
	Addons       []*SubscriptionPlanAddonRo      `p:"addons" json:"addons" dc:"订阅Addon"`
}
