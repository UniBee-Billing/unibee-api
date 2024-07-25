package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewEventReq struct {
	g.Meta           `path:"/event/new" tags:"Metric" method:"post" summary:"Merchant Metric Event"`
	MetricCode       string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `json:"metricProperties" dc:"MetricProperties"`
}

type NewEventRes struct {
	MerchantMetricEvent *bean.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DeleteEventReq struct {
	g.Meta          `path:"/event/delete" tags:"Metric" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId string `json:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DeleteEventRes struct {
}
