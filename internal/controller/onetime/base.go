package onetime

import (
	"context"
	"strings"
	"unibee/api/onetime/payment"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func currencyNumberCheck(amount *payment.AmountVo) {
	utility.Assert(amount != nil, "amount is nil")
	if strings.Compare(amount.Currency, "JPY") == 0 {
		utility.Assert(amount.Amount%100 == 0, "this currency No decimals allowed，made it divisible by 100")
	}
}

func merchantCheck(ctx context.Context, merchantId uint64) (apiconfig *entity.OpenApiConfig, res *entity.Merchant) {
	openApiConfig := _interface.BizCtx().Get(ctx).OpenApiConfig
	utility.Assert(openApiConfig != nil, "api config not found")
	utility.Assert(openApiConfig.MerchantId == merchantId, "api config not found")
	utility.Assert(openApiConfig.MerchantId > 0, "api config not found")
	err := dao.Merchant.Ctx(ctx).Where(entity.Merchant{Id: openApiConfig.MerchantId}).OmitEmpty().Scan(&res)
	if err != nil {
		return openApiConfig, res
	}
	return openApiConfig, res
}

//func convertToGJson(target interface{}) (res *gjson.PortalJson) {
//	resBytes, err := gjson.Marshal(target)
//	if err != nil {
//		return nil
//	}
//	return string(resBytes)
//}
