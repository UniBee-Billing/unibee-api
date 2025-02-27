package task

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Task" method:"get,post" summary:"Get Task List"`
	Page   int `json:"page"  description:"Page, Start With 0" `
	Count  int `json:"count"  description:"Count Of Page"`
}

type ListRes struct {
	Downloads []*bean.MerchantBatchTask `json:"downloads" dc:"Merchant Member Task List"`
	Total     int                       `json:"total" dc:"Total"`
}

type ExportColumnListReq struct {
	g.Meta `path:"/export_column_list" tags:"Task" method:"post" summary:"Export Column List" description:""`
	Task   string `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
}

type ExportColumnListRes struct {
	Columns        []interface{}            `json:"columns" dc:"Export Column List"`
	GroupColumns   map[string][]interface{} `json:"groupColumns" dc:"Group Column List"`
	ColumnHeaders  map[string]string        `json:"columnHeaders" dc:"Export Column Headers"`
	ColumnComments map[string]string        `json:"columnComments" dc:"Export Column Comments"`
}

type NewReq struct {
	g.Meta        `path:"/new_export" tags:"Task" method:"post" summary:"New Export" description:""`
	Task          string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport" v:"required"`
	Payload       map[string]interface{} `json:"payload" dc:"Payload, Task query parameters, positive or negative 'timeZone' available for all task"`
	ExportColumns []string               `json:"exportColumns" dc:"ExportColumns, the export file column list, will export all columns if not specified"`
	Format        string                 `json:"format" dc:"The format of export file, xlsx|csv, will be xlsx if not specified"`
}

type NewRes struct {
}

type NewImportReq struct {
	g.Meta `path:"/new_import" method:"post" mime:"multipart/form-data" tags:"Task" summary:"New Import"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"File To Upload" v:"required"`
	Task   string            `json:"task" dc:"Task,UserImport|ActiveSubscriptionImport|HistorySubscriptionImport" v:"required"`
}
type NewImportRes struct {
}

type NewTemplateReq struct {
	g.Meta        `path:"/new_export_template" tags:"Task" method:"post" summary:"New Export Template" description:""`
	Name          string                 `json:"name"      v:"required"    description:"name"`
	Task          string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport" v:"required"`
	Payload       map[string]interface{} `json:"payload" dc:"Payload"`
	ExportColumns []string               `json:"exportColumns" dc:"ExportColumns, the export file column list, will export all columns if not specified, first char should lower case"`
	Format        string                 `json:"format" dc:"The format of export file, xlsx|csv, will be xlsx if not specified"`
}

type NewTemplateRes struct {
	Template *bean.MerchantBatchExportTemplate `json:"template" dc:"Merchant Member Export Template"`
}

type EditTemplateReq struct {
	g.Meta        `path:"/edit_export_template" tags:"Task" method:"post" summary:"Edit Export Template" description:""`
	TemplateId    int64                   `json:"templateId"    v:"required"      description:"templateId"`
	Name          *string                 `json:"name"          description:"name"`
	Task          *string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Payload       *map[string]interface{} `json:"payload" dc:"Payload"`
	ExportColumns *[]string               `json:"exportColumns" dc:"ExportColumns, the export file column list, will export all columns if not specified"`
	Format        *string                 `json:"format" dc:"The format of export file, xlsx|csv, will be xlsx if not specified"`
}

type EditTemplateRes struct {
	Template *bean.MerchantBatchExportTemplate `json:"template" dc:"Merchant Member Export Template"`
}

type DeleteTemplateReq struct {
	g.Meta     `path:"/delete_export_template" tags:"Task" method:"post" summary:"Delete Export Template" description:""`
	TemplateId int64 `json:"templateId"     v:"required"       description:"templateId"`
}

type DeleteTemplateRes struct {
}

type ExportTemplateListReq struct {
	g.Meta `path:"/export_template_list" tags:"Task" method:"get,post" summary:"Get Export Template List"`
	Task   string `json:"task" dc:"Filter Task, Optional, InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Page   int    `json:"page"  description:"Page, Start With 0" `
	Count  int    `json:"count"  description:"Count Of Page"`
}

type ExportTemplateListRes struct {
	Templates []*bean.MerchantBatchExportTemplate `json:"templates" dc:"Merchant Member Export Template List"`
	Total     int                                 `json:"total" dc:"Total"`
}
