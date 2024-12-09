package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) EditConfig(ctx context.Context, req *credit.EditConfigReq) (res *credit.EditConfigRes, err error) {
	utility.Assert(req.Id > 0, "invalid id")
	one := query.GetCreditConfigById(ctx, _interface.GetMerchantId(ctx), req.Id)
	utility.Assert(one != nil, "config not found by id")
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
		dao.CreditConfig.Columns().MetaData:              utility.MergeMetadata(one.MetaData, *req.MetaData),
		dao.CreditConfig.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.CreditConfig.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "Update Credit Config Error:%s\n", err.Error())
		return nil, err
	}

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("CreditConfig(%d)", req.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &credit.EditConfigRes{}, nil
}
