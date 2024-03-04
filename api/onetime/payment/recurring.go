package payment

import "github.com/gogf/gf/v2/frame/g"

type DisableRecurringDetailReq struct {
	g.Meta `path:"/disableRecurringDetail" tags:"OneTime-Payment" method:"post" summary:"Disable Recurring Details (Support Klarna）"`
}
type DisableRecurringDetailRes struct {
}

type RecurringDetailListReq struct {
	g.Meta `path:"/recurringDetailList" tags:"OneTime-Payment" method:"post" summary:"Query Recurring Detail List (Support Klarna）"`
}
type RecurringDetailListRes struct {
}
