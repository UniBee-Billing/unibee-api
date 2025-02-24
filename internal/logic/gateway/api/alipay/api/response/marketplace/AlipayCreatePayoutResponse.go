package responseMarketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayCreatePayoutResponse struct {
	response.AlipayResponse
	TransferId         string                    `json:"transferId,omitempty"`
	TransferRequestId  string                    `json:"transferRequestId,omitempty"`
	TransferFromDetail *model.TransferFromDetail `json:"transferFromDetail,omitempty"`
	TransferToDetail   *model.TransferToDetail   `json:"transferToDetail,omitempty"`
}
