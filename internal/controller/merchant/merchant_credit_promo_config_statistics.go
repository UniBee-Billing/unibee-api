package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) PromoConfigStatistics(ctx context.Context, req *credit.PromoConfigStatisticsReq) (res *credit.PromoConfigStatisticsRes, err error) {
	if len(req.Currency) <= 0 {
		return &credit.PromoConfigStatisticsRes{CreditConfigStatistics: &bean.CreditConfigStatistics{
			TotalDecrementAmount: 0,
			TotalIncrementAmount: 0,
		}}, nil
	}
	return &credit.PromoConfigStatisticsRes{CreditConfigStatistics: &bean.CreditConfigStatistics{
		TotalDecrementAmount: int64(bean.GetCreditConfigTotalDecrementAmount(ctx, _interface.GetMerchantId(ctx), req.Currency)),
		TotalIncrementAmount: int64(bean.GetCreditConfigTotalIncrementAmount(ctx, _interface.GetMerchantId(ctx), req.Currency)),
	}}, nil
}
