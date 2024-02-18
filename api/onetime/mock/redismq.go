package mock

import "github.com/gogf/gf/v2/frame/g"

type MockMessageSendReq struct {
	g.Meta  `path:"/message_mock_test" tags:"Open-Mock-Controller" method:"post" summary:"Mock Message Test"`
	Message string `p:"message" dc:"Message" v:"required"`
}
type MockMessageSendRes struct {
}
