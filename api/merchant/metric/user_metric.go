package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type UserMetricReq struct {
	g.Meta         `path:"/user/metric" tags:"User Metric" method:"get" summary:"Query User Metric"`
	UserId         int64  `json:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type UserMetricRes struct {
	UserMetric *detail.UserMetric `json:"userMetric" dc:"UserMetric"`
}

type UserSubscriptionMetricReq struct {
	g.Meta         `path:"/user/sub/metric" tags:"User Metric" method:"get" summary:"Query User Metric By Subscription"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId"`
}

type UserSubscriptionMetricRes struct {
	UserMetric *detail.UserMetric `json:"userMetric" dc:"UserMetric"`
}
