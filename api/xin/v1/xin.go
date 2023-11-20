package v1

import "github.com/gogf/gf/v2/frame/g"

type InsertReq struct {
	g.Meta `path:"/xinpost" tags:"Xin" method:"post" summary:"You first xin_service api"`
	Name   string `json:"name" dc:"名称"`
}
type InsertRes struct {
	g.Meta `mime:"text/html" example:"string"`
}

type GetReq struct {
	g.Meta `path:"/xinget" tags:"Xin" method:"get" summary:"You first xin_service get api"`
}
type GetRes struct {
	g.Meta `mime:"text/html; charset=UTF-8" example:"string"`
}
