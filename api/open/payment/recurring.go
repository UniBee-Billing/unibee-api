package payment

import "github.com/gogf/gf/v2/frame/g"

type DisableRecurringDetailsReq struct {
	g.Meta `path:"/disableRecurringDetails" tags:"Open-Payment-Controller" method:"post" summary:"1.7 终⽌账户绑定信息(仅Klarna支持）"`
}
type DisableRecurringDetailsRes struct {
}

type ListRecurringDetailsReq struct {
	g.Meta `path:"/listRecurringDetails" tags:"Open-Payment-Controller" method:"post" summary:"1.6 查询绑定账户信息(仅Klarna支持）"`
}
type ListRecurringDetailsRes struct {
}
