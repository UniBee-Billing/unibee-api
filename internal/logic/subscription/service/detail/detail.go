package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*detail.SubscriptionDetail, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	{
		one.Data = ""
		one.ResponseData = ""
	}
	user := query.GetUserAccountById(ctx, one.UserId)
	var addonParams []*bean.PlanAddonParam
	if len(one.AddonData) > 0 {
		err := utility.UnmarshalFromJsonString(one.AddonData, &addonParams)
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionDetail parse addon param:%s", err.Error())
		}
	}
	latestInvoiceOne := bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId))
	if latestInvoiceOne != nil {
		latestInvoiceOne.Discount = bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, one.MerchantId, latestInvoiceOne.DiscountCode))
		latestInvoiceOne.PromoCreditTransaction = bean.SimplifyCreditTransaction(ctx, query.GetPromoCreditTransactionByInvoiceId(ctx, latestInvoiceOne.UserId, latestInvoiceOne.InvoiceId))
	}
	return &detail.SubscriptionDetail{
		User:                                bean.SimplifyUserAccount(user),
		Subscription:                        bean.SimplifySubscription(ctx, one),
		Gateway:                             detail.ConvertGatewayDetail(ctx, query.GetGatewayById(ctx, one.GatewayId)),
		Plan:                                bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		AddonParams:                         addonParams,
		Addons:                              addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		LatestInvoice:                       latestInvoiceOne,
		Discount:                            bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, one.MerchantId, one.DiscountCode)),
		UnfinishedSubscriptionPendingUpdate: GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId(ctx, one.PendingUpdateId),
	}, nil
}

func GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId(ctx context.Context, pendingUpdateId string) *detail.SubscriptionPendingUpdateDetail {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdateId).
		Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).
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
			fmt.Printf("GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId Unmarshal Metadata error:%s", err.Error())
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
