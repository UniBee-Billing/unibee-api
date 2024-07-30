package user

type NewUserInternalReq struct {
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	Password       string `json:"password" dc:"Password"`
	Phone          string `json:"phone" dc:"Phone" `
	Address        string `json:"address" dc:"Address"`
	UserName       string `json:"userName" dc:"UserName"`
	CountryCode    string `json:"countryCode" dc:"CountryCode"`
	Type           int64  `json:"type" dc:"User type, 1-Individual|2-organization"`
	CompanyName    string `json:"companyName" dc:"company name"`
	VATNumber      string `json:"vATNumber" dc:"vat number"`
	City           string `json:"city" dc:"city"`
	ZipCode        string `json:"zipCode" dc:"zip_code"`
	Custom         string `json:"custom" dc:"Custom"`
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId"`
	Language       string `json:"language" dc:"Language"`
}
