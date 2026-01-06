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

func (c *ControllerMember) ConfirmTotpKey(ctx context.Context, req *member.ConfirmTotpKeyReq) (res *member.ConfirmTotpKeyRes, err error) {
	one := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(one != nil, "Merchant Member Not Found")
	utility.Assert(one.TotpValidatorType == 0, "Already setup 2FA")
	utility.Assert(req.TotpType > 0, "invalid 2FA Type")
	utility.Assert(len(req.TotpKey) > 0, "invalid 2FA Key")
	utility.Assert(len(req.TotpCode) > 0, "Invalid 2FA Code")
	utility.Assert(totp.ValidateMerchantMemberTotp(ctx, req.TotpType, one.Email, req.TotpKey, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity), "Invalid 2FA Code")
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().TotpValidatorSecret: req.TotpKey,
		dao.MerchantMember.Columns().TotpValidatorType:   req.TotpType,
		dao.MerchantMember.Columns().GmtModify:           gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).Update()
	return &member.ConfirmTotpKeyRes{}, nil
}
