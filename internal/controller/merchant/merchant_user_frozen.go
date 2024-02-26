package merchant

import (
	"context"
	"unibee/api/merchant/user"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/utility"
)

func (c *ControllerUser) Frozen(ctx context.Context, req *user.FrozenReq) (res *user.FrozenRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	auth.FrozenUser(ctx, req.UserId)
	return &user.FrozenRes{}, nil
}
