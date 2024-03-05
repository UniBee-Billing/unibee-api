package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type UserMetricReq struct {
	g.Meta         `path:"/user/metric" tags:"Merchant-User-Metric" method:"get" summary:"Query User Metric"`
	UserId         int64  `json:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
}

type UserMetricRes struct {
	UserMetric *ro.UserMetric `json:"userMetric" dc:"UserMetric"`
}
