package responseAuth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayAuthApplyTokenResponse struct {
	response.AlipayResponse
	AccessToken            string                `json:"accessToken"`
	AccessTokenExpiryTime  string                `json:"accessTokenExpiryTime"`
	RefreshToken           string                `json:"refreshToken"`
	RefreshTokenExpiryTime string                `json:"refreshTokenExpiryTime"`
	ExtendInfo             string                `json:"extendInfo"`
	UserLoginId            string                `json:"userLoginId"`
	PspCustomerInfo        model.PspCustomerInfo `json:"pspCustomerInfo"`
	GrantType              string                `json:"grant_type"`
}
