package merchantinfo

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type MerchantInfoReq struct {
	g.Meta `path:"/info" tags:"Merchant-Info-Controller" method:"get" summary:"Get Merchant Info"`
}

type MerchantInfoRes struct {
	MerchantInfo *entity.MerchantInfo `p:"merchantInfo" dc:"merchantInfo"`
}

type MerchantInfoUpdateReq struct {
	g.Meta      `path:"/update" tags:"Merchant-Info-Controller" method:"post" summary:"Update Merchant Info"`
	CompanyName string `p:"companyName" description:"company_name"`
	Email       string `p:"email"       description:"email"`
	Address     string `p:"address"     description:"address"`
	CompanyLogo string `p:"companyLogo" description:"company_logo"`
	Phone       string `p:"phone"       description:"phone"`
	TimeZone    string `p:"timeZone" description:"User TimeZone"`
	Host        string `p:"host" description:"User Portal Host"`
}

type MerchantInfoUpdateRes struct {
	MerchantInfo *entity.MerchantInfo `p:"merchantInfo" dc:"merchantInfo"`
}
