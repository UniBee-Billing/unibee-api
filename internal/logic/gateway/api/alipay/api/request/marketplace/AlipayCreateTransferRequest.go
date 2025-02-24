package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipayCreateTransferRequest struct {
	TransferRequestId  string                    `json:"transferRequestId,omitempty"`
	TransferFromDetail *model.TransferFromDetail `json:"transferFromDetail,omitempty"`
	TransferToDetail   *model.TransferToDetail   `json:"transferToDetail,omitempty"`
}

func NewAlipayCreateTransferRequest() (*request.AlipayRequest, *AlipayCreateTransferRequest) {
	alipayCreateTransferRequest := &AlipayCreateTransferRequest{}
	alipayRequest := request.NewAlipayRequest(alipayCreateTransferRequest, model.MARKETPLACE_CREATETRANSFER_PATH, &responseMarketplace.AlipayCreateTransferResponse{})
	return alipayRequest, alipayCreateTransferRequest
}
