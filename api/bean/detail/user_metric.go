package detail

import "unibee/api/bean"

type UserMerchantMetricLimitStat struct {
	MetricLimit      *PlanMetricLimitDetail `json:"metricLimit" dc:"MetricLimit"`
	CurrentUsedValue int64                  `json:"CurrentUsedValue" dc:"CurrentUsedValue"`
}

type UserMerchantMetricChargeStat struct {
	MetricId          uint64                              `json:"metricId" dc:"MetricId" v:"required"`
	Metric            *bean.MerchantMetric                `json:"merchantMetric"    description:"MerchantMetric"`
	CurrentUsedValue  int64                               `json:"CurrentUsedValue" dc:"CurrentUsedValue"`
	MaxEventId        uint64                              `json:"maxEventId"`
	MinEventId        uint64                              `json:"minEventId"`
	ChargePricing     *bean.PlanMetricMeteredChargeParam  `json:"chargePricing" dc:"ChargePricing"`
	TotalChargeAmount int64                               `json:"totalChargeAmount" dc:"TotalChargeAmount"`
	GraduatedStep     *bean.MetricPlanChargeGraduatedStep `json:"graduatedStep" dc:"GraduatedStep"`
}
