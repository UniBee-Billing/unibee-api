package bean

import entity "unibee/internal/model/entity/default"

type MerchantSimplify struct {
	Id          uint64 `json:"id"          description:"merchant_id"`                // merchant_id
	UserId      int64  `json:"userId"      description:"create_user_id"`             // create_user_id
	Type        int    `json:"type"        description:"type"`                       // type
	Email       string `json:"email"       description:"email"`                      // email
	BusinessNum string `json:"businessNum" description:"business_num"`               // business_num
	Name        string `json:"name"        description:"name"`                       // name
	Idcard      string `json:"idcard"      description:"idcard"`                     // idcard
	Location    string `json:"location"    description:"location"`                   // location
	Address     string `json:"address"     description:"address"`                    // address
	CompanyLogo string `json:"companyLogo" description:"company_logo"`               // company_logo
	HomeUrl     string `json:"homeUrl"     description:""`                           //
	Phone       string `json:"phone"       description:"phone"`                      // phone
	CreateTime  int64  `json:"createTime"  description:"create utc time"`            // create utc time
	TimeZone    string `json:"timeZone"    description:"merchant default time zone"` // merchant default time zone
	Host        string `json:"host"        description:"merchant user portal host"`  // merchant user portal host
	CompanyName string `json:"companyName" description:"company_name"`               // company_name
}

func SimplifyMerchant(one *entity.Merchant) *MerchantSimplify {
	if one == nil {
		return nil
	}
	return &MerchantSimplify{
		Id:          one.Id,
		UserId:      one.UserId,
		Type:        one.Type,
		Email:       one.Email,
		BusinessNum: one.BusinessNum,
		Name:        one.Name,
		Idcard:      one.Idcard,
		Location:    one.Location,
		Address:     one.Address,
		CompanyLogo: one.CompanyLogo,
		HomeUrl:     one.HomeUrl,
		Phone:       one.Phone,
		CreateTime:  one.CreateTime,
		TimeZone:    one.TimeZone,
		Host:        one.Host,
		CompanyName: one.CompanyName,
	}
}
