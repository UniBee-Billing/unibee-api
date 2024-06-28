package task

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Task" method:"get,post" summary:"GetTaskList"`
	Page   int `json:"page"  description:"Page, Start With 0" `
	Count  int `json:"count"  description:"Count Of Page"`
}

type ListRes struct {
	Downloads []*bean.MerchantBatchTaskSimplify `json:"downloads" dc:"Merchant Member Download List"`
	Total     int                               `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta  `path:"/new_export" tags:"Task" method:"post" summary:"NewExport" description:""`
	Task    string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Payload map[string]interface{} `json:"payload" dc:"Payload"`
}

type NewRes struct {
}
