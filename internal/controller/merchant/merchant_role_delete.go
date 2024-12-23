package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	role2 "unibee/internal/logic/role"
	"unibee/utility"

	"unibee/api/merchant/role"
)

func (c *ControllerRole) Delete(ctx context.Context, req *role.DeleteReq) (res *role.DeleteRes, err error) {
	//utility.Assert(_interface.Context().Get(ctx).MerchantMember.IsOwner, "only owner can edit permission")
	utility.Assert(req.Id > 0, "invalid roleId")
	err = role2.DeleteMerchantRole(ctx, _interface.GetMerchantId(ctx), req.Id)
	if err != nil {
		return nil, err
	}
	return &role.DeleteRes{}, nil
}
