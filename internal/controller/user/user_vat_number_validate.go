package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/vat"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {
	return service.VatNumberValidate(ctx, req)
}
