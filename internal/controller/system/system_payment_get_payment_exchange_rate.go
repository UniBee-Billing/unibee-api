package system

import (
	"context"
	"unibee/internal/logic/multi_currencies/currency_exchange"

	"unibee/api/system/payment"
)

func (c *ControllerPayment) GetPaymentExchangeRate(ctx context.Context, req *payment.GetPaymentExchangeRateReq) (res *payment.GetPaymentExchangeRateRes, err error) {
	cloud, err := currency_exchange.GetExchangeConversionRateFromClusterCloud(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, err
	}
	return &payment.GetPaymentExchangeRateRes{ExchangeRate: *cloud}, nil
}
