package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type TimeLineListReq struct {
	g.Meta    `path:"/timeline_list" tags:"Subscription Timeline" method:"get,post" summary:"Get Subscription TimeLine List"`
	UserId    uint64 `json:"userId" dc:"Filter UserId, Default All " `
	SortField string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start With 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	SubscriptionTimeLines []*detail.SubscriptionTimeLineDetail `json:"subscriptionTimeLines" description:"SubscriptionTimeLines" `
	Total                 int                                  `json:"total" dc:"Total"`
}
