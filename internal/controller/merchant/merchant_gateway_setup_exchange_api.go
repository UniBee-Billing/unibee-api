package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/fiat_exchange"
	"unibee/internal/logic/merchant_config/update"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) SetupExchangeApi(ctx context.Context, req *gateway.SetupExchangeApiReq) (res *gateway.SetupExchangeApiRes, err error) {
	if len(req.ExchangeRateApiKey) > 0 {
		rate, err := currency.GetExchangeConversionRates(ctx, req.ExchangeRateApiKey, "USD", "EUR")
		utility.AssertError(err, "invalid exchange api key")
		utility.Assert(rate != nil, "invalid exchange api key")
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), fiat_exchange.FiatExchangeApiKey, req.ExchangeRateApiKey)
		if err != nil {
			return nil, err
		}
	}
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
