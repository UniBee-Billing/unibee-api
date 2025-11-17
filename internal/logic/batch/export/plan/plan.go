package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/plan"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskPlanExport struct {
}

func (t TaskPlanExport) TaskName() string {
	return "PlanExport"
}

func (t TaskPlanExport) Header() interface{} {
	return ExportPlanEntity{}
}

func (t TaskPlanExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &plan.ListInternalReq{
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
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["productIds"].([]interface{}); ok {
			req.ProductIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["type"].([]interface{}); ok {
			req.Type = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["publishStatus"].(float64); ok {
			req.PublishStatus = int(value)
		}
		if value, ok := payload["searchKey"].(string); ok {
			req.SearchKey = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
	}
	req.SkipTotal = true
	result, _ := plan.PlanList(ctx, req)
	if result != nil {
		for _, one := range result {
			var operationLog *entity.MerchantOperationLog
			_ = dao.MerchantOperationLog.Ctx(ctx).
				Where(dao.MerchantOperationLog.Columns().MerchantId, req.MerchantId).
				Where(dao.MerchantOperationLog.Columns().OptContent, "New").
				Where(dao.MerchantOperationLog.Columns().PlanId, one.Plan.Id).
				Scan(&operationLog)
			var createBy = ""
			if operationLog != nil {
				member := query.GetMerchantMemberById(ctx, operationLog.MemberId)
				if member != nil {
					createBy = fmt.Sprintf("%s_%s(%s)", member.FirstName, member.LastName, member.Email)
				}
			}

			mainList = append(mainList, &ExportPlanEntity{
				Id:                fmt.Sprintf("%v", one.Plan.Id),
				ExternalPlanId:    one.Plan.ExternalPlanId,
				MerchantName:      merchant.Name,
				ProductName:       one.Product.ProductName,
				PlanName:          one.Plan.PlanName,
				InternalName:      one.Plan.InternalName,
				Description:       one.Plan.Description,
				Status:            consts.PLanStatusToEnum(one.Plan.Status).Description(),
				PlanType:          consts.PlanTypeToEnum(one.Plan.Type).Description(),
				IsPublish:         fmt.Sprintf("%t", one.Plan.PublishStatus == 2),
				PlanAmount:        utility.ConvertCentToDollarStr(one.Plan.Amount, one.Plan.Currency),
				Currency:          one.Plan.Currency,
				IntervalUnit:      one.Plan.IntervalUnit,
				IntervalCount:     fmt.Sprintf("%d", one.Plan.IntervalCount),
				Metadata:          utility.MarshalToJsonString(one.Plan.Metadata),
				TrialAmount:       utility.ConvertCentToDollarStr(one.Plan.TrialAmount, one.Plan.Currency),
				TrialDurationTime: fmt.Sprintf("%d", one.Plan.TrialDurationTime),
				TrialDemand:       one.Plan.TrialDemand,
				CancelAtTrialEnd:  fmt.Sprintf("%t", one.Plan.CancelAtTrialEnd == 1),
				CheckoutUrl:       one.Plan.CheckoutUrl,
				DisableAutoCharge: fmt.Sprintf("%t", one.Plan.DisableAutoCharge == 1),
				CreateTime:        gtime.NewFromTimeStamp(one.Plan.CreateTime + timeZone),
				CreateBy:          createBy,
				TimeZone:          timeZoneStr,
			})
		}
	}
	return mainList, nil
}

type ExportPlanEntity struct {
	Id                string      `json:"Id"                        comment:""`
	ExternalPlanId    string      `json:"ExternalPlanId"            comment:"external_user_id"`
	MerchantName      string      `json:"MerchantName"              comment:""`
	ProductName       string      `json:"ProductName"               comment:""`
	PlanName          string      `json:"PlanName"                  comment:""`
	InternalName      string      `json:"InternalName"              comment:"PlanInternalName"`
	Description       string      `json:"Description"               comment:"description"`
	Status            string      `json:"Status"                    comment:""`
	PlanType          string      `json:"PlanType"                  comment:""`
	IsPublish         string      `json:"IsPublish"                 comment:""`
	PlanAmount        string      `json:"PlanAmount"                comment:""`
	Currency          string      `json:"Currency"                  comment:""`
	IntervalUnit      string      `json:"IntervalUnit"              comment:"period unit,day|month|year|week"`
	IntervalCount     string      `json:"IntervalCount"             comment:"period unit count"`
	Metadata          string      `json:"Metadata"                  comment:""`
	TrialAmount       string      `json:"TrialAmount"               comment:"price of trial period"`
	TrialDurationTime string      `json:"TrialDurationTime"         comment:"duration of trial，seconds"`
	TrialDemand       string      `json:"TrialDemand"               comment:""`
	CancelAtTrialEnd  string      `json:"CancelAtTrialEnd"          comment:"whether cancel at subscription first trial end"`
	CreateTime        *gtime.Time `json:"CreateTime"                layout:"2006-01-02 15:04:05"  comment:""`
	CreateBy          string      `json:"CreateBy"                  comment:""`
	TimeZone          string      `json:"TimeZone"                  comment:""`
	CheckoutUrl       string      `json:"CheckoutUrl"               comment:"Checkout Link"`
	DisableAutoCharge string      `json:"DisableAutoCharge"         comment:"disable auto-charge"`
}
