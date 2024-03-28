package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	role2 "unibee/internal/logic/role"
	"unibee/utility"

	"unibee/api/merchant/role"
)

func (c *ControllerRole) Delete(ctx context.Context, req *role.DeleteReq) (res *role.DeleteRes, err error) {
	utility.Assert(len(req.Role) > 0, "invalid role")
	err = role2.DeleteMerchantRole(ctx, _interface.GetMerchantId(ctx), req.Role)
	if err != nil {
		return nil, err
	}
	return &role.DeleteRes{}, nil
}
