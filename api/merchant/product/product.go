package product

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewReq struct {
	g.Meta      `path:"/new" tags:"Product" method:"post" summary:"Create Product"`
	ProductName string                  `json:"productName" description:"ProductName"`                                // ProductName
	Description string                  `json:"description" description:"description"`                                // description
	ImageUrl    string                  `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     string                  `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      int                     `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	Metadata    *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type NewRes struct {
	Product *bean.Product `json:"product" dc:"Product Object"`
}

type EditReq struct {
	g.Meta      `path:"/edit" tags:"Product" method:"post" summary:"Edit Product" dc:"Edit exist product, product is editable for both active and inactive status "`
	ProductId   uint64                  `json:"productId" dc:"Id of product" v:"required"`
	ProductName *string                 `json:"productName" description:"ProductName"`                                // ProductName
	Description *string                 `json:"description" description:"description"`                                // description
	ImageUrl    *string                 `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     *string                 `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      *int                    `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	Metadata    *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type EditRes struct {
	Product *bean.Product `json:"product" dc:"Product Object"`
}

type ListReq struct {
	g.Meta    `path:"/list" tags:"Product" method:"get,post" summary:"Get Product List"`
	Status    []int  `json:"status" dc:"Filter, Default All，,Status，1-active，2-inactive" `
	SortField string `json:"sortField" dc:"Sort Field，id|create_time|gmt_modify，Default id" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start 0" `
	Count     int    `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Products []*bean.Product `json:"products" dc:"Product Object List"`
	Total    int             `json:"total" dc:"Total"`
}

type CopyReq struct {
	g.Meta    `path:"/copy" tags:"Product" method:"post" summary:"Copy Product"`
	ProductId uint64 `json:"productId" dc:"ProductId" v:"required"`
}
type CopyRes struct {
	Product *bean.Product `json:"product" dc:"Product Object"`
}

type ActivateReq struct {
	g.Meta    `path:"/activate" tags:"Product" method:"post" summary:"Activate Product"`
	ProductId uint64 `json:"productId" dc:"ProductId" v:"required"`
}
type ActivateRes struct {
}

type InactiveReq struct {
	g.Meta    `path:"/inactivate" tags:"Product" method:"post" summary:"Inactivate Product" `
	ProductId uint64 `json:"productId" dc:"ProductId" v:"required"`
}
type InactiveRes struct {
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Product" method:"get,post" summary:"Product Detail"`
	ProductId uint64 `json:"productId" dc:"ProductId" v:"required"`
}
type DetailRes struct {
	Product *bean.Product `json:"product" dc:"Product Object"`
}

type DeleteReq struct {
	g.Meta    `path:"/delete" tags:"Product" method:"post" summary:"Delete Product" dc:"Product can being deleted when has no plan linked"`
	ProductId uint64 `json:"productId" dc:"ProductId" v:"required"`
}
type DeleteRes struct {
}
