package merchant

import (
	"context"
	"unibee/api/merchant/user"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
)

func (c *ControllerUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	result, err := auth.UserList(ctx, &auth.UserListInternalReq{
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
	return &user.ListRes{UserAccounts: result.UserAccounts, Total: result.Total}, nil
}
