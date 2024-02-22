package merchant

import (
	"context"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/auth"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	searchUser, err := auth.SearchUser(ctx, _interface.GetMerchantId(ctx), req.SearchKey)
	if err != nil {
		return nil, err
	}
	return &user.SearchRes{UserAccounts: searchUser}, nil
}
