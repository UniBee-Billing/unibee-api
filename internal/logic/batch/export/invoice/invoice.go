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
	entity "unibee/internal/model/entity/oversea_pay"
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
	if task == nil && task.MerchantId <= 0 {
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
	if payload != nil {
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
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
		if value, ok := payload["reportTimeStart"].(float64); ok {
			req.ReportTimeStart = int64(value)
		}
		if value, ok := payload["reportTimeEnd"].(float64); ok {
			req.ReportTimeEnd = int64(value)
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
				one.UserAccount = &bean.UserAccountSimplify{}
			}
			if one.Subscription == nil {
				one.Subscription = &bean.SubscriptionSimplify{}
			}
			if one.Payment == nil {
				one.Payment = &bean.PaymentSimplify{}
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
						billingPeriod = plan.IntervalUnit
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
			mainList = append(mainList, &ExportInvoiceEntity{
				InvoiceId:                      one.InvoiceId,
				InvoiceNumber:                  fmt.Sprintf("%s%s", api.GatewayShortNameMapping[invoiceGateway], one.InvoiceId),
				UserId:                         fmt.Sprintf("%v", one.UserId),
				ExternalUserId:                 fmt.Sprintf("%v", one.UserAccount.ExternalUserId),
				FirstName:                      one.UserAccount.FirstName,
				LastName:                       one.UserAccount.LastName,
				FullName:                       fmt.Sprintf("%s %s", one.UserAccount.FirstName, one.UserAccount.LastName),
				UserType:                       userType,
				Email:                          one.UserAccount.Email,
				City:                           one.UserAccount.City,
				Address:                        one.UserAccount.Address,
				InvoiceName:                    one.InvoiceName,
				ProductName:                    one.ProductName,
				TaxCode:                        one.CountryCode,
				CountryCode:                    one.CountryCode,
				VatNumber:                      one.VatNumber,
				CountryName:                    countryName,
				IsEU:                           IsEu,
				InvoiceType:                    invoiceType,
				Gateway:                        invoiceGateway,
				MerchantName:                   merchant.Name,
				DiscountCode:                   one.DiscountCode,
				OriginAmount:                   utility.ConvertCentToDollarStr(one.OriginAmount, one.Currency),
				TotalAmount:                    utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
				DiscountAmount:                 utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency),
				TotalAmountExcludingTax:        utility.ConvertCentToDollarStr(one.TotalAmountExcludingTax, one.Currency),
				Currency:                       one.Currency,
				TaxAmount:                      utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency),
				TaxPercentage:                  utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
				SubscriptionAmount:             utility.ConvertCentToDollarStr(one.SubscriptionAmount, one.Currency),
				SubscriptionAmountExcludingTax: utility.ConvertCentToDollarStr(one.SubscriptionAmountExcludingTax, one.Currency),
				PeriodEnd:                      gtime.NewFromTimeStamp(one.PeriodEnd),
				PeriodStart:                    gtime.NewFromTimeStamp(one.PeriodStart),
				FinishTime:                     gtime.NewFromTimeStamp(one.FinishTime),
				DueDate:                        gtime.NewFromTimeStamp(dueTime),
				CreateTime:                     gtime.NewFromTimeStamp(one.CreateTime),
				PaidTime:                       gtime.NewFromTimeStamp(one.Payment.PaidTime),
				Status:                         consts.InvoiceStatusToEnum(one.Status).Description(),
				PaymentId:                      one.PaymentId,
				RefundId:                       one.RefundId,
				SubscriptionId:                 one.SubscriptionId,
				PlanId:                         fmt.Sprintf("%v", one.Subscription.PlanId),
				TrialEnd:                       gtime.NewFromTimeStamp(one.TrialEnd),
				BillingCycleAnchor:             gtime.NewFromTimeStamp(one.BillingCycleAnchor),
				OriginInvoiceId:                OriginInvoiceId,
				BillingPeriod:                  billingPeriod,
			})
		}
	}
	return mainList, nil
}

type ExportInvoiceEntity struct {
	InvoiceId                      string      `json:"InvoiceId"  comment:"The unique id of invoice, pure digital"`
	InvoiceNumber                  string      `json:"InvoiceNumber" comment:"The unique number of invoice, format: Gateway+InvoiceId"`
	UserId                         string      `json:"UserId"              comment:"The unique id of user"`
	ExternalUserId                 string      `json:"ExternalUserId"      comment:"The external unique id of user"`
	FirstName                      string      `json:"FirstName"           comment:"The first name of user"`
	LastName                       string      `json:"LastName"            comment:"The last name of user"`
	FullName                       string      `json:"FullName"            comment:"The full name of user, format: FirstName LastName"`
	UserType                       string      `json:"UserType"            comment:"The type of user, Individual or Business"`
	Email                          string      `json:"Email"               comment:"The email of user"`
	InvoiceName                    string      `json:"InvoiceName" comment:"The name of invoice"`
	ProductName                    string      `json:"ProductName" comment:"The product name of invoice"`
	InvoiceType                    string      `json:"InvoiceType" comment:"The type of invoice, Tax Invoice or Credit Note"`
	Address                        string      `json:"Address" comment:"The address of user"`
	City                           string      `json:"City" comment:"The city of user"`
	State                          string      `json:"State" comment:"The state of user"`
	CountryCode                    string      `json:"CountryCode"                    comment:"The country code of invoice"`
	VatNumber                      string      `json:"VatNumber"                    comment:"The vat number of invoice"`
	TaxCode                        string      `json:"TaxCode"                    comment:"The tax code of invoice"`
	CountryName                    string      `json:"CountryName" comment:"The country name of invoice"`
	IsEU                           string      `json:"EU/Non-EU" comment:"Is eu country or not"`
	Gateway                        string      `json:"Gateway"             comment:"The gateway name of invoice, stripe|paypal|changelly|wire_transfer"`
	MerchantName                   string      `json:"MerchantName"        comment:"The name of merchant"`
	DiscountCode                   string      `json:"DiscountCode" comment:"The code of discount which apply to invoice"`
	OriginAmount                   string      `json:"OriginAmount"                 comment:"The original amount of invoice"`
	TotalAmount                    string      `json:"TotalAmount" comment:"The total amount of invoice"`
	DiscountAmount                 string      `json:"DiscountAmount" comment:"The discount amount of invoice"`
	TotalAmountExcludingTax        string      `json:"TotalAmountExcludingTax" comment:"The total amount of invoice with tax excluded"`
	Currency                       string      `json:"Currency" comment:"The currency of invoice"`
	TaxAmount                      string      `json:"TaxAmount" comment:"The tax amount of invoice"`
	TaxPercentage                  string      `json:"TaxPercentage"            comment:"The tax percentage of invoice applied"`
	SubscriptionAmount             string      `json:"SubscriptionAmount" comment:"The amount of subscription if invoice is generated by subscription"`
	SubscriptionAmountExcludingTax string      `json:"SubscriptionAmountExcludingTax" comment:"The amount of subscription which excluded tax amount if invoice is generated by subscription"`
	PeriodEnd                      *gtime.Time `json:"PeriodEnd"  layout:"2006-01-02 15:04:05"   comment:"The end time of period, will apply to subscription if invoice paid"`
	PeriodStart                    *gtime.Time `json:"PeriodStart"  layout:"2006-01-02 15:04:05"   comment:"The start time of period, will apply to subscription if invoice paid"`
	FinishTime                     *gtime.Time `json:"FinishTime"  layout:"2006-01-02 15:04:05"   comment:"The time when invoice finished, invoice will not be editable after finished"`
	CreateTime                     *gtime.Time `json:"CreateTime"  layout:"2006-01-02 15:04:05"   comment:"The create time of invoice"`
	DueDate                        *gtime.Time `json:"DueDate"  layout:"2006-01-02 15:04:05"   comment:"The due date of invoice, invoice will expired after due date"`
	PaidTime                       *gtime.Time `json:"PaidTime"  layout:"2006-01-02 15:04:05"   comment:"The paid time of invoice"`
	Status                         string      `json:"Status"                         comment:"The status of invoice，will be Pending｜Processing｜Paid | Failed | Cancelled | Reversed" `
	PaymentId                      string      `json:"PaymentId"                      comment:"paymentId" comment:"The id of payment connected to invoice"`
	RefundId                       string      `json:"RefundId"                       comment:"refundId" comment:"The id of refund connected to invoice, invoice will be credit not who contains refundId"`
	SubscriptionId                 string      `json:"SubscriptionId"                 comment:"subscription_id" comment:"the id of subscription connected to invoice"`
	PlanId                         string      `json:"PlanId"              comment:"The id of plan connected to invoice"`
	TrialEnd                       *gtime.Time `json:"TrialEnd"                       layout:"2006-01-02 15:04:05"   comment:"The time of trial end, will apply to subscription when invoice paid "`
	BillingCycleAnchor             *gtime.Time `json:"BillingCycleAnchor"             layout:"2006-01-02 15:04:05"   comment:"The subscription anchor time of billing cycle connected to invoice "`
	OriginInvoiceId                string      `json:"OriginInvoiceId"                description:"" comment:"The origin invoiceId if invoice type is credit note"`
	BillingPeriod                  string      `json:"BillingPeriod"      comment:"The billing period type of invoice, will be day|month|year or one-time purchase"`
}
