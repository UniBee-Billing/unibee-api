package payment

import "github.com/gogf/gf/v2/frame/g"

type DisableRecurringDetailsReq struct {
	g.Meta `path:"/disableRecurringDetails" tags:"Open-Payment-Controller" method:"post" summary:"Disable Recurring Details (Support Klarna）"`
}
type DisableRecurringDetailsRes struct {
}

type ListRecurringDetailsReq struct {
	g.Meta `path:"/listRecurringDetails" tags:"Open-Payment-Controller" method:"post" summary:"Query Recurring Detail List (Support Klarna）"`
}
type ListRecurringDetailsRes struct {
}
