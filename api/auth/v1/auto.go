package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type LoginReq struct {
	g.Meta `path:"/sso/login" tags:"Auth-Controller" method:"post" summary:"1.1 用户登录"`
}
type LoginRes struct {
}

type RegisterReq struct {
	g.Meta `path:"/sso/register" tags:"Auth-Controller" method:"post" summary:"1.2 用户注册"`
	Email  string `p:"email" dc:"email" v:"required"`
	Phone  string `p:"Phone" dc:"手机号" `
	Gender string `p:"gender" dc:"性别" `
}
type RegisterRes struct {
	User *entity.UserAccount `p:"user" dc:"user"`
}
