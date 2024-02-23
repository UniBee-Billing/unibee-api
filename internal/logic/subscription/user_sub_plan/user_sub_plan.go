package user_sub_plan

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	addon2 "unibee-api/internal/logic/subscription/addon"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

const (
	UserSubPlanCacheKeyPrefix = "UserSubPlanCacheKeyPrefix_"
	UserSubPlanCacheKeyExpire = 24 * 60 * 60
)

type UserSubPlan struct {
	MerchantId              uint64
	UserId                  int64
	PlanId                  uint64
	PlanType                int
	Quantity                int64
	SubscriptionIds         string
	SubscriptionPeriodStart int64
	SubscriptionPeriodEnd   int64
}

func UserSubPlanCachedList(ctx context.Context, merchantId uint64, userId int64, sub *entity.Subscription, reloadCache bool) []*UserSubPlan {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(userId > 0, "invalid userId")
	if sub == nil {
		return make([]*UserSubPlan, 0)
	}
	var list = make([]*UserSubPlan, 0)
	cacheKey := fmt.Sprintf("%s_%d_%d", UserSubPlanCacheKeyPrefix, merchantId, userId)
	if !reloadCache {
		get, err := g.Redis().Get(ctx, cacheKey)
		if err == nil && !get.IsNil() && !get.IsEmpty() {
			value := get.String()
			_ = utility.UnmarshalFromJsonString(value, &list)
			if len(list) > 0 {
				return list
			}
		}
	}
	if merchantId > 0 {
		plan := query.GetPlanById(ctx, sub.PlanId)
		if plan != nil {
			list = append(list, &UserSubPlan{
				MerchantId:              sub.MerchantId,
				UserId:                  userId,
				PlanId:                  sub.PlanId,
				PlanType:                plan.Type,
				Quantity:                sub.Quantity,
				SubscriptionIds:         sub.SubscriptionId,
				SubscriptionPeriodStart: sub.CurrentPeriodStart,
				SubscriptionPeriodEnd:   sub.CurrentPeriodEnd,
			})
		}
		//append addons
		addons := addon2.GetSubscriptionAddonsByAddonJson(ctx, sub.AddonData)
		for _, addon := range addons {
			list = append(list, &UserSubPlan{
				MerchantId:              sub.MerchantId,
				UserId:                  userId,
				PlanId:                  addon.AddonPlan.Id,
				PlanType:                addon.AddonPlan.Type,
				Quantity:                addon.Quantity,
				SubscriptionIds:         sub.SubscriptionId,
				SubscriptionPeriodStart: sub.CurrentPeriodStart,
				SubscriptionPeriodEnd:   sub.CurrentPeriodEnd,
			})
		}
	}
	if len(list) > 0 {
		_, _ = g.Redis().Set(ctx, cacheKey, utility.MarshalToJsonString(list))
		_, _ = g.Redis().Expire(ctx, cacheKey, UserSubPlanCacheKeyExpire) // one day cache expire time
	}
	return list
}

func ReloadUserSubPlanCacheListBackground(merchantId uint64, userId int64) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(ctx, "ReloadUserSubPlanCacheListBackground panic error:%s", err.Error())
				return
			}
		}()
		sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, userId, merchantId)
		if sub != nil {
			UserSubPlanCachedList(ctx, merchantId, userId, sub, true)
		}
	}()
}
