package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/gateway"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

func SubscriptionPlanChannelActivate(ctx context.Context, planId int64, channelId int64) (err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	utility.Assert(planId > 0, "invalid planId")
	utility.Assert(channelId > 0, "invalid channelId")
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetPlanChannel(ctx, planId, channelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	err = gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanActive(ctx, plan, planChannel)
	if err != nil {
		return
	}
	update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlanChannel.Columns().Status: consts.PlanChannelStatusActive,
		//dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: consts.PlanChannelStatusActive,// todo mark
		dao.SubscriptionPlanChannel.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).Update()
	if err != nil {
		return err
	}
	// todo mark update 值没变化会报错
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionPlanChannelActivate update err:%s", update)
	}
	return
}

func SubscriptionPlanChannelDeactivate(ctx context.Context, planId int64, channelId int64) (err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	utility.Assert(planId > 0, "invalid planId")
	utility.Assert(channelId > 0, "invalid channelId")
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetPlanChannel(ctx, planId, channelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	err = gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanDeactivate(ctx, plan, planChannel)
	if err != nil {
		return
	}
	update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlanChannel.Columns().Status: consts.PlanChannelStatusInActive,
		//dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: consts.PlanChannelStatusInActive,// todo mark
		dao.SubscriptionPlanChannel.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	// todo mark update 值没变化会报错
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionPlanChannelDeactivate update err:%s", update)
	}
	return
}

func SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (one *entity.SubscriptionPlan, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")

	utility.Assert(strings.HasPrefix(req.ImageUrl, "http"), "imageUrl should start with http")
	merchantInfo := query.GetMerchantInfoById(ctx, req.MerchantId)
	if len(req.ImageUrl) == 0 {
		req.ImageUrl = merchantInfo.CompanyLogo
	}
	if len(req.HomeUrl) == 0 {
		req.HomeUrl = merchantInfo.HomeUrl
	}
	utility.Assert(len(req.ImageUrl) > 0, "imageUrl should not be null")
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(req.Type == 1 || req.Type == 2, "type should be 1 or 2")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit 错误，day｜month｜year｜week\"")
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
		//检查 addonIds 类型
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, req.AddonIds).Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		}
	}

	one = &entity.SubscriptionPlan{
		CompanyId:                 merchantInfo.CompanyId,
		MerchantId:                req.MerchantId,
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
		ChannelProductName:        req.ProductName,
		ChannelProductDescription: req.ProductDescription,
		Status:                    consts.PlanStatusEditable,
	}
	result, err := dao.SubscriptionPlan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPlanCreate record insert failure %s`, err)
		one = nil
		return
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	//todo mark 是否直接发布
	return one, nil
}

func SubscriptionPlanEdit(ctx context.Context, req *v1.SubscriptionPlanEditReq) (one *entity.SubscriptionPlan, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")
	utility.Assert(len(req.ImageUrl) > 0, "imageUrl should not be null")
	utility.Assert(strings.HasPrefix(req.ImageUrl, "http"), "imageUrl should start with http")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit 错误，day｜month｜year｜week\"")
	if req.IntervalCount < 1 {
		req.IntervalCount = 1
	}
	utility.Assert(req.PlanId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", req.PlanId))
	utility.Assert(one.Status == consts.PlanStatusEditable, fmt.Sprintf("plan is not in edit status, id:%d", req.PlanId))

	if len(req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if len(req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}

	if len(req.AddonIds) > 0 {
		//检查 addonIds 类型
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, req.AddonIds).Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		}
	}

	_, err = dao.SubscriptionPlan.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlan.Columns().PlanName:                  req.PlanName,
		dao.SubscriptionPlan.Columns().Amount:                    req.Amount,
		dao.SubscriptionPlan.Columns().Currency:                  strings.ToUpper(req.Currency),
		dao.SubscriptionPlan.Columns().IntervalUnit:              req.IntervalUnit,
		dao.SubscriptionPlan.Columns().IntervalCount:             req.IntervalCount,
		dao.SubscriptionPlan.Columns().Description:               req.Description,
		dao.SubscriptionPlan.Columns().ImageUrl:                  req.ImageUrl,
		dao.SubscriptionPlan.Columns().HomeUrl:                   req.HomeUrl,
		dao.SubscriptionPlan.Columns().BindingAddonIds:           intListToString(req.AddonIds),
		dao.SubscriptionPlan.Columns().ChannelProductName:        req.ProductName,
		dao.SubscriptionPlan.Columns().ChannelProductDescription: req.ProductDescription,
	}).Where(dao.SubscriptionPlan.Columns().Id, req.PlanId).OmitNil().Update()
	if err != nil {
		err = gerror.Newf(`SubscriptionPlanEdit record insert failure %s`, err)
		one = nil
		return
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
	one.ChannelProductName = req.ProductName
	one.ChannelProductDescription = req.ProductDescription

	return one, nil
}

func SubscriptionPlanAddonsBinding(ctx context.Context, req *v1.SubscriptionPlanAddonsBindingReq) (one *entity.SubscriptionPlan, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
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
	var allAddonList []*entity.SubscriptionPlan
	err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, req.AddonIds).Scan(&allAddonList)
	for _, addonPlan := range allAddonList {
		utility.Assert(addonPlan.Type == consts.PlanTypeAddon, fmt.Sprintf("plan not addon type, id:%d", addonPlan.Id))
		utility.Assert(addonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
		//addon 周期校验
		utility.Assert(addonPlan.IntervalUnit == one.IntervalUnit && addonPlan.IntervalCount == one.IntervalCount, fmt.Sprintf("addon not match plan's recycle interval, id:%d", addonPlan.Id))
	}

	if req.Action == 0 {
		//覆盖
		addonIdsList = req.AddonIds
	} else if req.Action == 1 {
		//添加
		utility.Assert(len(req.AddonIds) > 0, "action add, addonids is empty")
		addonIdsList = mergeArrays(addonIdsList, req.AddonIds) // 添加到整数列表中
	} else if req.Action == 2 {
		//删除
		utility.Assert(len(req.AddonIds) > 0, "action delete, addonids is empty")
		addonIdsList = removeArrays(addonIdsList, req.AddonIds) // 添加到整数列表中
	}

	utility.Assert(len(addonIdsList) <= 10, "addon too much, should <= 10")

	newIds := intListToString(addonIdsList)
	one.BindingAddonIds = newIds
	update, err := dao.SubscriptionPlan.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlan.Columns().BindingAddonIds: one.BindingAddonIds,
		dao.SubscriptionPlan.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.SubscriptionPlan.Columns().Id, one.Id).Update()
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
	// 使用 map 跟踪已存在的元素
	seen := make(map[int64]bool)
	var result []int64

	// 将 arr1 中的元素添加到结果中
	for _, num := range arr1 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	// 将 arr2 中不重复的元素添加到结果中
	for _, num := range arr2 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

func removeArrays(arr, toRemove []int64) []int64 {
	// 创建一个 map 来存储要删除的元素
	removeMap := make(map[int64]bool)
	for _, num := range toRemove {
		removeMap[num] = true
	}

	// 遍历数组并删除要删除的元素
	var result []int64
	for _, num := range arr {
		if !removeMap[num] {
			result = append(result, num)
		}
	}

	return result
}

// 将整数数组转换为逗号分隔的字符串
func intListToString(arr []int64) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strArr, ",")
}
