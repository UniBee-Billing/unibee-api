package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type InvoiceDetail struct {
	Id                             uint64                                  `json:"id"                             description:""`
	MerchantId                     uint64                                  `json:"merchantId"                     description:"MerchantId"`
	UserId                         uint64                                  `json:"userId"                         description:"UserId"`
	SubscriptionId                 string                                  `json:"subscriptionId"                 description:"SubscriptionId"`
	InvoiceName                    string                                  `json:"invoiceName"                    description:"InvoiceName"`
	ProductName                    string                                  `json:"productName"`
	InvoiceId                      string                                  `json:"invoiceId"                      description:"InvoiceId"`
	GatewayPaymentType             string                                  `json:"gatewayPaymentType"               description:"GatewayPaymentType"`
	UniqueId                       string                                  `json:"uniqueId"                       description:"UniqueId"`
	GmtCreate                      *gtime.Time                             `json:"gmtCreate"                      description:"GmtCreate"`
	OriginAmount                   int64                                   `json:"originAmount"                   description:"OriginAmount,Cents"`
	TotalAmount                    int64                                   `json:"totalAmount"                    description:"TotalAmount,Cents"`
	DiscountCode                   string                                  `json:"discountCode"`
	DiscountAmount                 int64                                   `json:"discountAmount"                 description:"DiscountAmount,Cents"`
	TaxAmount                      int64                                   `json:"taxAmount"                      description:"TaxAmount,Cents"`
	SubscriptionAmount             int64                                   `json:"subscriptionAmount"             description:"SubscriptionAmount,Cents"`
	Currency                       string                                  `json:"currency"                       description:"Currency"`
	Lines                          []*bean.InvoiceItemSimplify             `json:"lines"                          description:"lines json data"`
	GatewayId                      uint64                                  `json:"gatewayId"                      description:"Id"`
	Status                         int                                     `json:"status"                         description:"Status，1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"`
	SendStatus                     int                                     `json:"sendStatus"                     description:"SendStatus，0-No | 1- YES"`
	SendEmail                      string                                  `json:"sendEmail"                      description:"SendEmail"`
	SendPdf                        string                                  `json:"sendPdf"                        description:"SendPdf"`
	GmtModify                      *gtime.Time                             `json:"gmtModify"                      description:"GmtModify"`
	IsDeleted                      int                                     `json:"isDeleted"                      description:""`
	Link                           string                                  `json:"link"                           description:"Link"`
	GatewayStatus                  string                                  `json:"gatewayStatus"                  description:"GatewayStatus，Stripe：https://stripe.com/docs/api/invoices/object"`
	GatewayPaymentId               string                                  `json:"gatewayPaymentId"               description:"GatewayPaymentId PaymentId"`
	GatewayUserId                  string                                  `json:"gatewayUserId"                  description:"GatewayUserId Id"`
	GatewayInvoicePdf              string                                  `json:"gatewayInvoicePdf"              description:"GatewayInvoicePdf pdf"`
	TaxPercentage                  int64                                   `json:"taxPercentage"                  description:"TaxPercentage，1000 = 10%"`
	SendNote                       string                                  `json:"sendNote"                       description:"SendNote"`
	TotalAmountExcludingTax        int64                                   `json:"totalAmountExcludingTax"        description:"TotalAmountExcludingTax,Cents"`
	SubscriptionAmountExcludingTax int64                                   `json:"subscriptionAmountExcludingTax" description:"SubscriptionAmountExcludingTax,Cents"`
	PeriodStart                    int64                                   `json:"periodStart"                    description:"period_start"`
	PeriodEnd                      int64                                   `json:"periodEnd"                      description:"period_end"`
	PaymentId                      string                                  `json:"paymentId"                      description:"PaymentId"`
	RefundId                       string                                  `json:"refundId"                       description:"refundId"`
	Gateway                        *Gateway                                `json:"gateway"                        description:"Gateway"`
	Merchant                       *bean.Merchant                          `json:"merchant"                       description:"Merchant"`
	UserAccount                    *bean.UserAccount                       `json:"userAccount"                    description:"UserAccount"`
	UserSnapshot                   *bean.UserAccount                       `json:"userSnapshot"                   description:"UserSnapshot"`
	Subscription                   *bean.Subscription                      `json:"subscription"                   description:"Subscription"`
	SubscriptionPendingUpdate      *SubscriptionPendingUpdateDetail        `json:"subscriptionPendingUpdate"     description:"SubscriptionPendingUpdate"`
	Payment                        *bean.Payment                           `json:"payment"                        description:"Payment"`
	Refund                         *bean.Refund                            `json:"refund"                         description:"Refund"`
	Discount                       *bean.MerchantDiscountCode              `json:"discount"                       description:"Discount"`
	CryptoAmount                   int64                                   `json:"cryptoAmount"                   description:"crypto_amount, cent"` // crypto_amount, cent
	CryptoCurrency                 string                                  `json:"cryptoCurrency"                 description:"crypto_currency"`
	DayUtilDue                     int64                                   `json:"dayUtilDue"                     description:"day util due after finish"` // day util due after finish
	BillingCycleAnchor             int64                                   `json:"billingCycleAnchor"             description:"billing_cycle_anchor"`      // billing_cycle_anchor
	CreateFrom                     string                                  `json:"createFrom"                     description:"create from"`               // create from
	Metadata                       map[string]interface{}                  `json:"metadata" dc:"Metadata，Map"`
	CountryCode                    string                                  `json:"countryCode"                    description:""`
	VatNumber                      string                                  `json:"vatNumber"                      description:""`
	FinishTime                     int64                                   `json:"finishTime"`
	CreateTime                     int64                                   `json:"createTime"`
	BizType                        int                                     `json:"bizType"`
	ProrationDate                  int64                                   `json:"prorationDate"`
	TrialEnd                       int64                                   `json:"trialEnd"                       description:"trial_end, utc time"` // trial_end, utc time
	AutoCharge                     bool                                    `json:"autoCharge"                      description:""`
	OriginalPaymentInvoice         *bean.Invoice                           `json:"originalPaymentInvoice"                      description:""`
	PromoCreditDiscountAmount      int64                                   `json:"promoCreditDiscountAmount"      description:"promo credit discount amount"`
	PromoCreditTransaction         *bean.CreditTransaction                 `json:"promoCreditTransaction"               description:"promo credit transaction"`
	PartialCreditPaidAmount        int64                                   `json:"partialCreditPaidAmount"        description:"partial credit paid amount"`
	Message                        string                                  `json:"message"                      description:""`
	UserMetricChargeForInvoice     *bean.UserMetricChargeInvoiceItemEntity `json:"userMetricChargeForInvoice"`
}

func ConvertInvoiceToDetail(ctx context.Context, invoice *entity.Invoice) *InvoiceDetail {
	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
	for _, line := range lines {
		line.Currency = invoice.Currency
		line.TaxPercentage = invoice.TaxPercentage
	}
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	var metadata = make(map[string]interface{})
	if len(invoice.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(invoice.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
		}
	}
	var userSnapShot *entity.UserAccount
	if len(invoice.Data) > 0 {
		err = gjson.Unmarshal([]byte(invoice.Data), &userSnapShot)
		if err != nil {
			fmt.Printf("UserSnapshot Unmarshal Metadata error:%s", err.Error())
		}
	}
	autoCharge := false
	if invoice.CreateFrom == consts.InvoiceAutoChargeFlag {
		autoCharge = true
	}
	payment := bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, invoice.PaymentId))
	refund := bean.SimplifyRefund(query.GetRefundByRefundId(ctx, invoice.RefundId))
	var originalPaymentInvoice *bean.Invoice
	message := ""
	if refund != nil {
		originalPaymentInvoice = bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId))
		if invoice.Status == consts.InvoiceStatusFailed {
			message = refund.RefundCommentExplain
		}
	}
	var userMetricChargeForInvoice *bean.UserMetricChargeInvoiceItemEntity
	if len(invoice.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(invoice.MetricCharge, &userMetricChargeForInvoice)
	}
	return &InvoiceDetail{
		Id:                             invoice.Id,
		MerchantId:                     invoice.MerchantId,
		SubscriptionId:                 invoice.SubscriptionId,
		InvoiceId:                      invoice.InvoiceId,
		InvoiceName:                    invoice.InvoiceName,
		ProductName:                    invoice.ProductName,
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
		SendPdf:                        link.GetInvoicePdfLink(invoice.InvoiceId, invoice.SendTerms),
		UserId:                         invoice.UserId,
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           link.GetInvoiceLink(invoice.InvoiceId, invoice.SendTerms),
		GatewayStatus:                  invoice.GatewayStatus,
		GatewayPaymentType:             invoice.GatewayInvoiceId,
		GatewayInvoicePdf:              invoice.GatewayInvoicePdf,
		TaxPercentage:                  invoice.TaxPercentage,
		SendNote:                       invoice.SendNote,
		DiscountCode:                   invoice.DiscountCode,
		DiscountAmount:                 invoice.DiscountAmount,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		Gateway:                        ConvertGatewayDetail(ctx, query.GetGatewayById(ctx, invoice.GatewayId)),
		Merchant:                       bean.SimplifyMerchant(query.GetMerchantById(ctx, invoice.MerchantId)),
		UserAccount:                    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, invoice.UserId)),
		UserSnapshot:                   bean.SimplifyUserAccount(userSnapShot),
		Subscription:                   bean.SimplifySubscription(ctx, query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)),
		SubscriptionPendingUpdate:      ConvertSubscriptionPendingUpdateDetailByInvoiceId(ctx, invoice.InvoiceId),
		Payment:                        payment,
		Refund:                         refund,
		Discount:                       bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, invoice.MerchantId, invoice.DiscountCode)),
		CryptoCurrency:                 invoice.CryptoCurrency,
		CryptoAmount:                   invoice.CryptoAmount,
		DayUtilDue:                     invoice.DayUtilDue,
		BillingCycleAnchor:             invoice.BillingCycleAnchor,
		CreateFrom:                     invoice.CreateFrom,
		CountryCode:                    invoice.CountryCode,
		VatNumber:                      invoice.VatNumber,
		Metadata:                       metadata,
		FinishTime:                     invoice.FinishTime,
		TrialEnd:                       invoice.TrialEnd,
		CreateTime:                     invoice.CreateTime,
		BizType:                        invoice.BizType,
		PaymentId:                      invoice.PaymentId,
		RefundId:                       invoice.RefundId,
		AutoCharge:                     autoCharge,
		OriginalPaymentInvoice:         originalPaymentInvoice,
		PromoCreditDiscountAmount:      invoice.PromoCreditDiscountAmount,
		PromoCreditTransaction:         bean.SimplifyCreditTransaction(ctx, query.GetPromoCreditTransactionByInvoiceId(ctx, invoice.UserId, invoice.InvoiceId)),
		PartialCreditPaidAmount:        invoice.PartialCreditPaidAmount,
		Message:                        message,
		UserMetricChargeForInvoice:     userMetricChargeForInvoice,
	}
}
