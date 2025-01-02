package subscription

import "github.com/gogf/gf/v2/frame/g"

type NewAdminNoteReq struct {
	g.Meta         `path:"/new_admin_note" tags:"Subscription Note" method:"post" summary:"New Subscription Note"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Note           string `json:"note" dc:"Note" v:"required"`
}

type NewAdminNoteRes struct {
}

type AdminNoteRo struct {
	Id             uint64 `json:"id"               description:"Id"`
	Note           string `json:"note"             description:"Note"`
	CreateTime     int64  `json:"createTime"       description:"CreateTime, UTC Time"`
	SubscriptionId string `json:"subscriptionId" description:"SubscriptionId"`
	UserName       string `json:"userName"   description:"UserName"`
	Mobile         string `json:"mobile"     description:"Mobile"`
	Email          string `json:"email"      description:"Email"`
	FirstName      string `json:"firstName"  description:"FirstName"`
	LastName       string `json:"lastName"   description:"LastName"`
}

type AdminNoteListReq struct {
	g.Meta         `path:"/admin_note_list" tags:"Subscription Note" method:"get,post" summary:"Get Subscription Note List"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	Page           int    `json:"page"  dc:"Page, Start With 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

type AdminNoteListRes struct {
	NoteLists []*AdminNoteRo `json:"noteLists"   description:""`
}
