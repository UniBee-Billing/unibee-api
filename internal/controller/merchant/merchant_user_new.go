package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) New(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error) {
	one, err := auth.QueryOrCreateUser(ctx, &auth.NewReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Phone:          req.Phone,
		Address:        req.Address,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	return &user.NewRes{User: bean.SimplifyUserAccount(one)}, nil
}
