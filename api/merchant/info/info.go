package info

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"Merchant-Info-Controller" method:"get" summary:"Get Merchant Info"`
}

type GetRes struct {
	Merchant *entity.Merchant `json:"merchant" dc:"Merchant"`
}

type UpdateReq struct {
	g.Meta      `path:"/update" tags:"Merchant-Info-Controller" method:"post" summary:"Update Merchant Info"`
	CompanyName string `p:"companyName" description:"company_name"`
	Email       string `p:"email"       description:"email"`
	Address     string `p:"address"     description:"address"`
	CompanyLogo string `p:"companyLogo" description:"company_logo"`
	Phone       string `p:"phone"       description:"phone"`
	TimeZone    string `p:"timeZone" description:"User TimeZone"`
	Host        string `p:"host" description:"User Portal Host"`
}

type UpdateRes struct {
	Merchant *entity.Merchant `json:"merchant" dc:"Merchant"`
}
