package oss

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileUploadReq struct {
	g.Meta `path:"/file" method:"post" mime:"multipart/form-data" tags:"File-Controller" summary:"Upload"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"File To Upload"`
}
type FileUploadRes struct {
	Url string `json:"url"  dc:"URL Of File Or Image"`
}
