package user

import "github.com/gogf/gf/v2/frame/g"

type LogoutReq struct {
	g.Meta `path:"/user_logout" tags:"User-Controller" method:"post" summary:"User Logout"`
}

type LogoutRes struct {
}
