package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"
)

func PlanActivate(ctx context.Context, planId uint64) error {
	utility.Assert(planId > 0, "invalid planId")
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found, invalid planId")
	PlanOrAddonIntervalVerify(ctx, planId)
	if one.Status == consts.PlanStatusActive {
		return nil
	}
	update, err := dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().Status:    consts.PlanStatusActive,
		dao.Plan.Columns().IsDeleted: 0,
		dao.Plan.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Plan.Columns().Id, planId).OmitNil().Update()
	if err != nil {
		return err
	}
	affected, err := update.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return gerror.New("internal err, publish count != 1")
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "Activate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	return nil
}

func PlanOrAddonIntervalVerify(ctx context.Context, planId uint64) {
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "plan not found")
	if plan.Type != consts.PlanTypeOnetimeAddon {
		intervals := []string{"day", "month", "year", "week"}
		utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(plan.IntervalUnit)), "IntervalUnit Error，Must One Of day｜month｜year｜week")
		if strings.ToLower(plan.IntervalUnit) == "day" {
			utility.Assert(plan.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
		} else if strings.ToLower(plan.IntervalUnit) == "month" {
			utility.Assert(plan.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
		} else if strings.ToLower(plan.IntervalUnit) == "year" {
			utility.Assert(plan.IntervalCount <= 1, "IntervalCount Must Lower Then 2 While IntervalUnit is year")
		} else if strings.ToLower(plan.IntervalUnit) == "week" {
			utility.Assert(plan.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
		}
	}
}
