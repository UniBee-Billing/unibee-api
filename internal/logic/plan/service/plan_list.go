package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unibee/api/merchant/plan"
	dao "unibee/internal/dao/oversea_pay"
	ro2 "unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/metric"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPlanListInternalReq struct {
	MerchantId    uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	Type          []int  `p:"type"  d:"1"  dc:"Default All，,1-main plan，2-addon plan" `
	Status        []int  `p:"status" dc:"Default All，,Status，1-Editing，2-Active，3-NonActive，4-Expired" `
	PublishStatus int    `p:"publishStatus" dc:"Default All，,Status，1-UnPublished，2-Published" `
	Currency      string `p:"currency" dc:"Currency"  `
	SortField     string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `p:"page" dc:"Page, Start WIth 0" `
	Count         int    `p:"count" dc:"Count Of Page" `
}

func SubscriptionPlanDetail(ctx context.Context, planId uint64) (*plan.SubscriptionPlanDetailRes, error) {
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found")
	var addonIds []int64
	if len(one.BindingAddonIds) > 0 {
		//初始化
		strList := strings.Split(one.BindingAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
			} else {
				addonIds = append(addonIds, num) // 添加到整数列表中
			}
		}
	}
	return &plan.SubscriptionPlanDetailRes{
		Plan: &ro2.PlanDetailRo{
			Plan:             ro2.SimplifyPlan(one),
			MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, one.Id, false),
			Addons:           ro2.SimplifyPlanList(query.GetPlanBindingAddonsByPlanId(ctx, planId)),
			AddonIds:         addonIds,
		},
	}, nil
}

func SubscriptionPlanList(ctx context.Context, req *SubscriptionPlanListInternalReq) (list []*ro2.PlanDetailRo) {
	var mainList []*entity.SubscriptionPlan
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
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
	q := dao.SubscriptionPlan.Ctx(ctx).
		Where(dao.SubscriptionPlan.Columns().MerchantId, req.MerchantId)
	if len(req.Type) > 0 {
		q = q.Where(dao.SubscriptionPlan.Columns().Type, req.Type)
	}
	if len(req.Status) > 0 {
		q = q.Where(dao.SubscriptionPlan.Columns().Status, req.Status)
	}
	err := q.Where(dao.SubscriptionPlan.Columns().PublishStatus, req.PublishStatus).
		Where(dao.SubscriptionPlan.Columns().Currency, strings.ToLower(req.Currency)).
		WhereIn(dao.SubscriptionPlan.Columns().IsDeleted, []int{0}).
		OmitEmpty().
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalAddonIds []int64
	var totalPlanIds []uint64
	for _, p := range mainList {
		totalPlanIds = append(totalPlanIds, p.Id)
		if p.Type != 1 {
			//非主 Plan 不查询 addons
			list = append(list, &ro2.PlanDetailRo{
				Plan:             ro2.SimplifyPlan(p),
				MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
				Addons:           nil,
				AddonIds:         nil,
			})
			continue
		}
		var addonIds []int64
		if len(p.BindingAddonIds) > 0 {
			//初始化
			strList := strings.Split(p.BindingAddonIds, ",")

			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
				if err != nil {
					fmt.Println("Internal Error converting string to int:", err)
				} else {
					totalAddonIds = append(totalAddonIds, num) // 添加到整数列表中
					addonIds = append(addonIds, num)           // 添加到整数列表中
				}
			}
		}
		list = append(list, &ro2.PlanDetailRo{
			Plan:             ro2.SimplifyPlan(p),
			MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
			Addons:           nil,
			AddonIds:         addonIds,
		})
	}
	if len(totalAddonIds) > 0 {
		//主 Plan 查询 addons
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, planRo := range list {
				if len(planRo.AddonIds) > 0 {
					for _, id := range planRo.AddonIds {
						if mapPlans[id] != nil {
							planRo.Addons = append(planRo.Addons, ro2.SimplifyPlan(mapPlans[id]))
						}
					}
				}
			}
		}
	}
	return list
}
