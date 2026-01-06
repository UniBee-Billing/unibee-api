package user

import (
	"context"
	"fmt"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	preload2 "unibee/internal/logic/preload"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskUserV2Export struct {
}

func (t TaskUserV2Export) TaskName() string {
	return "UserExport"
}

func (t TaskUserV2Export) Header() interface{} {
	return ExportUserEntity{}
}

func (t TaskUserV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &userListInternalReq{
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
			req.UserId = int64(value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
		}
		if value, ok := payload["firstName"].(string); ok {
			req.FirstName = value
		}
		if value, ok := payload["lastName"].(string); ok {
			req.LastName = value
		}
		if value, ok := payload["subscriptionId"].(string); ok {
			req.SubscriptionId = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["subStatus"].([]interface{}); ok {
			req.SubStatus = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["gatewayIds"].([]interface{}); ok {
			req.GatewayIds = export.JsonArrayTypeConvertInt64(ctx, value)
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
	result, _ := userList(ctx, req)
	if result != nil && result.UserAccounts != nil {
		// Preload all related data
		preload := preload2.UserListPreload(ctx, result.UserAccounts)

		for _, one := range result.UserAccounts {
			var userGateway = ""
			if gateway, ok := preload.Gateways[one.GatewayId]; ok {
				userGateway = gateway.GatewayName
			}

			var taxPercentage = one.TaxPercentage
			if vatRate, ok := preload.UserTaxPercentages[one.Id]; ok {
				taxPercentage = vatRate
			}

			mainList = append(mainList, &ExportUserEntity{
				Id:                 fmt.Sprintf("%v", one.Id),
				FirstName:          one.FirstName,
				LastName:           one.LastName,
				Email:              one.Email,
				MerchantName:       merchant.Name,
				Phone:              one.Phone,
				Address:            one.Address,
				VatNumber:          one.VATNumber,
				CountryCode:        one.CountryCode,
				CountryName:        one.CountryName,
				SubscriptionName:   one.SubscriptionName,
				SubscriptionId:     one.SubscriptionId,
				SubscriptionStatus: consts.SubStatusToEnum(one.SubscriptionStatus).Description(),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				ExternalUserId:     one.ExternalUserId,
				Status:             consts.UserStatusToEnum(one.Status).Description(),
				TaxPercentage:      utility.ConvertTaxPercentageToPercentageString(taxPercentage),
				Type:               consts.UserTypeToEnum(one.Type).Description(),
				Gateway:            userGateway,
				City:               one.City,
				ZipCode:            one.ZipCode,
				TimeZone:           timeZoneStr,
			})
		}
	}
	return mainList, nil
}

// userListInternalReq is a local copy of UserListInternalReq to avoid import cycles
type userListInternalReq struct {
	MerchantId      uint64  `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId          int64   `json:"userId" dc:"Filter UserId, Default All" `
	Email           string  `json:"email" dc:"Search Email" `
	FirstName       string  `json:"firstName" dc:"Search FirstName" `
	LastName        string  `json:"lastName" dc:"Search LastName" `
	SubscriptionId  string  `json:"subscriptionId" dc:"Search Filter SubscriptionId" `
	SubStatus       []int   `json:"subStatus" dc:"Filter, Default All，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing | 9-Failed" `
	Status          []int   `json:"status" dc:"Status, 0-Active｜2-Frozen" `
	PlanIds         []int   `json:"planIds" dc:"PlanIds, Search Filter PlanIds" `
	GatewayIds      []int64 `json:"gatewayIds" dc:"GatewayIds, Search Filter GatewayIds" `
	DeleteInclude   bool    `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField       string  `json:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType        string  `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int     `json:"page"  dc:"Page,Start 0" `
	Count           int     `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64   `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64   `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

// userListInternalRes is a local copy of UserListInternalRes to avoid import cycles
type userListInternalRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
	Total        int                   `json:"total" dc:"Total"`
}

// userList is a local copy of UserList to avoid import cycles
func userList(ctx context.Context, req *userListInternalReq) (res *userListInternalRes, err error) {
	var mainList []*entity.UserAccount
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	var isDeletes = []int{0}
	if req.DeleteInclude {
		isDeletes = append(isDeletes, 1)
	}
	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type", req.SortField), "sortField should one of user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().MerchantId, req.MerchantId).
		WhereIn(dao.UserAccount.Columns().IsDeleted, isDeletes)
	if req.UserId > 0 {
		q = q.Where(dao.UserAccount.Columns().Id, req.UserId)
	}
	if len(req.SubscriptionId) > 0 {
		q = q.Where(dao.UserAccount.Columns().SubscriptionId, req.SubscriptionId)
	}
	if len(req.Email) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().Email, "%"+req.Email+"%")
	}
	if len(req.FirstName) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
	}
	if len(req.LastName) > 0 {
		q = q.WhereLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
	}
	if req.SubStatus != nil && len(req.SubStatus) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().SubscriptionStatus, req.SubStatus)
	}
	if req.Status != nil && len(req.Status) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().Status, req.Status)
	}
	if req.PlanIds != nil && len(req.PlanIds) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().PlanId, req.PlanIds)
	}
	if req.GatewayIds != nil && len(req.GatewayIds) > 0 {
		q = q.WhereIn(dao.UserAccount.Columns().GatewayId, req.GatewayIds)
	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeEnd)
	}
	q = q.Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	if req.SkipTotal {
		err = q.Scan(&mainList)
	} else {
		err = q.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return nil, err
	}
	return &userListInternalRes{UserAccounts: mainList, Total: total}, nil
}
