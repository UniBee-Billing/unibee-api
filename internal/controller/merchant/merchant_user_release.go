package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/auth"
	"unibee-api/utility"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) Release(ctx context.Context, req *user.ReleaseReq) (res *user.ReleaseRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	auth.ReleaseUser(ctx, req.UserId)
	return &user.ReleaseRes{}, nil
}
