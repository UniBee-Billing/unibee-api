package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type TaskInvoice struct {
}

func (t TaskInvoice) TaskName() string {
	return fmt.Sprintf("InvoiceExport")
}

func (t TaskInvoice) Header() interface{} {
	return ExportInvoiceEntity{}
}

func (t TaskInvoice) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
		if value, ok := payload["AmountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
	}
	result, _ := service.InvoiceList(ctx, req)
	if result != nil && result.Invoices != nil {
		for _, one := range result.Invoices {
			var invoiceGateway = ""
			if one.Gateway != nil {
				invoiceGateway = one.Gateway.GatewayName
			}
			var firstName = ""
			var lastName = ""
			var email = ""
			if one.UserAccount != nil {
				firstName = one.UserAccount.FirstName
				lastName = one.UserAccount.LastName
				email = one.UserAccount.Email
			}
			mainList = append(mainList, &ExportInvoiceEntity{
				InvoiceId:                      one.InvoiceId,
				UserId:                         fmt.Sprintf("%v", one.UserId),
				FirstName:                      firstName,
				LastName:                       lastName,
				Email:                          email,
				InvoiceName:                    one.InvoiceName,
				ProductName:                    one.ProductName,
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
				Status:                         consts.InvoiceStatusToEnum(one.Status).Description(),
				PaymentId:                      one.PaymentId,
				RefundId:                       one.RefundId,
				SubscriptionId:                 one.SubscriptionId,
				TrialEnd:                       gtime.NewFromTimeStamp(one.TrialEnd),
				BillingCycleAnchor:             gtime.NewFromTimeStamp(one.BillingCycleAnchor),
				CreateFrom:                     one.CreateFrom,
				CountryCode:                    one.CountryCode,
			})
		}
	}
	return mainList, nil
}

type ExportInvoiceEntity struct {
	InvoiceId                      string      `json:"InvoiceId"`
	UserId                         string      `json:"UserId"                 `
	FirstName                      string      `json:"FirstName"          `
	LastName                       string      `json:"LastName"           `
	Email                          string      `json:"Email"              `
	InvoiceName                    string      `json:"InvoiceName"`
	ProductName                    string      `json:"ProductName"`
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
	Status                         string      `json:"Status"                         description:"status，1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	PaymentId                      string      `json:"PaymentId"                      description:"paymentId"`                                                     // paymentId
	RefundId                       string      `json:"RefundId"                       description:"refundId"`                                                      // refundId
	SubscriptionId                 string      `json:"SubscriptionId"                 description:"subscription_id"`                                               // subscription_id
	TrialEnd                       *gtime.Time `json:"TrialEnd"                       layout:"2006-01-02 15:04:05"  `                                              // trial_end, utc time
	BillingCycleAnchor             *gtime.Time `json:"BillingCycleAnchor"             layout:"2006-01-02 15:04:05"  `                                              // billing_cycle_anchor
	CreateFrom                     string      `json:"CreateFrom"                     description:"create from"`                                                   // create from
	CountryCode                    string      `json:"CountryCode"                    description:""`                                                              //
}
