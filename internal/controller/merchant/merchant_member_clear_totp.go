package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) ClearTotp(ctx context.Context, req *member.ClearTotpReq) (res *member.ClearTotpRes, err error) {
	admin := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(admin != nil, "Merchant Admin Not Found")
	utility.Assert(admin.Role == "Owner", "Only Owner can clear member 2FA")
	one := query.GetMerchantMemberById(ctx, req.MemberId)
	utility.Assert(one.MerchantId == admin.MerchantId, "no permission")
	utility.Assert(len(req.TotpCode) > 0, "Invalid 2FA Code")
	utility.Assert(totp.ValidateMerchantMemberTotp(ctx, admin.TotpValidatorType, admin.Email, admin.TotpValidatorSecret, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity), "Invalid 2FA Code")
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().TotpValidatorSecret: "",
		dao.MerchantMember.Columns().TotpValidatorType:   0,
		dao.MerchantMember.Columns().GmtModify:           gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).Update()
	return &member.ClearTotpRes{}, nil
}
