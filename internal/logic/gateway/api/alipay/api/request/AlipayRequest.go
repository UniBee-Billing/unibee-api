package request

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayRequest struct {
	Param          any    `json:"param,omitempty"`
	ClientId       string `json:"clientId,omitempty"`
	Path           string `json:"path,omitempty"`
	KeyVersion     string `json:"keyVersion,omitempty"`
	HttpMethod     string `json:"httpMethod,omitempty"`
	AlipayResponse any    `json:"response,omitempty"`
}

func NewAlipayRequest(param any, path string, alipayResponse any) *AlipayRequest {
	return &AlipayRequest{
		Param:          param,
		ClientId:       "",
		Path:           path,
		KeyVersion:     "1",
		HttpMethod:     model.MethodPost,
		AlipayResponse: alipayResponse,
	}
}

func NEWAlipayRequest(param any, clientId string, path string, keyVersion string, httpMethod string, alipayResponse any) *AlipayRequest {
	return &AlipayRequest{
		Param:          param,
		ClientId:       clientId,
		Path:           path,
		KeyVersion:     keyVersion,
		HttpMethod:     httpMethod,
		AlipayResponse: alipayResponse,
	}
}
