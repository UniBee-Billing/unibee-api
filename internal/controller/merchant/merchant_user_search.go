package merchant

import (
	"context"
	"unibee/api/merchant/user"
	_interface "unibee/internal/interface/context"
	user2 "unibee/internal/logic/user"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	searchUser, err := user2.SearchUser(ctx, _interface.GetMerchantId(ctx), req.SearchKey)
	if err != nil {
		return nil, err
	}
	return &user.SearchRes{UserAccounts: searchUser}, nil
}
