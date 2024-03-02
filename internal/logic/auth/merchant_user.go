package auth

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func ChangeMerchantUserPasswordWithOutOldVerify(ctx context.Context, email string, newPassword string) {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "user not found")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

func ChangeMerchantUserPassword(ctx context.Context, merchantId uint64, email string, oldPassword string, newPassword string) {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "user not found")
	utility.Assert(utility.ComparePasswords(one.Password, oldPassword), "password not match")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}
