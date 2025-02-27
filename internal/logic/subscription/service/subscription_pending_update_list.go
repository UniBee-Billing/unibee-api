package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPendingUpdateListInternalReq struct {
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	SubscriptionId string `json:"subscriptionId" `
	SortField      string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify" `
	SortType       string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count"  dc:"Count Of Page"`
}

type SubscriptionPendingUpdateListInternalRes struct {
	SubscriptionPendingUpdateDetails []*detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdateDetails" dc:"SubscriptionPendingUpdateDetails"`
	Total                            int                                       `json:"total" dc:"Total"`
}

func GetSubscriptionPendingUpdateEventByPendingUpdateId(ctx context.Context, pendingUpdateId string) *detail.SubscriptionPendingUpdateEvent {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdateId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("GetSubscriptionPendingUpdateDetailByPendingUpdateId Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &detail.SubscriptionPendingUpdateEvent{
		MerchantId:      one.MerchantId,
		User:            bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		Subscription:    bean.SimplifySubscription(ctx, query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)),
		Invoice:         bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)),
		PendingUpdateId: one.PendingUpdateId,
		GmtCreate:       one.GmtCreate,
		Amount:          one.Amount,
		Status:          one.Status,
		UpdateAmount:    one.UpdateAmount,
		Currency:        one.Currency,
		UpdateCurrency:  one.UpdateCurrency,
		PlanId:          one.PlanId,
		UpdatePlanId:    one.UpdatePlanId,
		Quantity:        one.Quantity,
		UpdateQuantity:  one.UpdateQuantity,
		AddonData:       one.AddonData,
		UpdateAddonData: one.UpdateAddonData,
		ProrationAmount: one.ProrationAmount,
		GatewayId:       one.GatewayId,
		GmtModify:       one.GmtModify,
		Paid:            one.Paid,
		Link:            one.Link,
		MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
		EffectImmediate: one.EffectImmediate,
		EffectTime:      one.EffectTime,
		Note:            one.Note,
		Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
		UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		Metadata:        metadata,
	}
}

func GetSubscriptionPendingUpdateDetailByPendingUpdateId(ctx context.Context, pendingUpdateId string) *detail.SubscriptionPendingUpdateDetail {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdateId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("GetSubscriptionPendingUpdateDetailByPendingUpdateId Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &detail.SubscriptionPendingUpdateDetail{
		MerchantId:      one.MerchantId,
		SubscriptionId:  one.SubscriptionId,
		PendingUpdateId: one.PendingUpdateId,
		GmtCreate:       one.GmtCreate,
		Amount:          one.Amount,
		Status:          one.Status,
		UpdateAmount:    one.UpdateAmount,
		Currency:        one.Currency,
		UpdateCurrency:  one.UpdateCurrency,
		PlanId:          one.PlanId,
		UpdatePlanId:    one.UpdatePlanId,
		Quantity:        one.Quantity,
		UpdateQuantity:  one.UpdateQuantity,
		AddonData:       one.AddonData,
		UpdateAddonData: one.UpdateAddonData,
		ProrationAmount: one.ProrationAmount,
		GatewayId:       one.GatewayId,
		UserId:          one.UserId,
		InvoiceId:       one.InvoiceId,
		GmtModify:       one.GmtModify,
		Paid:            one.Paid,
		Link:            one.Link,
		MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
		EffectImmediate: one.EffectImmediate,
		EffectTime:      one.EffectTime,
		Note:            one.Note,
		Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
		UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		Metadata:        metadata,
	}
}

func SubscriptionPendingUpdateList(ctx context.Context, req *SubscriptionPendingUpdateListInternalReq) (res *SubscriptionPendingUpdateListInternalRes, err error) {
	var mainList []*entity.SubscriptionPendingUpdate
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
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
		OmitEmpty().ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, err
	}

	var updateList []*detail.SubscriptionPendingUpdateDetail
	for _, one := range mainList {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SubscriptionPendingUpdateList Unmarshal Metadata error:%s", err.Error())
			}
		}
		updateList = append(updateList, &detail.SubscriptionPendingUpdateDetail{
			MerchantId:      one.MerchantId,
			SubscriptionId:  one.SubscriptionId,
			PendingUpdateId: one.PendingUpdateId,
			GmtCreate:       one.GmtCreate,
			Amount:          one.Amount,
			Status:          one.Status,
			UpdateAmount:    one.UpdateAmount,
			Currency:        one.Currency,
			UpdateCurrency:  one.UpdateCurrency,
			PlanId:          one.PlanId,
			UpdatePlanId:    one.UpdatePlanId,
			Quantity:        one.Quantity,
			UpdateQuantity:  one.UpdateQuantity,
			AddonData:       one.AddonData,
			UpdateAddonData: one.UpdateAddonData,
			ProrationAmount: one.ProrationAmount,
			GatewayId:       one.GatewayId,
			UserId:          one.UserId,
			InvoiceId:       one.InvoiceId,
			GmtModify:       one.GmtModify,
			Paid:            one.Paid,
			Link:            one.Link,
			MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate: one.EffectImmediate,
			EffectTime:      one.EffectTime,
			Note:            one.Note,
			Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:        metadata,
		})
	}

	return &SubscriptionPendingUpdateListInternalRes{SubscriptionPendingUpdateDetails: updateList, Total: total}, nil
}
