package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	v1 "unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/metric"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func SubscriptionPlanPublish(ctx context.Context, planId uint64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan.Status == consts.PlanStatusActive, "plan not activate")
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().PublishStatus: consts.PlanPublishStatusPublished,
		dao.Plan.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.Plan.Columns().Id, planId).Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionPlanUnPublish(ctx context.Context, planId uint64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan.Status == consts.PlanStatusActive, "plan not activate")
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().PublishStatus: consts.PlanPublishStatusUnPublished,
		dao.Plan.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.Plan.Columns().Id, planId).Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionPlanCreate(ctx context.Context, req *v1.NewReq) (one *entity.Plan, err error) {
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")
	utility.Assert(len(req.PlanName) > 0, "plan name should not blank")
	utility.Assert(len(req.Description) > 0, "description should not blank")

	//check metricLimitList
	if len(req.MetricLimits) > 0 {
		for _, ml := range req.MetricLimits {
			utility.Assert(ml.MetricId > 0, "invalid metricId")
			utility.Assert(ml.MetricLimit > 0, "invalid MetricLimit")
			me := query.GetMerchantMetric(ctx, ml.MetricId)
			utility.Assert(me != nil, "metric not found")
			utility.Assert(me.Type == metric.MetricTypeLimitMetered, "metric type invalid")
		}
	}

	utility.Assert(strings.HasPrefix(req.ImageUrl, "http"), "imageUrl should start with http")
	merchantInfo := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if len(req.ImageUrl) == 0 {
		req.ImageUrl = merchantInfo.CompanyLogo
	}
	if len(req.HomeUrl) == 0 {
		req.HomeUrl = merchantInfo.HomeUrl
	}
	utility.Assert(len(req.ImageUrl) > 0, "imageUrl should not be null")
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(req.Type == 1 || req.Type == 2, "type should be 1 or 2")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit Error， must one of day｜month｜year｜week\"")
	if strings.ToLower(req.IntervalUnit) == "day" {
		utility.Assert(req.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
	} else if strings.ToLower(req.IntervalUnit) == "month" {
		utility.Assert(req.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
	} else if strings.ToLower(req.IntervalUnit) == "year" {
		utility.Assert(req.IntervalCount <= 1, "IntervalCount Must Lower Then 2 While IntervalUnit is year")
	} else if strings.ToLower(req.IntervalUnit) == "week" {
		utility.Assert(req.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
	}
	if req.IntervalCount < 1 {
		req.IntervalCount = 1
	}

	if len(req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if len(req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}

	if len(req.AddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		}
	}

	one = &entity.Plan{
		CompanyId:                 merchantInfo.CompanyId,
		MerchantId:                _interface.GetMerchantId(ctx),
		PlanName:                  req.PlanName,
		Amount:                    req.Amount,
		Currency:                  strings.ToUpper(req.Currency),
		IntervalUnit:              strings.ToLower(req.IntervalUnit),
		IntervalCount:             req.IntervalCount,
		Type:                      req.Type,
		Description:               req.Description,
		ImageUrl:                  req.ImageUrl,
		HomeUrl:                   req.HomeUrl,
		BindingAddonIds:           intListToString(req.AddonIds),
		GatewayProductName:        req.ProductName,
		GatewayProductDescription: req.ProductDescription,
		Status:                    consts.PlanStatusEditable,
		CreateTime:                gtime.Now().Timestamp(),
	}
	result, err := dao.Plan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`SubscriptionPlanCreate record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	if len(req.MetricLimits) > 0 {
		err = metric.BulkMetricLimitPlanBindingReplace(ctx, one, req.MetricLimits)
		if err != nil {
			return nil, gerror.Newf(`BulkMetricLimitPlanBindingReplace %s`, err)
		}
	}

	return one, nil
}

func SubscriptionPlanEdit(ctx context.Context, req *v1.EditReq) (one *entity.Plan, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")
	utility.Assert(len(req.ImageUrl) > 0, "imageUrl should not be null")
	utility.Assert(len(req.PlanName) > 0, "plan name should not blank")
	utility.Assert(len(req.Description) > 0, "description should not blank")
	utility.Assert(strings.HasPrefix(req.ImageUrl, "http"), "imageUrl should start with http")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit Error， must one of day｜month｜year｜week\"")
	if strings.ToLower(req.IntervalUnit) == "day" {
		utility.Assert(req.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
	} else if strings.ToLower(req.IntervalUnit) == "month" {
		utility.Assert(req.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
	} else if strings.ToLower(req.IntervalUnit) == "year" {
		utility.Assert(req.IntervalCount <= 1, "IntervalCount Must Lower Then 2 While IntervalUnit is year")
	} else if strings.ToLower(req.IntervalUnit) == "week" {
		utility.Assert(req.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
	}
	if req.IntervalCount < 1 {
		req.IntervalCount = 1
	}
	utility.Assert(req.PlanId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", req.PlanId))
	utility.Assert(one.Status == consts.PlanStatusEditable, fmt.Sprintf("plan is not in edit status, id:%d", req.PlanId))

	//check metricLimitList
	if len(req.MetricLimits) > 0 {
		for _, ml := range req.MetricLimits {
			utility.Assert(ml.MetricId > 0, "invalid metricId")
			utility.Assert(ml.MetricLimit > 0, "invalid MetricLimit")
			me := query.GetMerchantMetric(ctx, ml.MetricId)
			utility.Assert(me != nil, "metric not found")
			utility.Assert(me.Type == metric.MetricTypeLimitMetered, "metric type invalid")
		}
	}

	if len(req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if len(req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}

	if len(req.AddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		}
	}

	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().PlanName:                  req.PlanName,
		dao.Plan.Columns().Amount:                    req.Amount,
		dao.Plan.Columns().Currency:                  strings.ToUpper(req.Currency),
		dao.Plan.Columns().IntervalUnit:              req.IntervalUnit,
		dao.Plan.Columns().IntervalCount:             req.IntervalCount,
		dao.Plan.Columns().Description:               req.Description,
		dao.Plan.Columns().ImageUrl:                  req.ImageUrl,
		dao.Plan.Columns().HomeUrl:                   req.HomeUrl,
		dao.Plan.Columns().BindingAddonIds:           intListToString(req.AddonIds),
		dao.Plan.Columns().GatewayProductName:        req.ProductName,
		dao.Plan.Columns().GatewayProductDescription: req.ProductDescription,
	}).Where(dao.Plan.Columns().Id, req.PlanId).OmitNil().Update()
	if err != nil {
		return nil, gerror.Newf(`SubscriptionPlanEdit record insert failure %s`, err)
	}

	one.PlanName = req.PlanName
	one.Amount = req.Amount
	one.Currency = strings.ToUpper(req.Currency)
	one.IntervalUnit = strings.ToLower(req.IntervalUnit)
	one.IntervalCount = req.IntervalCount
	one.Description = req.Description
	one.ImageUrl = req.ImageUrl
	one.HomeUrl = req.HomeUrl
	one.BindingAddonIds = intListToString(req.AddonIds)
	one.GatewayProductName = req.ProductName
	one.GatewayProductDescription = req.ProductDescription

	if len(req.MetricLimits) > 0 {
		err = metric.BulkMetricLimitPlanBindingReplace(ctx, one, req.MetricLimits)
		if err != nil {
			return nil, gerror.Newf(`BulkMetricLimitPlanBindingReplace %s`, err)
		}
	}

	return one, nil
}

func SubscriptionPlanDelete(ctx context.Context, planId uint64) (one *entity.Plan, err error) {
	utility.Assert(planId > 0, "planId invlaid")
	one = query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", planId))
	utility.Assert(one.Status == consts.PlanStatusEditable, fmt.Sprintf("plan is not in edit status, id:%d", planId))
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.Plan.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Plan.Columns().Id, one.Id).Update()
	if err != nil {
		return nil, err
	}
	one.IsDeleted = int(gtime.Now().Timestamp())
	return one, nil
}

func SubscriptionPlanAddonsBinding(ctx context.Context, req *v1.AddonsBindingReq) (one *entity.Plan, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Action >= 0 && req.Action <= 2, "action should 0-2")
	utility.Assert(req.PlanId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", req.PlanId))
	utility.Assert(one.Type == consts.PlanTypeMain, fmt.Sprintf("plan not type main, id:%d", req.PlanId))

	var addonIdsList []int64
	if len(one.BindingAddonIds) > 0 {
		//初始化
		strList := strings.Split(one.BindingAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
				return nil, err
			}
			addonIdsList = append(addonIdsList, num) // 添加到整数列表中
		}
	}
	//检查 addonIds 类型
	var allAddonList []*entity.Plan
	err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).OmitEmpty().Scan(&allAddonList)
	for _, addonPlan := range allAddonList {
		utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
		utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		//addon 周期校验
		utility.Assert(addonPlan.IntervalUnit == one.IntervalUnit && addonPlan.IntervalCount == one.IntervalCount, fmt.Sprintf("addon not match plan's recycle interval, id:%d", addonPlan.Id))
	}

	if req.Action == 0 {
		//覆盖
		addonIdsList = req.AddonIds
	} else if req.Action == 1 {
		//添加
		utility.Assert(len(req.AddonIds) > 0, "action add, addon ids is empty")
		addonIdsList = mergeArrays(addonIdsList, req.AddonIds)
	} else if req.Action == 2 {
		//删除
		utility.Assert(len(req.AddonIds) > 0, "action delete, addon ids is empty")
		addonIdsList = removeArrays(addonIdsList, req.AddonIds)
	}

	utility.Assert(len(addonIdsList) <= 10, "addon too much, should <= 10")

	newIds := intListToString(addonIdsList)
	one.BindingAddonIds = newIds
	update, err := dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().BindingAddonIds: one.BindingAddonIds,
		dao.Plan.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Plan.Columns().Id, one.Id).Update()
	if err != nil {
		return nil, err
	}
	affected, err := update.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, gerror.New("internal err, publish count != 1")
	}
	return one, nil
}

func mergeArrays(arr1, arr2 []int64) []int64 {
	seen := make(map[int64]bool)
	var result []int64
	for _, num := range arr1 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	for _, num := range arr2 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

func removeArrays(arr, toRemove []int64) []int64 {
	removeMap := make(map[int64]bool)
	for _, num := range toRemove {
		removeMap[num] = true
	}
	var result []int64
	for _, num := range arr {
		if !removeMap[num] {
			result = append(result, num)
		}
	}

	return result
}

func intListToString(arr []int64) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strArr, ",")
}
