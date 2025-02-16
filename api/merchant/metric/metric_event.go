package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewEventReq struct {
	g.Meta           `path:"/event/new" tags:"User Metric" method:"post" summary:"New Merchant Metric Event"`
	MetricCode       string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `json:"metricProperties" dc:"MetricProperties"`
	ProductId        int64       `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type NewEventRes struct {
	MerchantMetricEvent *bean.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DeleteEventReq struct {
	g.Meta          `path:"/event/delete" tags:"User Metric" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId string `json:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DeleteEventRes struct {
}
