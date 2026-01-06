package invoice

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskCreditNoteExport struct {
}

func (t TaskCreditNoteExport) TaskName() string {
	return fmt.Sprintf("CreditNoteExport")
}

func (t TaskCreditNoteExport) Header() interface{} {
	return ExportCreditNoteEntity{}
}

func (t TaskCreditNoteExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.CreditNoteListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZone = zone
			}
		}
		if value, ok := payload["gatewayIds"].([]interface{}); ok {
			req.GatewayIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["searchKey"].(string); ok {
			req.SearchKey = value
		}
		if value, ok := payload["emails"].(string); ok {
			emails := make([]string, 0)

			// 1. Process directly passed emails parameter
			if len(value) > 0 {
				cleanedEmails := strings.ReplaceAll(value, ";", ",")
				emails = strings.Split(cleanedEmails, ",")
				// Clean each email, remove spaces
				for i, email := range emails {
					emails[i] = strings.TrimSpace(email)
				}
				// Filter empty strings
				var filteredEmails []string
				for _, email := range emails {
					if email != "" {
						filteredEmails = append(filteredEmails, email)
					}
				}
				emails = filteredEmails
			}
			req.Emails = emails
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
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
	result, _ := service.CreditNoteList(ctx, req)
	if result != nil && result.CreditNotes != nil {
		for _, one := range result.CreditNotes {
			var creditNoteGateway = ""
			if one.Gateway != nil {
				creditNoteGateway = one.Gateway.GatewayName
			}
			if one.UserSnapshot == nil {
				one.UserSnapshot = &bean.UserAccount{}
			}
			if one.Refund == nil {
				one.Refund = &bean.Refund{}
			}
			var planIdStr string
			var planNameStr string
			if one.PlanSnapshot != nil && one.PlanSnapshot.Plan != nil {
				planIdStr = fmt.Sprintf("%v", one.PlanSnapshot.Plan.Id)
				planNameStr = one.PlanSnapshot.Plan.PlanName
			}
			mainList = append(mainList, &ExportCreditNoteEntity{
				CreditNoteId:        one.InvoiceId,
				UserId:              fmt.Sprintf("%v", one.UserId),
				Email:               one.UserSnapshot.Email,
				FirstName:           one.UserSnapshot.FirstName,
				LastName:            one.UserSnapshot.LastName,
				CreditNoteName:      one.InvoiceName,
				ProductName:         one.ProductName,
				Currency:            one.Currency,
				TotalAmount:         utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
				TaxAmount:           utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency),
				Status:              consts.InvoiceStatusToEnum(one.Status).Description(),
				Gateway:             creditNoteGateway,
				CreateTime:          gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				FinishTime:          gtime.NewFromTimeStamp(one.FinishTime + timeZone),
				RefundId:            one.RefundId,
				PaymentId:           one.PaymentId,
				PlanId:              planIdStr,
				PlanName:            planNameStr,
				SubscriptionId:      one.SubscriptionId,
				RefundReason:        one.Refund.RefundComment,
				PartialCreditAmount: utility.ConvertCentToDollarStr(one.PartialCreditPaidAmount, one.Currency),
			})
		}
	}
	return mainList, nil
}

type ExportCreditNoteEntity struct {
	CreditNoteId        string      `json:"CreditNoteId" comment:"The unique id of credit note" group:"Credit Note"`
	UserId              string      `json:"UserId" comment:"The unique id of user" group:"User Information"`
	Email               string      `json:"Email" comment:"The email of user" group:"User Information"`
	FirstName           string      `json:"FirstName" comment:"The first name of user" group:"User Information"`
	LastName            string      `json:"LastName" comment:"The last name of user" group:"User Information"`
	CreditNoteName      string      `json:"CreditNoteName" comment:"The name of credit note" group:"Credit Note"`
	ProductName         string      `json:"ProductName" comment:"The product name of credit note" group:"Product Information"`
	PlanId              string      `json:"PlanId"              comment:"The id of plan connected to invoice" group:"Product and Subscription"`
	PlanName            string      `json:"PlanName"              comment:"The name of plan connected to invoice" group:"Product and Subscription"`
	Currency            string      `json:"Currency" comment:"The currency of credit note" group:"Transaction"`
	TotalAmount         string      `json:"TotalAmount" comment:"The total amount of credit note (negative value)" group:"Transaction"`
	TaxAmount           string      `json:"TaxAmount" comment:"The tax amount of credit note" group:"Transaction"`
	Status              string      `json:"Status" comment:"The status of credit note" group:"Credit Note"`
	Gateway             string      `json:"Gateway" comment:"The gateway name of credit note" group:"Transaction"`
	CreateTime          *gtime.Time `json:"CreateTime" layout:"2006-01-02 15:04:05" comment:"The create time of credit note" group:"Credit Note"`
	FinishTime          *gtime.Time `json:"FinishTime" layout:"2006-01-02 15:04:05" comment:"The finish time of credit note" group:"Credit Note"`
	RefundId            string      `json:"RefundId" comment:"The refund ID connected to credit note" group:"Transaction"`
	PaymentId           string      `json:"PaymentId" comment:"The payment ID connected to credit note" group:"Transaction"`
	SubscriptionId      string      `json:"SubscriptionId" comment:"The subscription ID connected to credit note" group:"Subscription"`
	RefundReason        string      `json:"RefundReason" comment:"The reason for the refund/credit note" group:"Transaction"`
	PartialCreditAmount string      `json:"PartialCreditAmount" comment:"The partial credit amount if applicable" group:"Transaction"`
}
