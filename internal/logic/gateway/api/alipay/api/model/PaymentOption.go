package model

type PaymentOption struct {
	PaymentMethodType     string                     `json:"paymentMethodType,omitempty"`
	PaymentMethodCategory PaymentMethodCategoryType  `json:"paymentMethodCategory,omitempty"`
	PaymentMethodRegion   []string                   `json:"paymentMethodRegion,omitempty"`
	Enabled               bool                       `json:"enabled,omitempty"`
	Preferred             bool                       `json:"preferred,omitempty"`
	DisableReason         DisableReasonType          `json:"disableReason,omitempty"`
	AmountLimitInfoMap    map[string]AmountLimitInfo `json:"amountLimitInfoMap,omitempty"`
	SupportedCurrencies   []string                   `json:"supportedCurrencies,omitempty"`
	PaymentOptionDetail   PaymentOptionDetail        `json:"paymentOptionDetail,omitempty"`
	ExtendInfo            string                     `json:"extendInfo,omitempty"`
	Logo                  Logo                       `json:"logo,omitempty"`
	PromoNames            []string                   `json:"promoNames,omitempty"`
	Installment           Installment                `json:"installment,omitempty"`
	PromotionInfos        []PromotionInfo            `json:"promotionInfos,omitempty"`
}

type PaymentMethodCategoryType string

const (
	//ALIPAY_PLUS, WALLET, MOBILE_BANKING_APP, BANK_TRANSFER, ONLINE_BANKING, CARD, OTC;
	PaymentMethodCategoryTypeAlipayPlus       PaymentMethodCategoryType = "ALIPAY_PLUS"
	PaymentMethodCategoryTypeWallet           PaymentMethodCategoryType = "WALLET"
	PaymentMethodCategoryTypeMOBILEBANKINGAPP PaymentMethodCategoryType = "MOBILE_BANKING_APP"
	PaymentMethodCategoryTypeBankTransfer     PaymentMethodCategoryType = "BANK_TRANSFER"
	PaymentMethodCategoryTypeOnlineBanking    PaymentMethodCategoryType = "ONLINE_BANKING"
	PaymentMethodCategoryTypeCard             PaymentMethodCategoryType = "CARD"
	PaymentMethodCategoryTypeOTC              PaymentMethodCategoryType = "OTC"
)

type AmountLimitInfo struct {
	SingleLimit AmountLimit `json:"singleLimit,omitempty"`
	DayLimit    AmountLimit `json:"dayLimit,omitempty"`
	MonthLimit  AmountLimit `json:"monthLimit,omitempty"`
}
type AmountLimit struct {
	MinAmount    Amount `json:"minAmount,omitempty"`
	MaxAmount    Amount `json:"maxAmount,omitempty"`
	RemainAmount Amount `json:"remainAmount,omitempty"`
}

type PaymentOptionDetail struct {
	SupportCardBrands []SupportCardBrand `json:"supportCardBrands,omitempty"`
	Funding           []string           `json:"funding,omitempty"`
	SupportBanks      []SupportBank      `json:"supportBanks,omitempty"`
}

type SupportBank struct {
	BankIdentifierCode string `json:"bankIdentifierCode,omitempty"`
	BankShortName      string `json:"bankShortName,omitempty"`
	BankLogo           Logo   `json:"bankLogo,omitempty"`
}

type Installment struct {
	SupportCardBrands []SupportCardBrand `json:"supportCardBrands,omitempty"`
	Plans             []Plan             `json:"plans,omitempty"`
}

type Plan struct {
	InterestRate         string `json:"interestRate,omitempty"`
	MinInstallmentAmount Amount `json:"minInstallmentAmount,omitempty"`
	MaxInstallmentAmount Amount `json:"maxInstallmentAmount,omitempty"`
	InstallmentNum       string `json:"installmentNum,omitempty"`
	Interval             string `json:"interval,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
}
type PromotionInfo struct {
	PromotionType PromotionType `json:"promotionType,omitempty"`
	Discount      Discount      `json:"discount,omitempty"`
	InterestFree  InterestFree  `json:"interestFree,omitempty"`
}

type InterestFree struct {
	Provider            string `json:"provider,omitempty"`
	ExpireTime          string `json:"expireTime,omitempty"`
	InstallmentFreeNums []int  `json:"installmentFreeNums,omitempty"`
	MinPaymentAmount    Amount `json:"minPaymentAmount,omitempty"`
	MaxPaymentAmount    Amount `json:"maxPaymentAmount,omitempty"`
	FreePercentage      int    `json:"freePercentage,omitempty"`
}
