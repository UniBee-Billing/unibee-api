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
	PayMethod                int                `json:"payMethod"`
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
	GatewayUpdateId string `json:"gatewayUpdateId" description:""`
	Data            string `json:"data"`
	Link            string `json:"link" description:""`
	Paid            bool   `json:"paid" description:""`
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
	PayMethod    int                    `json:"payMethod"` // 1-Automatic， 2-Send Payment Link By Invoice
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
	TaxAmount                      int64                    `json:"taxAmount"          `
	SubscriptionAmount             int64                    `json:"subscriptionAmount" `
	SubscriptionAmountExcludingTax int64                    `json:"subscriptionAmountExcludingTax" `
	Currency                       string                   `json:"currency"           `
	Lines                          []*InvoiceItemDetailRo   `json:"lines"              `
	GatewayId                      int64                    `json:"gatewayId"          `
	Status                         consts.InvoiceStatusEnum `json:"status"             `
	Reason                         string                   `json:"reason"             `
	GatewayUserId                  string                   `json:"gatewayUserId"             `
	Link                           string                   `json:"link"               `
	GatewayStatus                  string                   `json:"gatewayStatus"      `
	GatewayInvoiceId               string                   `json:"gatewayInvoiceId"   `
	GatewayInvoicePdf              string                   `json:"GatewayInvoicePdf"   `
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
	MerchantId                     int64                  `json:"merchantId"                     description:"MerchantId"`
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
	GatewayId                      int64                  `json:"gatewayId"                      description:"GatewayId"`
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
}

type PlanDetailRo struct {
	Plan     *entity.SubscriptionPlan   `p:"plan" json:"plan" dc:"Plan"`
	Gateways []*OutGatewayRo            `p:"gateways" json:"gateways" dc:"Gateways"`
	Addons   []*entity.SubscriptionPlan `p:"addons" json:"addons" dc:"Addons"`
	AddonIds []int64                    `p:"addonIds" json:"addonIds" dc:"AddonIds"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int64 `p:"quantity" json:"quantity" dc:"Quantity，Default 1" `
	AddonPlanId int64 `p:"addonPlanId" json:"addonPlanId" dc:"AddonPlanId"`
}

type SubscriptionPlanAddonRo struct {
	Quantity         int64                    `p:"quantity"  json:"quantity" dc:"Quantity" `
	AddonPlan        *entity.SubscriptionPlan `p:"addonPlan"  json:"addonPlan" dc:"addonPlan" `
	AddonGatewayPlan *entity.GatewayPlan      `p:"addonGatewayPlan"   json:"addonGatewayPlan" dc:"AddonGatewayPlan" `
}

type SubscriptionDetailRo struct {
	User                                *entity.UserAccount              `json:"user" dc:"user"`
	Subscription                        *entity.Subscription             `p:"subscription" json:"subscription" dc:"Subscription"`
	Plan                                *entity.SubscriptionPlan         `p:"plan" json:"plan" dc:"Plan"`
	Gateway                             *OutGatewayRo                    `p:"gateway" json:"gateway" dc:"Gateway"`
	AddonParams                         []*SubscriptionPlanAddonParamRo  `p:"addonParams" json:"addonParams" dc:"AddonParams"`
	Addons                              []*SubscriptionPlanAddonRo       `p:"addons" json:"addons" dc:"Addon"`
	UnfinishedSubscriptionPendingUpdate *SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionPendingUpdateDetail struct {
	MerchantId           int64                       `json:"merchantId"           description:"MerchantId"`
	SubscriptionId       string                      `json:"subscriptionId"       description:"SubscriptionId"`
	UpdateSubscriptionId string                      `json:"updateSubscriptionId" description:"UpdateSubscriptionId"`
	GmtCreate            *gtime.Time                 `json:"gmtCreate"            description:"GmtCreate"`
	Amount               int64                       `json:"amount"               description:"Amount, Cent"`
	Status               int                         `json:"status"               description:"Status，0-Init | 1-Create｜2-Finished｜3-Cancelled"`
	UpdateAmount         int64                       `json:"updateAmount"         description:"UpdateAmount, Cents"`
	ProrationAmount      int64                       `json:"prorationAmount"      description:"ProrationAmount,Cents"`
	Currency             string                      `json:"currency"             description:"Currency"`
	UpdateCurrency       string                      `json:"updateCurrency"       description:"UpdateCurrency"`
	PlanId               int64                       `json:"planId"               description:"PlanId"`
	UpdatePlanId         int64                       `json:"updatePlanId"         description:"UpdatePlanId"`
	Quantity             int64                       `json:"quantity"             description:"quantity"`
	UpdateQuantity       int64                       `json:"updateQuantity"       description:"UpdateQuantity"`
	AddonData            string                      `json:"addonData"            description:"plan addon json data"`
	UpdateAddonData      string                      `json:"updateAddonData"     description:"UpdateAddonData"`
	GatewayId            int64                       `json:"gatewayId"            description:"GatewayId"`
	UserId               int64                       `json:"userId"               description:"UserId"`
	GmtModify            *gtime.Time                 `json:"gmtModify"            description:"GmtModify"`
	Paid                 int                         `json:"paid"                 description:"Paid"`
	Link                 string                      `json:"link"                 description:"Link"`
	MerchantUser         *entity.MerchantUserAccount `json:"merchantUser"       description:"merchant_user"`
	EffectImmediate      int                         `json:"effectImmediate"      description:"EffectImmediate"`
	EffectTime           int64                       `json:"effectTime"           description:"effect_immediate=0, EffectTime unit_time"`
	Note                 string                      `json:"note"            description:"Update Note"`
	Plan                 *entity.SubscriptionPlan    `json:"plan" dc:"Plan"`
	Addons               []*SubscriptionPlanAddonRo  `json:"addons" dc:"Addons"`
	UpdatePlan           *entity.SubscriptionPlan    `json:"updatePlan" dc:"UpdatePlan"`
	UpdateAddons         []*SubscriptionPlanAddonRo  `json:"updateAddons" dc:"UpdateAddons"`
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
