package checkout

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta    `path:"/list" tags:"Checkout Setup" method:"get,post" summary:"Get Merchant Checkout list"`
	SearchKey string `q:"searchKey" description:"Search checkout id|name|description"`
}

type ListRes struct {
	MerchantCheckouts []*bean.MerchantCheckout `json:"merchantCheckouts" dc:"MerchantCheckout List Object"`
}

type DetailReq struct {
	g.Meta     `path:"/detail" tags:"Checkout Setup" method:"get,post" summary:"Get Merchant Checkout Detail"`
	CheckoutId int64 `json:"checkoutId"               description:"checkout id" v:"required#checkoutId"`
}

type DetailRes struct {
	MerchantCheckout *bean.MerchantCheckout `json:"merchantCheckout" dc:"MerchantCheckout Object"`
}

type NewReq struct {
	g.Meta      `path:"/new_checkout" tags:"Checkout Setup" method:"post" summary:"New Merchant Checkout"`
	Name        string                 `json:"name"                  description:"name"`        // name
	Description string                 `json:"description"           description:"description"` // description
	Data        map[string]interface{} `json:"data"              description:"checkout config data"`
	StagingData map[string]interface{} `json:"stagingData"              description:"checkout staging config data"`
}

type NewRes struct {
	MerchantCheckout *bean.MerchantCheckout `json:"merchantCheckout" dc:"MerchantCheckout Object"`
}

type EditReq struct {
	g.Meta      `path:"/edit_checkout" tags:"Checkout Setup" method:"post" summary:"Edit Merchant Checkout"`
	CheckoutId  int64                  `json:"checkoutId"               description:"checkout id" v:"required#checkoutId"`
	Name        string                 `json:"name"                  description:"name"`        // name
	Description string                 `json:"description"           description:"description"` // description
	Data        map[string]interface{} `json:"data"              description:"checkout config data"`
	StagingData map[string]interface{} `json:"stagingData"              description:"checkout staging config data"`
}

type EditRes struct {
	MerchantCheckout *bean.MerchantCheckout `json:"merchantCheckout" dc:"MerchantCheckout Object"`
}

type GetLinkReq struct {
	g.Meta     `path:"/get_link" tags:"Checkout Setup" method:"get,post" summary:"Get Merchant Checkout Link"`
	CheckoutId int64 `json:"checkoutId"               description:"checkout id" v:"required#checkoutId"`
	PlanId     int64 `json:"planId" v:"required#checkoutPlanId"`
}

type GetLinkRes struct {
	Link string `json:"link" dc:"Merchant Checkout Link"`
}

type ArchiveReq struct {
	g.Meta     `path:"/archive" tags:"Checkout Setup" method:"post" summary:"Archive Merchant Checkout"`
	CheckoutId int64 `json:"checkoutId"               description:"checkout id" v:"required#checkoutId"`
}

type ArchiveRes struct {
}
