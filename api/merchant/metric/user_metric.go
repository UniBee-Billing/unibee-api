package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/metric_event"
)

type UserMerchantMetricStatReq struct {
	g.Meta         `path:"/user_merchant_metric_stat" tags:"Merchant-Metric-Controller" method:"post" summary:"User Merchant Metric Stat"`
	UserId         int64  `p:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `p:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
}

type UserMerchantMetricStatRes struct {
	UserMerchantMetricStats []*metric_event.UserMerchantMetricStat `json:"userMerchantMetricStats" dc:"UserMerchantMetricStats"`
}
