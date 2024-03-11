package merchant

import (
	"context"
	"unibee/api/merchant/user"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	searchUser, err := auth.SearchUser(ctx, _interface.GetMerchantId(ctx), req.SearchKey)
	if err != nil {
		return nil, err
	}
	return &user.SearchRes{UserAccounts: searchUser}, nil
}
