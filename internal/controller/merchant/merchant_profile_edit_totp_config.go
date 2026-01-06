package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) EditTotpConfig(ctx context.Context, req *profile.EditTotpConfigReq) (res *profile.EditTotpConfigRes, err error) {
	totp.UpdateMerchantTotpGlobalConfig(ctx, _interface.GetMerchantId(ctx), req.Activate)
	return &profile.EditTotpConfigRes{}, nil
}
