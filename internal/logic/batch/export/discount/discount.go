package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/discount"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskDiscountExport struct {
}

func (t TaskDiscountExport) TaskName() string {
	return "DiscountExport"
}

func (t TaskDiscountExport) Header() interface{} {
	return ExportDiscountEntity{}
}

func (t TaskDiscountExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &discount.ListInternalReq{
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
		if value, ok := payload["discountType"].([]interface{}); ok {
			req.DiscountType = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["billingType"].([]interface{}); ok {
			req.BillingType = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["code"].(string); ok {
			req.Code = value
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
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
		}
	}
	req.SkipTotal = true
	result, _ := discount.MerchantDiscountCodeList(ctx, req)
	if result != nil {
		for _, one := range result {
			totalUsed, err := dao.MerchantUserDiscountCode.Ctx(ctx).
				Where(dao.MerchantUserDiscountCode.Columns().MerchantId, one.MerchantId).
				//Where(dao.MerchantUserDiscountCode.Columns().Code, one.Code).
				Where("LOWER(code) = LOWER(?)", one.Code). // case_insensitive
				Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
				Count()
			if err != nil {
				totalUsed = 0
			}
			var operationLog *entity.MerchantOperationLog
			_ = dao.MerchantOperationLog.Ctx(ctx).
				Where(dao.MerchantOperationLog.Columns().MerchantId, req.MerchantId).
				Where(dao.MerchantOperationLog.Columns().OptContent, "New").
				Where(dao.MerchantOperationLog.Columns().DiscountCode, one.Code).
				Scan(&operationLog)
			var createBy = ""
			if operationLog != nil {
				member := query.GetMerchantMemberById(ctx, operationLog.MemberId)
				if member != nil {
					createBy = fmt.Sprintf("%s_%s(%s)", member.FirstName, member.LastName, member.Email)
				}
			}

			mainList = append(mainList, &ExportDiscountEntity{
				Id:                 fmt.Sprintf("%v", one.Id),
				MerchantName:       merchant.Name,
				Name:               one.Name,
				Code:               one.Code,
				Status:             consts.DiscountStatusToEnum(one.Status).Description(),
				BillingType:        consts.DiscountBillingTypeToEnum(one.BillingType).Description(),
				DiscountType:       consts.DiscountTypeToEnum(one.DiscountType).Description(),
				DiscountAmount:     utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency),
				DiscountPercentage: utility.ConvertTaxPercentageToPercentageString(one.DiscountPercentage),
				Currency:           one.Currency,
				CycleLimit:         fmt.Sprintf("%v", one.CycleLimit),
				StartTime:          gtime.NewFromTimeStamp(one.StartTime + timeZone),
				EndTime:            gtime.NewFromTimeStamp(one.EndTime + timeZone),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				CreateBy:           createBy,
				TotalUsed:          fmt.Sprintf("%v", totalUsed),
				TimeZone:           timeZoneStr,
			})
		}
	}
	return mainList, nil
}

type ExportDiscountEntity struct {
	Id                 string      `json:"Id"                 comment:""`
	MerchantName       string      `json:"MerchantName"           comment:""`
	Name               string      `json:"Name"               comment:""`
	Code               string      `json:"Code"               comment:""`
	Status             string      `json:"Status"             comment:""`
	BillingType        string      `json:"BillingType"        comment:""`
	DiscountType       string      `json:"DiscountType"       comment:""`
	DiscountAmount     string      `json:"DiscountAmount"     comment:""`
	DiscountPercentage string      `json:"DiscountPercentage" comment:""`
	Currency           string      `json:"Currency"           comment:""`
	CycleLimit         string      `json:"CycleLimit"         comment:""`
	StartTime          *gtime.Time `json:"StartTime"         layout:"2006-01-02 15:04:05"  comment:""`
	EndTime            *gtime.Time `json:"EndTime"           layout:"2006-01-02 15:04:05"  comment:""`
	CreateTime         *gtime.Time `json:"CreateTime"        layout:"2006-01-02 15:04:05"  comment:""`
	CreateBy           string      `json:"CreateBy"         comment:""`
	TotalUsed          string      `json:"TotalUsed"         comment:""`
	TimeZone           string      `json:"TimeZone"         comment:""`
}
