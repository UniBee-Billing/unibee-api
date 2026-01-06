package translater

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TranslateReq struct {
	g.Meta     `path:"/translate" tags:"Checkout" method:"post" summary:"Translate"`
	MerchantId uint64   `json:"merchantId" description:""  v:"required"`
	TargetLang string   `json:"targetLang" description:""  v:"required"`
	Texts      []string `json:"texts" description:""  v:"required"`
}
type TranslateRes struct {
	Translates map[string]string `json:"translates"`
	Success    bool              `json:"success"`
}
