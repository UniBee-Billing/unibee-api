package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/api/merchant/credit"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerCredit) EditPromoConfig(ctx context.Context, req *credit.EditPromoConfigReq) (res *credit.EditPromoConfigRes, err error) {
	one := query.GetCreditConfig(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, req.Currency)
	utility.Assert(one != nil, "Config not found, please setup")
	//if one == nil {
	//	one = &entity.CreditConfig{
	//		Type:                  consts.CreditAccountTypePromo,
	//		Currency:              strings.ToUpper(req.Currency),
	//		ExchangeRate:          *req.ExchangeRate,
	//		CreateTime:            gtime.Now().Timestamp(),
	//		MerchantId:            _interface.GetMerchantId(ctx),
	//		Recurring:             *req.Recurring,
	//		DiscountCodeExclusive: *req.DiscountCodeExclusive,
	//		Logo:                  *req.Logo,
	//		Name:                  *req.Name,
	//		Description:           *req.Description,
	//		LogoUrl:               *req.LogoUrl,
	//		MetaData:              unibee.StringValue(utility.MarshalMetadataToJsonString(req.MetaData)),
	//		RechargeEnable:        0,
	//		PayoutEnable:          *req.PayoutEnable,
	//		PreviewDefaultUsed:    *req.PreviewDefaultUsed,
	//	}
	//	result, err := dao.CreditConfig.Ctx(ctx).Data(one).OmitNil().Insert(one)
	//	if err != nil {
	//		return nil, gerror.Newf(`create credit config record failure %s`, err)
	//	}
	//	id, _ := result.LastInsertId()
	//	one.Id = uint64(id)
	//	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
	//		MerchantId:     one.MerchantId,
	//		Target:         fmt.Sprintf("CreditConfig(%d)", one.Id),
	//		Content:        "New",
	//		UserId:         0,
	//		SubscriptionId: "",
	//		InvoiceId:      "",
	//		PlanId:         0,
	//		DiscountCode:   "",
	//	}, err)
	//}
	//if req.ExchangeRate != nil {
	//	utility.Assert(*req.ExchangeRate == one.ExchangeRate, "ExchangeRate can't change after setup")
	//}
	_, err = dao.CreditConfig.Ctx(ctx).Data(g.Map{
		dao.CreditConfig.Columns().Name:                  req.Name,
		dao.CreditConfig.Columns().Description:           req.Description,
		dao.CreditConfig.Columns().Logo:                  req.Logo,
		dao.CreditConfig.Columns().LogoUrl:               req.LogoUrl,
		dao.CreditConfig.Columns().Recurring:             req.Recurring,
		dao.CreditConfig.Columns().DiscountCodeExclusive: req.DiscountCodeExclusive,
		dao.CreditConfig.Columns().RechargeEnable:        req.RechargeEnable,
		dao.CreditConfig.Columns().PayoutEnable:          req.PayoutEnable,
		dao.CreditConfig.Columns().PreviewDefaultUsed:    req.PreviewDefaultUsed,
		dao.CreditConfig.Columns().MetaData:              utility.MergeMetadata(one.MetaData, req.MetaData),
		dao.CreditConfig.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.CreditConfig.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "Update Credit Config Error:%s\n", err.Error())
		return nil, err
	}
	one = query.GetCreditConfig(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, req.Currency)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("CreditConfig(%d)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &credit.EditPromoConfigRes{CreditConfig: bean.SimplifyCreditConfig(one)}, nil
}
