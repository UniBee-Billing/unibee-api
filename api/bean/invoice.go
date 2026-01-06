package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Invoice struct {
	Id                             uint64                             `json:"id"                             description:""`
	UserId                         uint64                             `json:"userId"                         description:"UserId"`
	InvoiceId                      string                             `json:"invoiceId"`
	InvoiceName                    string                             `json:"invoiceName"`
	ProductName                    string                             `json:"productName"`
	DiscountCode                   string                             `json:"discountCode"`
	OriginAmount                   int64                              `json:"originAmount"                `
	TotalAmount                    int64                              `json:"totalAmount"`
	DiscountAmount                 int64                              `json:"discountAmount"`
	TotalAmountExcludingTax        int64                              `json:"totalAmountExcludingTax"`
	Currency                       string                             `json:"currency"`
	TaxAmount                      int64                              `json:"taxAmount"`
	TaxPercentage                  int64                              `json:"taxPercentage"                  description:"TaxPercentage，1000 = 10%"`
	SubscriptionAmount             int64                              `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                              `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemSimplify             `json:"lines"`
	PeriodEnd                      int64                              `json:"periodEnd"`
	PeriodStart                    int64                              `json:"periodStart"`
	FinishTime                     int64                              `json:"finishTime"`
	ProrationDate                  int64                              `json:"prorationDate"`
	ProrationScale                 int64                              `json:"prorationScale"`
	SendNote                       string                             `json:"sendNote"                       description:"send_note"`    // send_note
	Link                           string                             `json:"link"                           description:"invoice link"` // invoice link
	PaymentLink                    string                             `json:"paymentLink"                    description:"invoice payment link"`
	Status                         int                                `json:"status"                         description:"status，1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	PaymentId                      string                             `json:"paymentId"                      description:"paymentId"`                                                     // paymentId
	RefundId                       string                             `json:"refundId"                       description:"refundId"`                                                      // refundId
	SubscriptionId                 string                             `json:"subscriptionId"                 description:"subscription_id"`                                               // subscription_id
	BizType                        int                                `json:"bizType"                        description:"biz type from payment 1-onetime payment, 3-subscription"`       // biz type from payment 1-single payment, 3-subscription
	CryptoAmount                   int64                              `json:"cryptoAmount"                   description:"crypto_amount, cent"`                                           // crypto_amount, cent
	CryptoCurrency                 string                             `json:"cryptoCurrency"                 description:"crypto_currency"`
	SendStatus                     int                                `json:"sendStatus"                     description:"email send status，0-No | 1- YES| 2-Unnecessary"` // email send status，0-No | 1- YES| 2-Unnecessary
	DayUtilDue                     int64                              `json:"dayUtilDue"                     description:"day util due after finish"`                      // day util due after finish
	Discount                       *MerchantDiscountCode              `json:"discount" dc:"Discount"`
	TrialEnd                       int64                              `json:"trialEnd"                       description:"trial_end, utc time"`  // trial_end, utc time
	BillingCycleAnchor             int64                              `json:"billingCycleAnchor"             description:"billing_cycle_anchor"` // billing_cycle_anchor
	CreateFrom                     string                             `json:"createFrom"                     description:"create from"`          // create from
	Metadata                       map[string]interface{}             `json:"metadata" dc:"Metadata，Map"`
	CountryCode                    string                             `json:"countryCode"                    description:""`
	VatNumber                      string                             `json:"vatNumber"                      description:""`
	Data                           string                             `json:"data"                      description:""`
	AutoCharge                     bool                               `json:"autoCharge"                      description:""`
	PromoCreditAccount             *CreditAccount                     `json:"promoCreditAccount"                      description:""`
	PromoCreditPayout              *CreditPayout                      `json:"promoCreditPayout"                      description:""`
	PromoCreditDiscountAmount      int64                              `json:"promoCreditDiscountAmount"      description:"promo credit discount amount"`
	CreditAccount                  *CreditAccount                     `json:"creditAccount"                      description:""`
	CreditPayout                   *CreditPayout                      `json:"creditPayout"                      description:""`
	PartialCreditPaidAmount        int64                              `json:"partialCreditPaidAmount"        description:"partial credit paid amount"`
	PromoCreditTransaction         *CreditTransaction                 `json:"promoCreditTransaction"               description:"promo credit transaction"`
	UserMetricChargeForInvoice     *UserMetricChargeInvoiceItemEntity `json:"userMetricChargeForInvoice"`
	PaymentType                    string                             `json:"paymentType"               description:""`
	PaymentMethodId                string                             `json:"PaymentMethodId"               description:""`
	PlanSnapshot                   *InvoicePlanSnapshot               `json:"planSnapshot" description:"Snapshot of the plan and addons at the time of billing. Includes both the current and previous plans when applicable (e.g., upgrade or downgrade)."`
}

type InvoiceItemSimplify struct {
	Currency                   string                       `json:"currency"`
	OriginAmount               int64                        `json:"originAmount"`
	OriginUnitAmountExcludeTax int64                        `json:"originUnitAmountExcludeTax"`
	DiscountAmount             int64                        `json:"discountAmount"`
	Amount                     int64                        `json:"amount"`
	Tax                        int64                        `json:"tax"`
	AmountExcludingTax         int64                        `json:"amountExcludingTax"`
	TaxPercentage              int64                        `json:"taxPercentage"                  description:"Tax Percentage，1000 = 10%"`
	UnitAmountExcludingTax     int64                        `json:"unitAmountExcludingTax"`
	Name                       string                       `json:"name"`
	Description                string                       `json:"description"`
	PdfDescription             string                       `json:"pdfDescription"`
	Proration                  bool                         `json:"proration"`
	Quantity                   int64                        `json:"quantity"`
	PeriodEnd                  int64                        `json:"periodEnd"`
	PeriodStart                int64                        `json:"periodStart"`
	Plan                       *Plan                        `json:"plan"`
	MetricCharge               *UserMetricChargeInvoiceItem `json:"metricCharge"`
}

type InvoiceSnapshotChargeType int

const (
	InvoiceChargeTypeOneTime InvoiceSnapshotChargeType = iota
	InvoiceChargeTypeSubscriptionCreate
	InvoiceChargeTypeSubscriptionUpgrade
	InvoiceChargeTypeSubscriptionDowngrade
	InvoiceChargeTypeSubscriptionRenew
	InvoiceChargeTypeSubscriptionCycle
)

type InvoicePlanSnapshot struct {
	ChargeType     InvoiceSnapshotChargeType `json:"chargeType" dc:"Billing charge type. 0: One-time, 1: New Subscription, 2: Upgrade, 3: Downgrade, 4: Renewal, 5: Billing Cycle Charge."`
	AutoCharge     bool                      `json:"autoCharge" dc:"Billing charge"`
	Plan           *Plan                     `json:"plan" dc:"Current plan snapshot at the time of billing."`
	Addons         []*PlanAddonDetail        `json:"addons" dc:"Addons associated with the current plan."`
	PreviousPlan   *Plan                     `json:"previousPlan" dc:"Previous plan before upgrade or downgrade. Available only when paidType = 2 or 3."`
	PreviousAddons []*PlanAddonDetail        `json:"previousAddons" dc:"Addons from the previous plan, relevant for upgrade or downgrade (paidType = 2 or 3)."`
}

func SimplifyInvoice(one *entity.Invoice) *Invoice {
	if one == nil {
		return nil
	}
	var lines []*InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	if err != nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
		}
	}
	autoCharge := false
	if one.CreateFrom == consts.InvoiceAutoChargeFlag {
		autoCharge = true
	}
	var userMetricChargeForInvoice *UserMetricChargeInvoiceItemEntity
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &userMetricChargeForInvoice)
	}
	return &Invoice{
		Id:                             one.Id,
		UserId:                         one.UserId,
		InvoiceName:                    one.InvoiceName,
		ProductName:                    one.ProductName,
		InvoiceId:                      one.InvoiceId,
		OriginAmount:                   one.TotalAmount + one.DiscountAmount + one.PromoCreditDiscountAmount,
		TotalAmount:                    one.TotalAmount,
		DiscountCode:                   one.DiscountCode,
		DiscountAmount:                 one.DiscountAmount,
		TotalAmountExcludingTax:        one.TotalAmountExcludingTax,
		Currency:                       one.Currency,
		TaxAmount:                      one.TaxAmount,
		SubscriptionAmount:             one.SubscriptionAmount,
		SubscriptionAmountExcludingTax: one.SubscriptionAmountExcludingTax,
		Lines:                          lines,
		PeriodEnd:                      one.PeriodEnd,
		PeriodStart:                    one.PeriodStart,
		FinishTime:                     one.FinishTime,
		SendNote:                       one.SendNote,
		Link:                           link.GetInvoiceLink(one.InvoiceId, one.SendTerms),
		PaymentLink:                    one.PaymentLink,
		Status:                         one.Status,
		PaymentId:                      one.PaymentId,
		RefundId:                       one.RefundId,
		SubscriptionId:                 one.SubscriptionId,
		BizType:                        one.BizType,
		CryptoCurrency:                 one.CryptoCurrency,
		CryptoAmount:                   one.CryptoAmount,
		SendStatus:                     one.SendStatus,
		DayUtilDue:                     one.DayUtilDue,
		TaxPercentage:                  one.TaxPercentage,
		TrialEnd:                       one.TrialEnd,
		BillingCycleAnchor:             one.BillingCycleAnchor,
		CreateFrom:                     one.CreateFrom,
		CountryCode:                    one.CountryCode,
		Metadata:                       metadata,
		VatNumber:                      one.VatNumber,
		Data:                           one.Data,
		AutoCharge:                     autoCharge,
		PromoCreditDiscountAmount:      one.PromoCreditDiscountAmount,
		PartialCreditPaidAmount:        one.PartialCreditPaidAmount,
		UserMetricChargeForInvoice:     userMetricChargeForInvoice,
		PaymentType:                    one.GatewayInvoiceId,
	}
}
