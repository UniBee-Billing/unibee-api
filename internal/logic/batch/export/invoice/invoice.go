package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/invoice/service"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type TaskInvoiceExport struct {
}

func (t TaskInvoiceExport) TaskName() string {
	return fmt.Sprintf("InvoiceExport")
}

func (t TaskInvoiceExport) Header() interface{} {
	return ExportInvoiceEntity{}
}

func (t TaskInvoiceExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.InvoiceListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZoneStr = value
				timeZone = zone
			}
		}
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = uint64(value)
		}
		if value, ok := payload["firstName"].(string); ok {
			req.FirstName = value
		}
		if value, ok := payload["lastName"].(string); ok {
			req.LastName = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		//if value, ok := payload["deleteInclude"].(bool); ok {
		//	req.DeleteInclude = value
		//}
		if value, ok := payload["type"].(float64); ok {
			req.Type = unibee.Int(int(value))
		}
		if value, ok := payload["sendEmail"].(string); ok {
			req.SendEmail = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["amountStart"].(float64); ok {
			req.AmountStart = unibee.Int64(int64(value))
		}
		if value, ok := payload["amountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeStart"].(float64); ok {
			req.ReportTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeEnd"].(float64); ok {
			req.ReportTimeEnd = int64(value) - timeZone
		}
	}
	req.SkipTotal = true
	result, _ := service.InvoiceList(ctx, req)
	if result != nil && result.Invoices != nil {
		for _, one := range result.Invoices {
			var invoiceGateway = ""
			if one.Gateway != nil {
				invoiceGateway = one.Gateway.GatewayName
			}
			if one.UserAccount == nil {
				one.UserAccount = &bean.UserAccount{}
			}
			if one.UserSnapshot == nil {
				one.UserSnapshot = one.UserAccount
			}
			if one.Subscription == nil {
				one.Subscription = &bean.Subscription{}
			}
			if one.Payment == nil {
				one.Payment = &bean.Payment{}
			}
			invoiceType := "Tax invoice"
			OriginInvoiceId := ""
			if one.Refund != nil {
				invoiceType = "Credit Note"
				if one.Payment != nil {
					OriginInvoiceId = one.Payment.InvoiceId
				}
			}
			userType := "Individual"
			if one.UserAccount.Type == 2 {
				userType = "Business"
			}
			var dueTime int64 = 0
			if one.FinishTime > 0 {
				dueTime = one.FinishTime + one.DayUtilDue*86400
			}
			var billingPeriod = ""
			if one.BizType == consts.BizTypeSubscription {
				if one.Subscription != nil && one.Subscription.PlanId > 0 {
					plan := query.GetPlanById(ctx, one.Subscription.PlanId)
					if plan != nil {
						if plan.IntervalCount <= 1 {
							billingPeriod = plan.IntervalUnit
						} else {
							billingPeriod = fmt.Sprintf("%d x %s", plan.IntervalCount, plan.IntervalUnit)
						}
					}
				}
			} else {
				billingPeriod = "one time purchase"
			}
			countryName := ""
			IsEu := ""
			if vat_gateway.GetDefaultVatGateway(ctx, one.MerchantId) != nil {
				vatCountryRate, _ := vat_gateway.QueryVatCountryRateByMerchant(ctx, one.MerchantId, one.CountryCode)
				if vatCountryRate != nil {
					countryName = vatCountryRate.CountryName
					if vatCountryRate.IsEU {
						IsEu = "EU"
					} else {
						IsEu = "Non-EU"
					}
				}
			}
			var productId int64 = 0
			if one.Subscription.PlanId > 0 {
				onePlan := query.GetPlanById(ctx, one.Subscription.PlanId)
				if onePlan != nil {
					productId = onePlan.ProductId
				}
			}
			var transactionType = "payment"
			var transactionId = one.PaymentId
			var externalTransactionId = one.Payment.GatewayPaymentId
			if one.Refund != nil {
				transactionType = "refund"
				transactionId = one.RefundId
				externalTransactionId = one.Refund.GatewayRefundId
			}
			mainList = append(mainList, &ExportInvoiceEntity{
				InvoiceId:                      one.InvoiceId,
				InvoiceNumber:                  fmt.Sprintf("%s%s", api.GatewayShortNameMapping[invoiceGateway], one.InvoiceId),
				UserId:                         fmt.Sprintf("%v", one.UserId),
				ExternalUserId:                 fmt.Sprintf("%v", one.UserAccount.ExternalUserId),
				FirstName:                      one.UserSnapshot.FirstName,
				LastName:                       one.UserSnapshot.LastName,
				FullName:                       fmt.Sprintf("%s %s", one.UserSnapshot.FirstName, one.UserSnapshot.LastName),
				UserType:                       userType,
				Email:                          one.UserSnapshot.Email,
				City:                           one.UserSnapshot.City,
				Address:                        one.UserSnapshot.Address,
				InvoiceName:                    one.InvoiceName,
				ProductName:                    one.ProductName,
				TaxCode:                        one.CountryCode,
				CountryCode:                    one.CountryCode,
				PostCode:                       one.UserSnapshot.ZipCode,
				VatNumber:                      one.VatNumber,
				CountryName:                    countryName,
				IsEU:                           IsEu,
				InvoiceType:                    invoiceType,
				Gateway:                        invoiceGateway,
				CompanyName:                    one.UserSnapshot.CompanyName,
				MerchantName:                   merchant.CompanyName,
				DiscountCode:                   one.DiscountCode,
				OriginalAmount:                 utility.ConvertCentToDollarStr(one.OriginAmount, one.Currency),
				TotalAmount:                    utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
				DiscountAmount:                 utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency),
				TotalAmountExcludingTax:        utility.ConvertCentToDollarStr(one.TotalAmountExcludingTax, one.Currency),
				Currency:                       one.Currency,
				TaxAmount:                      utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency),
				TaxPercentage:                  utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
				SubscriptionAmount:             utility.ConvertCentToDollarStr(one.SubscriptionAmount, one.Currency),
				SubscriptionAmountExcludingTax: utility.ConvertCentToDollarStr(one.SubscriptionAmountExcludingTax, one.Currency),
				PeriodEnd:                      gtime.NewFromTimeStamp(one.PeriodEnd + timeZone),
				PeriodStart:                    gtime.NewFromTimeStamp(one.PeriodStart + timeZone),
				FinishTime:                     gtime.NewFromTimeStamp(one.FinishTime + timeZone),
				DueDate:                        gtime.NewFromTimeStamp(dueTime + timeZone),
				CreateTime:                     gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				PaidTime:                       gtime.NewFromTimeStamp(one.Payment.PaidTime + timeZone),
				Status:                         consts.InvoiceStatusToEnum(one.Status).Description(),
				PaymentId:                      one.PaymentId,
				RefundId:                       one.RefundId,
				SubscriptionId:                 one.SubscriptionId,
				PlanId:                         fmt.Sprintf("%v", one.Subscription.PlanId),
				ProductId:                      fmt.Sprintf("%v", productId),
				TrialEnd:                       gtime.NewFromTimeStamp(one.TrialEnd + timeZone),
				BillingCycleAnchor:             gtime.NewFromTimeStamp(one.BillingCycleAnchor + timeZone),
				OriginalInvoiceId:              OriginInvoiceId,
				BillingPeriod:                  billingPeriod,
				TransactionType:                transactionType,
				TransactionId:                  transactionId,
				ExternalTransactionId:          externalTransactionId,
				TimeZone:                       timeZoneStr,
			})
		}
	}
	return mainList, nil
}

type ExportInvoiceEntity struct {
	InvoiceId                      string      `json:"InvoiceId"  comment:"The unique id of invoice, pure digital" group:"Invoice"`
	InvoiceNumber                  string      `json:"InvoiceNumber" comment:"The unique number of invoice, format: Gateway+InvoiceId" group:"Invoice"`
	UserId                         string      `json:"UserId"              comment:"The unique id of user" group:"User Information"`
	ExternalUserId                 string      `json:"ExternalUserId"      comment:"The external unique id of user" group:"User Information"`
	FirstName                      string      `json:"FirstName"           comment:"The first name of user" group:"User Information"`
	LastName                       string      `json:"LastName"            comment:"The last name of user" group:"User Information"`
	FullName                       string      `json:"FullName"            comment:"The full name of user, format: FirstName LastName" group:"User Information"`
	UserType                       string      `json:"UserType"            comment:"The type of user, Individual or Business" group:"User Information"`
	Email                          string      `json:"Email"               comment:"The email of user" group:"User Information"`
	InvoiceName                    string      `json:"InvoiceName" comment:"The name of invoice" group:"Invoice"`
	ProductName                    string      `json:"ProductName" comment:"The product name of invoice" group:"Product and Subscription"`
	InvoiceType                    string      `json:"InvoiceType" comment:"The type of invoice, Tax Invoice or Credit note" group:"Invoice"`
	CompanyName                    string      `json:"CompanyName"        comment:"The CompanyName of user" group:"Transaction"`
	Address                        string      `json:"Address" comment:"The address of user" group:"User Information"`
	City                           string      `json:"City" comment:"The city of user" group:"User Information"`
	State                          string      `json:"State" comment:"The state of user" group:"User Information"`
	CountryCode                    string      `json:"CountryCode"                    comment:"The country code of invoice" group:"User Information"`
	PostCode                       string      `json:"PostCode"                    comment:"The post code of invoice" group:"User Information"`
	VatNumber                      string      `json:"VatNumber"                    comment:"The vat number of invoice" group:"User Information"`
	TaxCode                        string      `json:"TaxCode"                    comment:"The tax code of invoice" group:"Transaction"`
	CountryName                    string      `json:"CountryName" comment:"The country name of invoice" group:"User Information"`
	IsEU                           string      `json:"EU/Non-EU" comment:"Is eu country or not" group:"User Information"`
	Gateway                        string      `json:"Gateway"             comment:"The gateway name of invoice, stripe|paypal|changelly|wire_transfer" group:"Transaction"`
	MerchantName                   string      `json:"MerchantName"        comment:"The name of merchant" group:"Transaction"`
	DiscountCode                   string      `json:"DiscountCode" comment:"The code of discount which apply to invoice" group:"Transaction"`
	OriginalAmount                 string      `json:"OriginalAmount"                 comment:"The original amount of invoice" group:"Transaction"`
	TotalAmount                    string      `json:"TotalAmount" comment:"The total amount of invoice" group:"Transaction"`
	DiscountAmount                 string      `json:"DiscountAmount" comment:"The discount amount of invoice" group:"Transaction"`
	TotalAmountExcludingTax        string      `json:"TotalAmountExcludingTax" comment:"The total amount of invoice with tax excluded" group:"Transaction"`
	Currency                       string      `json:"Currency" comment:"The currency of invoice" group:"Transaction"`
	TaxAmount                      string      `json:"TaxAmount" comment:"The tax amount of invoice" group:"Transaction"`
	TaxPercentage                  string      `json:"TaxPercentage"            comment:"The tax percentage of invoice applied" group:"Transaction"`
	SubscriptionAmount             string      `json:"SubscriptionAmount" comment:"The amount of subscription if invoice is generated by subscription" group:"Product and Subscription"`
	SubscriptionAmountExcludingTax string      `json:"SubscriptionAmountExcludingTax" comment:"The amount of subscription which excluded tax amount if invoice is generated by subscription" group:"Product and Subscription"`
	PeriodEnd                      *gtime.Time `json:"PeriodEnd"  layout:"2006-01-02 15:04:05"   comment:"The end time of period, will apply to subscription if invoice paid" group:"Product and Subscription"`
	PeriodStart                    *gtime.Time `json:"PeriodStart"  layout:"2006-01-02 15:04:05"   comment:"The start time of period, will apply to subscription if invoice paid" group:"Product and Subscription"`
	FinishTime                     *gtime.Time `json:"FinishTime"  layout:"2006-01-02 15:04:05"   comment:"The time when invoice finished, invoice will not be editable after finished" group:"Invoice"`
	CreateTime                     *gtime.Time `json:"CreateTime"  layout:"2006-01-02 15:04:05"   comment:"The create time of invoice" group:"Invoice"`
	DueDate                        *gtime.Time `json:"DueDate"  layout:"2006-01-02 15:04:05"   comment:"The due date of invoice, invoice will expired after due date" group:"Invoice"`
	PaidTime                       *gtime.Time `json:"PaidTime"  layout:"2006-01-02 15:04:05"   comment:"The paid time of invoice" group:"Invoice"`
	Status                         string      `json:"InvoiceStatus"                         comment:"The status of invoice，will be Pending｜Processing｜Paid | Failed | Cancelled | Reversed"  group:"Invoice"`
	PaymentId                      string      `json:"PaymentId"                      comment:"paymentId" comment:"The id of payment connected to invoice" group:"Transaction"`
	RefundId                       string      `json:"RefundId"                       comment:"refundId" comment:"The id of refund connected to invoice, invoice will be credit not who contains refundId" group:"Transaction"`
	SubscriptionId                 string      `json:"SubscriptionId"                 comment:"subscription_id" comment:"the id of subscription connected to invoice" group:"Product and Subscription"`
	ProductId                      string      `json:"ProductId"             comment:"" group:"Product and Subscription"`
	PlanId                         string      `json:"PlanId"              comment:"The id of plan connected to invoice" group:"Product and Subscription"`
	TrialEnd                       *gtime.Time `json:"TrialEnd"                       layout:"2006-01-02 15:04:05"   comment:"The time of trial end, will apply to subscription when invoice paid " group:"Product and Subscription"`
	BillingCycleAnchor             *gtime.Time `json:"BillingCycleAnchor"             layout:"2006-01-02 15:04:05"   comment:"The subscription anchor time of billing cycle connected to invoice " group:"Product and Subscription"`
	OriginalInvoiceId              string      `json:"OriginalInvoiceId"                description:"" comment:"The origin invoiceId if invoice type is credit note" group:"Transaction"`
	BillingPeriod                  string      `json:"BillingPeriod"      comment:"The billing period type of invoice, will be day|month|year or one-time purchase" group:"Product and Subscription"`
	TransactionType                string      `json:"TransactionType"    comment:"" group:"Transaction"`
	TransactionId                  string      `json:"TransactionId"       comment:"" group:"Transaction"`
	ExternalTransactionId          string      `json:"ExternalTransactionId"   comment:"" group:"Transaction"`
	TimeZone                       string      `json:"TimeZone"         comment:"" group:"Transaction"`
}
