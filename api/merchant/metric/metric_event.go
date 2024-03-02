package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type NewEventReq struct {
	g.Meta           `path:"/event/new" tags:"Merchant-Metric-Controller" method:"post" summary:"Merchant Metric Event"`
	MetricCode       string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `json:"metricProperties" dc:"MetricProperties"`
}

type NewEventRes struct {
	MerchantMetricEvent *entity.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DeleteEventReq struct {
	g.Meta          `path:"/event/delete" tags:"Merchant-Metric-Controller" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId string `json:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DeleteEventRes struct {
}
