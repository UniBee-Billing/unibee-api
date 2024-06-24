package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Profile(ctx context.Context, req *member.ProfileReq) (res *member.ProfileRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Merchant Member Not Found")
	one := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(one != nil, "Merchant Member Not Found")
	return &member.ProfileRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one)}, nil
}
