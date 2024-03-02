package merchant

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"User-Merchant-Info-Controller" method:"get" summary:"Get Merchant Info"`
}

type GetRes struct {
	Merchant *entity.Merchant `json:"merchant" dc:"Merchant"`
}
