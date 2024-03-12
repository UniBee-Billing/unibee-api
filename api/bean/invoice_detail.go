package bean

import (
	"github.com/gogf/gf/v2/os/gtime"
)

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
	Lines                          []*InvoiceItemSimplify `json:"lines"                          description:"lines json data"`
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
	Merchant                       *MerchantSimplify      `json:"merchant"                       description:"Merchant"`
	UserAccount                    *UserAccountSimplify   `json:"userAccount"                       description:"UserAccount"`
	Subscription                   *SubscriptionSimplify  `json:"subscription"                       description:"Subscription"`
	Payment                        *PaymentSimplify       `json:"payment"                       description:"Payment"`
	Refund                         *RefundSimplify        `json:"refund"                       description:"Refund"`
}
