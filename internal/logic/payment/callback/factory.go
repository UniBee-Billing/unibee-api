package callback

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
)

func GetPaymentCallbackServiceProvider(ctx context.Context, bizType int) (channelService _interface.PaymentBizCallbackInterface) {
	proxy := &proxy{}
	proxy.BizType = bizType
	return proxy
}
