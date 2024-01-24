package ro

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type CreatePayContext struct {
	OpenApiId                int64                     `json:"openApiId"`
	AppId                    string                    `json:"appId"`
	Desc                     string                    `json:"desc"`
	Pay                      *entity.Payment           `json:"pay"`
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
	PaymentId      string      `json:"payOrderNo"`
	AlreadyPaid    bool        `json:"alreadyPaid"`
	OrderString    string      `json:"orderString"`
	Message        string      `json:"message"`
	TppOrderNo     string      `json:"tppOrderNo"`
	TppPayId       string      `json:"tppPayId"`
	ChannelId      int64       `json:"payChannel"`
	PayChannelType string      `json:"payChannelType"`
	Action         *gjson.Json `json:"action"`
	AdditionalData *gjson.Json `json:"additionalData"`
}

// OutPayCaptureRo is the golang structure for table oversea_pay.
type OutPayCaptureRo struct {
	MerchantId       string          `json:"merchantId"         `          // 商户ID
	ChannelCaptureId string          `json:"channelCaptureId"            ` // 业务类型。1-订单
	Reference        string          `json:"reference"              `      // 业务id-即商户订单号
	Amount           *v1.PayAmountVo `json:"amount"`
	Status           string          `json:"status"`
}

// OutPayCancelRo is the golang structure for table oversea_pay.
type OutPayCancelRo struct {
	MerchantId      string `json:"merchantId"         `         // 商户ID
	ChannelCancelId string `json:"channelCancelId"            ` // 业务类型。1-订单
	Reference       string `json:"reference"              `     // 业务id-即商户订单号
	Status          string `json:"status"`
}

// OutPayRefundRo is the golang structure for table oversea_pay.
type OutPayRefundRo struct {
	MerchantId       string      `json:"merchantId"         ` // 商户ID
	RefundId         string      `json:"refundId"              `
	ChannelRefundId  string      `json:"channelRefundId"            `  // 渠道退款订单
	ChannelPaymentId string      `json:"channelPaymentId"            ` // 渠道支付订单
	Status           int         `json:"status"`
	Reason           string      `json:"reason"              `    // 业务id-即商户订单号
	RefundFee        int64       `json:"refundFee"              ` // 业务id-即商户订单号
	Currency         string      `json:"currency"              `
	RefundTime       *gtime.Time `json:"refundTime" `
}

type ChannelPaymentListReq struct {
	ChannelUserId string `json:"channelUserId"         `
}

// OutPayRo is the golang structure for table oversea_pay.
type OutPayRo struct {
	MerchantId                string                                 `json:"merchantId"         `
	ChannelId                 int64                                  `json:"channelId"         `
	ChannelPaymentId          string                                 `json:"channelPaymentId"              ` // 业务id-即渠道支付单号
	ChannelUserId             string                                 `json:"channelUserId"         `
	Status                    int                                    `json:"status"`
	CaptureStatus             int                                    `json:"captureStatus"`
	Reason                    string                                 `json:"reason"              `
	PayFee                    int64                                  `json:"PayFee"              `
	ReceiptFee                int64                                  `json:"receiptFee"              `
	Currency                  string                                 `json:"currency"              `
	PayTime                   *gtime.Time                            `json:"payTime" `
	CreateTime                *gtime.Time                            `json:"createTime" `
	CancelTime                *gtime.Time                            `json:"cancelTime" `
	CancelReason              string                                 `json:"cancelReason" `
	TotalRefundFee            int64                                  `json:"totalRefundFee"              `
	ChannelInvoiceId          string                                 `json:"channelInvoiceId"         `
	ChannelSubscriptionId     string                                 `json:"channelSubscriptionId"         `
	ChannelUpdateId           string                                 `json:"channelUpdateId"`
	Subscription              *entity.Subscription                   `json:"subscription"         `
	ChannelUser               *entity.SubscriptionUserChannel        `json:"channelUser"         `
	ChannelInvoiceDetail      *ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"              `
	ChannelSubscriptionDetail *ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"              `
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
	VatCountryRate *VatCountryRate                 `json:"vatCountryRate"`
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
	Data              string                            `json:"data"`
	TotalAmount       int64                             `json:"totalAmount"`
	Currency          string                            `json:"currency"`
	Invoice           *ChannelDetailInvoiceInternalResp `json:"invoice"`
	NextPeriodInvoice *ChannelDetailInvoiceInternalResp `json:"nextPeriodInvoice"`
	ProrationDate     int64                             `json:"prorationDate"`
}

type ChannelUpdateSubscriptionInternalResp struct {
	ChannelUpdateId string `json:"channelUpdateId" description:"渠道更新单Id"`
	Data            string `json:"data"`
	Link            string `json:"link" description:"需要支付情况下，提供支付链接"`
	Paid            bool   `json:"paid" description:"是否已支付，false-未支付，需要支付，true-已支付或不需要支付"`
}

type ChannelDetailSubscriptionInternalResp struct {
	Status                 consts.SubscriptionStatusEnum `json:"status"`
	ChannelSubscriptionId  string                        `json:"channelSubscriptionId"`
	ChannelStatus          string                        `json:"channelStatus"                  `
	Data                   string                        `json:"data"`
	ChannelItemData        string                        `json:"channelItemData"`
	ChannelLatestInvoiceId string                        `json:"channelLatestInvoiceId"`
	ChannelLatestPaymentId string                        `json:"channelLatestPaymentId"`
	CancelAtPeriodEnd      bool                          `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd       int64                         `json:"currentPeriodEnd"`
	CurrentPeriodStart     int64                         `json:"currentPeriodStart"`
	BillingCycleAnchor     int64                         `json:"billingCycleAnchor"`
	TrialEnd               int64                         `json:"trialEnd"`
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
	ChannelSubscriptionId          string                      `json:"channelSubscriptionId"           `
	TotalAmount                    int64                       `json:"totalAmount"        `
	TotalAmountExcludingTax        int64                       `json:"totalAmountExcludingTax"        `
	TaxAmount                      int64                       `json:"taxAmount"          `             // Tax金额,单位：分
	SubscriptionAmount             int64                       `json:"subscriptionAmount" `             // Sub金额,单位：分
	SubscriptionAmountExcludingTax int64                       `json:"subscriptionAmountExcludingTax" ` // Sub金额,单位：分
	Currency                       string                      `json:"currency"           `
	Lines                          []*ChannelDetailInvoiceItem `json:"lines"              `        // lines json data
	ChannelId                      int64                       `json:"channelId"          `        // 支付渠道Id
	Status                         consts.InvoiceStatusEnum    `json:"status"             `        // 订阅单状态，0-Init | 1-Pending ｜2-Processing｜3-paid | 4-failed | 5-cancelled
	ChannelUserId                  string                      `json:"channelUserId"             ` // channelUserId
	Link                           string                      `json:"link"               `        //
	ChannelStatus                  string                      `json:"channelStatus"      `        // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelInvoiceId               string                      `json:"channelInvoiceId"   `        // 关联渠道发票 Id
	ChannelInvoicePdf              string                      `json:"ChannelInvoicePdf"   `       // 关联渠道发票 Pdf
	PeriodEnd                      int64                       `json:"periodEnd"`
	PeriodStart                    int64                       `json:"periodStart"`
	ChannelPaymentId               string                      `json:"channelPaymentId"`
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
	Id                             uint64                      `json:"id"                             description:""`                                                       //
	MerchantId                     int64                       `json:"merchantId"                     description:"商户Id"`                                                   // 商户Id
	UserId                         int64                       `json:"userId"                         description:"userId"`                                                 // userId
	SubscriptionId                 string                      `json:"subscriptionId"                 description:"订阅id（内部编号）"`                                             // 订阅id（内部编号）
	InvoiceName                    string                      `json:"invoiceName"                    description:"发票名称"`                                                   // 发票名称
	InvoiceId                      string                      `json:"invoiceId"                      description:"发票ID（内部编号）"`                                             // 发票ID（内部编号）
	ChannelInvoiceId               string                      `json:"channelInvoiceId"               description:"关联渠道发票 Id"`                                              // 关联渠道发票 Id
	UniqueId                       string                      `json:"uniqueId"                       description:"唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键"` // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	GmtCreate                      *gtime.Time                 `json:"gmtCreate"                      description:"创建时间"`                                                   // 创建时间
	TotalAmount                    int64                       `json:"totalAmount"                    description:"金额,单位：分"`
	TaxAmount                      int64                       `json:"taxAmount"                      description:"Tax金额,单位：分"` // Tax金额,单位：分
	SubscriptionAmount             int64                       `json:"subscriptionAmount"             description:"Sub金额,单位：分"` // Sub金额,单位：分
	Currency                       string                      `json:"currency"                       description:"货币"`
	Lines                          []*ChannelDetailInvoiceItem `json:"lines"                          description:"lines json data"`                                                       // lines json data
	ChannelId                      int64                       `json:"channelId"                      description:"支付渠道Id"`                                                                // 支付渠道Id
	Status                         int                         `json:"status"                         description:"订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     int                         `json:"sendStatus"                     description:"邮件发送状态，0-No | 1- YES"`                                                  // 邮件发送状态，0-No | 1- YES
	SendEmail                      string                      `json:"sendEmail"                      description:"email 发送地址，取自 UserAccount 表 email"`                                     // email 发送地址，取自 UserAccount 表 email
	SendPdf                        string                      `json:"sendPdf"                        description:"pdf 文件地址"`                                                              // pdf 文件地址
	Data                           string                      `json:"data"                           description:"渠道额外参数，JSON格式"`                                                         // 渠道额外参数，JSON格式
	GmtModify                      *gtime.Time                 `json:"gmtModify"                      description:"修改时间"`                                                                  // 修改时间
	IsDeleted                      int                         `json:"isDeleted"                      description:""`                                                                      //
	Link                           string                      `json:"link"                           description:"invoice 链接（可用于支付）"`                                                     // invoice 链接（可用于支付）
	ChannelStatus                  string                      `json:"channelStatus"                  description:"渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object"`             // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelPaymentId               string                      `json:"channelPaymentId"               description:"关联渠道 PaymentId"`                                                        // 关联渠道 PaymentId
	ChannelUserId                  string                      `json:"channelUserId"                  description:"渠道用户 Id"`                                                               // 渠道用户 Id
	ChannelInvoicePdf              string                      `json:"channelInvoicePdf"              description:"关联渠道发票 pdf"`                                                            // 关联渠道发票 pdf
	TaxPercentage                  int64                       `json:"taxPercentage"                  description:"Tax税率，万分位，1000 表示 10%"`                                                 // Tax税率，万分位，1000 表示 10%
	SendNote                       string                      `json:"sendNote"                       description:"send_note"`                                                             // send_note
	SendTerms                      string                      `json:"sendTerms"                      description:"send_terms"`                                                            // send_terms
	TotalAmountExcludingTax        int64                       `json:"totalAmountExcludingTax"        description:"金额(不含税）,单位：分"`                                                          // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax int64                       `json:"subscriptionAmountExcludingTax" description:"Sub金额(不含税）,单位：分"`                                                       // Sub金额(不含税）,单位：分
	PeriodStart                    int64                       `json:"periodStart"                    description:"period_start"`                                                          // period_start
	PeriodEnd                      int64                       `json:"periodEnd"                      description:"period_end"`                                                            // period_end
	PaymentId                      string                      `json:"paymentId"                      description:"PaymentId"`                                                             // PaymentId
	RefundId                       string                      `json:"refundId"                       description:"refundId"`
}

type PlanDetailRo struct {
	Plan     *entity.SubscriptionPlan   `p:"plan" json:"plan" dc:"订阅计划"`
	Channels []*OutChannelRo            `p:"channels" json:"channels" dc:"订阅计划 Channel 开通明细"`
	Addons   []*entity.SubscriptionPlan `p:"addons" json:"addons" dc:"订阅计划 Addons 明细"`
	AddonIds []int64                    `p:"addonIds" json:"addonIds" dc:"订阅计划 Addon Ids"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int64 `p:"quantity" json:"quantity" dc:"数量，Default 1" `
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

type SubscriptionPendingUpdateDetail struct {
	MerchantId           int64                       `json:"merchantId"           description:"商户Id"`        // 商户Id
	SubscriptionId       string                      `json:"subscriptionId"       description:"订阅id（内部编号）"`  // 订阅id（内部编号）
	UpdateSubscriptionId string                      `json:"updateSubscriptionId" description:"升级单ID（内部编号）"` // 升级单ID（内部编号）
	GmtCreate            *gtime.Time                 `json:"gmtCreate"            description:"创建时间"`        // 创建时间
	Amount               int64                       `json:"amount"               description:"金额,单位：分"`
	Status               int                         `json:"status"               description:"订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled"` // 订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled
	UpdateAmount         int64                       `json:"updateAmount"         description:"升级到金额,单位：分"`                                     // 升级到金额,单位：分
	Currency             string                      `json:"currency"             description:"货币"`
	UpdateCurrency       string                      `json:"updateCurrency"       description:"升级到货币"`                                // 升级到货币
	PlanId               int64                       `json:"planId"               description:"计划ID"`                                 // 计划ID
	UpdatePlanId         int64                       `json:"updatePlanId"         description:"升级到计划ID"`                              // 升级到计划ID
	Quantity             int64                       `json:"quantity"             description:"quantity"`                             // quantity
	UpdateQuantity       int64                       `json:"updateQuantity"       description:"升级到quantity"`                          // 升级到quantity
	AddonData            string                      `json:"addonData"            description:"plan addon json data"`                 // plan addon json data
	UpdateAddonData      string                      `json:"updateAddonData"     description:"升级到plan addon json data"`               // 升级到plan addon json data
	ChannelId            int64                       `json:"channelId"            description:"支付渠道Id"`                               // 支付渠道Id
	UserId               int64                       `json:"userId"               description:"userId"`                               // userId
	GmtModify            *gtime.Time                 `json:"gmtModify"            description:"修改时间"`                                 // 修改时间
	Paid                 int                         `json:"paid"                 description:"是否已支付，0-否，1-是"`                        // 是否已支付，0-否，1-是
	Link                 string                      `json:"link"                 description:"支付链接"`                                 // 支付链接
	MerchantUser         *entity.MerchantUserAccount `json:"merchantUser"       description:"merchant_user"`                          // merchant_user_id
	EffectImmediate      int                         `json:"effectImmediate"      description:"是否马上生效，0-否，1-是"`                       // 是否马上生效，0-否，1-是
	EffectTime           int64                       `json:"effectTime"           description:"effect_immediate=0, 预计生效时间 unit_time"` // effect_immediate=0, 预计生效时间 unit_time
	AdminNote            string                      `json:"adminNote"            description:"Admin 修改备注"`
	Plan                 *entity.SubscriptionPlan    `json:"plan" dc:"旧订阅计划"`
	Addons               []*SubscriptionPlanAddonRo  `json:"addons" dc:"旧订阅Addon"`
	UpdatePlan           *entity.SubscriptionPlan    `json:"updatePlan" dc:"更新订阅计划"`
	UpdateAddons         []*SubscriptionPlanAddonRo  `json:"updateAddons" dc:"更新订阅Addon"`
}

type ValidResult struct {
	Valid           bool   `json:"valid"           `
	VatNumber       string `json:"vatNumber"           `
	CountryCode     string `json:"countryCode"           `
	CompanyName     string `json:"companyName"           `
	CompanyAddress  string `json:"companyAddress"           `
	ValidateMessage string `json:"validateMessage"           `
}

type VatCountryRate struct {
	Id                    uint64 `json:"id"  dc:"TaxId"`
	Gateway               string `json:"channel"           `                                          // channel
	CountryCode           string `json:"countryCode"           `                                      // country_code
	CountryName           string `json:"countryName"           `                                      // country_name
	VatSupport            bool   `json:"vatSupport"          dc:"vat support,true or false"         ` // vat support true or false
	StandardTaxPercentage int64  `json:"standardTaxPercentage"  dc:"Tax税率，万分位，1000 表示 10%"`
}
