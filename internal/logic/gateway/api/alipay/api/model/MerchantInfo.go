package model

type MerchantInfo struct {
	ReferenceMerchantId string                `json:"referenceMerchantId,omitempty"`
	LoginId             string                `json:"loginId,omitempty"`
	LegalEntityType     LegalEntityType       `json:"legalEntityType,omitempty"`
	Company             *Company              `json:"company,omitempty"`
	BusinessInfo        *BusinessInfo         `json:"businessInfo,omitempty"`
	EntityAssociations  []*EntityAssociations `json:"entityAssociations,omitempty"`
}

type LegalEntityType string

const (
	LegalEntityType_Company    LegalEntityType = "COMPANY"
	LegalEntityType_INDIVIDUAL LegalEntityType = "INDIVIDUAL"
)

type Company struct {
	LegalName         string           `json:"legalName,omitempty"`
	CompanyType       CompanyType      `json:"companyType,omitempty"`
	RegisteredAddress *Address         `json:"registeredAddress,omitempty"`
	OperatingAddress  *Address         `json:"operatingAddress,omitempty"`
	IncorporationDate string           `json:"incorporationDate,omitempty"`
	StockInfo         *StockInfo       `json:"stockInfo,omitempty"`
	Certificates      *Certificate     `json:"certificates,omitempty"`
	Attachments       []*Attachment    `json:"attachments,omitempty"`
	CompanyUnit       *CompanyUnitType `json:"companyUnit,omitempty"`
	Contacts          []*Contact       `json:"contacts,omitempty"`
	VatNo             string           `json:"vatNo,omitempty"`
}

type BusinessInfo struct {
	Mcc                string     `json:"mcc,omitempty"`
	Websites           []*WebSite `json:"websites,omitempty"`
	EnglishName        string     `json:"englishName,omitempty"`
	DoingBusinessAs    string     `json:"doingBusinessAs,omitempty"`
	MainSalesCountry   string     `json:"mainSalesCountry,omitempty"`
	AppName            string     `json:"appName,omitempty"`
	ServiceDescription string     `json:"serviceDescription,omitempty"`
}

type EntityAssociations struct {
	AssociationType   AssociationType `json:"associationType,omitempty"`
	LegalEntityType   LegalEntityType `json:"legalEntityType,omitempty"`
	Company           *Company        `json:"company,omitempty"`
	Individual        *Individual     `json:"individual,omitempty"`
	ShareholdingRatio string          `json:"shareholdingRatio,omitempty"`
}

type Individual struct {
	Name         *UserName    `json:"name,omitempty"`
	EnglishName  *UserName    `json:"englishName,omitempty"`
	DateOfBirth  string       `json:"dateOfBirth,omitempty"`
	PlaceOfBirth *Address     `json:"placeOfBirth,omitempty"`
	Certificates *Certificate `json:"certificates,omitempty"`
	Nationality  string       `json:"nationality,omitempty"`
	Contacts     []*Contact   `json:"contacts,omitempty"`
}

type AssociationType string

const (
	AssociationType_LEGAL_REPRESENTATIVE AssociationType = "LEGAL_REPRESENTATIVE"
	AssociationType_UBO                  AssociationType = "UBO"
	AssociationType_CONTACT              AssociationType = "CONTACT"
	AssociationType_DIRECTOR             AssociationType = "DIRECTOR"
	AssociationType_AUTHORIZER           AssociationType = "AUTHORIZER"
	AssociationType_BOARD_MEMBER         AssociationType = "BOARD_MEMBER"
)

type WebSite struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
	Desc string `json:"desc,omitempty"`
	Type string `json:"type,omitempty"`
}

type StockInfo struct {
	ListedRegion string `json:"listedRegion,omitempty"`
	TickerSymbol string `json:"tickerSymbol,omitempty"`
}

type Attachment struct {
	AttachmentType AttachmentType `json:"attachmentType,omitempty"`
	File           string         `json:"file,omitempty"`
	AttachmentName string         `json:"attachmentName,omitempty"`
	FileKey        string         `json:"fileKey,omitempty"`
}

type CompanyUnitType string

const (
	CompanyUnitType_HEADQUARTER CompanyUnitType = "HEADQUARTER"
	CompanyUnitType_BRANCH      CompanyUnitType = "BRANCH"
)

type CompanyType string

const (
	CompanyType_ENTERPRISE               CompanyType = "ENTERPRISE"
	CompanyType_SOLE_PROPRIETORSHIP      CompanyType = "SOLE_PROPRIETORSHIP"
	CompanyType_PARTNERSHIP              CompanyType = "PARTNERSHIP"
	CompanyType_STATE_OWNED_BUSINESS     CompanyType = "STATE_OWNED_BUSINESS"
	CompanyType_PRIVATELY_OWNED_BUSINESS CompanyType = "PRIVATELY_OWNED_BUSINESS"
	CompanyType_PUBLICLY_LISTED_BUSINESS CompanyType = "PUBLICLY_LISTED_BUSINESS"
	CompanyType_LTDA                     CompanyType = "LTDA"
	CompanyType_SA                       CompanyType = "SA"
	CompanyType_EIRELI                   CompanyType = "EIRELI"
	CompanyType_BOFC                     CompanyType = "BOFC"
	CompanyType_MEI                      CompanyType = "MEI"
	CompanyType_EI                       CompanyType = "EI"
)

type Contact struct {
	Type ContactType `json:"type,omitempty"`
	Info string      `json:"info,omitempty"`
}

type ContactType string

const (
	ContactType_EMAIL               ContactType = "EMAIL"
	ContactType_PHONE_NO            ContactType = "PHONE_NO"
	ContactType_COMMERCIAL_PHONE_NO ContactType = "COMMERCIAL_PHONE_NO"
)
