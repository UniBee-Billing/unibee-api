package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	role2 "unibee/internal/logic/role"

	"unibee/api/merchant/role"
)

func (c *ControllerRole) List(ctx context.Context, req *role.ListReq) (res *role.ListRes, err error) {
	list, total := role2.MerchantRoleList(ctx, _interface.GetMerchantId(ctx))
	return &role.ListRes{MerchantRoles: list, Total: total}, nil
}
