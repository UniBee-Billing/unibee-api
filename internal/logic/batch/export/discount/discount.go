package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/discount"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskDiscount struct {
}

func (t TaskDiscount) TaskName() string {
	return "DiscountExport"
}

func (t TaskDiscount) Header() interface{} {
	return ExportDiscountEntity{}
}

func (t TaskDiscount) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &discount.ListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	if payload != nil {
		if value, ok := payload["discountType"].([]int); ok {
			req.DiscountType = value
		}
		if value, ok := payload["billingType"].([]int); ok {
			req.BillingType = value
		}
		if value, ok := payload["status"].([]int); ok {
			req.Status = value
		}
		if value, ok := payload["code"].(string); ok {
			req.Code = value
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
		if value, ok := payload["createTimeStart"].(int64); ok {
			req.CreateTimeStart = value
		}
		if value, ok := payload["createTimeEnd"].(int64); ok {
			req.CreateTimeEnd = value
		}
	}
	result, _ := discount.MerchantDiscountCodeList(ctx, req)
	if result != nil {
		for _, one := range result {
			totalUsed, err := dao.MerchantUserDiscountCode.Ctx(ctx).
				Where(dao.MerchantUserDiscountCode.Columns().MerchantId, one.MerchantId).
				Where(dao.MerchantUserDiscountCode.Columns().Code, one.Code).
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
				StartTime:          gtime.NewFromTimeStamp(one.StartTime),
				EndTime:            gtime.NewFromTimeStamp(one.EndTime),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime),
				CreateBy:           createBy,
				TotalUsed:          fmt.Sprintf("%v", totalUsed),
			})
		}
	}
	return mainList, nil
}

type ExportDiscountEntity struct {
	Id                 string      `json:"Id"                `
	MerchantName       string      `json:"MerchantName"          `
	Name               string      `json:"Name"              `
	Code               string      `json:"Code"              `
	Status             string      `json:"Status"            `
	BillingType        string      `json:"BillingType"       `
	DiscountType       string      `json:"DiscountType"      `
	DiscountAmount     string      `json:"DiscountAmount"    `
	DiscountPercentage string      `json:"DiscountPercentage"`
	Currency           string      `json:"Currency"          `
	CycleLimit         string      `json:"CycleLimit"        `
	StartTime          *gtime.Time `json:"StartTime"         layout:"2006-01-02 15:04:05" `
	EndTime            *gtime.Time `json:"EndTime"           layout:"2006-01-02 15:04:05" `
	CreateTime         *gtime.Time `json:"CreateTime"        layout:"2006-01-02 15:04:05" `
	CreateBy           string      `json:"CreateBy"        `
	TotalUsed          string      `json:"TotalUsed"        `
}