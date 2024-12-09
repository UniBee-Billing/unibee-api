package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/discount"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskUserDiscountExport struct {
}

func (t TaskUserDiscountExport) TaskName() string {
	return "UserDiscountExport"
}

func (t TaskUserDiscountExport) Header() interface{} {
	return ExportUserDiscountEntity{}
}

func (t TaskUserDiscountExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	var id int64
	if value, ok := payload["id"].(float64); ok {
		id = int64(value)
	}
	if id <= 0 {
		return mainList, nil
	}
	req := &discount.UserDiscountListInternalReq{
		MerchantId: task.MerchantId,
		Id:         uint64(id),
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
		if value, ok := payload["userIds"].([]interface{}); ok {
			req.UserIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
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
	}
	req.SkipTotal = true
	result, _ := discount.MerchantUserDiscountCodeList(ctx, req)
	if result != nil {
		for _, one := range result {
			var firstName = ""
			var lastName = ""
			var email = ""
			if one.User != nil {
				firstName = one.User.FirstName
				lastName = one.User.LastName
				email = one.User.Email
			} else {
				one.User = &bean.UserAccount{}
			}
			if one.Plan == nil {
				one.Plan = &bean.Plan{}
			}
			recurring := "No"
			if one.Recurring == 1 {
				recurring = "Yes"
			}
			statusStr := "Finished"
			if one.Status == 2 {
				statusStr = "Rollback"
			}
			mainList = append(mainList, &ExportUserDiscountEntity{
				Id:             fmt.Sprintf("%v", one.Id),
				UserId:         fmt.Sprintf("%v", one.User.Id),
				ExternalUserId: fmt.Sprintf("%v", one.User.ExternalUserId),
				MerchantName:   merchant.Name,
				FirstName:      firstName,
				LastName:       lastName,
				Email:          email,
				PlanId:         fmt.Sprintf("%v", one.Plan.Id),
				ExternalPlanId: fmt.Sprintf("%v", one.Plan.ExternalPlanId),
				PlanName:       one.Plan.PlanName,
				Code:           one.Code,
				Status:         statusStr,
				SubscriptionId: one.SubscriptionId,
				TransactionId:  one.PaymentId,
				InvoiceId:      one.InvoiceId,
				CreateTime:     gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				ApplyAmount:    utility.ConvertCentToDollarStr(one.ApplyAmount, one.Currency),
				Currency:       one.Currency,
				Recurring:      recurring,
				TimeZone:       timeZoneStr,
			})
		}
	}
	return mainList, nil
}

type ExportUserDiscountEntity struct {
	Id             string      `json:"Id"                  comment:""`
	UserId         string      `json:"UserId"              comment:""`
	ExternalUserId string      `json:"ExternalUserId"      comment:""`
	PlanId         string      `json:"PlanId"              comment:""`
	ExternalPlanId string      `json:"ExternalPlanId"      comment:""`
	FirstName      string      `json:"FirstName"           comment:""`
	LastName       string      `json:"LastName"            comment:""`
	Email          string      `json:"Email"               comment:""`
	MerchantName   string      `json:"MerchantName"        comment:""`
	PlanName       string      `json:"PlanName"            comment:""`
	Status         string      `json:"Status"            comment:""`
	Code           string      `json:"Code"            comment:""`
	SubscriptionId string      `json:"SubscriptionId"  comment:""`
	TransactionId  string      `json:"TransactionId"       comment:""`
	InvoiceId      string      `json:"InvoiceId"       comment:""`
	CreateTime     *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05" comment:""`
	ApplyAmount    string      `json:"ApplyAmount"     comment:""`
	Currency       string      `json:"Currency"        comment:""`
	Recurring      string      `json:"recurring"      description:"is recurring apply, Yes|No"` // is recurring apply, 0-no, 1-yes
	TimeZone       string      `json:"TimeZone"         comment:""`
}
