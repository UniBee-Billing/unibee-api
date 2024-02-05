package service

import (
	"context"
	"fmt"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	ro2 "go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

type SubscriptionPlanListInternalReq struct {
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Type          int    `p:"type"  d:"1"  dc:"Default All，,1-main plan，2-addon plan" `
	Status        int    `p:"status" dc:"Default All，,Status，1-Editing，2-Active，3-NonActive，4-Expired" `
	PublishStatus int    `p:"publishStatus" dc:"Default All，,Status，1-UnPublished，2-Published" `
	Currency      string `p:"currency" dc:"Currency"  `
	SortField     string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `p:"page" d:"0"  dc:"Page, Start WIth 0" `
	Count         int    `p:"count" d:"20"  dc:"Count Of Page" `
}

func SubscriptionPlanDetail(ctx context.Context, planId int64) (*plan.SubscriptionPlanDetailRes, error) {
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found")
	return &plan.SubscriptionPlanDetailRes{
		Plan: &ro2.PlanDetailRo{
			Plan:     one,
			Channels: query.GetListActiveOutChannelRos(ctx, planId),
			Addons:   query.GetPlanBindingAddonsByPlanId(ctx, planId),
		},
	}, nil
}

func SubscriptionPlanList(ctx context.Context, req *SubscriptionPlanListInternalReq) (list []*ro2.PlanDetailRo) {
	var mainList []*entity.SubscriptionPlan
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
	err := dao.SubscriptionPlan.Ctx(ctx).
		Where(dao.SubscriptionPlan.Columns().MerchantId, req.MerchantId).
		Where(dao.SubscriptionPlan.Columns().Type, req.Type).
		Where(dao.SubscriptionPlan.Columns().Status, req.Status).
		Where(dao.SubscriptionPlan.Columns().PublishStatus, req.PublishStatus).
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
		list = append(list, &ro2.PlanDetailRo{
			Plan:     p,
			Channels: []*ro2.OutChannelRo{},
			Addons:   nil,
			AddonIds: addonIds,
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
							planRo.Addons = append(planRo.Addons, mapPlans[id])
						}
					}
				}
			}
		}
	}
	//添加 Channel 信息
	var allPlanChannelList []*entity.ChannelPlan
	err = dao.ChannelPlan.Ctx(ctx).WhereIn(dao.ChannelPlan.Columns().PlanId, totalPlanIds).OmitEmpty().Scan(&allPlanChannelList)
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
