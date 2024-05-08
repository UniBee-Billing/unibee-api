package member

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/oversea_pay"
	email2 "unibee/internal/logic/email"
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

func UpdateMemberRole(ctx context.Context, merchantId uint64, memberId uint64, roleName string) error {
	member := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(member != nil, "member not found")
	utility.Assert(strings.Compare(member.Role, "Owner") == 0, "Cannot Update Owner Role")
	role := query.GetRoleByName(ctx, merchantId, roleName)
	utility.Assert(role != nil, "role not found")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Role:      role,
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, memberId).Where(dao.MerchantMember.Columns().MerchantId, merchantId).OmitNil().Update()
	return err
}

func TransferOwnerMember(ctx context.Context, merchantId uint64, memberId uint64) error {
	member := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(member != nil, "member not found")
	if strings.Compare(member.Role, "Owner") == 0 {
		return nil
	}
	role := query.GetRoleByName(ctx, merchantId, "Owner")
	utility.Assert(role != nil, "owner role not found")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Role:      "",
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().Role, "Owner").
		OmitNil().Update()
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Role:      "Owner",
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).
		Where(dao.MerchantMember.Columns().Id, memberId).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		OmitNil().Update()
	return err
}

func AddMerchantMember(ctx context.Context, merchantId uint64, email string, firstName string, lastName string, role string) error {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one == nil, "email exist")
	utility.Assert(strings.Compare(role, "Owner") != 0, "can not add owner member")
	merchantRole := query.GetRoleByName(ctx, merchantId, role)
	utility.Assert(merchantRole != nil, "role not found")

	merchantMasterMember := &entity.MerchantMember{
		MerchantId: merchantId,
		Email:      email,
		CreateTime: gtime.Now().Timestamp(),
		FirstName:  firstName,
		LastName:   lastName,
		Role:       role,
	}
	_, err := dao.MerchantMember.Ctx(ctx).Data(merchantMasterMember).OmitNil().Insert(merchantMasterMember)
	if err != nil {
		return err
	}
	err = email2.SendTemplateEmail(ctx, merchantId, email, "", email2.TemplateMerchantMemberInvite, "", &email2.TemplateVariable{
		UserName: merchantMasterMember.FirstName + " " + merchantMasterMember.LastName,
		Link:     "<a href=\"" + config.GetConfigInstance().Server.GetServerPath() + "\">Link</a>",
	})
	utility.AssertError(err, "AddMerchantMember")
	return nil
}
