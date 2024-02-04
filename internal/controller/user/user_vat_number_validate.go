package user

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/user/vat"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {

	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
	return service.VatNumberValidate(ctx, req, int64(_interface.BizCtx().Get(ctx).User.Id))
}
