package checkout

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type GetReq struct {
	g.Meta     `path:"/get" tags:"Checkout" method:"post" summary:"Get Merchant Checkout"`
	CheckoutId int64 `json:"checkoutId" description:"checkout id" v:"required#checkoutId"`
}

type GetRes struct {
	MerchantCheckout *bean.MerchantCheckout `json:"merchantCheckout" dc:"MerchantCheckout Object"`
}
