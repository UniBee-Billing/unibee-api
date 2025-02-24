package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipayCreatePayoutRequest struct {
	TransferRequestId  string                    `json:"transferRequestId,omitempty"`
	TransferFromDetail *model.TransferFromDetail `json:"transferFromDetail,omitempty"`
	TransferToDetail   *model.TransferToDetail   `json:"transferToDetail,omitempty"`
}

func NewAlipayCreatePayoutRequest() (*request.AlipayRequest, *AlipayCreatePayoutRequest) {
	alipayCreatePayoutRequest := &AlipayCreatePayoutRequest{}
	alipayRequest := request.NewAlipayRequest(alipayCreatePayoutRequest, model.MARKETPLACE_CREATEPAYOUT_PATH, &responseMarketplace.AlipayCreatePayoutResponse{})
	return alipayRequest, alipayCreatePayoutRequest
}
