package bean

import entity "unibee/internal/model/entity/default"

type MerchantMetricEvent struct {
	Id                        uint64 `json:"id"                          description:"Id"`                                                                     // Id
	MerchantId                uint64 `json:"merchantId"                  description:"merchantId"`                                                             // merchantId
	MetricId                  uint64 `json:"metricId"                    description:"metric_id"`                                                              // metric_id
	ExternalEventId           string `json:"externalEventId"             description:"external_event_id, should be unique"`                                    // external_event_id, should be unique
	UserId                    int64  `json:"userId"                      description:"user_id"`                                                                // user_id
	AggregationPropertyInt    uint64 `json:"aggregationPropertyInt"      description:"aggregation property int, use for metric of max|sum type"`               // aggregation property int, use for metric of max|sum type
	AggregationPropertyString string `json:"aggregationPropertyString"   description:"aggregation property string, use for metric of count|count_unique type"` // aggregation property string, use for metric of count|count_unique type
	CreateTime                int64  `json:"createTime"                  description:"create utc time"`                                                        // create utc time
	AggregationPropertyData   string `json:"aggregationPropertyData"     description:"aggregation property data (Json)"`                                       // aggregation property data (Json)
	SubscriptionIds           string `json:"subscriptionIds"             description:""`                                                                       //
	SubscriptionPeriodStart   int64  `json:"subscriptionPeriodStart"     description:"matched subscription's current_period_start"`                            // matched subscription's current_period_start
	SubscriptionPeriodEnd     int64  `json:"subscriptionPeriodEnd"       description:"matched subscription's current_period_end"`                              // matched subscription's current_period_end
	MetricLimit               uint64 `json:"metricLimit"                 description:""`                                                                       //
	Used                      uint64 `json:"used"                        description:""`                                                                       //
}

func SimplifyMerchantMetricEvent(one *entity.MerchantMetricEvent) *MerchantMetricEvent {
	if one == nil {
		return nil
	}
	return &MerchantMetricEvent{
		Id:                        one.Id,
		MerchantId:                one.MerchantId,
		MetricId:                  one.MetricId,
		ExternalEventId:           one.ExternalEventId,
		UserId:                    one.UserId,
		AggregationPropertyInt:    one.AggregationPropertyInt,
		AggregationPropertyString: one.AggregationPropertyString,
		CreateTime:                one.CreateTime,
		AggregationPropertyData:   one.AggregationPropertyData,
		SubscriptionIds:           one.SubscriptionIds,
		SubscriptionPeriodStart:   one.SubscriptionPeriodStart,
		SubscriptionPeriodEnd:     one.SubscriptionPeriodEnd,
		MetricLimit:               one.MetricLimit,
		Used:                      one.Used,
	}
}
