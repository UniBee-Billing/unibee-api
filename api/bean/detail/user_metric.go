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

type UserMetric struct {
	IsPaid               bool                            `json:"isPaid" dc:"IsPaid"`
	Product              *bean.Product                   `json:"product" dc:"Product"`
	User                 *bean.UserAccount               `json:"user" dc:"user"`
	Subscription         *bean.Subscription              `json:"subscription" dc:"Subscription"`
	Plan                 *bean.Plan                      `json:"plan" dc:"Plan"`
	Addons               []*bean.PlanAddonDetail         `json:"addons" dc:"Addon"`
	LimitStats           []*UserMerchantMetricLimitStat  `json:"limitStats" dc:"LimitStats"`
	MeteredChargeStats   []*UserMerchantMetricChargeStat `json:"meteredChargeStats" dc:"MeteredChargeStats"`
	RecurringChargeStats []*UserMerchantMetricChargeStat `json:"recurringChargeStats" dc:"RecurringChargeStats"`
	Description          string                          `json:"description" dc:"description"`
}
