package merchant

import (
	"context"
	"unibee-api/internal/logic/auth"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	searchUser, err := auth.SearchUser(ctx, req.SearchKey)
	if err != nil {
		return nil, err
	}
	return &user.SearchRes{UserAccounts: searchUser}, nil
}
