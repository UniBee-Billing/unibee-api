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
				CreateFrom:                     one.CreateFrom,
				OriginInvoiceId:                OriginInvoiceId,
				BillingPeriod:                  billingPeriod,
			})
		}
	}
	return mainList, nil
}

type ExportInvoiceEntity struct {
	InvoiceId                      string      `json:"InvoiceId"`
	InvoiceNumber                  string      `json:"InvoiceNumber"`
	UserId                         string      `json:"UserId"             `
	ExternalUserId                 string      `json:"ExternalUserId"     `
	FirstName                      string      `json:"FirstName"          `
	LastName                       string      `json:"LastName"           `
	FullName                       string      `json:"FullName"           `
	UserType                       string      `json:"UserType"           `
	Email                          string      `json:"Email"              `
	InvoiceName                    string      `json:"InvoiceName"`
	ProductName                    string      `json:"ProductName"`
	InvoiceType                    string      `json:"InvoiceType"`
	Address                        string      `json:"Address"`
	City                           string      `json:"City"`
	State                          string      `json:"State"`
	CountryCode                    string      `json:"CountryCode"                    description:""`
	TaxCode                        string      `json:"TaxCode"                    description:""`
	CountryName                    string      `json:"CountryName"`
	IsEU                           string      `json:"EU/Non-EU"`
	Gateway                        string      `json:"Gateway"            `
	MerchantName                   string      `json:"MerchantName"       `
	DiscountCode                   string      `json:"DiscountCode"`
	OriginAmount                   string      `json:"OriginAmount"                `
	TotalAmount                    string      `json:"TotalAmount"`
	DiscountAmount                 string      `json:"DiscountAmount"`
	TotalAmountExcludingTax        string      `json:"TotalAmountExcludingTax"`
	Currency                       string      `json:"Currency"`
	TaxAmount                      string      `json:"TaxAmount"`
	TaxPercentage                  string      `json:"TaxPercentage"           `
	SubscriptionAmount             string      `json:"SubscriptionAmount"`
	SubscriptionAmountExcludingTax string      `json:"SubscriptionAmountExcludingTax"`
	PeriodEnd                      *gtime.Time `json:"PeriodEnd"  layout:"2006-01-02 15:04:05"  `
	PeriodStart                    *gtime.Time `json:"PeriodStart"  layout:"2006-01-02 15:04:05"  `
	FinishTime                     *gtime.Time `json:"FinishTime"  layout:"2006-01-02 15:04:05"  `
	CreateTime                     *gtime.Time `json:"CreateTime"  layout:"2006-01-02 15:04:05"  `
	DueDate                        *gtime.Time `json:"DueDate"  layout:"2006-01-02 15:04:05"  `
	PaidTime                       *gtime.Time `json:"PaidTime"  layout:"2006-01-02 15:04:05"  `
	Status                         string      `json:"Status"                         description:"status，1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	PaymentId                      string      `json:"PaymentId"                      description:"paymentId"`                                                     // paymentId
	RefundId                       string      `json:"RefundId"                       description:"refundId"`                                                      // refundId
	SubscriptionId                 string      `json:"SubscriptionId"                 description:"subscription_id"`                                               // subscription_id
	PlanId                         string      `json:"PlanId"             `
	TrialEnd                       *gtime.Time `json:"TrialEnd"                       layout:"2006-01-02 15:04:05"  ` // trial_end, utc time
	BillingCycleAnchor             *gtime.Time `json:"BillingCycleAnchor"             layout:"2006-01-02 15:04:05"  ` // billing_cycle_anchor
	CreateFrom                     string      `json:"CreateFrom"                     description:"create from"`      // create from

	OriginInvoiceId string `json:"OriginInvoiceId"                description:""`
	BillingPeriod   string `json:"BillingPeriod"     `
}
