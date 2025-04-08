package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	user2 "unibee/internal/logic/user"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) Count(ctx context.Context, req *user.CountReq) (res *user.CountRes, err error) {
	total, err := user2.UserCount(ctx, &user2.UserListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	if err != nil {
		return nil, err
	}
	return &user.CountRes{Total: total}, nil
}
