package customs

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseCustoms "unibee/internal/logic/gateway/api/alipay/api/response/customs"
)

type AlipayCustomsQueryRequest struct {
	DeclarationRequestIds []string `json:"declarationRequestIds,omitempty"`
}

func (alipayCustomsQueryRequest *AlipayCustomsQueryRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayCustomsQueryRequest, model.INQUIRY_DECLARE_PATH, &responseCustoms.AlipayCustomsQueryResponse{})
}

func NewAlipayCustomsQueryRequest() (*request.AlipayRequest, *AlipayCustomsQueryRequest) {
	alipayCustomsQueryRequest := &AlipayCustomsQueryRequest{}
	alipayRequest := request.NewAlipayRequest(alipayCustomsQueryRequest, model.INQUIRY_DECLARE_PATH, &responseCustoms.AlipayCustomsQueryResponse{})
	return alipayRequest, alipayCustomsQueryRequest
}
