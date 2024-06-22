package download

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Download" method:"get,post" summary:"GetDownloadList"`
	Page   int `json:"page"  description:"Page, Start With 0" `
	Count  int `json:"count"  description:"Count Of Page"`
}

type ListRes struct {
	Downloads []*bean.MerchantBatchTaskSimplify `json:"downloads" dc:"Merchant Member Download List"`
	Total     int                               `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta  `path:"/new" tags:"Download" method:"post" summary:"NewDownload" description:""`
	Task    string            `json:"task" dc:"Task,InvoiceExport"`
	Payload map[string]string `json:"payload" dc:"Payload"`
}

type NewRes struct {
}
