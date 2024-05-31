package subscription

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consts"
	"unibee/internal/query"
	"unibee/utility"
)

func GetPeriodEndFromStart(ctx context.Context, start int64, billingCycleAnchor int64, planId uint64) int64 {
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "Plan Not Active")
	var periodEnd = gtime.NewFromTimeStamp(start)
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		periodEnd = periodEnd.AddDate(0, 0, plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		periodEnd = periodEnd.AddDate(0, 0, 7*plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		//periodEnd = periodEnd.AddDate(0, plan.IntervalCount, 0)
		periodEnd = periodEnd.AddDate(0, plan.IntervalCount, -periodEnd.Day()+1)
		periodEnd = periodEnd.AddDate(0, 0, utility.MinInt(gtime.NewFromTimeStamp(billingCycleAnchor).Day(), periodEnd.EndOfMonth().Day())-1)
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		//periodEnd = periodEnd.AddDate(plan.IntervalCount, 0, 0)
		periodEnd = periodEnd.AddDate(plan.IntervalCount, 0, -periodEnd.Day()+1)
		periodEnd = periodEnd.AddDate(0, 0, utility.MinInt(gtime.NewFromTimeStamp(billingCycleAnchor).Day(), periodEnd.EndOfMonth().Day())-1)
	}
	return periodEnd.Timestamp()
}

func GetDunningTimeFromEnd(ctx context.Context, end int64, planId uint64) int64 {
	if end == 0 {
		return 0
	}
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "Plan Not Active")
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		return end - 60*60 // one hour
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		return end - 24*60*60 // 24h
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		return end - 3*24*60*60 // 3 day
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		return end - 15*24*60*60 // 15 day
	}
	return end - 30*60 // half hour
}

func GetDunningTimeCap(ctx context.Context, planId uint64) int64 {
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "Plan Not Active")
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		return 60 * 60 // one hour
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		return 24 * 60 * 60 // 24h
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		return 3 * 24 * 60 * 60 // 3 day
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		return 15 * 24 * 60 * 60 // 15 day
	}
	return 30 * 60 // half hour
}
