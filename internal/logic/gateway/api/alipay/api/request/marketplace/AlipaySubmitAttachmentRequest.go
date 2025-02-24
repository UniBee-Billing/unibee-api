package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipaySubmitAttachmentRequest struct {
	SubmitAttachmentRequestId string               `json:"submitAttachmentRequestId,omitempty"`
	AttachmentType            model.AttachmentType `json:"attachmentType,omitempty"`
	FileSha256                string               `json:"fileSha256,omitempty"`
}

func NewAlipaySubmitAttachmentRequest() (*request.AlipayRequest, *AlipaySubmitAttachmentRequest) {
	alipaySubmitAttachmentRequest := &AlipaySubmitAttachmentRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySubmitAttachmentRequest, model.MARKETPLACE_SUBMITATTACHMENT_PATH, &responseMarketplace.AlipaySubmitAttachmentResponse{})
	return alipayRequest, alipaySubmitAttachmentRequest
}
