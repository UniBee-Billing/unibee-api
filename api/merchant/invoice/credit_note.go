package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/api/bean/detail"
)

type CreditNoteListReq struct {
	g.Meta          `path:"/credit_note/list" tags:"Invoice" method:"post" summary:"Bulk CreditNote(Refund Invoice) Invoice List" dc:"Bulk credit note invoice list"`
	SearchKey       string            `json:"searchKey" dc:"The search key of invoice" `
	Emails          string            `json:"emails" dc:"The email list of invoice user, split by commas or semicolons" `
	File            *ghttp.UploadFile `json:"file" type:"file" dc:"Email CSV File To Search"`
	Status          []int             `json:"status" dc:"The status of invoice, 2-processing｜3-paid | 4-failed | 5-cancelled" `
	GatewayIds      []int64           `json:"gatewayIds" dc:"GatewayIds, Search Filter GatewayIds" `
	PlanIds         []int64           `json:"planIds" dc:"PlanIds, Search Filter PlanIds" `
	Currency        string            `json:"currency" dc:"The currency of invoice" `
	SortField       string            `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType        string            `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	Page            int               `json:"page"  dc:"Page, Start 0" `
	Count           int               `json:"count"  dc:"Count" dc:"Count By Page" `
	CreateTimeStart int64             `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64             `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
}

type CreditNoteListRes struct {
	CreditNotes []*detail.CreditNoteDetail `json:"creditNotes" dc:"CreditNote Detail Object List"`
	Total       int                        `json:"total" dc:"Total"`
}
