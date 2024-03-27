package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId"`
	UserId     int64  `json:"userId"  dc:"UserId" `
	Status     []int  `json:"status" dc:"Default All，,Status，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `json:"page" dc:"Page, Start WIth 0" `
	Count      int    `json:"count" dc:"Count Of Page" `
}

func SubscriptionDetail(ctx context.Context, subscriptionId string) (*detail.SubscriptionDetail, error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	{
		one.Data = ""
		one.ResponseData = ""
	}
	user := query.GetUserAccountById(ctx, one.UserId)
	var addonParams []*bean.PlanAddonParam
	if len(one.AddonData) > 0 {
		err := utility.UnmarshalFromJsonString(one.AddonData, &addonParams)
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionDetail parse addon param:%s", err.Error())
		}
	}
	return &detail.SubscriptionDetail{
		User:                                bean.SimplifyUserAccount(user),
		Subscription:                        bean.SimplifySubscription(one),
		Gateway:                             bean.SimplifyGateway(query.GetGatewayById(ctx, one.GatewayId)),
		Plan:                                bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		AddonParams:                         addonParams,
		Addons:                              addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		LatestInvoice:                       bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId)),
		UnfinishedSubscriptionPendingUpdate: GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId(ctx, one.PendingUpdateId),
	}, nil
}

func SubscriptionList(ctx context.Context, req *SubscriptionListInternalReq) (list []*detail.SubscriptionDetail) {
	var mainList []*entity.Subscription
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
	baseQuery := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, req.MerchantId).
		Where(dao.Subscription.Columns().UserId, req.UserId)
	if req.Status != nil && len(req.Status) > 0 {
		baseQuery = baseQuery.WhereIn(dao.Subscription.Columns().Status, req.Status)
	}
	err := baseQuery.Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	var totalPlanIds []uint64
	for _, sub := range mainList {
		totalPlanIds = append(totalPlanIds, sub.PlanId)
		var addonParams []*bean.PlanAddonParam
		if len(sub.AddonData) > 0 {
			err = utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
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
		user := query.GetUserAccountById(ctx, sub.UserId)
		if user != nil {
			user.Password = ""
		}
		list = append(list, &detail.SubscriptionDetail{
			User:          bean.SimplifyUserAccount(user),
			Subscription:  bean.SimplifySubscription(sub),
			Gateway:       bean.SimplifyGateway(query.GetGatewayById(ctx, sub.GatewayId)),
			Plan:          nil,
			Addons:        nil,
			LatestInvoice: bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)),
			AddonParams:   addonParams,
		})
	}
	if len(totalPlanIds) > 0 {
		var allPlanList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalPlanIds).OmitEmpty().Scan(&allPlanList)
		if err == nil {
			mapPlans := make(map[uint64]*entity.Plan)
			for _, pair := range allPlanList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, subRo := range list {
				subRo.Plan = bean.SimplifyPlan(mapPlans[subRo.Subscription.PlanId])
				if len(subRo.AddonParams) > 0 {
					for _, param := range subRo.AddonParams {
						if mapPlans[param.AddonPlanId] != nil {
							subRo.Addons = append(subRo.Addons, &bean.PlanAddonDetail{
								Quantity:  param.Quantity,
								AddonPlan: bean.SimplifyPlan(mapPlans[param.AddonPlanId]),
							})
						}
					}
				}
			}
		}
	}
	return list
}
