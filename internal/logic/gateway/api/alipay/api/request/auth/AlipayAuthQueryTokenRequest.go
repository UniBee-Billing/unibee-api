package auth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
)

type AlipayAuthQueryTokenRequest struct {
	AccessToken string `json:"accessToken"`
}

func NewAlipayAuthQueryTokenRequest() (*request.AlipayRequest, *AlipayAuthQueryTokenRequest) {
	alipayAuthQueryTokenRequest := &AlipayAuthQueryTokenRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAuthQueryTokenRequest, model.AUTH_QUERY_PATH, &responseAuth.AlipayAuthQueryTokenResponse{})
	return alipayRequest, alipayAuthQueryTokenRequest
}
