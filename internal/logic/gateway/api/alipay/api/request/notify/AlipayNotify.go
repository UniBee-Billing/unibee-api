package notify

import (
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayNotify struct {
	NotifyType string          `json:"notifyType,omitempty"`
	Result     response.Result `json:"result,omitempty"`
}
