package ro

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type CreatePayContext struct {
	CheckoutMode  bool                    `json:"checkoutMode"`
	OpenApiId     int64                   `json:"openApiId"`
	AppId         string                  `json:"appId"`
	Desc          string                  `json:"desc"`
	Pay           *entity.Payment         `json:"pay"`
	Gateway       *entity.MerchantGateway `json:"gateway"`
	TerminalIp    string                  `json:"terminalIp"`
	ShopperUserId string                  `json:"merchantUserId"`
	ShopperEmail  string                  `json:"shopperEmail"`
	ShopperLocale string                  `json:"shopperLocale"`
	Mobile        string                  `json:"mobile"`
	MediaData     map[string]string       `json:"mediaInfo"`
	Invoice       *InvoiceDetailSimplify  `json:"invoice"`
	//BillingDetails           *v1.OutPayAddress         `json:"billingDetails"`
	//ShippingDetails          *v1.OutPayAddress         `json:"shippingDetails"`
	ShopperName              *v1.OutShopperName `json:"shopperName"`
	ShopperInteraction       string             `json:"shopperInteraction"`
	RecurringProcessingModel string             `json:"recurringProcessingModel"`
	TokenId                  string             `json:"tokenId"`
	MerchantOrderReference   string             `json:"merchantOrderReference"`
	DateOfBirth              *gtime.Time        `json:"dateOfBirth"`
	Platform                 string             `json:"platform"`
	DeviceType               string             `json:"deviceType"`
	PayMethod                int                `json:"payMethod"` // 1-自动支付， 2-发送邮件支付
	DaysUtilDue              int                `json:"daysUtilDue"`
	GatewayPaymentMethod     string             `json:"gatewayPaymentMethod"`
	PayImmediate             bool               `json:"payImmediate"`
}

type CreatePayInternalResp struct {
	Status                 consts.PayStatusEnum `json:"status"`
	PaymentId              string               `json:"paymentId"`
	GatewayPaymentId       string               `json:"gatewayPaymentId"`
	GatewayPaymentIntentId string               `json:"gatewayPaymentIntentId"`
	Link                   string               `json:"link"`
	Action                 *gjson.Json          `json:"action"`
}

// OutPayCaptureRo is the golang structure for table oversea_pay.
type OutPayCaptureRo struct {
	MerchantId       string       `json:"merchantId"         `
	GatewayCaptureId string       `json:"gatewayCaptureId"            `
	Reference        string       `json:"reference"              `
	Amount           *v1.AmountVo `json:"amount"`
	Status           string       `json:"status"`
}

// OutPayCancelRo is the golang structure for table oversea_pay.
type OutPayCancelRo struct {
	MerchantId      string `json:"merchantId"         `         // 商户ID
	GatewayCancelId string `json:"gatewayCancelId"            ` // 业务类型。1-订单
	Reference       string `json:"reference"              `     // 业务id-即商户订单号
	Status          string `json:"status"`
}

// OutPayRefundRo is the golang structure for table oversea_pay.
type OutPayRefundRo struct {
	MerchantId       string                  `json:"merchantId"         `          // 商户ID
	GatewayRefundId  string                  `json:"gatewayRefundId"            `  // 渠道退款订单
	GatewayPaymentId string                  `json:"gatewayPaymentId"            ` // 渠道支付订单
	Status           consts.RefundStatusEnum `json:"status"`
	Reason           string                  `json:"reason"              `    // 业务id-即商户订单号
	RefundAmount     int64                   `json:"refundFee"              ` // 业务id-即商户订单号
	Currency         string                  `json:"currency"              `
	RefundTime       *gtime.Time             `json:"refundTime" `
}

type GatewayPaymentListReq struct {
	UserId int64 `json:"userId"         `
}

// GatewayPaymentRo is the golang structure for table oversea_pay.
type GatewayPaymentRo struct {
	MerchantId                  int64                                  `json:"merchantId"         `
	Status                      int                                    `json:"status"`
	AuthorizeStatus             int                                    `json:"captureStatus"`
	AuthorizeReason             string                                 `json:"authorizeReason" `
	Currency                    string                                 `json:"currency"              `
	TotalAmount                 int64                                  `json:"totalAmount"              `
	PaymentAmount               int64                                  `json:"paymentAmount"              `
	BalanceAmount               int64                                  `json:"balanceAmount"              `
	RefundAmount                int64                                  `json:"refundAmount"              `
	BalanceStart                int64                                  `json:"balanceStart"              `
	BalanceEnd                  int64                                  `json:"balanceEnd"              `
	Reason                      string                                 `json:"reason"              `
	UniqueId                    string                                 `json:"uniqueId"              `
	PayTime                     *gtime.Time                            `json:"payTime" `
	CreateTime                  *gtime.Time                            `json:"createTime" `
	CancelTime                  *gtime.Time                            `json:"cancelTime" `
	CancelReason                string                                 `json:"cancelReason" `
	PaymentData                 string                                 `json:"paymentData" `
	GatewayId                   int64                                  `json:"gatewayId"         `
	GatewayUserId               string                                 `json:"gatewayUserId"         `
	GatewayPaymentId            string                                 `json:"gatewayPaymentId"              `
	GatewayPaymentMethod        string                                 `json:"gatewayPaymentMethod"              `
	GatewayInvoiceId            string                                 `json:"gatewayInvoiceId"         `
	GatewaySubscriptionId       string                                 `json:"gatewaySubscriptionId"         `
	GatewaySubscriptionUpdateId string                                 `json:"gatewaySubscriptionUpdateId" `
	GatewayInvoiceDetail        *GatewayDetailInvoiceInternalResp      `json:"gatewayInvoiceDetail"  `
	GatewaySubscriptionDetail   *GatewayDetailSubscriptionInternalResp `json:"gatewaySubscriptionDetail"              `
}

type OutGatewayRo struct {
	GatewayId   uint64 `json:"gatewayId"`
	GatewayName string `json:"gatewayName"`
}

type GatewayCreateProductInternalResp struct {
	GatewayProductId     string `json:"gatewayProductId"`
	GatewayProductStatus string `json:"gatewayProductStatus"`
}

type GatewayCreatePlanInternalResp struct {
	GatewayPlanId     string                                   `json:"gatewayPlanId"`
	GatewayPlanStatus string                                   `json:"gatewayPlanStatus"`
	Data              string                                   `json:"data"`
	Status            consts.SubscriptionGatewayPlanStatusEnum `json:"status"`
}

type GatewayCreateSubscriptionInternalResp struct {
	GatewayUserId             string                                   `json:"gatewayUserId"`
	GatewaySubscriptionId     string                                   `json:"gatewaySubscriptionId"`
	GatewaySubscriptionStatus string                                   `json:"gatewaySubscriptionStatus"`
	Data                      string                                   `json:"data"`
	Link                      string                                   `json:"link"`
	Status                    consts.SubscriptionGatewayPlanStatusEnum `json:"status"`
	Paid                      bool                                     `json:"paid"`
}

type GatewayCreateSubscriptionInternalReq struct {
	Plan           *entity.SubscriptionPlan   `json:"plan"`
	AddonPlans     []*SubscriptionPlanAddonRo `json:"addonPlans"`
	GatewayPlan    *entity.GatewayPlan        `json:"gatewayPlan"`
	Subscription   *entity.Subscription       `json:"subscription"`
	VatCountryRate *VatCountryRate            `json:"vatCountryRate"`
}

type GatewayUpdateSubscriptionInternalReq struct {
	Plan            *entity.SubscriptionPlan   `json:"plan"`
	Quantity        int64                      `json:"quantity" dc:"数量" `
	AddonPlans      []*SubscriptionPlanAddonRo `json:"addonPlans"`
	GatewayPlan     *entity.GatewayPlan        `json:"gatewayPlan"`
	Subscription    *entity.Subscription       `json:"subscription"`
	ProrationDate   int64                      `json:"prorationDate"`
	EffectImmediate bool                       `json:"EffectImmediate"`
}

type GatewayCancelSubscriptionInternalReq struct {
	Plan         *entity.SubscriptionPlan `json:"plan"`
	GatewayPlan  *entity.GatewayPlan      `json:"gatewayPlan"`
	Subscription *entity.Subscription     `json:"subscription"`
	InvoiceNow   bool                     `json:"invoiceNow"`
	Prorate      bool                     `json:"prorate"`
}

type GatewayCancelSubscriptionInternalResp struct {
}

type GatewayCancelAtPeriodEndSubscriptionInternalResp struct {
}

type GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp struct {
}

type GatewayUpdateSubscriptionPreviewInternalResp struct {
	Data              string                            `json:"data"`
	TotalAmount       int64                             `json:"totalAmount"`
	Currency          string                            `json:"currency"`
	Invoice           *GatewayDetailInvoiceInternalResp `json:"invoice"`
	NextPeriodInvoice *GatewayDetailInvoiceInternalResp `json:"nextPeriodInvoice"`
	ProrationDate     int64                             `json:"prorationDate"`
}

type GatewayUpdateSubscriptionInternalResp struct {
	GatewayUpdateId string `json:"gatewayUpdateId" description:"渠道更新单Id"`
	Data            string `json:"data"`
	Link            string `json:"link" description:"需要支付情况下，提供支付链接"`
	Paid            bool   `json:"paid" description:"是否已支付，false-未支付，需要支付，true-已支付或不需要支付"`
}

type GatewayDetailSubscriptionInternalResp struct {
	Status                      consts.SubscriptionStatusEnum `json:"status"`
	GatewaySubscriptionId       string                        `json:"gatewaySubscriptionId"`
	GatewayStatus               string                        `json:"gatewayStatus"                  `
	Data                        string                        `json:"data"`
	GatewayItemData             string                        `json:"gatewayItemData"`
	GatewayLatestInvoiceId      string                        `json:"gatewayLatestInvoiceId"`
	GatewayLatestPaymentId      string                        `json:"gatewayLatestPaymentId"`
	GatewayDefaultPaymentMethod string                        `json:"gatewayDefaultPaymentMethod"`
	CancelAtPeriodEnd           bool                          `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd            int64                         `json:"currentPeriodEnd"`
	CurrentPeriodStart          int64                         `json:"currentPeriodStart"`
	BillingCycleAnchor          int64                         `json:"billingCycleAnchor"`
	TrialEnd                    int64                         `json:"trialEnd"`
}

type GatewayBalance struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type GatewayUserCreateInternalResp struct {
	GatewayUserId string `json:"gatewayUserId"`
}

type GatewayUserDetailQueryInternalResp struct {
	GatewayUserId        string            `json:"gatewayUserId"`
	DefaultPaymentMethod string            `json:"defaultPaymentMethod"`
	Balance              *GatewayBalance   `json:"balance"`
	CashBalance          []*GatewayBalance `json:"cashBalance"`
	InvoiceCreditBalance []*GatewayBalance `json:"invoiceCreditBalance"`
	Email                string            `json:"email"`
	Description          string            `json:"description"`
}

type GatewayUserPaymentMethodListInternalResp struct {
	PaymentMethods []string `json:"paymentMethods"`
}

type GatewayMerchantBalanceQueryInternalResp struct {
	AvailableBalance       []*GatewayBalance `json:"available"`
	ConnectReservedBalance []*GatewayBalance `json:"connectReserved"`
	PendingBalance         []*GatewayBalance `json:"pending"`
}

type GatewayWebhookSubscriptionInternalResp struct {
}

type GatewayRedirectInternalResp struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ReturnUrl string `json:"returnUrl"`
	QueryPath string `json:"queryPath"`
}

type GatewayCreateInvoiceInternalReq struct {
	Invoice      *entity.Invoice        `json:"invoice"`
	InvoiceLines []*InvoiceItemDetailRo `json:"invoiceLines"`
	PayMethod    int                    `json:"payMethod"` // 1-自动支付， 2-发送邮件支付
	DaysUtilDue  int                    `json:"daysUtilDue"`
}

type GatewayPayInvoiceInternalReq struct {
	GatewayInvoiceId string `json:"gatewayInvoiceId"`
}

type GatewayCancelInvoiceInternalReq struct {
	GatewayInvoiceId string `json:"gatewayInvoiceId"`
}

type GatewayDetailInvoiceInternalResp struct {
	GatewayDefaultPaymentMethod    string                   `json:"gatewayDefaultPaymentMethod"`
	GatewaySubscriptionId          string                   `json:"gatewaySubscriptionId"           `
	SubscriptionId                 string                   `json:"subscriptionId"           `
	TotalAmount                    int64                    `json:"totalAmount"        `
	PaymentAmount                  int64                    `json:"paymentAmount"              `
	BalanceAmount                  int64                    `json:"balanceAmount"              `
	BalanceStart                   int64                    `json:"balanceStart"              `
	BalanceEnd                     int64                    `json:"balanceEnd"              `
	TotalAmountExcludingTax        int64                    `json:"totalAmountExcludingTax"        `
	TaxAmount                      int64                    `json:"taxAmount"          `             // Tax金额,单位：分
	SubscriptionAmount             int64                    `json:"subscriptionAmount" `             // Sub金额,单位：分
	SubscriptionAmountExcludingTax int64                    `json:"subscriptionAmountExcludingTax" ` // Sub金额,单位：分
	Currency                       string                   `json:"currency"           `
	Lines                          []*InvoiceItemDetailRo   `json:"lines"              `        // lines json data
	GatewayId                      int64                    `json:"gatewayId"          `        // 支付渠道Id
	Status                         consts.InvoiceStatusEnum `json:"status"             `        // 订阅单状态，0-Init | 1-Pending ｜2-Processing｜3-paid | 4-failed | 5-cancelled
	Reason                         string                   `json:"reason"             `        // reason
	GatewayUserId                  string                   `json:"gatewayUserId"             ` // gatewayUserId
	Link                           string                   `json:"link"               `        //
	GatewayStatus                  string                   `json:"gatewayStatus"      `        // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	GatewayInvoiceId               string                   `json:"gatewayInvoiceId"   `        // 关联渠道发票 Id
	GatewayInvoicePdf              string                   `json:"GatewayInvoicePdf"   `       // 关联渠道发票 Pdf
	PeriodEnd                      int64                    `json:"periodEnd"`
	PeriodStart                    int64                    `json:"periodStart"`
	GatewayPaymentId               string                   `json:"gatewayPaymentId"`
	PaymentTime                    int64                    `json:"paymentTime"        `
	CreateTime                     int64                    `json:"createTime"        `
	CancelTime                     int64                    `json:"cancelTime"        `
}

type InvoiceItemDetailRo struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	Tax                    int64  `json:"tax"`
	TaxScale               int64  `json:"taxScale"                  description:"Tax税率，万分位，1000 表示 10%"` // Tax税率，万分位，1000 表示 10%
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Proration              bool   `json:"proration"`
	Quantity               int64  `json:"quantity"`
	PeriodEnd              int64  `json:"periodEnd"`
	PeriodStart            int64  `json:"periodStart"`
}

type InvoiceDetailSimplify struct {
	InvoiceId                      string                 `json:"invoiceId"`
	TotalAmount                    int64                  `json:"totalAmount"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"`
	Currency                       string                 `json:"currency"`
	TaxAmount                      int64                  `json:"taxAmount"`
	TaxScale                       int64                  `json:"taxScale"                  description:"Tax税率，万分位，1000 表示 10%"` // Tax税率，万分位，1000 表示 10%
	SubscriptionAmount             int64                  `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemDetailRo `json:"lines"`
	PeriodEnd                      int64                  `json:"periodEnd"`
	PeriodStart                    int64                  `json:"periodStart"`
	ProrationDate                  int64                  `json:"prorationDate"`
	ProrationScale                 int64                  `json:"prorationScale"`
}

type InvoiceDetailRo struct {
	Id                             uint64                 `json:"id"                             description:""`                                                       //
	MerchantId                     int64                  `json:"merchantId"                     description:"商户Id"`                                                   // 商户Id
	UserId                         int64                  `json:"userId"                         description:"userId"`                                                 // userId
	SubscriptionId                 string                 `json:"subscriptionId"                 description:"订阅id（内部编号）"`                                             // 订阅id（内部编号）
	InvoiceName                    string                 `json:"invoiceName"                    description:"发票名称"`                                                   // 发票名称
	InvoiceId                      string                 `json:"invoiceId"                      description:"发票ID（内部编号）"`                                             // 发票ID（内部编号）
	GatewayInvoiceId               string                 `json:"gatewayInvoiceId"               description:"关联渠道发票 Id"`                                              // 关联渠道发票 Id
	UniqueId                       string                 `json:"uniqueId"                       description:"唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键"` // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	GmtCreate                      *gtime.Time            `json:"gmtCreate"                      description:"创建时间"`                                                   // 创建时间
	TotalAmount                    int64                  `json:"totalAmount"                    description:"金额,单位：分"`
	DiscountAmount                 int64                  `json:"discountAmount"                    description:"优惠金额,单位：分"`
	TaxAmount                      int64                  `json:"taxAmount"                      description:"Tax金额,单位：分"` // Tax金额,单位：分
	SubscriptionAmount             int64                  `json:"subscriptionAmount"             description:"Sub金额,单位：分"` // Sub金额,单位：分
	Currency                       string                 `json:"currency"                       description:"货币"`
	Lines                          []*InvoiceItemDetailRo `json:"lines"                          description:"lines json data"`                                                       // lines json data
	GatewayId                      int64                  `json:"gatewayId"                      description:"支付渠道Id"`                                                                // 支付渠道Id
	Status                         int                    `json:"status"                         description:"订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     int                    `json:"sendStatus"                     description:"邮件发送状态，0-No | 1- YES"`                                                  // 邮件发送状态，0-No | 1- YES
	SendEmail                      string                 `json:"sendEmail"                      description:"email 发送地址，取自 UserAccount 表 email"`                                     // email 发送地址，取自 UserAccount 表 email
	SendPdf                        string                 `json:"sendPdf"                        description:"pdf 文件地址"`                                                              // pdf 文件地址
	Data                           string                 `json:"data"                           description:"渠道额外参数，JSON格式"`                                                         // 渠道额外参数，JSON格式
	GmtModify                      *gtime.Time            `json:"gmtModify"                      description:"修改时间"`                                                                  // 修改时间
	IsDeleted                      int                    `json:"isDeleted"                      description:""`                                                                      //
	Link                           string                 `json:"link"                           description:"invoice 链接（可用于支付）"`                                                     // invoice 链接（可用于支付）
	GatewayStatus                  string                 `json:"gatewayStatus"                  description:"渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object"`             // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	GatewayPaymentId               string                 `json:"gatewayPaymentId"               description:"关联渠道 PaymentId"`                                                        // 关联渠道 PaymentId
	GatewayUserId                  string                 `json:"gatewayUserId"                  description:"渠道用户 Id"`                                                               // 渠道用户 Id
	GatewayInvoicePdf              string                 `json:"gatewayInvoicePdf"              description:"关联渠道发票 pdf"`                                                            // 关联渠道发票 pdf
	TaxScale                       int64                  `json:"taxScale"                  description:"Tax税率，万分位，1000 表示 10%"`                                                      // Tax税率，万分位，1000 表示 10%
	SendNote                       string                 `json:"sendNote"                       description:"send_note"`                                                             // send_note
	SendTerms                      string                 `json:"sendTerms"                      description:"send_terms"`                                                            // send_terms
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"        description:"金额(不含税）,单位：分"`                                                          // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax" description:"Sub金额(不含税）,单位：分"`                                                       // Sub金额(不含税）,单位：分
	PeriodStart                    int64                  `json:"periodStart"                    description:"period_start"`                                                          // period_start
	PeriodEnd                      int64                  `json:"periodEnd"                      description:"period_end"`                                                            // period_end
	PaymentId                      string                 `json:"paymentId"                      description:"PaymentId"`                                                             // PaymentId
	RefundId                       string                 `json:"refundId"                       description:"refundId"`
}

type PlanDetailRo struct {
	Plan     *entity.SubscriptionPlan   `p:"plan" json:"plan" dc:"订阅计划"`
	Gateways []*OutGatewayRo            `p:"gateways" json:"gateways" dc:"订阅计划 Gateway 开通明细"`
	Addons   []*entity.SubscriptionPlan `p:"addons" json:"addons" dc:"订阅计划 Addons 明细"`
	AddonIds []int64                    `p:"addonIds" json:"addonIds" dc:"订阅计划 Addon Ids"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int64 `p:"quantity" json:"quantity" dc:"数量，Default 1" `
	AddonPlanId int64 `p:"addonPlanId" json:"addonPlanId" dc:"订阅计划Addon ID"`
}

type SubscriptionPlanAddonRo struct {
	Quantity         int64                    `p:"quantity"  json:"quantity" dc:"Quantity" `
	AddonPlan        *entity.SubscriptionPlan `p:"addonPlan"  json:"addonPlan" dc:"addonPlan" `
	AddonGatewayPlan *entity.GatewayPlan      `p:"addonGatewayPlan"   json:"addonGatewayPlan" dc:"AddonGatewayPlan" `
}

type SubscriptionDetailRo struct {
	User                                *entity.UserAccount              `json:"user" dc:"user"`
	Subscription                        *entity.Subscription             `p:"subscription" json:"subscription" dc:"订阅"`
	Plan                                *entity.SubscriptionPlan         `p:"plan" json:"plan" dc:"订阅计划"`
	Gateway                             *OutGatewayRo                    `p:"gateway" json:"gateway" dc:"订阅渠道"`
	AddonParams                         []*SubscriptionPlanAddonParamRo  `p:"addonParams" json:"addonParams" dc:"订阅Addon参数"`
	Addons                              []*SubscriptionPlanAddonRo       `p:"addons" json:"addons" dc:"订阅Addon"`
	UnfinishedSubscriptionPendingUpdate *SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionPendingUpdateDetail struct {
	MerchantId           int64                       `json:"merchantId"           description:"商户Id"`        // 商户Id
	SubscriptionId       string                      `json:"subscriptionId"       description:"订阅id（内部编号）"`  // 订阅id（内部编号）
	UpdateSubscriptionId string                      `json:"updateSubscriptionId" description:"升级单ID（内部编号）"` // 升级单ID（内部编号）
	GmtCreate            *gtime.Time                 `json:"gmtCreate"            description:"创建时间"`        // 创建时间
	Amount               int64                       `json:"amount"               description:"金额,单位：分"`
	Status               int                         `json:"status"               description:"订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled"` // 订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled
	UpdateAmount         int64                       `json:"updateAmount"         description:"升级到金额,单位：分"`                                     // 升级到金额,单位：分
	ProrationAmount      int64                       `json:"prorationAmount"      description:"按比例金额,单位：分"`                                     // 升级到金额,单位：分
	Currency             string                      `json:"currency"             description:"货币"`
	UpdateCurrency       string                      `json:"updateCurrency"       description:"升级到货币"`                                // 升级到货币
	PlanId               int64                       `json:"planId"               description:"计划ID"`                                 // 计划ID
	UpdatePlanId         int64                       `json:"updatePlanId"         description:"升级到计划ID"`                              // 升级到计划ID
	Quantity             int64                       `json:"quantity"             description:"quantity"`                             // quantity
	UpdateQuantity       int64                       `json:"updateQuantity"       description:"升级到quantity"`                          // 升级到quantity
	AddonData            string                      `json:"addonData"            description:"plan addon json data"`                 // plan addon json data
	UpdateAddonData      string                      `json:"updateAddonData"     description:"升级到plan addon json data"`               // 升级到plan addon json data
	GatewayId            int64                       `json:"gatewayId"            description:"支付渠道Id"`                               // 支付渠道Id
	UserId               int64                       `json:"userId"               description:"userId"`                               // userId
	GmtModify            *gtime.Time                 `json:"gmtModify"            description:"修改时间"`                                 // 修改时间
	Paid                 int                         `json:"paid"                 description:"是否已支付，0-否，1-是"`                        // 是否已支付，0-否，1-是
	Link                 string                      `json:"link"                 description:"支付链接"`                                 // 支付链接
	MerchantUser         *entity.MerchantUserAccount `json:"merchantUser"       description:"merchant_user"`                          // merchant_user_id
	EffectImmediate      int                         `json:"effectImmediate"      description:"是否马上生效，0-否，1-是"`                       // 是否马上生效，0-否，1-是
	EffectTime           int64                       `json:"effectTime"           description:"effect_immediate=0, 预计生效时间 unit_time"` // effect_immediate=0, 预计生效时间 unit_time
	Note                 string                      `json:"note"            description:"Update Note"`
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
	Gateway               string `json:"gateway"           `                                          // gateway
	CountryCode           string `json:"countryCode"           `                                      // country_code
	CountryName           string `json:"countryName"           `                                      // country_name
	VatSupport            bool   `json:"vatSupport"          dc:"vat support,true or false"         ` // vat support true or false
	StandardTaxPercentage int64  `json:"standardTaxPercentage"  dc:"Tax税率，万分位，1000 表示 10%"`
}
