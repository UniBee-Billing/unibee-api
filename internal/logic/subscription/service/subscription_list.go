package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionListInternalReq struct {
	MerchantId int64  `p:"merchantId" dc:"MerchantId"`
	UserId     int64  `p:"userId"  dc:"UserId" `
	Status     int    `p:"status" dc:"Default All，,Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page" d:"0"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count" d:"20" dc:"Count Of Page" `
}

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*subscription.SubscriptionDetailRes, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")

	if one.Type == consts.SubTypeDefault {
		go func() {
			defer func() {
				if exception := recover(); exception != nil {
					var err error
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						err = v
					} else {
						err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
					}
					g.Log().Errorf(context.Background(), "SubscriptionDetail Background panic error:%s\n", err.Error())
					return
				}
			}()
			backgroundCtx := context.Background()
			plan := query.GetPlanById(backgroundCtx, one.PlanId)
			utility.Assert(plan != nil, "invalid planId")
			utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
			planChannel := query.GetPlanChannel(backgroundCtx, one.PlanId, one.ChannelId)
			details, err := out.GetPayChannelServiceProvider(backgroundCtx, one.ChannelId).DoRemoteChannelSubscriptionDetails(backgroundCtx, plan, planChannel, one)
			if err == nil {
				err := handler.UpdateSubWithChannelDetailBack(backgroundCtx, one, details)
				if err != nil {
					fmt.Printf("SubscriptionDetail Background Fetch error%s", err)
					return
				}
			}
		}()
	}
	//删减返回值
	{
		one.Data = ""
		one.ResponseData = ""
	}
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	if user != nil {
		user.Password = ""
	}
	return &subscription.SubscriptionDetailRes{
		User:                                user,
		Subscription:                        one,
		Channel:                             ConvertChannelToRo(query.GetPayChannelById(ctx, one.ChannelId)),
		Plan:                                query.GetPlanById(ctx, one.PlanId),
		Addons:                              query.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UnfinishedSubscriptionPendingUpdate: GetUnfinishedSubscriptionPendingUpdateDetailByUpdateSubscriptionId(ctx, one.PendingUpdateId),
	}, nil
}

func SubscriptionList(ctx context.Context, req *SubscriptionListInternalReq) (list []*ro.SubscriptionDetailRo) {
	var mainList []*entity.Subscription
	if req.Count <= 0 {
		req.Count = 10 //每页数量Default 10
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
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, req.MerchantId).
		Where(dao.Subscription.Columns().UserId, req.UserId).
		Where(dao.Subscription.Columns().Status, req.Status).
		Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalPlanIds []int64
	for _, sub := range mainList {
		totalPlanIds = append(totalPlanIds, sub.PlanId)
		var addonParams []*ro.SubscriptionPlanAddonParamRo
		if len(sub.AddonData) > 0 {
			err := utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
			if err == nil {
				for _, s := range addonParams {
					totalPlanIds = append(totalPlanIds, s.AddonPlanId) // 添加到整数列表中
				}
			}
		}
		{
			sub.Data = ""
			sub.ResponseData = ""
		}
		user := query.GetUserAccountById(ctx, uint64(sub.UserId))
		if user != nil {
			user.Password = ""
		}
		list = append(list, &ro.SubscriptionDetailRo{
			User:         user,
			Subscription: sub,
			Channel:      ConvertChannelToRo(query.GetPayChannelById(ctx, sub.ChannelId)),
			Plan:         nil,
			Addons:       nil,
			AddonParams:  addonParams,
		})
	}
	if len(totalPlanIds) > 0 {
		//查询所有 Plan
		var allPlanList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalPlanIds).Scan(&allPlanList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allPlanList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, subRo := range list {
				subRo.Plan = mapPlans[subRo.Subscription.PlanId]
				if len(subRo.AddonParams) > 0 {
					for _, param := range subRo.AddonParams {
						if mapPlans[param.AddonPlanId] != nil {
							subRo.Addons = append(subRo.Addons, &ro.SubscriptionPlanAddonRo{
								Quantity:  param.Quantity,
								AddonPlan: mapPlans[param.AddonPlanId],
							})
						}
					}
				}
			}
		}
	}
	return list
}
