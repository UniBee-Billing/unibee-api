package plan

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/metric"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type ListInternalReq struct {
	MerchantId    uint64  `json:"merchantId" dc:"MerchantId" v:"required"`
	PlanIds       []int64 `json:"planIds"  dc:"filter id list of plan, default all" `
	ProductIds    []int64 `json:"productIds"  dc:"filter id list of product, default all" `
	Type          []int   `json:"type" dc:"Default All，,1-main plan，2-recurring addon, 3-one time addon" `
	Status        []int   `json:"status" dc:"Default All，,Status，1-Editing，2-Active，3-NonActive，4-Expired" `
	PublishStatus int     `json:"publishStatus" dc:"Default All，,Status，1-UnPublished，2-Published" `
	Currency      string  `json:"currency" dc:"Currency"  `
	SearchKey     string  `json:"searchKey" dc:"Search Key, plan name or description"  `
	SortField     string  `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string  `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int     `json:"page" dc:"Page, Start With 0" `
	Count         int     `json:"count" dc:"Count Of Page" `
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
			Product:               bean.SimplifyProduct(query.GetProductById(ctx, uint64(one.ProductId), merchantId)),
			Plan:                  bean.SimplifyPlan(one),
			Addons:                bean.SimplifyPlanList(query.GetAddonsByIds(ctx, addonIds)),
			AddonIds:              addonIds,
			OnetimeAddons:         bean.SimplifyPlanList(query.GetAddonsByIds(ctx, oneTimeAddonIds)),
			OnetimeAddonIds:       oneTimeAddonIds,
			MetricPlanLimits:      metric.MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, one.Id, false),
			MetricMeteredCharge:   detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(one).MetricMeteredCharge),
			MetricRecurringCharge: detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(one).MetricRecurringCharge),
		},
	}, nil
}

func PlanList(ctx context.Context, req *ListInternalReq) (list []*detail.PlanDetail, total int) {
	var mainList []*entity.Plan
	if req.Count <= 0 {
		req.Count = 100
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("plan_name|gmt_create|gmt_modify", req.SortField), "sortField should one of plan_name|gmt_create|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.Plan.Ctx(ctx).
		Where(dao.Plan.Columns().MerchantId, req.MerchantId)
	if len(req.ProductIds) > 0 {
		if isInt64InArray(req.ProductIds, 0) {
			q = q.Where(q.Builder().WhereOrIn(dao.Plan.Columns().ProductId, req.ProductIds).WhereOrNull(dao.Plan.Columns().ProductId))
		} else {
			q = q.WhereIn(dao.Plan.Columns().ProductId, req.ProductIds)
		}
	}
	if len(req.PlanIds) > 0 {
		q = q.WhereIn(dao.Plan.Columns().Id, req.PlanIds)
	}
	if len(req.Type) > 0 {
		q = q.WhereIn(dao.Plan.Columns().Type, req.Type)
	}
	if len(req.Status) > 0 {
		q = q.WhereIn(dao.Plan.Columns().Status, req.Status)
	}
	if len(req.SearchKey) > 0 {
		q = q.Where(q.Builder().WhereOrLike(dao.Plan.Columns().PlanName, "%"+req.SearchKey+"%").
			WhereOrLike(dao.Plan.Columns().Description, "%"+req.SearchKey+"%"))
	}
	err := q.Where(dao.Plan.Columns().PublishStatus, req.PublishStatus).
		Where(dao.Plan.Columns().Currency, strings.ToUpper(req.Currency)).
		WhereLTE(dao.Plan.Columns().IsDeleted, 0).
		OmitEmpty().
		Order(fmt.Sprintf("is_deleted desc, %s, status asc", sortKey)).
		Limit(req.Page*req.Count, req.Count).
		ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, 0
	}
	var totalAddonIds []int64
	var totalOneTimeAddonIds []int64
	for _, p := range mainList {
		if p.Type != consts.PlanTypeMain {
			list = append(list, &detail.PlanDetail{
				Product:               bean.SimplifyProduct(query.GetProductById(ctx, uint64(p.ProductId), req.MerchantId)),
				Plan:                  bean.SimplifyPlan(p),
				MetricPlanLimits:      metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
				Addons:                nil,
				AddonIds:              nil,
				MetricMeteredCharge:   detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(p).MetricMeteredCharge),
				MetricRecurringCharge: detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(p).MetricRecurringCharge),
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
			Product:               bean.SimplifyProduct(query.GetProductById(ctx, uint64(p.ProductId), req.MerchantId)),
			Plan:                  bean.SimplifyPlan(p),
			MetricPlanLimits:      metric.MerchantMetricPlanLimitCachedList(ctx, p.MerchantId, p.Id, false),
			Addons:                nil,
			AddonIds:              addonIds,
			OnetimeAddons:         nil,
			OnetimeAddonIds:       oneTimeAddonIds,
			MetricMeteredCharge:   detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(p).MetricMeteredCharge),
			MetricRecurringCharge: detail.ConvertMetricPlanChargeDetailArrayFromParam(ctx, bean.ConvertMetricPlanBindingEntityFromPlan(p).MetricRecurringCharge),
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
	return list, total
}

func isInt64InArray(arr []int64, target int64) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}
