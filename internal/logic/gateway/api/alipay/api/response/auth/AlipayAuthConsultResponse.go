package responseAuth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayAuthConsultResponse struct {
	response.AlipayResponse
	AuthUrl       string             `json:"authUrl"`
	ExtendInfo    string             `json:"extendInfo"`
	NormalUrl     string             `json:"normalUrl"`
	SchemeUrl     string             `json:"schemeUrl"`
	ApplinkUrl    string             `json:"applinkUrl"`
	AppIdentifier string             `json:"appIdentifier"`
	AuthCodeForm  model.AuthCodeForm `json:"authCodeForm"`
	GrantType     string             `json:"grant_type"`
}
