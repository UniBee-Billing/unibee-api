package bean

type VatCountryRate struct {
	Id                    uint64 `json:"id"`
	Gateway               string `json:"gateway"           `        // gateway
	CountryCode           string `json:"countryCode"           `    // country_code
	CountryName           string `json:"countryName"           `    // country_name
	VatSupport            bool   `json:"vatSupport"               ` // vat support true or false
	IsEU                  bool   `json:"isEU"                  `
	StandardTaxPercentage int64  `json:"standardTaxPercentage" `
}

type ValidResult struct {
	Valid           bool   `json:"valid"           `
	VatNumber       string `json:"vatNumber"           `
	CountryCode     string `json:"countryCode"           `
	CompanyName     string `json:"companyName"           `
	CompanyAddress  string `json:"companyAddress"           `
	ValidateMessage string `json:"validateMessage"           `
}
