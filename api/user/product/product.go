package product

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta    `path:"/list" tags:"Product" method:"get,post" summary:"ProductList"`
	SortField string `json:"sortField" dc:"Sort Field，id|create_time|gmt_modify，Default id" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start 0" `
	Count     int    `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Products []*bean.Product `json:"products" dc:"Product Object List"`
	Total    int             `json:"total" dc:"Total"`
}
