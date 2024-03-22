package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/jwt"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func PasswordLogin(ctx context.Context, email string, password string) (one *entity.MerchantMember, token string) {
	one = query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "Email Not Found")
	utility.Assert(utility.ComparePasswords(one.Password, password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEMERCHANTMember, one.MerchantId, one.Id, one.Email)
	fmt.Println("logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", one.Id)), "Cache Error")
	return one, token
}

func ChangeMerchantMemberPasswordWithOutOldVerify(ctx context.Context, email string, newPassword string) {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "user not found")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

func ChangeMerchantMemberPassword(ctx context.Context, email string, oldPassword string, newPassword string) {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "user not found")
	utility.Assert(utility.ComparePasswords(one.Password, oldPassword), "password not match")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}
