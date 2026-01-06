package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Merchant struct {
	Id                  uint64 `json:"id"          description:"merchant_id"`                // merchant_id
	UserId              int64  `json:"userId"      description:"create_user_id"`             // create_user_id
	Type                int    `json:"type"        description:"type"`                       // type
	Email               string `json:"email"       description:"email"`                      // email
	Name                string `json:"name"        description:"name"`                       // name
	Location            string `json:"location"    description:"location"`                   // location
	Address             string `json:"address"     description:"address"`                    // address
	CompanyLogo         string `json:"companyLogo" description:"company_logo"`               // company_logo
	HomeUrl             string `json:"homeUrl"     description:""`                           //
	Phone               string `json:"phone"       description:"phone"`                      // phone
	CreateTime          int64  `json:"createTime"  description:"create utc time"`            // create utc time
	TimeZone            string `json:"timeZone"    description:"merchant default time zone"` // merchant default time zone
	Host                string `json:"host"        description:"merchant user portal host"`  // merchant user portal host
	CompanyName         string `json:"companyName" description:"company_name"`               // company_name
	CountryCode         string `json:"countryCode" dc:"Country Code"`
	CountryName         string `json:"countryName" dc:"Country Name"`
	CompanyVatNumber    string `json:"companyVatNumber" dc:"Country Vat Number"`
	CompanyRegistryCode string `json:"companyRegistryCode" dc:"Country Registry Code"`
}

func SimplifyMerchant(one *entity.Merchant) *Merchant {
	if one == nil {
		return nil
	}
	createTime := one.CreateTime
	if createTime == 0 && one.GmtCreate != nil {
		createTime = one.GmtCreate.Timestamp()
	}
	return &Merchant{
		Id:                  one.Id,
		UserId:              one.UserId,
		Type:                one.Type,
		Email:               one.Email,
		CompanyVatNumber:    one.BusinessNum,
		Name:                one.Name,
		CompanyRegistryCode: one.Idcard,
		Location:            one.Location,
		Address:             one.Address,
		CompanyLogo:         one.CompanyLogo,
		HomeUrl:             one.HomeUrl,
		Phone:               one.Phone,
		CreateTime:          createTime,
		TimeZone:            one.TimeZone,
		Host:                one.Host,
		CompanyName:         one.CompanyName,
		CountryCode:         one.CountryCode,
		CountryName:         one.CountryName,
	}
}

type MerchantMember struct {
	Id            uint64                 `json:"id"         description:"userId"`          // userId
	MerchantId    uint64                 `json:"merchantId" description:"merchant id"`     // merchant id
	Email         string                 `json:"email"      description:"email"`           // email
	FirstName     string                 `json:"firstName"  description:"first name"`      // first name
	LastName      string                 `json:"lastName"   description:"last name"`       // last name
	CreateTime    int64                  `json:"createTime" description:"create utc time"` // create utc time
	Mobile        string                 `json:"mobile"     description:"mobile"`          // mobile
	IsOwner       bool                   `json:"isOwner" description:"Check Member is Owner" `
	IsBlankPasswd bool                   `json:"isBlankPasswd" description:"is blank password"`
	TotpType      int                    `json:"totpType"   description:"0-Inactive, 1-General, Google Authenticator | 2-Microsoft Authenticator | 3-Authy | 4-1Password | 5-LastPass | 6-FreeOTP | 7-Other TOTP"`
	OAuthAccounts []*Oauth               `json:"oauthAccounts" description:"List of connected OAuth accounts"`
	Metadata      map[string]interface{} `json:"metadata"                  description:""`
}

func SimplifyMerchantMember(one *entity.MerchantMember) *MerchantMember {
	if one == nil {
		return nil
	}
	isOwner := false
	if strings.Contains(one.Role, "Owner") {
		isOwner = true
	}
	var identityData = &OauthIdentity{}
	_ = utility.UnmarshalFromJsonString(one.AuthJs, &identityData)
	oauthAccounts := make([]*Oauth, 0)
	for _, account := range identityData.OAuthAccountMap {
		oauthAccounts = append(oauthAccounts, account)
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyUserAccount Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &MerchantMember{
		Id:            one.Id,
		MerchantId:    one.MerchantId,
		Email:         one.Email,
		FirstName:     one.FirstName,
		LastName:      one.LastName,
		CreateTime:    one.CreateTime,
		Mobile:        one.Mobile,
		IsBlankPasswd: len(one.Password) == 0,
		IsOwner:       isOwner,
		TotpType:      one.TotpValidatorType,
		OAuthAccounts: oauthAccounts,
		Metadata:      metadata,
	}
}

type MerchantMemberDevice struct {
	Name                     string `json:"name"         description:"Name"` // userId
	Identity                 string `json:"identity"     description:"Identity"`
	LastLoginTime            int64  `json:"lastLoginTime" description:"Last Login Time"`
	LastActiveTime           int64  `json:"lastActiveTime" description:"Last Active Time"`
	LastTotpVerificationTime int64  `json:"lastTotpVerificationTime" description:"Last Totp Verification Time"`
	Status                   bool   `json:"status" description:"true-Active, false-Offline"`
	IPAddress                string `json:"ipAddress"     description:"IP Address"`
	CurrentDevice            bool   `json:"currentDevice" description:"Is CurrentDevice"`
}

type Oauth struct {
	Provider      string `json:"provider"`
	ProviderId    string `json:"providerId"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
}

type OauthIdentity struct {
	Identity        string            `json:"identity"`
	OAuthAccountMap map[string]*Oauth `json:"oauthAccountsMap" description:"Map of connected OAuth accounts"`
}

type MerchantMultiCurrencyConfig struct {
	Name            string                    `json:"name"`
	DefaultCurrency string                    `json:"defaultCurrency"`
	MultiCurrencies []*MerchantCurrencyConfig `json:"currencyConfigs"`
	LastUpdateTime  int64                     `json:"lastUpdateTime" description:"Last Update UTC Time"`
}

type MerchantCurrencyConfig struct {
	Currency     string  `json:"currency" description:"target currency"`
	AutoExchange bool    `json:"autoExchange" description:"using https://app.exchangerate-api.com/ to update exchange rate if true, the exchange APIKey need setup first"`
	ExchangeRate float64 `json:"exchangeRate"  description:"the exchange rate of gateway, no setup required if AutoExchange is true"`
}
