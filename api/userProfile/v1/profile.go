package v1

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Auth-Controller" method:"get" summary:"1.1 get user profile"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}
// with token to be implemented in the future
type ProfileRes struct {
	User *entity.UserAccount `p:"user" dc:"user"`
	// Token string `p:"token" dc:"token string"`
}
