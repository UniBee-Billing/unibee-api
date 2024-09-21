package track

import "github.com/gogf/gf/v2/frame/g"

type SetupSegmentReq struct {
	g.Meta           `path:"/setup_segment" tags:"Track" method:"post" summary:"SegmentSetup"`
	ServerSideSecret string `json:"serverSideSecret" dc:"ServerSideSecret" v:"required"`
	UserPortalSecret string `json:"userPortalSecret" dc:"UserPortalSecret" v:"required"`
}
type SetupSegmentRes struct {
}
