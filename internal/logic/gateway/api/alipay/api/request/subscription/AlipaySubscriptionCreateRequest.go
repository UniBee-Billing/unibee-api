package subscription

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseSubscription "unibee/internal/logic/gateway/api/alipay/api/response/subscription"
)

type AlipaySubscriptionCreateRequest struct {
	SubscriptionRequestId       string                    `json:"subscriptionRequestId,omitempty"`
	SubscriptionDescription     string                    `json:"subscriptionDescription,omitempty"`
	SubscriptionRedirectUrl     string                    `json:"subscriptionRedirectUrl,omitempty"`
	SubscriptionStartTime       string                    `json:"subscriptionStartTime,omitempty"`
	SubscriptionEndTime         string                    `json:"subscriptionEndTime,omitempty"`
	PeriodRule                  *model.PeriodRule         `json:"periodRule,omitempty"`
	SubscriptionExpiryTime      string                    `json:"subscriptionExpiryTime,omitempty"`
	PaymentMethod               *model.PaymentMethod      `json:"paymentMethod,omitempty"`
	SubscriptionNotificationUrl string                    `json:"subscriptionNotificationUrl,omitempty"`
	PaymentNotificationUrl      string                    `json:"paymentNotificationUrl,omitempty"`
	OrderInfo                   *model.OrderInfo          `json:"orderInfo,omitempty"`
	PaymentAmount               *model.Amount             `json:"paymentAmount,omitempty"`
	SettlementStrategy          *model.SettlementStrategy `json:"settlementStrategy,omitempty"`
	Env                         *model.Env                `json:"env,omitempty"`
	Trials                      []*model.Trial            `json:"trials,omitempty"`
}

func (alipaySubscriptionCreateRequest *AlipaySubscriptionCreateRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipaySubscriptionCreateRequest, model.SUBSCRIPTION_CREATE_PATH, &responseSubscription.AlipaySubscriptionCreateResponse{})
}

func NewAlipaySubscriptionCreateRequest() (*request.AlipayRequest, *AlipaySubscriptionCreateRequest) {
	alipaySubscriptionCreateRequest := &AlipaySubscriptionCreateRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySubscriptionCreateRequest, model.SUBSCRIPTION_CREATE_PATH, &responseSubscription.AlipaySubscriptionCreateResponse{})
	return alipayRequest, alipaySubscriptionCreateRequest
}
