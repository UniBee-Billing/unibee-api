package bean

import entity "unibee/internal/model/entity/default"

type MerchantMetric struct {
	Id                  uint64 `json:"id"            description:"id"`                                                                                // id
	MerchantId          uint64 `json:"merchantId"          description:"merchantId"`                                                                  // merchantId
	Code                string `json:"code"                description:"code"`                                                                        // code
	MetricName          string `json:"metricName"          description:"metric name"`                                                                 // metric name
	MetricDescription   string `json:"metricDescription"   description:"metric description"`                                                          // metric description
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int    `json:"aggregationType"     description:"1-count，2-count unique, 3-latest, 4-max, 5-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	UpdateTime          int64  `json:"gmtModify"     description:"update time"`     // update time
	CreateTime          int64  `json:"createTime"    description:"create utc time"` // create utc time
}

func SimplifyMerchantMetric(one *entity.MerchantMetric) *MerchantMetric {
	if one == nil {
		return nil
	}
	var updateTime int64
	if one.GmtModify != nil {
		updateTime = one.GmtModify.Timestamp()
	}
	return &MerchantMetric{
		Id:                  one.Id,
		MerchantId:          one.MerchantId,
		Code:                one.Code,
		MetricName:          one.MetricName,
		MetricDescription:   one.MetricDescription,
		Type:                one.Type,
		AggregationType:     one.AggregationType,
		AggregationProperty: one.AggregationProperty,
		UpdateTime:          updateTime,
		CreateTime:          one.CreateTime,
	}
}
