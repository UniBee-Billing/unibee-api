package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/metric_event"
)

type UserStatReq struct {
	g.Meta         `path:"/user/stat" tags:"Merchant-Metric-Controller" method:"post" summary:"User Merchant Metric Stat"`
	UserId         int64  `p:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `p:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
}

type UserStatRes struct {
	UserMerchantMetricStats []*metric_event.UserMerchantMetricStat `json:"userMerchantMetricStats" dc:"UserMerchantMetricStats"`
}
