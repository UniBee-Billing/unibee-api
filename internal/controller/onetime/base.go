package onetime

import (
	"context"
	"strings"
	"unibee-api/api/onetime/payment"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/interface"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

func currencyNumberCheck(amount *payment.AmountVo) {
	utility.Assert(amount != nil, "amount is nil")
	if strings.Compare(amount.Currency, "JPY") == 0 {
		utility.Assert(amount.Amount%100 == 0, "this currency No decimals allowed，made it divisible by 100")
	}
}

func merchantCheck(ctx context.Context, merchantId uint64) (apiconfig *entity.OpenApiConfig, res *entity.MerchantInfo) {
	openApiConfig := _interface.BizCtx().Get(ctx).OpenApiConfig
	utility.Assert(openApiConfig != nil, "api config not found")
	utility.Assert(openApiConfig.MerchantId == merchantId, "api config not found")
	utility.Assert(openApiConfig.MerchantId > 0, "api config not found")
	err := dao.MerchantInfo.Ctx(ctx).Where(entity.MerchantInfo{Id: openApiConfig.MerchantId}).OmitEmpty().Scan(&res)
	if err != nil {
		return openApiConfig, res
	}
	return openApiConfig, res
}

//func convertToGJson(target interface{}) (res *gjson.Json) {
//	resBytes, err := gjson.Marshal(target)
//	if err != nil {
//		return nil
//	}
//	return string(resBytes)
//}
