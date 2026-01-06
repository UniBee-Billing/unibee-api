package bean

import (
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type MerchantMetric struct {
	Id                  uint64                 `json:"id"            description:"id"`                       // id
	MerchantId          uint64                 `json:"merchantId"          description:"merchantId"`         // merchantId
	Code                string                 `json:"code"                description:"code"`               // code
	MetricName          string                 `json:"metricName"          description:"metric name"`        // metric name
	MetricDescription   string                 `json:"metricDescription"   description:"metric description"` // metric description
	Type                int                    `json:"type"                description:"1-limit_metered，2-charge_metered,3-charge_recurring"`
	AggregationType     int                    `json:"aggregationType"     description:"1-count，2-count unique, 3-latest, 4-max, 5-sum"` // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string                 `json:"aggregationProperty" description:"aggregation property"`
	UpdateTime          int64                  `json:"gmtModify"     description:"update time"`     // update time
	CreateTime          int64                  `json:"createTime"    description:"create utc time"` // create utc time
	MetaData            map[string]interface{} `json:"metaData"            description:"meta_data(json)"`
	Unit                string                 `json:"unit"                description:"unit"`
	Archived            bool                   `json:"archived"                description:"archived"`
}

func SimplifyMerchantMetric(one *entity.MerchantMetric) *MerchantMetric {
	if one == nil {
		return nil
	}
	var updateTime int64
	if one.GmtModify != nil {
		updateTime = one.GmtModify.Timestamp()
	}
	metadata := make(map[string]interface{}, 0)
	if len(one.MetaData) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetaData, &metadata)
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
		Unit:                one.Unit,
		MetaData:            metadata,
		Archived:            one.IsDeleted != 0,
	}
}
