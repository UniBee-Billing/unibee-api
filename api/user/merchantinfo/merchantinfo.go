package merchantinfo

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type MerchantInfoReq struct {
	g.Meta `path:"/info" tags:"User-Merchant-Info-Controller" method:"get" summary:"Get Merchant Info"`
}

type MerchantInfoRes struct {
	MerchantInfo *entity.MerchantInfo `p:"merchantInfo" dc:"merchantInfo"`
}
