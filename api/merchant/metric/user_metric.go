package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type UserMetricReq struct {
	g.Meta         `path:"/user/metric" tags:"User Metric" method:"get" summary:"Query User Metric"`
	UserId         int64  `json:"userId" dc:"UserId, One Of UserId|ExternalUserId Needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, One Of UserId|ExternalUserId Needed"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type UserMetricRes struct {
	UserMetric *UserMetric `json:"userMetric" dc:"UserMetric"`
}

type UserMetric struct {
	IsPaid               bool                                   `json:"isPaid" dc:"IsPaid"`
	Product              *bean.Product                          `json:"product" dc:"Product"`
	User                 *bean.UserAccount                      `json:"user" dc:"user"`
	Subscription         *bean.Subscription                     `json:"subscription" dc:"Subscription"`
	Plan                 *bean.Plan                             `json:"plan" dc:"Plan"`
	Addons               []*bean.PlanAddonDetail                `json:"addons" dc:"Addon"`
	LimitStats           []*detail.UserMerchantMetricLimitStat  `json:"limitStats" dc:"LimitStats"`
	MeteredChargeStats   []*detail.UserMerchantMetricChargeStat `json:"meteredChargeStats" dc:"MeteredChargeStats"`
	RecurringChargeStats []*detail.UserMerchantMetricChargeStat `json:"recurringChargeStats" dc:"RecurringChargeStats"`
	Description          string                                 `json:"description" dc:"description"`
}
