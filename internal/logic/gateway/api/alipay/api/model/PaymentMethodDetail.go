package model

type PaymentMethodDetail struct {
	PaymentMethodDetailType PaymentMethodDetailType      `json:"paymentMethodDetailType,omitempty"`
	Card                    *CardPaymentMethodDetail     `json:"card,omitempty"`
	ExternalAccount         *ExternalPaymentMethodDetail `json:"externalAccount,omitempty"`
	Discount                *DiscountPaymentMethodDetail `json:"discount,omitempty"`
	Coupon                  *CouponPaymentMethodDetail   `json:"coupon,omitempty"`
	PaymentMethodType       string                       `json:"paymentMethodType,omitempty"`
	ExtendInfo              string                       `json:"extendInfo,omitempty"`
}

type PaymentMethodDetailType string

const (
	PaymentMethodDetailType_CARD            PaymentMethodDetailType = "CARD"
	PaymentMethodDetailType_EXTERNALACCOUNT PaymentMethodDetailType = "EXTERNALACCOUNT"
	PaymentMethodDetailType_COUPON          PaymentMethodDetailType = "COUPON"
	PaymentMethodDetailType_DISCOUNT        PaymentMethodDetailType = "DISCOUNT"
)

type CardPaymentMethodDetail struct {
	CardToken                   string    `json:"cardToken,omitempty"`
	CardNo                      string    `json:"cardNo,omitempty"`
	Brand                       CardBrand `json:"brand,omitempty"`
	SelectedCardBrand           CardBrand `json:"selectedCardBrand,omitempty"`
	CardIssuer                  string    `json:"cardIssuer,omitempty"`
	CountryIssue                string    `json:"countryIssue,omitempty"`
	CardHolderName              *UserName `json:"cardHolderName,omitempty"`
	ExpiryYear                  string    `json:"expiryYear,omitempty"`
	ExpiryMonth                 string    `json:"expiryMonth,omitempty"`
	BillingAddress              *Address  `json:"billingAddress,omitempty"`
	Mask                        string    `json:"mask,omitempty"`
	Last4                       string    `json:"last4,omitempty"`
	PaymentMethodDetailMetadata string    `json:"paymentMethodDetailMetadata,omitempty"`
	MaskedCardNo                string    `json:"maskedCardNo,omitempty"`
	Fingerprint                 string    `json:"fingerprint,omitempty"`
	AuthenticationFlow          string    `json:"authenticationFlow,omitempty"`
	Funding                     string    `json:"funding,omitempty"`
	AvsResultRaw                string    `json:"avsResultRaw,omitempty"`
	Bin                         string    `json:"bin,omitempty"`
	IssuerName                  string    `json:"issuerName,omitempty"`
	IssuingCountry              string    `json:"issuingCountry,omitempty"`
	LastFour                    string    `json:"lastFour,omitempty"`
	CardholderName              *UserName `json:"cardholderName,omitempty"`
	Cvv                         string    `json:"cvv,omitempty"`
	DateOfBirth                 string    `json:"dateOfBirth,omitempty"`
	BusinessNo                  string    `json:"businessNo,omitempty"`
	CardPasswordDigest          string    `json:"cardPasswordDigest,omitempty"`
	Cpf                         string    `json:"cpf,omitempty"`
	PayerEmail                  string    `json:"payerEmail,omitempty"`
	NetworkTransactionId        string    `json:"networkTransactionId,omitempty"`
	Is3DSAuthentication         bool      `json:"is3DSAuthentication,omitempty"`
}

type ExternalPaymentMethodDetail struct {
	AssetToken                  string `json:"assetToken,omitempty"`
	AccountDisplayName          string `json:"accountDisplayName,omitempty"`
	DisableReason               string `json:"disableReason,omitempty"`
	PaymentMethodDetailMetadata string `json:"paymentMethodDetailMetadata,omitempty"`
}

type DiscountPaymentMethodDetail struct {
	DiscountId                  string  `json:"discountId,omitempty"`
	AvailableAmount             *Amount `json:"availableAmount,omitempty"`
	DiscountName                string  `json:"discountName,omitempty"`
	DiscountDescription         string  `json:"discountDescription,omitempty"`
	PaymentMethodDetailMetadata string  `json:"paymentMethodDetailMetadata,omitempty"`
}

type CouponPaymentMethodDetail struct {
	CouponId                    string  `json:"couponId,omitempty"`
	AvailableAmount             *Amount `json:"availableAmount,omitempty"`
	CouponName                  string  `json:"couponName,omitempty"`
	CouponDescription           string  `json:"couponDescription,omitempty"`
	CouponExpireTime            string  `json:"couponExpireTime,omitempty"`
	PaymentMethodDetailMetadata string  `json:"paymentMethodDetailMetadata,omitempty"`
}
