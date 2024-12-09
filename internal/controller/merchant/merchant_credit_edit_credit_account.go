package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) EditCreditAccount(ctx context.Context, req *credit.EditCreditAccountReq) (res *credit.EditCreditAccountRes, err error) {
	utility.Assert(req.Id > 0, "Invalid id")
	one := query.GetCreditAccountById(ctx, req.Id)
	utility.Assert(one != nil, "Account not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "Invalid Merchant")
	if req.PayoutEnable != nil {
		utility.Assert(*req.PayoutEnable == 0 || *req.PayoutEnable == 1, "payout should be 0 or 1")
	}
	if req.RechargeEnable != nil {
		utility.Assert(*req.RechargeEnable == 0 || *req.RechargeEnable == 1, "rechargeEnable should be 0 or 1")
	}
	_, err = dao.CreditAccount.Ctx(ctx).Data(g.Map{
		dao.CreditAccount.Columns().PayoutEnable:   req.PayoutEnable,
		dao.CreditAccount.Columns().RechargeEnable: req.RechargeEnable,
		dao.CreditAccount.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.CreditAccount.Columns().Id, req.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("CreditAccount(%d)", req.Id),
		Content:        "Edit",
		UserId:         one.UserId,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	one = query.GetCreditAccountById(ctx, req.Id)
	return &credit.EditCreditAccountRes{UserCreditAccount: bean.SimplifyCreditAccount(one)}, nil
}
