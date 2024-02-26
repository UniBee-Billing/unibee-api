package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type MerchantMetricEventReq struct {
	g.Meta           `path:"/merchant_metric_event" tags:"Merchant-Metric-Controller" method:"post" summary:"Merchant Metric Event"`
	MetricCode       string      `p:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `p:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `p:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `p:"metricProperties" dc:"MetricProperties"`
}

type MerchantMetricEventRes struct {
	MerchantMetricEvent *entity.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DelMerchantMetricEventReq struct {
	g.Meta          `path:"/del_merchant_metric_event" tags:"Merchant-Metric-Controller" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `p:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `p:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId string `p:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DelMerchantMetricEventRes struct {
}
