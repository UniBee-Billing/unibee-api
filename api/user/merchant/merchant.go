package merchant

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"User-Merchant-Info" method:"get" summary:"Get Merchant Info"`
}

type GetRes struct {
	Merchant *bean.Merchant `json:"merchant" dc:"Merchant"`
}
