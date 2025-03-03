package detail

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMetricEventDetail struct {
	Id                        uint64                  `json:"id"                          description:"Id"`                                  // Id
	MerchantId                uint64                  `json:"merchantId"                  description:"merchantId"`                          // merchantId
	MetricId                  uint64                  `json:"metricId"                    description:"metric_id"`                           // metric_id
	Metric                    *bean.MerchantMetric    `json:"merchantMetric"    description:"MerchantMetric"`                                // metricId
	ExternalEventId           string                  `json:"externalEventId"             description:"external_event_id, should be unique"` // external_event_id, should be unique
	UserId                    int64                   `json:"userId"                      description:"user_id"`                             // user_id
	User                      *bean.UserAccount       `json:"user" dc:"user"`
	AggregationPropertyInt    uint64                  `json:"aggregationPropertyInt"      description:"aggregation property int, use for metric of max|sum type"`               // aggregation property int, use for metric of max|sum type
	AggregationPropertyString string                  `json:"aggregationPropertyString"   description:"aggregation property string, use for metric of count|count_unique type"` // aggregation property string, use for metric of count|count_unique type
	CreateTime                int64                   `json:"createTime"                  description:"create utc time"`                                                        // create utc time
	AggregationPropertyData   string                  `json:"aggregationPropertyData"     description:"aggregation property data (Json)"`                                       // aggregation property data (Json)
	SubscriptionIds           string                  `json:"subscriptionIds"             description:""`                                                                       //
	SubscriptionPeriodStart   int64                   `json:"subscriptionPeriodStart"     description:"matched subscription's current_period_start"`                            // matched subscription's current_period_start
	SubscriptionPeriodEnd     int64                   `json:"subscriptionPeriodEnd"       description:"matched subscription's current_period_end"`                              // matched subscription's current_period_end
	MetricLimit               uint64                  `json:"metricLimit"                 description:""`                                                                       //
	Used                      uint64                  `json:"used"                        description:""`                                                                       //
	ChargeInvoiceId           string                  `json:"chargeInvoiceId"             description:"charge invoice id"`                                                      // charge invoice id
	ChargeInvoice             *bean.Invoice           `json:"chargeInvoice" dc:"chargeInvoice"`
	EventCharge               *bean.EventMetricCharge `json:"eventCharge"                  description:"event charge"`
	ChargeStatus              int                     `json:"chargeStatus"                description:"0-Uncharged，1-charged"` // 0-Uncharged，1-charged
}

func ConvertMerchantMetricEventDetail(ctx context.Context, one *entity.MerchantMetricEvent) *MerchantMetricEventDetail {
	if one == nil {
		return nil
	}
	eventCharge := &bean.EventMetricCharge{}
	if len(one.ChargeData) > 0 {
		_ = utility.UnmarshalFromJsonString(one.ChargeData, &eventCharge)
	}
	return &MerchantMetricEventDetail{
		Id:                        one.Id,
		MerchantId:                one.MerchantId,
		MetricId:                  one.MetricId,
		ExternalEventId:           one.ExternalEventId,
		Metric:                    bean.SimplifyMerchantMetric(query.GetMerchantMetric(ctx, one.MetricId)),
		UserId:                    one.UserId,
		User:                      bean.SimplifyUserAccount(query.GetUserAccountById(ctx, uint64(one.UserId))),
		AggregationPropertyInt:    one.AggregationPropertyInt,
		AggregationPropertyString: one.AggregationPropertyString,
		CreateTime:                one.CreateTime,
		AggregationPropertyData:   one.AggregationPropertyData,
		SubscriptionIds:           one.SubscriptionIds,
		SubscriptionPeriodStart:   one.SubscriptionPeriodStart,
		SubscriptionPeriodEnd:     one.SubscriptionPeriodEnd,
		MetricLimit:               one.MetricLimit,
		Used:                      one.Used,
		ChargeInvoiceId:           one.ChargeInvoiceId,
		ChargeInvoice:             bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, one.ChargeInvoiceId)),
		EventCharge:               eventCharge,
		ChargeStatus:              one.ChargeStatus,
	}
}
