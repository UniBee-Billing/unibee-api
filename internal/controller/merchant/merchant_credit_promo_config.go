package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/credit"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/currency"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func (c *ControllerCredit) PromoConfig(ctx context.Context, req *credit.PromoConfigReq) (res *credit.PromoConfigRes, err error) {
	utility.Assert(len(req.Currency) > 0, "Invalid currency")
	utility.Assert(currency.IsCurrencySupport(req.Currency), "invalid currency")
	var one *entity.CreditConfig
	query := dao.CreditConfig.Ctx(ctx).
		Where(dao.CreditConfig.Columns().MerchantId, _interface.GetMerchantId(ctx)).
		Where(dao.CreditConfig.Columns().Type, consts.CreditAccountTypePromo)
	query = query.Where(dao.CreditConfig.Columns().Currency, req.Currency)
	_ = query.Scan(&one)

	if one == nil {
		return &credit.PromoConfigRes{}, nil
		//one = &entity.CreditConfig{
		//	Type:                  consts.CreditAccountTypePromo,
		//	Currency:              strings.ToUpper(req.Currency),
		//	ExchangeRate:          0,
		//	CreateTime:            gtime.Now().Timestamp(),
		//	MerchantId:            _interface.GetMerchantId(ctx),
		//	Recurring:             0,
		//	DiscountCodeExclusive: 0,
		//	Logo:                  "",
		//	Name:                  "",
		//	Description:           "",
		//	LogoUrl:               "",
		//	MetaData:              "",
		//	RechargeEnable:        1,
		//	PayoutEnable:          0,
		//	PreviewDefaultUsed:    0,
		//}
		//result, err := dao.CreditConfig.Ctx(ctx).Data(one).OmitNil().Insert(one)
		//if err != nil {
		//	return nil, gerror.Newf(`create credit config record insert failure %s`, err)
		//}
		//id, _ := result.LastInsertId()
		//one.Id = uint64(id)
		//operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		//	MerchantId:     one.MerchantId,
		//	Target:         fmt.Sprintf("CreditConfig(%d)", one.Id),
		//	Content:        "SystemAutoSetup",
		//	UserId:         0,
		//	SubscriptionId: "",
		//	InvoiceId:      "",
		//	PlanId:         0,
		//	DiscountCode:   "",
		//}, err)
	}

	return &credit.PromoConfigRes{CreditConfig: bean.SimplifyCreditConfig(one)}, nil
}
