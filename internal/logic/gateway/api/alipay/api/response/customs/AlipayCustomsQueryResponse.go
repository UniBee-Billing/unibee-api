package responseCustoms

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayCustomsQueryResponse struct {
	response.AlipayResponse
	DeclarationRequestsNotFound []string                  `json:"declarationRequestsNotFound"`
	DeclarationRecords          []model.DeclarationRecord `json:"declarationRecords"`
}
