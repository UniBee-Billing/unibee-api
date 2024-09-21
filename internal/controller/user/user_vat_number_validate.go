package user

import (
	"context"
	"unibee/api/user/vat"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {
	return service.VatNumberValidate(ctx, req)
}
