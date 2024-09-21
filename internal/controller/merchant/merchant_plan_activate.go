package merchant

import (
	"context"
	"strconv"
	"strings"
	_plan "unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	plan2 "unibee/internal/logic/plan"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPlan) Activate(ctx context.Context, req *_plan.ActivateReq) (res *_plan.ActivateRes, err error) {
	utility.Assert(req.PlanId > 0, "plan should > 0")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	plan2.PlanOrAddonIntervalVerify(ctx, req.PlanId)

	//Activate Plan
	err = plan2.PlanActivate(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}

	if len(plan.BindingAddonIds) > 0 {
		var addonIds []int64
		var addonIdsList []int64
		if len(plan.BindingAddonIds) > 0 {
			strList := strings.Split(plan.BindingAddonIds, ",")
			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return nil, err
				}
				addonIdsList = append(addonIdsList, num)
				addonIds = append(addonIds, num)
			}
		}
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, addonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			if addonPlan.Status != consts.PlanStatusActive {
				plan2.PlanOrAddonIntervalVerify(ctx, addonPlan.Id)
				err = plan2.PlanActivate(ctx, addonPlan.Id)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &_plan.ActivateRes{}, nil
}
