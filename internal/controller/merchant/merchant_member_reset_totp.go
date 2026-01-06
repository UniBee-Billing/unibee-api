package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) ResetTotp(ctx context.Context, req *member.ResetTotpReq) (res *member.ResetTotpRes, err error) {
	one := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(one != nil, "Merchant Member Not Found")
	if one.TotpValidatorType <= 0 {
		return &member.ResetTotpRes{}, nil
	}
	utility.Assert(len(req.TotpCode) > 0, "Invalid 2FA Code")
	utility.Assert(totp.ValidateMerchantMemberTotp(ctx, one.TotpValidatorType, one.Email, one.TotpValidatorSecret, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity), "Invalid 2FA Code")
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().TotpValidatorSecret: "",
		dao.MerchantMember.Columns().TotpValidatorType:   0,
		dao.MerchantMember.Columns().GmtModify:           gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).Update()
	one = query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	return &member.ResetTotpRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one)}, nil
}
