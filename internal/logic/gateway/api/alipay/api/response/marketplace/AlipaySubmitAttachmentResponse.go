package responseMarketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipaySubmitAttachmentResponse struct {
	response.AlipayResponse
	SubmitAttachmentRequestId string               `json:"submitAttachmentRequestId,omitempty"`
	AttachmentType            model.AttachmentType `json:"attachmentType,omitempty"`
	AttachmentKey             string               `json:"attachmentKey,omitempty"`
}
