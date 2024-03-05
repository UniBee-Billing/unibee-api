package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type UserStatReq struct {
	g.Meta         `path:"/user/stat" tags:"Merchant-Metric" method:"get,post" summary:"User Merchant Metric Stat"`
	UserId         int64  `json:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
}

type UserStatRes struct {
	UserMerchantMetricStats []*ro.UserMerchantMetricStat `json:"userMerchantMetricStats" dc:"UserMerchantMetricStats"`
}
