package responseAuth

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayAuthRevokeTokenResponse struct {
	response.AlipayResponse
	ExtendInfo string `json:"extendInfo"`
}
