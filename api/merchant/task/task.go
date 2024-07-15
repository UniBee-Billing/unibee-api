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
	Downloads []*bean.MerchantBatchTaskSimplify `json:"downloads" dc:"Merchant Member Task List"`
	Total     int                               `json:"total" dc:"Total"`
}

type ExportColumnListReq struct {
	g.Meta `path:"/export_column_list" tags:"Task" method:"post" summary:"ExportColumnList" description:""`
	Task   string `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
}

type ExportColumnListRes struct {
	Columns []interface{} `json:"columns" dc:"Export Column List"`
}

type NewReq struct {
	g.Meta        `path:"/new_export" tags:"Task" method:"post" summary:"NewExport" description:""`
	Task          string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Payload       map[string]interface{} `json:"payload" dc:"Payload"`
	ExportColumns []string               `json:"exportColumns" dc:"ExportColumns, the export file column list, will export all columns if not specified"`
}

type NewRes struct {
}

type NewImportReq struct {
	g.Meta `path:"/new_import" method:"post" mime:"multipart/form-data" tags:"Task" summary:"NewImport"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"File To Upload"`
	Task   string            `json:"task" dc:"Task,UserImport|ActiveSubscriptionImport|HistorySubscriptionImport"`
}
type NewImportRes struct {
}
