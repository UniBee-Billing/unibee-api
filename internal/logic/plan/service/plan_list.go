package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/metric"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPlanListInternalReq struct {
	MerchantId    uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	Type          []int  `json:"type" dc:"Default All，,1-main plan，2-addon plan" `
	Status        []int  `json:"status" dc:"Default All，,Status，1-Editing，2-Active，3-NonActive，4-Expired" `
	PublishStatus int    `json:"publishStatus" dc:"Default All，,Status，1-UnPublished，2-Published" `
	Currency      string `json:"currency" dc:"Currency"  `
	SortField     string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `json:"page" dc:"Page, Start With 0" `
	Count         int    `json:"count" dc:"Count Of Page" `
}

func PlanDetail(ctx context.Context, merchantId uint64, planId uint64) (*plan.DetailRes, error) {
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found")
	utility.Assert(one.MerchantId == merchantId, "wrong merchant account")
	var addonIds = make([]int64, 0)
	if len(one.BindingAddonIds) > 0 {
		strList := strings.Split(one.BindingAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
			} else {
				addonIds = append(addonIds, num)
			}
		}
	}
	var oneTimeAddonIds = make([]int64, 0)
	if len(one.BindingOnetimeAddonIds) > 0 {
		strList := strings.Split(one.BindingOnetimeAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
			} else {
				oneTimeAddonIds = append(oneTimeAddonIds, num)
			}
		}
	}
	return &plan.DetailRes{
		Plan: &detail.PlanDetail{
			Plan:             bean.SimplifyPlan(one),
			MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, one.Id, false),
			Addons:           bean.SimplifyPlanList(query.GetAddonsByIds(ctx, addonIds)),
			AddonIds:         addonIds,
			OnetimeAddons:    bean.SimplifyPlanList(query.GetAddonsByIds(ctx, oneTimeAddonIds)),
			OnetimeAddonIds:  oneTimeAddonIds,
		},
	}, nil
}

func PlanList(ctx context.Context, req *SubscriptionPlanListInternalReq) (list []*detail.PlanDetail) {
	var mainList []*entity.Plan
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
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
	q := dao.Plan.Ctx(ctx).
		Where(dao.Plan.Columns().MerchantId, req.MerchantId)
	if len(req.Type) > 0 {
		q = q.Where(dao.Plan.Columns().Type, req.Type)
	}
	if len(req.Status) > 0 {
		q = q.Where(dao.Plan.Columns().Status, req.Status)
	}
	err := q.Where(dao.Plan.Columns().PublishStatus, req.PublishStatus).
		Where(dao.Plan.Columns().Currency, strings.ToLower(req.Currency)).
		WhereIn(dao.Plan.Columns().IsDeleted, []int{0}).
		OmitEmpty().
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalAddonIds []int64
	var totalOneTimeAddonIds []int64
	for _, p := range mainList {
		if p.Type != consts.PlanTypeMain {
			list = append(list, &detail.PlanDetail{
				Plan:             bean.SimplifyPlan(p),
				MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
				Addons:           nil,
				AddonIds:         nil,
			})
			continue
		}
		var addonIds []int64
		if len(p.BindingAddonIds) > 0 {
			strList := strings.Split(p.BindingAddonIds, ",")

			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					fmt.Println("Internal Error converting string to int:", err)
				} else {
					totalAddonIds = append(totalAddonIds, num)
					addonIds = append(addonIds, num)
				}
			}
		}
		var oneTimeAddonIds = make([]int64, 0)
		if len(p.BindingOnetimeAddonIds) > 0 {
			strList := strings.Split(p.BindingOnetimeAddonIds, ",")

			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					fmt.Println("Internal Error converting string to int:", err)
				} else {
					totalOneTimeAddonIds = append(totalOneTimeAddonIds, num)
					oneTimeAddonIds = append(oneTimeAddonIds, num)
				}
			}
		}
		list = append(list, &detail.PlanDetail{
			Plan:             bean.SimplifyPlan(p),
			MetricPlanLimits: metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
			Addons:           nil,
			AddonIds:         addonIds,
			OnetimeAddons:    nil,
			OnetimeAddonIds:  oneTimeAddonIds,
		})
	}
	if len(totalAddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			mapPlans := make(map[int64]*entity.Plan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, planRo := range list {
				if len(planRo.AddonIds) > 0 {
					for _, id := range planRo.AddonIds {
						if mapPlans[id] != nil {
							planRo.Addons = append(planRo.Addons, bean.SimplifyPlan(mapPlans[id]))
						}
					}
				}
			}
		}
	}
	if len(totalOneTimeAddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalOneTimeAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			mapPlans := make(map[int64]*entity.Plan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, planRo := range list {
				if len(planRo.OnetimeAddonIds) > 0 {
					for _, id := range planRo.OnetimeAddonIds {
						if mapPlans[id] != nil {
							planRo.OnetimeAddons = append(planRo.OnetimeAddons, bean.SimplifyPlan(mapPlans[id]))
						}
					}
				}
			}
		}
	}
	return list
}
