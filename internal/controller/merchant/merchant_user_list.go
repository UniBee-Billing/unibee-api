package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/utility"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	result, err := auth.UserAccountList(ctx, &auth.UserListInternalReq{
		MerchantId:    _interface.GetMerchantId(ctx),
		UserId:        req.UserId,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Status:        req.Status,
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
