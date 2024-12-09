package bean

import (
	"strings"
	entity "unibee/internal/model/entity/default"
)

type Merchant struct {
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

func SimplifyMerchant(one *entity.Merchant) *Merchant {
	if one == nil {
		return nil
	}
	return &Merchant{
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

type MerchantMember struct {
	Id         uint64 `json:"id"         description:"userId"`          // userId
	MerchantId uint64 `json:"merchantId" description:"merchant id"`     // merchant id
	Email      string `json:"email"      description:"email"`           // email
	FirstName  string `json:"firstName"  description:"first name"`      // first name
	LastName   string `json:"lastName"   description:"last name"`       // last name
	CreateTime int64  `json:"createTime" description:"create utc time"` // create utc time
	Mobile     string `json:"mobile"     description:"mobile"`          // mobile
	IsOwner    bool   `json:"isOwner" description:"Check Member is Owner" `
}

func SimplifyMerchantMember(one *entity.MerchantMember) *MerchantMember {
	if one == nil {
		return nil
	}
	isOwner := false
	if strings.Contains(one.Role, "Owner") {
		isOwner = true
	}
	return &MerchantMember{
		Id:         one.Id,
		MerchantId: one.MerchantId,
		Email:      one.Email,
		FirstName:  one.FirstName,
		LastName:   one.LastName,
		CreateTime: one.CreateTime,
		Mobile:     one.Mobile,
		IsOwner:    isOwner,
	}
}
