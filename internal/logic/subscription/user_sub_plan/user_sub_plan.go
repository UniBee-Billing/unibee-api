package user_sub_plan

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
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
	PlanId                  int64
	PlanType                int
	Quantity                int64
	SubscriptionIds         string
	SubscriptionPeriodStart int64
	SubscriptionPeriodEnd   int64
}

func UserSubPlanCachedList(ctx context.Context, merchantId uint64, userId int64, reloadCache bool) []*UserSubPlan {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(userId > 0, "invalid userId")
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
		var entities []*entity.Subscription
		var status = []int{consts.SubStatusActive, consts.SubStatusIncomplete, consts.SubStatusPendingInActive}
		err := dao.Subscription.Ctx(ctx).
			Where(dao.Subscription.Columns().MerchantId, merchantId).
			Where(dao.Subscription.Columns().UserId, userId).
			WhereIn(dao.Subscription.Columns().Status, status).
			Where(dao.Subscription.Columns().IsDeleted, 0).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				plan := query.GetPlanById(ctx, one.PlanId)
				if plan != nil {
					list = append(list, &UserSubPlan{
						MerchantId:              one.MerchantId,
						UserId:                  userId,
						PlanId:                  one.PlanId,
						PlanType:                plan.Type,
						Quantity:                one.Quantity,
						SubscriptionIds:         one.SubscriptionId,
						SubscriptionPeriodStart: one.CurrentPeriodStart,
						SubscriptionPeriodEnd:   one.CurrentPeriodEnd,
					})
				}
				//append addons
				addons := addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData)
				for _, addon := range addons {
					list = append(list, &UserSubPlan{
						MerchantId:              one.MerchantId,
						UserId:                  userId,
						PlanId:                  int64(addon.AddonPlan.Id),
						PlanType:                addon.AddonPlan.Type,
						Quantity:                addon.Quantity,
						SubscriptionIds:         one.SubscriptionId,
						SubscriptionPeriodStart: one.CurrentPeriodStart,
						SubscriptionPeriodEnd:   one.CurrentPeriodEnd,
					})
				}
			}
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
		UserSubPlanCachedList(ctx, merchantId, userId, true)
	}()
}
