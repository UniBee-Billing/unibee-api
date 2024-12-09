package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	credit2 "unibee/internal/logic/credit/recharge"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) NewCreditRecharge(ctx context.Context, req *credit.NewCreditRechargeReq) (res *credit.NewCreditRechargeRes, err error) {
	utility.Assert(req.UserId > 0, "invalid UserId")
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.MerchantId == _interface.GetMerchantId(ctx), "user not match")
	utility.Assert(req.GatewayId > 0, "invalid GatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(req.RechargeAmount > 0, "invalid Amount")
	utility.Assert(len(req.Currency) > 0, "invalid currency")
	if len(req.Name) == 0 {
		req.Name = "Credit Recharge"
	}
	if len(req.Description) == 0 {
		req.Description = "Credit Recharge"
	}
	createRes, err := credit2.CreateRechargePayment(ctx, &credit2.CreditRechargeInternalReq{
		UserId:         req.UserId,
		MerchantId:     _interface.GetMerchantId(ctx),
		GatewayId:      req.GatewayId,
		RechargeAmount: req.RechargeAmount,
		Currency:       req.Currency,
		Name:           req.Name,
		Description:    req.Description,
		ReturnUrl:      req.ReturnUrl,
		CancelUrl:      req.CancelUrl,
	})
	if err != nil {
		return nil, err
	}
	return &credit.NewCreditRechargeRes{
		User:           createRes.User,
		Merchant:       createRes.Merchant,
		Gateway:        createRes.Gateway,
		CreditAccount:  createRes.CreditAccount,
		CreditRecharge: createRes.CreditRecharge,
		Invoice:        createRes.Invoice,
		Payment:        createRes.Payment,
		Link:           createRes.Link,
		Paid:           createRes.Paid,
	}, nil
}
