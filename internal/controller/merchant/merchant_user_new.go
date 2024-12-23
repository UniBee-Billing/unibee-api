package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) New(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error) {
	if config.GetConfigInstance().IsProd() {
		existOne := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
		utility.Assert(existOne == nil, "same email exist")
	}
	one, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Phone:          req.Phone,
		Address:        req.Address,
		Language:       req.Language,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	return &user.NewRes{User: bean.SimplifyUserAccount(one)}, nil
}
