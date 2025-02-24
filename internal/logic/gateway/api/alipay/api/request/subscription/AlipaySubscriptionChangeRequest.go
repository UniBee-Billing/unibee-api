package subscription

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseSubscription "unibee/internal/logic/gateway/api/alipay/api/response/subscription"
)

type AlipaySubscriptionChangeRequest struct {
	SubscriptionChangeRequestId string            `json:"subscriptionChangeRequestId,omitempty"`
	SubscriptionId              string            `json:"subscriptionId,omitempty"`
	SubscriptionDescription     string            `json:"subscriptionDescription,omitempty"`
	SubscriptionStartTime       string            `json:"subscriptionStartTime,omitempty"`
	SubscriptionEndTime         string            `json:"subscriptionEndTime,omitempty"`
	PeriodRule                  *model.PeriodRule `json:"periodRule,omitempty"`
	SubscriptionExpiryTime      string            `json:"subscriptionExpiryTime,omitempty"`
	OrderInfo                   *model.OrderInfo  `json:"orderInfo,omitempty"`
	PaymentAmount               *model.Amount     `json:"paymentAmount,omitempty"`
	PaymentAmountDifference     *model.Amount     `json:"paymentAmountDifference,omitempty"`
}

func (alipaySubscriptionChangeRequest *AlipaySubscriptionChangeRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipaySubscriptionChangeRequest, model.SUBSCRIPTION_CHANGE_PATH, &responseSubscription.AlipaySubscriptionChangeResponse{})
}

func NewAlipaySubscriptionChangeRequest() (*request.AlipayRequest, *AlipaySubscriptionChangeRequest) {
	alipaySubscriptionChangeRequest := &AlipaySubscriptionChangeRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySubscriptionChangeRequest, model.SUBSCRIPTION_CHANGE_PATH, &responseSubscription.AlipaySubscriptionChangeResponse{})
	return alipayRequest, alipaySubscriptionChangeRequest
}
