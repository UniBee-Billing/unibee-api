package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) DeleteDevice(ctx context.Context, req *member.DeleteDeviceReq) (res *member.DeleteDeviceRes, err error) {
	admin := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	targetMember := query.GetMerchantMemberById(ctx, req.MemberId)
	utility.Assert(targetMember != nil, "Merchant Member Not Found")
	utility.Assert(admin != nil, "Merchant Admin Not Found")
	utility.Assert(admin.Role == "Owner" || admin.Id == targetMember.Id, "Only Owner or yourself can delete device")
	utility.Assert(len(req.DeviceIdentity) > 0, "Invalid Device Identity")
	utility.Assert(totp.IsClientIdentityExist(ctx, targetMember.Email, req.DeviceIdentity), "device not found")
	totp.DeleteClientIdentity(ctx, targetMember.Email, req.DeviceIdentity)

	return &member.DeleteDeviceRes{DeviceList: detail.GetClientIdentityList(ctx, targetMember.Email)}, nil
}
