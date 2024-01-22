package service

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionPendingUpdateListInternalReq struct {
	MerchantId     int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	SubscriptionId string `p:"subscriptionId" `
	SortField      string `p:"sortField" dc:"排序字段，gmt_create|gmt_modify" `
	SortType       string `p:"sortType" dc:"排序类型，asc|desc" `
	Page           int    `p:"page"  dc:"分页页码,0开始" `
	Count          int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}

type SubscriptionPendingUpdateListInternalRes struct {
	SubscriptionPendingUpdateDetails []*ro.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
}

func GetUnfinishedSubscriptionPendingUpdateDetailBySubscriptionId(ctx context.Context, subscriptionId string) *ro.SubscriptionPendingUpdateDetail {
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, subscriptionId).
		Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).
		OrderDesc(dao.SubscriptionPendingUpdate.Columns().Id).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	if one == nil {
		return nil
	}
	return &ro.SubscriptionPendingUpdateDetail{
		MerchantId:           one.MerchantId,
		SubscriptionId:       one.SubscriptionId,
		UpdateSubscriptionId: one.UpdateSubscriptionId,
		GmtCreate:            one.GmtCreate,
		Amount:               one.Amount,
		Status:               one.Status,
		UpdateAmount:         one.UpdateAmount,
		Currency:             one.Currency,
		UpdateCurrency:       one.UpdateCurrency,
		PlanId:               one.PlanId,
		UpdatePlanId:         one.UpdatePlanId,
		Quantity:             one.Quantity,
		UpdateQuantity:       one.UpdateQuantity,
		AddonData:            one.AddonData,
		UpdateAddonData:      one.UpdateAddonData,
		ChannelId:            one.ChannelId,
		UserId:               one.UserId,
		GmtModify:            one.GmtModify,
		Paid:                 one.Paid,
		Link:                 one.Link,
		MerchantUser:         query.GetMerchantAccountById(ctx, uint64(one.MerchantUserId)),
		EffectImmediate:      one.EffectImmediate,
		EffectTime:           one.EffectTime,
		AdminNote:            one.AdminNote,
		Plan:                 query.GetPlanById(ctx, one.PlanId),
		Addons:               query.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:           query.GetPlanById(ctx, one.UpdatePlanId),
		UpdateAddons:         query.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
	}
}

func SubscriptionPendingUpdateList(ctx context.Context, req *SubscriptionPendingUpdateListInternalReq) (res *SubscriptionPendingUpdateListInternalRes, err error) {
	var mainList []*entity.SubscriptionPendingUpdate
	if req.Count <= 0 {
		req.Count = 10 //每页数量默认 10
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_modify desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("gmt_create|gmt_modify", req.SortField), "sortField should one of gmt_create|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	err = dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().MerchantId, req.MerchantId).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, req.SubscriptionId).
		WhereNotNull(dao.SubscriptionPendingUpdate.Columns().MerchantUserId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}

	var updateList []*ro.SubscriptionPendingUpdateDetail
	for _, one := range mainList {
		updateList = append(updateList, &ro.SubscriptionPendingUpdateDetail{
			MerchantId:           one.MerchantId,
			SubscriptionId:       one.SubscriptionId,
			UpdateSubscriptionId: one.UpdateSubscriptionId,
			GmtCreate:            one.GmtCreate,
			Amount:               one.Amount,
			Status:               one.Status,
			UpdateAmount:         one.UpdateAmount,
			Currency:             one.Currency,
			UpdateCurrency:       one.UpdateCurrency,
			PlanId:               one.PlanId,
			UpdatePlanId:         one.UpdatePlanId,
			Quantity:             one.Quantity,
			UpdateQuantity:       one.UpdateQuantity,
			AddonData:            one.AddonData,
			UpdateAddonData:      one.UpdateAddonData,
			ChannelId:            one.ChannelId,
			UserId:               one.UserId,
			GmtModify:            one.GmtModify,
			Paid:                 one.Paid,
			Link:                 one.Link,
			MerchantUser:         query.GetMerchantAccountById(ctx, uint64(one.MerchantUserId)),
			EffectImmediate:      one.EffectImmediate,
			EffectTime:           one.EffectTime,
			AdminNote:            one.AdminNote,
			Plan:                 query.GetPlanById(ctx, one.PlanId),
			Addons:               query.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:           query.GetPlanById(ctx, one.UpdatePlanId),
			UpdateAddons:         query.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		})
	}

	return &SubscriptionPendingUpdateListInternalRes{SubscriptionPendingUpdateDetails: updateList}, nil
}
