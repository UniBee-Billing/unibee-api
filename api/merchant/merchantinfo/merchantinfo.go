package merchantinfo

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type MerchantInfoReq struct {
	g.Meta `path:"/info" tags:"Merchant-Info-Controller" method:"get" summary:"Get Merchant Info"`
}

type MerchantInfoRes struct {
	MerchantInfo *entity.MerchantInfo `p:"merchantInfo" dc:"merchantInfo"`
}

type MerchantInfoUpdateReq struct {
	g.Meta      `path:"/update" tags:"Merchant-Info-Controller" method:"post" summary:"Update Merchant Info"`
	CompanyName string `p:"companyName" description:"company_name"` // company_name
	Email       string `p:"email"       description:"email"`        // email
	Address     string `p:"address"     description:"address"`      // address
	CompanyLogo string `p:"companyLogo" description:"company_logo"` // company_logo
	FirstName   string `p:"firstName"   description:"first_name"`   // first_name
	LastName    string `p:"lastName"    description:"last_name"`    // last_name
	Phone       string `p:"phone"       description:"phone"`        // phone
}

type MerchantInfoUpdateRes struct {
	MerchantInfo *entity.MerchantInfo `p:"merchantInfo" dc:"merchantInfo"`
}
