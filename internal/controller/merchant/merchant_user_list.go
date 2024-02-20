package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/auth"
	"unibee-api/utility"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	result, err := auth.UserAccountList(ctx, &auth.UserListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Status:     req.Status,
		//UserName:           req.UserName,
		//SubscriptionName:   req.SubscriptionName,
		//SubscriptionStatus: req.SubscriptionStatus,
		//PaymentMethod:      req.PaymentMethod,
		//BillingType:        req.BillingType,
		DeleteInclude: req.DeleteInclude,
		SortField:     req.SortField,
		SortType:      req.SortType,
		Page:          req.Page,
		Count:         req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &user.ListRes{UserAccounts: result.UserAccounts}, nil
}
