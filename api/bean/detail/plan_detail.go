package detail

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/query"
)

type PlanDetail struct {
	Product               *bean.Product                     `json:"product" dc:"Product"`
	Plan                  *bean.Plan                        `json:"plan" dc:"Plan"`
	Addons                []*bean.Plan                      `json:"addons" dc:"Addons"`
	AddonIds              []int64                           `json:"addonIds" dc:"AddonIds"`
	OnetimeAddons         []*bean.Plan                      `json:"onetimeAddons" dc:"OneTimeAddons"`
	OnetimeAddonIds       []int64                           `json:"onetimeAddonIds" dc:"OneTimeAddonIds"`
	MetricPlanLimits      []*MerchantMetricPlanLimitDetail  `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	MetricMeteredCharge   []*MerchantMetricPlanChargeDetail `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*MerchantMetricPlanChargeDetail `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
}

type MerchantMetricPlanChargeDetail struct {
	MetricId           uint64                                `json:"metricId" dc:"MetricId" v:"required"`
	Metric             *bean.MerchantMetric                  `json:"merchantMetric"    description:"MerchantMetric"` // metricId
	ChargeType         int                                   `json:"chargeType" dc:"ChargeType,0-standard pricing 1-graduated pricing"`
	StandardAmount     int64                                 `json:"standardAmount" dc:"StandardAmount, used for standard pricing"`
	StandardStartValue int64                                 `json:"standardStartValue" dc:"StandardStartValue, used for standard pricing"`
	GraduatedAmounts   []*bean.MetricPlanChargeGraduatedStep `json:"graduatedAmounts" dc:"GraduatedAmounts"`
}

func ConvertMetricPlanChargeDetailArrayFromParam(ctx context.Context, params []*bean.PlanMetricMeteredChargeParam) []*MerchantMetricPlanChargeDetail {
	var list = make([]*MerchantMetricPlanChargeDetail, 0)
	for _, param := range params {
		list = append(list, &MerchantMetricPlanChargeDetail{
			MetricId:           param.MetricId,
			Metric:             bean.SimplifyMerchantMetric(query.GetMerchantMetric(ctx, param.MetricId)),
			ChargeType:         param.ChargeType,
			StandardAmount:     param.StandardAmount,
			StandardStartValue: param.StandardStartValue,
			GraduatedAmounts:   param.GraduatedAmounts,
		})
	}
	return list
}

type MerchantMetricPlanLimitDetail struct {
	Id          uint64               `json:"id"            description:"id"`                 // id
	MerchantId  uint64               `json:"merchantId"          description:"merchantId"`   // merchantId
	MetricId    uint64               `json:"metricId"    description:"metricId"`             // metricId
	Metric      *bean.MerchantMetric `json:"merchantMetric"    description:"MerchantMetric"` // metricId
	PlanId      uint64               `json:"planId"      description:"plan_id"`              // plan_id
	MetricLimit uint64               `json:"metricLimit" description:"plan metric limit"`    // plan metric limit
	UpdateTime  int64                `json:"gmtModify"     description:"update time"`        // update time
	CreateTime  int64                `json:"createTime"    description:"create utc time"`    // create utc time
}

type PlanMetricLimitDetail struct {
	MerchantId          uint64
	UserId              uint64
	MetricId            uint64
	Code                string `json:"code"                description:"code"`        // code
	MetricName          string `json:"metricName"          description:"metric name"` // metric name
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered,3-charge_recurring"`
	AggregationType     int    `json:"aggregationType"     description:"0-count，1-count unique, 2-latest, 3-max, 4-sum"` // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	TotalLimit          uint64
	PlanLimits          []*MerchantMetricPlanLimitDetail
}
