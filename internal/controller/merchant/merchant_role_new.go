package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	role2 "unibee/internal/logic/role"
	"unibee/utility"

	"unibee/api/merchant/role"
)

func (c *ControllerRole) New(ctx context.Context, req *role.NewReq) (res *role.NewRes, err error) {
	utility.Assert(len(req.Role) > 0, "invalid role")
	utility.Assert(_interface.Context().Get(ctx).MerchantMember.IsOwner, "only owner can edit permission")
	err = role2.NewMerchantRole(ctx, &role2.CreateRoleInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		Role:           req.Role,
		PermissionData: req.Permissions,
	})
	if err != nil {
		return nil, err
	}
	return &role.NewRes{}, nil
}
