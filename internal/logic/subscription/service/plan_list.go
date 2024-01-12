package service

import (
	"context"
	"fmt"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	ro2 "go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

type SubscriptionPlanListInternalReq struct {
	MerchantId int64  `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	Type       int    `p:"type"  d:"1"  dc:"不填查询所有类型，,1-main plan，2-addon plan" `
	Status     int    `p:"status" dc:"不填查询所有状态，,状态，1-编辑中，2-活跃，3-非活跃，4-过期" `
	Currency   string `p:"currency" d:"usd"  dc:"订阅计划货币"  `
	SortField  string `p:"sortField" dc:"排序字段，gmt_create|gmt_modify，默认 gmt_modify" `
	SortType   string `p:"sortType" dc:"排序类型，asc|desc，默认 desc" `
	Page       int    `p:"page" d:"0"  dc:"分页页码,0开始" `
	Count      int    `p:"count" d:"20"  dc:"订阅计划货币" dc:"每页数量" `
}

func SubscriptionPlanDetail(ctx context.Context, planId int64) (*plan.SubscriptionPlanDetailRes, error) {
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found")
	return &plan.SubscriptionPlanDetailRes{
		Plan: &ro.PlanDetailRo{
			Plan:     one,
			Channels: query.GetListActiveOutChannelRos(ctx, planId),
			Addons:   query.GetPlanBindingAddonsByPlanId(ctx, planId),
		},
	}, nil
}

func SubscriptionPlanList(ctx context.Context, req *SubscriptionPlanListInternalReq) (list []*ro.PlanDetailRo) {
	var mainList []*entity.SubscriptionPlan
	if req.Count <= 0 {
		req.Count = 10 //每页数量默认 10
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
	err := dao.SubscriptionPlan.Ctx(ctx).
		Where(dao.SubscriptionPlan.Columns().MerchantId, req.MerchantId).
		Where(dao.SubscriptionPlan.Columns().Type, req.Type).
		Where(dao.SubscriptionPlan.Columns().Status, req.Status).
		Where(dao.SubscriptionPlan.Columns().Currency, strings.ToLower(req.Currency)).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalAddonIds []int64
	var totalPlanIds []uint64
	for _, p := range mainList {
		totalPlanIds = append(totalPlanIds, p.Id)
		if p.Type != 1 {
			//非主 Plan 不查询 addons
			list = append(list, &ro.PlanDetailRo{
				Plan:     p,
				Channels: []*ro2.OutChannelRo{},
				Addons:   nil,
				AddonIds: nil,
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
		list = append(list, &ro.PlanDetailRo{
			Plan:     p,
			Channels: []*ro2.OutChannelRo{},
			Addons:   nil,
			AddonIds: addonIds,
		})
	}
	if len(totalAddonIds) > 0 {
		//主 Plan 查询 addons
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).Scan(&allAddonList)
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
							planRo.Addons = append(planRo.Addons, mapPlans[id])
						}
					}
				}
			}
		}
	}
	//添加 Channel 信息
	var allPlanChannelList []*entity.SubscriptionPlanChannel
	err = dao.SubscriptionPlanChannel.Ctx(ctx).WhereIn(dao.SubscriptionPlanChannel.Columns().PlanId, totalPlanIds).Scan(&allPlanChannelList)
	if err == nil {
		for _, planChannel := range allPlanChannelList {
			for _, planRo := range list {
				if int64(planRo.Plan.Id) == planChannel.PlanId && planChannel.Status == consts.PlanChannelStatusActive {
					outChannel := query.GetPayChannelById(ctx, planChannel.ChannelId)
					planRo.Channels = append(planRo.Channels, &ro2.OutChannelRo{
						ChannelId:   outChannel.Id,
						ChannelName: outChannel.Name,
					})
				}
			}
		}
	}
	return list
}
