package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) NewConfig(ctx context.Context, req *credit.NewConfigReq) (res *credit.NewConfigRes, err error) {
	utility.Assert(len(req.Currency) > 0, "invalid currency")
	utility.Assert(currency.IsCurrencySupport(req.Currency), "invalid currency")
	utility.Assert(req.Type == consts.CreditAccountTypeMain || req.Type == consts.CreditAccountTypePromo, "invalid type, should be 1-main account, 2-promo credit account")
	utility.Assert(req.ExchangeRate > 0, "invalid exchange rate")
	utility.Assert(query.GetCreditConfig(ctx, _interface.GetMerchantId(ctx), req.Type, req.Currency) == nil, "already exist config with currency")
	utility.Assert(req.RechargeEnable != nil, "rechargeEnable should setup")
	utility.Assert(req.PayoutEnable != nil, "payoutEnable should setup")
	utility.Assert(req.PreviewDefaultUsed != nil, "previewDefaultUsed should setup")
	if req.Type == 1 && req.ExchangeRate != 100 {
		req.ExchangeRate = 100
	}
	one := &entity.CreditConfig{
		Type:                  req.Type,
		Currency:              strings.ToUpper(req.Currency),
		ExchangeRate:          req.ExchangeRate,
		CreateTime:            gtime.Now().Timestamp(),
		MerchantId:            _interface.GetMerchantId(ctx),
		Recurring:             req.Recurring,
		DiscountCodeExclusive: req.DiscountCodeExclusive,
		Logo:                  req.Logo,
		Name:                  req.Name,
		Description:           req.Description,
		LogoUrl:               req.LogoUrl,
		MetaData:              unibee.StringValue(utility.MarshalMetadataToJsonString(req.MetaData)),
		RechargeEnable:        *req.RechargeEnable,
		PayoutEnable:          *req.PayoutEnable,
		PreviewDefaultUsed:    *req.PreviewDefaultUsed,
	}
	result, err := dao.CreditConfig.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`create credit config record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("CreditConfig(%d)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &credit.NewConfigRes{CreditConfig: bean.SimplifyCreditConfig(one)}, nil
}
