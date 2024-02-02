package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/auth"

	"go-oversea-pay/api/merchant/user"
)

func (c *ControllerUser) Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error) {
	searchUser, err := auth.SearchUser(ctx, req.SearchKey)
	if err != nil {
		return nil, err
	}
	return &user.SearchRes{UserAccounts: searchUser}, nil
}
