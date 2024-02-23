package service

import (
	"context"
	"strings"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	addon2 "unibee-api/internal/logic/subscription/addon"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type SubscriptionPendingUpdateListInternalReq struct {
	MerchantId     uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	SubscriptionId string `p:"subscriptionId" `
	SortField      string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify" `
	SortType       string `p:"sortType" dc:"Sort Type，asc|desc" `
	Page           int    `p:"page"  dc:"Page, Start WIth 0" `
	Count          int    `p:"count"  dc:"Count Of Page"`
}

type SubscriptionPendingUpdateListInternalRes struct {
	SubscriptionPendingUpdateDetails []*ro.SubscriptionPendingUpdateDetailVo `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
}

func GetUnfinishedSubscriptionPendingUpdateDetailByUpdateSubscriptionId(ctx context.Context, pendingUpdateId string) *ro.SubscriptionPendingUpdateDetailVo {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	if one == nil {
		return nil
	}
	return &ro.SubscriptionPendingUpdateDetailVo{
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
		ProrationAmount:      one.ProrationAmount,
		GatewayId:            one.GatewayId,
		UserId:               one.UserId,
		GmtModify:            one.GmtModify,
		Paid:                 one.Paid,
		Link:                 one.Link,
		MerchantUser:         ro.SimplifyMerchantUserAccount(query.GetMerchantUserAccountById(ctx, uint64(one.MerchantUserId))),
		EffectImmediate:      one.EffectImmediate,
		EffectTime:           one.EffectTime,
		Note:                 one.Note,
		Plan:                 ro.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:           ro.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
		UpdateAddons:         addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
	}
}

func SubscriptionPendingUpdateList(ctx context.Context, req *SubscriptionPendingUpdateListInternalReq) (res *SubscriptionPendingUpdateListInternalRes, err error) {
	var mainList []*entity.SubscriptionPendingUpdate
	if req.Count <= 0 {
		req.Count = 20
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

	var updateList []*ro.SubscriptionPendingUpdateDetailVo
	for _, one := range mainList {
		updateList = append(updateList, &ro.SubscriptionPendingUpdateDetailVo{
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
			ProrationAmount:      one.ProrationAmount,
			GatewayId:            one.GatewayId,
			UserId:               one.UserId,
			GmtModify:            one.GmtModify,
			Paid:                 one.Paid,
			Link:                 one.Link,
			MerchantUser:         ro.SimplifyMerchantUserAccount(query.GetMerchantUserAccountById(ctx, uint64(one.MerchantUserId))),
			EffectImmediate:      one.EffectImmediate,
			EffectTime:           one.EffectTime,
			Note:                 one.Note,
			Plan:                 ro.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:           ro.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:         addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		})
	}

	return &SubscriptionPendingUpdateListInternalRes{SubscriptionPendingUpdateDetails: updateList}, nil
}
