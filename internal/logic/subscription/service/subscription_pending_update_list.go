package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPendingUpdateListInternalReq struct {
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	SubscriptionId string `json:"subscriptionId" `
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page           int    `json:"page"  dc:"Page, Start WIth 0" `
	Count          int    `json:"count"  dc:"Count Of Page"`
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
	var metadata = make(map[string]string)
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("GetUnfinishedSubscriptionPendingUpdateDetailByUpdateSubscriptionId Unmarshal Metadata error:%s", err.Error())
		}
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
		MerchantMember:       ro.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
		EffectImmediate:      one.EffectImmediate,
		EffectTime:           one.EffectTime,
		Note:                 one.Note,
		Plan:                 ro.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:           ro.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
		UpdateAddons:         addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		Metadata:             metadata,
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
		WhereNotNull(dao.SubscriptionPendingUpdate.Columns().MerchantMemberId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}

	var updateList []*ro.SubscriptionPendingUpdateDetailVo
	for _, one := range mainList {
		var metadata = make(map[string]string)
		if len(one.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SubscriptionPendingUpdateList Unmarshal Metadata error:%s", err.Error())
			}
		}
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
			MerchantMember:       ro.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate:      one.EffectImmediate,
			EffectTime:           one.EffectTime,
			Note:                 one.Note,
			Plan:                 ro.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:           ro.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:         addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:             metadata,
		})
	}

	return &SubscriptionPendingUpdateListInternalRes{SubscriptionPendingUpdateDetails: updateList}, nil
}
