package out

import (
	"context"
	"go-oversea-pay/api/out/vo"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/service"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

func currencyNumberCheck(amount *vo.PayAmountVo) {
	utility.Assert(amount != nil, "amount is nil")
	if strings.Compare(amount.Currency, "JPY") == 0 {
		utility.Assert(amount.Value%100 == 0, "this currency No decimals allowed，made it divisible by 100")
	}
}

func merchantCheck(ctx context.Context, merchantAccount string) (apiconfig *entity.OpenApiConfig, res *entity.MerchantInfo) {
	openApiConfig := service.BizCtx().Get(ctx).OpenApiConfig
	utility.Assert(openApiConfig != nil, "api config not found")
	utility.Assert(strings.Compare(strconv.FormatInt(openApiConfig.MerchantId, 10), merchantAccount) == 0, "api config not found")
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
