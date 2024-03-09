package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"Profile" method:"get" summary:"Get Merchant Info"`
}

type GetRes struct {
	Merchant *entity.Merchant      `json:"merchant" dc:"Merchant"`
	Env      string                `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd   bool                  `json:"isProd" description:"Check System Env Is Prod, true|false" `
	TimeZone []string              `json:"TimeZone" description:"TimeZone List" `
	Currency []*ro.Currency        `json:"Currency" description:"Currency List" `
	Gateway  []*ro.GatewaySimplify `json:"gateway" description:"Gateway List" `
}

type UpdateReq struct {
	g.Meta      `path:"/update" tags:"Profile" method:"post" summary:"Update Merchant Info"`
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
