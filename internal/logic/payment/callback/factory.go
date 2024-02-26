package callback

import (
	"context"
	_interface "unibee/internal/interface"
)

func GetPaymentCallbackServiceProvider(ctx context.Context, bizType int) (one _interface.PaymentBizCallbackInterface) {
	proxy := &proxy{}
	proxy.BizType = bizType
	return proxy
}
