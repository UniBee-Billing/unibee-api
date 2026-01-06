package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) ClearTotp(ctx context.Context, req *auth.ClearTotpReq) (res *auth.ClearTotpRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.TotpResumeCode != "", "Resume Code Cannot Be Empty")
	one := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(one != nil, "Email Not Found")
	utility.Assert(one.TotpValidatorType != 0, "Member 2FA Already Clear")
	utility.Assert(utility.MD5(one.TotpValidatorSecret) == req.TotpResumeCode, "Invalid Resume Code")
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().TotpValidatorSecret: "",
		dao.MerchantMember.Columns().TotpValidatorType:   0,
		dao.MerchantMember.Columns().GmtModify:           gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).Update()

	return &auth.ClearTotpRes{}, nil
}
