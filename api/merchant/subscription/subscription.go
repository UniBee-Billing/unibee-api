package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/ro"
)

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"User-Subscription-Controller" method:"post" summary:"订阅列表"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int64 `p:"userId"  dc:"UserId" `
	Status     int   `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	Page       int   `p:"page" dc:"分页页码,0开始" `
	Count      int   `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `p:"subscriptions" dc:"订阅明细"`
}
