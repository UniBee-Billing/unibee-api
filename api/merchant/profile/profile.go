package profile

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"User-Profile-Controller" method:"get" summary:"get merchant user profile"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}

// with token to be implemented in the future
type ProfileRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"merchant user"`
	// Token string `p:"token" dc:"token string"`
}
