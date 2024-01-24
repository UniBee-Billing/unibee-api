package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/auth"

	"go-oversea-pay/api/merchant/user"
)

func (c *ControllerUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	result, err := auth.UserAccountList(ctx, &auth.UserListInternalReq{
		MerchantId:         req.MerchantId,
		UserId:             req.UserId,
		Email:              req.Email,
		UserName:           req.UserName,
		SubscriptionName:   req.SubscriptionName,
		SubscriptionStatus: req.SubscriptionStatus,
		PaymentMethod:      req.PaymentMethod,
		BillingType:        req.BillingType,
		DeleteInclude:      req.DeleteInclude,
		SortField:          req.SortField,
		SortType:           req.SortType,
		Page:               req.Page,
		Count:              req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &user.ListRes{UserAccounts: result.UserAccounts}, nil
}
