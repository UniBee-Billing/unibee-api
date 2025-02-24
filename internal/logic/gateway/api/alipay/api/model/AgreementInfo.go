package model

type AgreementInfo struct {
	AuthState          string `json:"authState,omitempty"`
	UserLoginId        string `json:"userLoginId,omitempty"`
	UserLoginType      string `json:"userLoginType,omitempty"`
	DisplayUserLoginId string `json:"displayUserLoginId,omitempty"`
}

type RiskData struct {
	Order                  *RiskOrder              `json:"order,omitempty"`
	Buyer                  *RiskBuyer              `json:"buyer,omitempty"`
	Env                    *RiskEnv                `json:"env,omitempty"`
	RiskSignal             *RiskSignal             `json:"riskSignal,omitempty"`
	Address                *RiskAddress            `json:"address,omitempty"`
	CardVerificationResult *CardVerificationResult `json:"cardVerificationResult,omitempty"`
}

type RiskOrder struct {
	OrderType     string `json:"orderType,omitempty"`
	ReferringSite string `json:"referringSite,omitempty"`
}

type RiskBuyer struct {
	NoteToMerchant  string `json:"noteToMerchant,omitempty"`
	NoteToShipping  string `json:"noteToShipping,omitempty"`
	OrderCountIn1H  string `json:"orderCountIn1H,omitempty"`
	OrderCountIn24H string `json:"orderCountIn24H,omitempty"`
}

type RiskEnv struct {
	IpAddressType string `json:"ipAddressType,omitempty"`
}

type RiskSignal struct {
	RiskCode   string `json:"riskCode,omitempty"`
	RiskReason string `json:"riskReason,omitempty"`
}

type RiskAddress struct {
	ShippingPhoneType             string `json:"shippingPhoneType,omitempty"`
	IsBillShipStateSame           bool   `json:"isBillShipStateSame,omitempty"`
	IsPreviousStateSame           bool   `json:"isPreviousStateSame,omitempty"`
	LocToShipDistance             int    `json:"locToShipDistance,omitempty"`
	MinPreviousShipToBillDistance int    `json:"minPreviousShipToBillDistance,omitempty"`
}

type CardVerificationResult struct {
	AuthenticationType   string             `json:"authenticationType,omitempty"`
	AuthenticationMethod string             `json:"authenticationMethod,omitempty"`
	CvvResult            string             `json:"cvvResult,omitempty"`
	AvsResult            string             `json:"avsResult,omitempty"`
	AuthorizationCode    string             `json:"authorizationCode,omitempty"`
	ThreeDSResult        *RiskThreeDSResult `json:"threeDSResult,omitempty"`
}

type RiskThreeDSResult struct {
	ThreeDSVersion         string `json:"threeDSVersion,omitempty"`
	ThreeDSInteractionMode string `json:"threeDSInteractionMode,omitempty"`
	Eci                    string `json:"eci,omitempty"`
	Cavv                   string `json:"cavv,omitempty"`
}
