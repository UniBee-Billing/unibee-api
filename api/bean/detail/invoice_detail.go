package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type InvoiceDetail struct {
	Id                             uint64                             `json:"id"                             description:""`
	MerchantId                     uint64                             `json:"merchantId"                     description:"MerchantId"`
	UserId                         uint64                             `json:"userId"                         description:"UserId"`
	SubscriptionId                 string                             `json:"subscriptionId"                 description:"SubscriptionId"`
	InvoiceName                    string                             `json:"invoiceName"                    description:"InvoiceName"`
	InvoiceId                      string                             `json:"invoiceId"                      description:"InvoiceId"`
	GatewayInvoiceId               string                             `json:"gatewayInvoiceId"               description:"GatewayInvoiceId"`
	UniqueId                       string                             `json:"uniqueId"                       description:"UniqueId"`
	GmtCreate                      *gtime.Time                        `json:"gmtCreate"                      description:"GmtCreate"`
	OriginAmount                   int64                              `json:"originAmount"                    description:"OriginAmount,Cents"`
	TotalAmount                    int64                              `json:"totalAmount"                    description:"TotalAmount,Cents"`
	DiscountAmount                 int64                              `json:"discountAmount"                 description:"DiscountAmount,Cents"`
	TaxAmount                      int64                              `json:"taxAmount"                      description:"TaxAmount,Cents"`
	SubscriptionAmount             int64                              `json:"subscriptionAmount"             description:"SubscriptionAmount,Cents"`
	Currency                       string                             `json:"currency"                       description:"Currency"`
	Lines                          []*bean.InvoiceItemSimplify        `json:"lines"                          description:"lines json data"`
	GatewayId                      uint64                             `json:"gatewayId"                      description:"Id"`
	Status                         int                                `json:"status"                         description:"Status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"`
	SendStatus                     int                                `json:"sendStatus"                     description:"SendStatus，0-No | 1- YES"`
	SendEmail                      string                             `json:"sendEmail"                      description:"SendEmail"`
	SendPdf                        string                             `json:"sendPdf"                        description:"SendPdf"`
	GmtModify                      *gtime.Time                        `json:"gmtModify"                      description:"GmtModify"`
	IsDeleted                      int                                `json:"isDeleted"                      description:""`
	Link                           string                             `json:"link"                           description:"Link"`
	GatewayStatus                  string                             `json:"gatewayStatus"                  description:"GatewayStatus，Stripe：https://stripe.com/docs/api/invoices/object"`
	GatewayPaymentId               string                             `json:"gatewayPaymentId"               description:"GatewayPaymentId PaymentId"`
	GatewayUserId                  string                             `json:"gatewayUserId"                  description:"GatewayUserId Id"`
	GatewayInvoicePdf              string                             `json:"gatewayInvoicePdf"              description:"GatewayInvoicePdf pdf"`
	TaxPercentage                  int64                              `json:"taxPercentage"                       description:"TaxPercentage，1000 = 10%"`
	SendNote                       string                             `json:"sendNote"                       description:"SendNote"`
	TotalAmountExcludingTax        int64                              `json:"totalAmountExcludingTax"        description:"TotalAmountExcludingTax,Cents"`
	SubscriptionAmountExcludingTax int64                              `json:"subscriptionAmountExcludingTax" description:"SubscriptionAmountExcludingTax,Cents"`
	PeriodStart                    int64                              `json:"periodStart"                    description:"period_start"`
	PeriodEnd                      int64                              `json:"periodEnd"                      description:"period_end"`
	PaymentId                      string                             `json:"paymentId"                      description:"PaymentId"`
	RefundId                       string                             `json:"refundId"                       description:"refundId"`
	Gateway                        *bean.GatewaySimplify              `json:"gateway"                        description:"Gateway"`
	Merchant                       *bean.MerchantSimplify             `json:"merchant"                       description:"Merchant"`
	UserAccount                    *bean.UserAccountSimplify          `json:"userAccount"                    description:"UserAccount"`
	Subscription                   *bean.SubscriptionSimplify         `json:"subscription"                   description:"Subscription"`
	Payment                        *bean.PaymentSimplify              `json:"payment"                        description:"Payment"`
	Refund                         *bean.RefundSimplify               `json:"refund"                         description:"Refund"`
	Discount                       *bean.MerchantDiscountCodeSimplify `json:"discount"                       description:"Discount"`
	CryptoAmount                   int64                              `json:"cryptoAmount"                   description:"crypto_amount, cent"` // crypto_amount, cent
	CryptoCurrency                 string                             `json:"cryptoCurrency"                 description:"crypto_currency"`
	DayUtilDue                     int64                              `json:"dayUtilDue"                     description:"day util due after finish"` // day util due after finish
}

func ConvertInvoiceToDetail(ctx context.Context, invoice *entity.Invoice) *InvoiceDetail {
	var lines []*bean.InvoiceItemSimplify
	err := bean.UnmarshalFromJsonString(invoice.Lines, &lines)
	for _, line := range lines {
		line.Currency = invoice.Currency
		line.TaxPercentage = invoice.TaxPercentage
	}
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	return &InvoiceDetail{
		Id:                             invoice.Id,
		MerchantId:                     invoice.MerchantId,
		SubscriptionId:                 invoice.SubscriptionId,
		InvoiceId:                      invoice.InvoiceId,
		InvoiceName:                    invoice.InvoiceName,
		GmtCreate:                      invoice.GmtCreate,
		OriginAmount:                   invoice.TotalAmount + invoice.DiscountAmount,
		TotalAmount:                    invoice.TotalAmount,
		TaxAmount:                      invoice.TaxAmount,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		Currency:                       invoice.Currency,
		Lines:                          lines,
		GatewayId:                      invoice.GatewayId,
		Status:                         invoice.Status,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      invoice.SendEmail,
		SendPdf:                        link.GetInvoicePdfLink(ctx, invoice.InvoiceId, invoice.SendTerms),
		UserId:                         invoice.UserId,
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           invoice.Link,
		GatewayStatus:                  invoice.GatewayStatus,
		GatewayInvoiceId:               invoice.GatewayInvoiceId,
		GatewayInvoicePdf:              invoice.GatewayInvoicePdf,
		TaxPercentage:                  invoice.TaxPercentage,
		SendNote:                       invoice.SendNote,
		DiscountAmount:                 invoice.DiscountAmount,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		Gateway:                        bean.SimplifyGateway(query.GetGatewayById(ctx, invoice.GatewayId)),
		Merchant:                       bean.SimplifyMerchant(query.GetMerchantById(ctx, invoice.MerchantId)),
		UserAccount:                    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, invoice.UserId)),
		Subscription:                   bean.SimplifySubscription(query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)),
		Payment:                        bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, invoice.PaymentId)),
		Refund:                         bean.SimplifyRefund(query.GetRefundByRefundId(ctx, invoice.RefundId)),
		Discount:                       bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, invoice.MerchantId, invoice.DiscountCode)),
		CryptoCurrency:                 invoice.CryptoCurrency,
		CryptoAmount:                   invoice.CryptoAmount,
		DayUtilDue:                     invoice.DayUtilDue,
	}
}
