package ro

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "unibee/api/onetime/payment"
	"unibee/internal/consts"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/oversea_pay"
)

type CreatePayContext struct {
	CheckoutMode             bool                    `json:"checkoutMode"`
	OpenApiId                int64                   `json:"openApiId"`
	AppId                    string                  `json:"appId"`
	Desc                     string                  `json:"desc"`
	Pay                      *entity.Payment         `json:"pay"`
	Gateway                  *entity.MerchantGateway `json:"gateway"`
	TerminalIp               string                  `json:"terminalIp"`
	ShopperUserId            string                  `json:"shopperUserId"`
	ShopperEmail             string                  `json:"shopperEmail"`
	ShopperLocale            string                  `json:"shopperLocale"`
	Mobile                   string                  `json:"mobile"`
	MediaData                map[string]string       `json:"mediaInfo"`
	Invoice                  *InvoiceDetailSimplify  `json:"invoice"`
	ShopperName              *v1.OutShopperName      `json:"shopperName"`
	ShopperInteraction       string                  `json:"shopperInteraction"`
	RecurringProcessingModel string                  `json:"recurringProcessingModel"`
	TokenId                  string                  `json:"tokenId"`
	MerchantOrderReference   string                  `json:"merchantOrderReference"`
	DateOfBirth              *gtime.Time             `json:"dateOfBirth"`
	Platform                 string                  `json:"platform"`
	DeviceType               string                  `json:"deviceType"`
	DaysUtilDue              int                     `json:"daysUtilDue"`
	GatewayPaymentMethod     string                  `json:"gatewayPaymentMethod"`
	PayImmediate             bool                    `json:"payImmediate"`
}

type CreatePayInternalResp struct {
	Status                 consts.PayStatusEnum `json:"status"`
	PaymentId              string               `json:"paymentId"`
	GatewayPaymentId       string               `json:"gatewayPaymentId"`
	GatewayPaymentIntentId string               `json:"gatewayPaymentIntentId"`
	GatewayPaymentMethod   string               `json:"gatewayPaymentMethod"`
	Link                   string               `json:"link"`
	Action                 *gjson.Json          `json:"action"`
	Invoice                *entity.Invoice      `json:"invoice"`
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
	MerchantId      string `json:"merchantId"         `
	GatewayCancelId string `json:"gatewayCancelId"            `
	Reference       string `json:"reference"              `
	Status          string `json:"status"`
}

// OutPayRefundRo is the golang structure for table oversea_pay.
type OutPayRefundRo struct {
	MerchantId       string                  `json:"merchantId"         `
	GatewayRefundId  string                  `json:"gatewayRefundId"            `
	GatewayPaymentId string                  `json:"gatewayPaymentId"            `
	Status           consts.RefundStatusEnum `json:"status"`
	Reason           string                  `json:"reason"              `
	RefundAmount     int64                   `json:"refundFee"              `
	Currency         string                  `json:"currency"              `
	RefundTime       *gtime.Time             `json:"refundTime" `
}

type GatewayPaymentListReq struct {
	UserId int64 `json:"userId"         `
}

// GatewayPaymentRo is the golang structure for table oversea_pay.
type GatewayPaymentRo struct {
	MerchantId           uint64      `json:"merchantId"         `
	Status               int         `json:"status"`
	AuthorizeStatus      int         `json:"captureStatus"`
	AuthorizeReason      string      `json:"authorizeReason" `
	Currency             string      `json:"currency"              `
	TotalAmount          int64       `json:"totalAmount"              `
	PaymentAmount        int64       `json:"paymentAmount"              `
	BalanceAmount        int64       `json:"balanceAmount"              `
	RefundAmount         int64       `json:"refundAmount"              `
	BalanceStart         int64       `json:"balanceStart"              `
	BalanceEnd           int64       `json:"balanceEnd"              `
	Reason               string      `json:"reason"              `
	UniqueId             string      `json:"uniqueId"              `
	PayTime              *gtime.Time `json:"payTime" `
	CreateTime           *gtime.Time `json:"createTime" `
	CancelTime           *gtime.Time `json:"cancelTime" `
	CancelReason         string      `json:"cancelReason" `
	PaymentData          string      `json:"paymentData" `
	GatewayId            uint64      `json:"gatewayId"         `
	GatewayUserId        string      `json:"gatewayUserId"         `
	GatewayPaymentId     string      `json:"gatewayPaymentId"              `
	GatewayPaymentMethod string      `json:"gatewayPaymentMethod"              `
}

type GatewayCreateSubscriptionInternalResp struct {
	GatewayUserId         string                                   `json:"gatewayUserId"`
	GatewaySubscriptionId string                                   `json:"gatewaySubscriptionId"`
	Data                  string                                   `json:"data"`
	Link                  string                                   `json:"link"`
	Status                consts.SubscriptionGatewayPlanStatusEnum `json:"status"`
	Paid                  bool                                     `json:"paid"`
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

type GatewayUserAttachPaymentMethodInternalResp struct {
}

type GatewayUserDeAttachPaymentMethodInternalResp struct {
}

type PaymentMethod struct {
	Id   string      `json:"id"`
	Type string      `json:"type"`
	Data *gjson.Json `json:"data"`
}

type GatewayUserPaymentMethodListInternalResp struct {
	PaymentMethods []*PaymentMethod `json:"paymentMethods"`
}

type GatewayUserPaymentMethodCreateAndBindInternalResp struct {
	PaymentMethod *PaymentMethod `json:"paymentMethod"`
}

type GatewayMerchantBalanceQueryInternalResp struct {
	AvailableBalance       []*GatewayBalance `json:"available"`
	ConnectReservedBalance []*GatewayBalance `json:"connectReserved"`
	PendingBalance         []*GatewayBalance `json:"pending"`
}

type GatewayRedirectInternalResp struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ReturnUrl string `json:"returnUrl"`
	QueryPath string `json:"queryPath"`
}

type InvoiceItemDetailRo struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	Tax                    int64  `json:"tax"`
	TaxScale               int64  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Proration              bool   `json:"proration"`
	Quantity               int64  `json:"quantity"`
	PeriodEnd              int64  `json:"periodEnd"`
	PeriodStart            int64  `json:"periodStart"`
}

type InvoiceDetailSimplify struct {
	InvoiceId                      string                 `json:"invoiceId"`
	InvoiceName                    string                 `json:"invoiceName"`
	TotalAmount                    int64                  `json:"totalAmount"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"`
	Currency                       string                 `json:"currency"`
	TaxAmount                      int64                  `json:"taxAmount"`
	TaxScale                       int64                  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	SubscriptionAmount             int64                  `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemDetailRo `json:"lines"`
	PeriodEnd                      int64                  `json:"periodEnd"`
	PeriodStart                    int64                  `json:"periodStart"`
	ProrationDate                  int64                  `json:"prorationDate"`
	ProrationScale                 int64                  `json:"prorationScale"`
}

type InvoiceDetailRo struct {
	Id                             uint64                 `json:"id"                             description:""`
	MerchantId                     uint64                 `json:"merchantId"                     description:"MerchantId"`
	UserId                         int64                  `json:"userId"                         description:"UserId"`
	SubscriptionId                 string                 `json:"subscriptionId"                 description:"SubscriptionId"`
	InvoiceName                    string                 `json:"invoiceName"                    description:"InvoiceName"`
	InvoiceId                      string                 `json:"invoiceId"                      description:"InvoiceId"`
	GatewayInvoiceId               string                 `json:"gatewayInvoiceId"               description:"GatewayInvoiceId"`
	UniqueId                       string                 `json:"uniqueId"                       description:"UniqueId"`
	GmtCreate                      *gtime.Time            `json:"gmtCreate"                      description:"GmtCreate"`
	TotalAmount                    int64                  `json:"totalAmount"                    description:"TotalAmount,Cents"`
	DiscountAmount                 int64                  `json:"discountAmount"                    description:"DiscountAmount,Cents"`
	TaxAmount                      int64                  `json:"taxAmount"                      description:"TaxAmount,Cents"`
	SubscriptionAmount             int64                  `json:"subscriptionAmount"             description:"SubscriptionAmount,Cents"`
	Currency                       string                 `json:"currency"                       description:"Currency"`
	Lines                          []*InvoiceItemDetailRo `json:"lines"                          description:"lines json data"`
	GatewayId                      uint64                 `json:"gatewayId"                      description:"Id"`
	Status                         int                    `json:"status"                         description:"Status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"`
	SendStatus                     int                    `json:"sendStatus"                     description:"SendStatus，0-No | 1- YES"`
	SendEmail                      string                 `json:"sendEmail"                      description:"SendEmail"`
	SendPdf                        string                 `json:"sendPdf"                        description:"SendPdf"`
	Data                           string                 `json:"data"                           description:"Data"`
	GmtModify                      *gtime.Time            `json:"gmtModify"                      description:"GmtModify"`
	IsDeleted                      int                    `json:"isDeleted"                      description:""`
	Link                           string                 `json:"link"                           description:"Link"`
	GatewayStatus                  string                 `json:"gatewayStatus"                  description:"GatewayStatus，Stripe：https://stripe.com/docs/api/invoices/object"`
	GatewayPaymentId               string                 `json:"gatewayPaymentId"               description:"GatewayPaymentId PaymentId"`
	GatewayUserId                  string                 `json:"gatewayUserId"                  description:"GatewayUserId Id"`
	GatewayInvoicePdf              string                 `json:"gatewayInvoicePdf"              description:"GatewayInvoicePdf pdf"`
	TaxScale                       int64                  `json:"taxScale"                  description:"TaxScale，1000 = 10%"`
	SendNote                       string                 `json:"sendNote"                       description:"SendNote"`
	SendTerms                      string                 `json:"sendTerms"                      description:"SendTerms"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"        description:"TotalAmountExcludingTax,Cents"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax" description:"SubscriptionAmountExcludingTax,Cents"`
	PeriodStart                    int64                  `json:"periodStart"                    description:"period_start"`
	PeriodEnd                      int64                  `json:"periodEnd"                      description:"period_end"`
	PaymentId                      string                 `json:"paymentId"                      description:"PaymentId"`
	RefundId                       string                 `json:"refundId"                       description:"refundId"`
	Gateway                        *GatewaySimplify       `json:"gateway"                       description:"Gateway"`
	Merchant                       *entity.Merchant       `json:"merchant"                       description:"Merchant"`
	UserAccount                    *entity.UserAccount    `json:"userAccount"                       description:"UserAccount"`
	Subscription                   *entity.Subscription   `json:"subscription"                       description:"Subscription"`
	Payment                        *entity.Payment        `json:"payment"                       description:"Payment"`
	Refund                         *entity.Refund         `json:"refund"                       description:"Refund"`
}

type PlanDetailRo struct {
	Plan             *PlanSimplify                `json:"plan" dc:"Plan"`
	MetricPlanLimits []*MerchantMetricPlanLimitVo `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	Addons           []*PlanSimplify              `json:"addons" dc:"Addons"`
	AddonIds         []int64                      `json:"addonIds" dc:"AddonIds"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int64  `json:"quantity" dc:"Quantity，Default 1" `
	AddonPlanId uint64 `json:"addonPlanId" dc:"AddonPlanId"`
}

type PlanAddonVo struct {
	Quantity  int64         `json:"quantity" dc:"Quantity" `
	AddonPlan *PlanSimplify `json:"addonPlan" dc:"addonPlan" `
}

type SubscriptionDetailVo struct {
	User                                *UserAccountSimplify                   `json:"user" dc:"user"`
	Subscription                        *SubscriptionSimplify                  `json:"subscription" dc:"Subscription"`
	Plan                                *PlanSimplify                          `json:"plan" dc:"Plan"`
	Gateway                             *GatewaySimplify                       `json:"gateway" dc:"Gateway"`
	AddonParams                         []*SubscriptionPlanAddonParamRo        `json:"addonParams" dc:"AddonParams"`
	Addons                              []*PlanAddonVo                         `json:"addons" dc:"Addon"`
	UnfinishedSubscriptionPendingUpdate *SubscriptionPendingUpdateDetailVo     `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
	UserMerchantMetricStats             []*metric_event.UserMerchantMetricStat `json:"userMerchantMetricStats" dc:"UserMerchantMetricStats"`
}

type SubscriptionTimeLineDetailVo struct {
	MerchantId      uint64         `json:"merchantId"      description:"merchant id"`             // merchant id
	UserId          int64          `json:"userId"          description:"userId"`                  // userId
	SubscriptionId  string         `json:"subscriptionId"  description:"subscription id"`         // subscription id
	PeriodStart     int64          `json:"periodStart"     description:"period_start"`            // period_start
	PeriodEnd       int64          `json:"periodEnd"       description:"period_end"`              // period_end
	PeriodStartTime *gtime.Time    `json:"periodStartTime" description:"period start (datetime)"` // period start (datetime)
	PeriodEndTime   *gtime.Time    `json:"periodEndTime"   description:"period end (datatime)"`   // period end (datatime)
	InvoiceId       string         `json:"invoiceId"       description:"invoice id"`              // invoice id
	UniqueId        string         `json:"uniqueId"        description:"unique id"`               // unique id
	Currency        string         `json:"currency"        description:"currency"`                // currency
	PlanId          uint64         `json:"planId"          description:"PlanId"`                  // PlanId
	Plan            *PlanSimplify  `json:"plan" description:"Plan"`
	Quantity        int64          `json:"quantity"        description:"quantity"` // quantity
	Addons          []*PlanAddonVo `json:"addons" description:"Addon"`
	GatewayId       uint64         `json:"gatewayId"       description:"gateway_id"`      // gateway_id
	CreateTime      int64          `json:"createTime"      description:"create utc time"` // create utc time
}

type SubscriptionPendingUpdateDetailVo struct {
	MerchantId           uint64                  `json:"merchantId"           description:"MerchantId"`
	SubscriptionId       string                  `json:"subscriptionId"       description:"SubscriptionId"`
	UpdateSubscriptionId string                  `json:"updateSubscriptionId" description:"UpdateSubscriptionId"`
	GmtCreate            *gtime.Time             `json:"gmtCreate"            description:"GmtCreate"`
	Amount               int64                   `json:"amount"               description:"Amount, Cent"`
	Status               int                     `json:"status"               description:"Status，0-Init | 1-Create｜2-Finished｜3-Cancelled"`
	UpdateAmount         int64                   `json:"updateAmount"         description:"UpdateAmount, Cents"`
	ProrationAmount      int64                   `json:"prorationAmount"      description:"ProrationAmount,Cents"`
	Currency             string                  `json:"currency"             description:"Currency"`
	UpdateCurrency       string                  `json:"updateCurrency"       description:"UpdateCurrency"`
	PlanId               uint64                  `json:"planId"               description:"PlanId"`
	UpdatePlanId         uint64                  `json:"updatePlanId"         description:"UpdatePlanId"`
	Quantity             int64                   `json:"quantity"             description:"quantity"`
	UpdateQuantity       int64                   `json:"updateQuantity"       description:"UpdateQuantity"`
	AddonData            string                  `json:"addonData"            description:"plan addon json data"`
	UpdateAddonData      string                  `json:"updateAddonData"     description:"UpdateAddonData"`
	GatewayId            uint64                  `json:"gatewayId"            description:"Id"`
	UserId               int64                   `json:"userId"               description:"UserId"`
	GmtModify            *gtime.Time             `json:"gmtModify"            description:"GmtModify"`
	Paid                 int                     `json:"paid"                 description:"Paid"`
	Link                 string                  `json:"link"                 description:"Link"`
	MerchantMember       *MerchantMemberSimplify `json:"merchantMember"       description:"Merchant Member"`
	EffectImmediate      int                     `json:"effectImmediate"      description:"EffectImmediate"`
	EffectTime           int64                   `json:"effectTime"           description:"effect_immediate=0, EffectTime unit_time"`
	Note                 string                  `json:"note"            description:"Update Note"`
	Plan                 *PlanSimplify           `json:"plan" dc:"Plan"`
	Addons               []*PlanAddonVo          `json:"addons" dc:"Addons"`
	UpdatePlan           *PlanSimplify           `json:"updatePlan" dc:"UpdatePlan"`
	UpdateAddons         []*PlanAddonVo          `json:"updateAddons" dc:"UpdateAddons"`
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

type BulkMetricLimitPlanBindingParam struct {
	MetricId    int64  `json:"metricId" dc:"MetricId" v:"required"`
	MetricLimit uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type MerchantMetricVo struct {
	Id                  uint64 `json:"id"            description:"id"`                                                                                // id
	MerchantId          uint64 `json:"merchantId"          description:"merchantId"`                                                                  // merchantId
	Code                string `json:"code"                description:"code"`                                                                        // code
	MetricName          string `json:"metricName"          description:"metric name"`                                                                 // metric name
	MetricDescription   string `json:"metricDescription"   description:"metric description"`                                                          // metric description
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int    `json:"aggregationType"     description:"1-count，2-count unique, 3-latest, 4-max, 5-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	UpdateTime          int64  `json:"gmtModify"     description:"update time"`     // update time
	CreateTime          int64  `json:"createTime"    description:"create utc time"` // create utc time
}

type MerchantMetricPlanLimitVo struct {
	Id          uint64            `json:"id"            description:"id"`                     // id
	MerchantId  uint64            `json:"merchantId"          description:"merchantId"`       // merchantId
	MetricId    int64             `json:"metricId"    description:"metricId"`                 // metricId
	Metric      *MerchantMetricVo `json:"merchantMetricVo"    description:"MerchantMetricVo"` // metricId
	PlanId      uint64            `json:"planId"      description:"plan_id"`                  // plan_id
	MetricLimit uint64            `json:"metricLimit" description:"plan metric limit"`        // plan metric limit
	UpdateTime  int64             `json:"gmtModify"     description:"update time"`            // update time
	CreateTime  int64             `json:"createTime"    description:"create utc time"`        // create utc time
}
