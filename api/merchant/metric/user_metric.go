package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type UserMetricReq struct {
	g.Meta         `path:"/user/metric" tags:"User-Metric" method:"get" summary:"Query User Metric"`
	UserId         int64  `json:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
}

type UserMetricRes struct {
	UserMetric *UserMetric `json:"userMetric" dc:"UserMetric"`
}

type UserMetric struct {
	IsPaid                  bool                         `json:"isPaid" dc:"IsPaid"`
	User                    *ro.UserAccountSimplify      `json:"user" dc:"user"`
	Subscription            *ro.SubscriptionSimplify     `json:"subscription" dc:"Subscription"`
	Plan                    *ro.PlanSimplify             `json:"plan" dc:"Plan"`
	Addons                  []*ro.PlanAddonVo            `json:"addons" dc:"Addon"`
	UserMerchantMetricStats []*ro.UserMerchantMetricStat `json:"userMerchantMetricStats" dc:"UserMerchantMetricStats"`
}
