package subscription

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

func GetPeriodEndFromStart(ctx context.Context, start int64, planId uint64) int64 {
	plan := query.GetPlanById(ctx, int64(planId))
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "Plan Not Active")
	var periodEnd = gtime.NewFromTimeStamp(start)
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		periodEnd = periodEnd.AddDate(0, 0, plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		periodEnd = periodEnd.AddDate(0, 0, 7*plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		periodEnd = periodEnd.AddDate(0, plan.IntervalCount, 0)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		periodEnd = periodEnd.AddDate(plan.IntervalCount, 0, 0)
	}
	return periodEnd.Timestamp()
}
