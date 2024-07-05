package task

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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
	g.Meta            `path:"/new_export" tags:"Task" method:"post" summary:"NewExport" description:""`
	Task              string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Payload           map[string]interface{} `json:"payload" dc:"Payload"`
	SkipColumnIndexes []int                  `json:"skipColumnIndexes" dc:"SkipColumnIndexes, the column will be skipped in the export file if its index specified"`
}

type NewRes struct {
}

type NewImportReq struct {
	g.Meta `path:"/new_import" method:"post" mime:"multipart/form-data" tags:"Task" summary:"NewImport"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"File To Upload"`
	Task   string            `json:"task" dc:"Task,UserImport|ActiveSubscriptionImport"`
}
type NewImportRes struct {
}
