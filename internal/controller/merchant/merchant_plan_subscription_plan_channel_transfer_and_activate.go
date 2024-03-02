package merchant

import (
	"context"
	"strconv"
	"strings"
	_plan "unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPlan) SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *_plan.SubscriptionPlanChannelTransferAndActivateReq) (res *_plan.SubscriptionPlanChannelTransferAndActivateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	utility.Assert(req.PlanId > 0, "plan should > 0")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	service.PlanOrAddonIntervalVerify(ctx, req.PlanId)

	//Activate Plan
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
				service.PlanOrAddonIntervalVerify(ctx, addonPlan.Id)

				//Activate Plan
				err = service.SubscriptionPlanActivate(ctx, addonPlan.Id)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &_plan.SubscriptionPlanChannelTransferAndActivateRes{}, nil
}
