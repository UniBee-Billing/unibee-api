package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	role2 "unibee/internal/logic/role"
	"unibee/utility"

	"unibee/api/merchant/role"
)

func (c *ControllerRole) Edit(ctx context.Context, req *role.EditReq) (res *role.EditRes, err error) {
	//utility.Assert(_interface.Context().Get(ctx).MerchantMember.IsOwner, "only owner can edit permission")
	utility.Assert(len(req.Role) > 0, "invalid role")
	err = role2.EditMerchantRole(ctx, &role2.CreateRoleInternalReq{
		Id:             req.Id,
		MerchantId:     _interface.GetMerchantId(ctx),
		Role:           req.Role,
		PermissionData: req.Permissions,
	})
	if err != nil {
		return nil, err
	}
	return &role.EditRes{}, nil
}
