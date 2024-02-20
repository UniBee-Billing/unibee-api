package user_sub_plan

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

const (
	UserSubPlanCacheKeyPrefix = "UserSubPlanCacheKeyPrefix_"
	UserSubPlanCacheKeyExpire = 24 * 60 * 60
)

type UserSubPlan struct {
	MerchantId int64
	UserId     int64
	PlanId     int64
	PlanType   int
	Quantity   int64
}

func UserSubPlanCacheList(ctx context.Context, merchantId int64, userId int64, reloadCache bool) []*UserSubPlan {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(userId > 0, "invalid userId")
	var list = make([]*UserSubPlan, 0)
	cacheKey := fmt.Sprintf("%s_%d_%d", UserSubPlanCacheKeyPrefix, merchantId, userId)
	if !reloadCache {
		get, _ := g.Redis().Get(ctx, cacheKey)
		value := get.String()
		if len(value) > 0 {
			_ = utility.UnmarshalFromJsonString(value, &list)
			if len(list) > 0 {
				return list
			}
		}
	}
	if merchantId > 0 {
		var entities []*entity.Subscription
		var status = []int{consts.SubStatusActive, consts.SubStatusIncomplete, consts.SubStatusPendingInActive}
		err := dao.MerchantMetric.Ctx(ctx).
			Where(entity.Subscription{MerchantId: merchantId}).
			Where(entity.Subscription{UserId: userId}).
			WhereIn(dao.Subscription.Columns().Status, status).
			Where(entity.Subscription{IsDeleted: 0}).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				plan := query.GetPlanById(ctx, one.PlanId)
				if plan != nil {
					list = append(list, &UserSubPlan{
						MerchantId: one.MerchantId,
						UserId:     userId,
						PlanId:     one.PlanId,
						PlanType:   plan.Type,
						Quantity:   one.Quantity,
					})
				}
				//append addons
				addons := query.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData)
				for _, addon := range addons {
					list = append(list, &UserSubPlan{
						MerchantId: one.MerchantId,
						UserId:     userId,
						PlanId:     int64(addon.AddonPlan.Id),
						PlanType:   addon.AddonPlan.Type,
						Quantity:   addon.Quantity,
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

func ReloadUserSubPlanCacheListBackground(merchantId int64, userId int64) {
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
		UserSubPlanCacheList(ctx, merchantId, userId, true)
	}()
}
