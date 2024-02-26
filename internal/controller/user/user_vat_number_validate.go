package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/user/vat"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {

	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
	return service.VatNumberValidate(ctx, req, int64(_interface.BizCtx().Get(ctx).User.Id))
}
