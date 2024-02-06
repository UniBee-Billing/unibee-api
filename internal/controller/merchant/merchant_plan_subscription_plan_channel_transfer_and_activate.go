package merchant

import (
	"context"
	_plan "unibee-api/api/merchant/plan"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/plan/service"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
	"strconv"
	"strings"
)

func (c *ControllerPlan) SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *_plan.SubscriptionPlanChannelTransferAndActivateReq) (res *_plan.SubscriptionPlanChannelTransferAndActivateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	utility.Assert(req.PlanId > 0, "plan should > 0")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	//多个渠道Plan 创建并激活
	list := query.GetListSubscriptionTypeGateways(ctx) // todo mark 需改造成获取 merchantId 相关的 Gateway
	utility.Assert(len(list) > 0, "no gateway found, need at least one")
	for _, gateway := range list {
		err = service.SubscriptionPlanChannelTransferAndActivate(ctx, req.PlanId, int64(gateway.Id))
		if err != nil {
			return nil, err
		}
	}

	//发布 Plan
	err = service.SubscriptionPlanActivate(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}

	if len(plan.BindingAddonIds) > 0 {
		//addon 检查
		var addonIds []int64
		var addonIdsList []int64
		if len(plan.BindingAddonIds) > 0 {
			//初始化
			strList := strings.Split(plan.BindingAddonIds, ",")
			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
				if err != nil {
					return nil, err
				}
				addonIdsList = append(addonIdsList, num) // 添加到整数列表中
				addonIds = append(addonIds, num)         // 添加到整数列表中
			}
		}
		//检查 addonIds 类型
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, addonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			if addonPlan.Status != consts.PlanStatusActive {
				//发布 addonPlan
				for _, gateway := range list {
					err = service.SubscriptionPlanChannelTransferAndActivate(ctx, int64(addonPlan.Id), int64(gateway.Id))
					if err != nil {
						return nil, err
					}
				}

				//发布 Plan
				err = service.SubscriptionPlanActivate(ctx, int64(addonPlan.Id))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &_plan.SubscriptionPlanChannelTransferAndActivateRes{}, nil
}
