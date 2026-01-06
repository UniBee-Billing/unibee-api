package period

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetPlanById(ctx context.Context, id uint64) (one *entity.Plan) {
	if id <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPeriodEndFromStart(ctx context.Context, start int64, billingCycleAnchor int64, planId uint64) int64 {
	if billingCycleAnchor == 0 {
		billingCycleAnchor = start
	}
	plan := GetPlanById(ctx, planId)
	//utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	if plan == nil {
		g.Log().Errorf(ctx, "GetPeriodEndFromStart planId %d not found", planId)
		return start
	}
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
	plan := GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		return end - 60*60 // one hour
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		return end - 24*60*60 // 24h
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		return end - 3*24*60*60 // 3 day
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		return end - 15*24*60*60 // 15 day
	}
	return end - 30*60 // half an hour
}

func GetDunningTimeCap(ctx context.Context, planId uint64) int64 {
	plan := GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "GetPeriod Plan Not Found")
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "day") == 0 {
		return 60 * 60 // one hour
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "week") == 0 {
		return 24 * 60 * 60 // 24h
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "month") == 0 {
		return 3 * 24 * 60 * 60 // 3 day
	} else if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		return 15 * 24 * 60 * 60 // 15 day
	}
	return 30 * 60 // half an hour
}
