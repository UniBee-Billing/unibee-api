package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayDisputeNotify struct {
	AlipayNotify
	PaymentRequestId        string                        `json:"paymentRequestId,omitempty"`
	DisputeId               string                        `json:"disputeId,omitempty"`
	PaymentId               string                        `json:"paymentId,omitempty"`
	DisputeTime             string                        `json:"disputeTime,omitempty"`
	DisputeAmount           *model.Amount                 `json:"disputeAmount,omitempty"`
	DisputeNotificationType model.DisputeNotificationType `json:"disputeNotificationType,omitempty"`
	DisputeReasonMsg        string                        `json:"disputeReasonMsg,omitempty"`
	DisputeJudgedTime       string                        `json:"disputeJudgedTime,omitempty"`
	DisputeJudgedAmount     *model.Amount                 `json:"disputeJudgedAmount,omitempty"`
	DisputeJudgedResult     model.DisputeJudgedResult     `json:"disputeJudgedResult,omitempty"`
	DefenseDueTime          string                        `json:"defenseDueTime,omitempty"`
	DisputeReasonCode       string                        `json:"disputeReasonCode,omitempty"`
	DisputeSource           string                        `json:"disputeSource,omitempty"`
	Rrn                     string                        `json:"arn,omitempty"`
	DisputeAcceptReason     model.DisputeAcceptReasonType `json:"disputeAcceptReason,omitempty"`
	DisputeAcceptTime       string                        `json:"disputeAcceptTime,omitempty"`
	DisputeType             string                        `json:"disputeType,omitempty"`
	Defendable              bool                          `json:"defendable,omitempty"`
}
