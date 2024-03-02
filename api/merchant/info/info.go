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
	CompanyName string `json:"companyName" description:"company_name"`
	Email       string `json:"email"       description:"email"`
	Address     string `json:"address"     description:"address"`
	CompanyLogo string `json:"companyLogo" description:"company_logo"`
	Phone       string `json:"phone"       description:"phone"`
	TimeZone    string `json:"timeZone" description:"User TimeZone"`
	Host        string `json:"host" description:"User Portal Host"`
}

type UpdateRes struct {
	Merchant *entity.Merchant `json:"merchant" dc:"Merchant"`
}
