package bean

type BulkMetricLimitPlanBindingParam struct {
	MetricId    uint64 `json:"metricId" dc:"MetricId" v:"required"`
	MetricLimit uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type MerchantMetricPlanLimit struct {
	Id          uint64          `json:"id"            description:"id"`                   // id
	MerchantId  uint64          `json:"merchantId"          description:"merchantId"`     // merchantId
	MetricId    uint64          `json:"metricId"    description:"metricId"`               // metricId
	Metric      *MerchantMetric `json:"merchantMetricVo"    description:"MerchantMetric"` // metricId
	PlanId      uint64          `json:"planId"      description:"plan_id"`                // plan_id
	MetricLimit uint64          `json:"metricLimit" description:"plan metric limit"`      // plan metric limit
	UpdateTime  int64           `json:"gmtModify"     description:"update time"`          // update time
	CreateTime  int64           `json:"createTime"    description:"create utc time"`      // create utc time
}

type PlanMetricLimitDetail struct {
	MerchantId          uint64
	UserId              uint64
	MetricId            uint64
	Code                string `json:"code"                description:"code"`                                                                        // code
	MetricName          string `json:"metricName"          description:"metric name"`                                                                 // metric name
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int    `json:"aggregationType"     description:"0-count，1-count unique, 2-latest, 3-max, 4-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	TotalLimit          uint64
	PlanLimits          []*MerchantMetricPlanLimit
}
