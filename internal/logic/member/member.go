package member

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/platform"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func NewSession(ctx context.Context, memberId int64, returnUrl string) (session string, err error) {
	utility.Assert(memberId > 0, "Invalid member id")
	one := query.GetMerchantMemberById(ctx, uint64(memberId))
	utility.Assert(one != nil, "Invalid Session, Member Not Found")
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	// create user session
	session = utility.GenerateRandomAlphanumeric(40)
	_, err = g.Redis().Set(ctx, fmt.Sprintf("MemberSessionKey:%s", session), one.Id)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, fmt.Sprintf("MemberSessionKey:%s", session), config.GetConfigInstance().Auth.Login.Expire*60)
	utility.AssertError(err, "Server Error")
	if len(returnUrl) > 0 {
		_, err = g.Redis().Set(ctx, fmt.Sprintf("MemberSessionReturnUrlKey:%s", session), returnUrl)
		utility.AssertError(err, "Server Error")
		_, err = g.Redis().Expire(ctx, fmt.Sprintf("MemberSessionReturnUrlKey:%s", session), config.GetConfigInstance().Auth.Login.Expire*60)
		utility.AssertError(err, "Server Error")
	}
	return session, nil
}

func SessionTransfer(ctx context.Context, session string) (*entity.MerchantMember, string) {
	utility.Assert(len(session) > 0, "Session Is Nil")
	id, err := g.Redis().Get(ctx, fmt.Sprintf("MemberSessionKey:%s", session))
	utility.AssertError(err, "System Error")
	utility.Assert(id != nil && !id.IsNil() && !id.IsEmpty(), "Session Expired")
	utility.Assert(len(id.String()) > 0, "Invalid Session")
	memberId, err := strconv.Atoi(id.String())
	utility.AssertError(err, "System Error")
	one := query.GetMerchantMemberById(ctx, uint64(memberId))
	utility.Assert(one != nil, "Invalid Session, Member Not Found")
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	var returnUrl = ""
	returnData, err := g.Redis().Get(ctx, fmt.Sprintf("MemberSessionReturnUrlKey:%s", session))
	if err == nil {
		returnUrl = returnData.String()
	}
	return one, returnUrl
}

func PasswordLogin(ctx context.Context, email string, password string) (one *entity.MerchantMember, token string) {
	one = query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one != nil, "Email Not Found")
	utility.Assert(utility.ComparePasswords(one.Password, password), "Login Failed, Password Not Match")

	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, one.MerchantId, one.Id, one.Email)
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Member(%v)", one.Id),
		Content:        "ChangePasswordByVerifyCode",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Member(%v)", one.Id),
		Content:        "ChangePassword",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
}

func UpdateMemberRole(ctx context.Context, merchantId uint64, memberId uint64, roleIds []uint64) error {
	member := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(member != nil, "member not found")
	utility.Assert(strings.Compare(member.Role, "Owner") != 0, "Can't Update Owner's Role")
	for _, roleId := range roleIds {
		role := query.GetRoleById(ctx, roleId)
		utility.Assert(role != nil, "roleId "+strconv.FormatUint(roleId, 10)+" not found")
	}
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Role:      utility.MarshalToJsonString(roleIds),
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, memberId).Where(dao.MerchantMember.Columns().MerchantId, merchantId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     merchantId,
		Target:         fmt.Sprintf("Member(%v)", member.Id),
		Content:        "UpdateRole",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	ReloadMemberCacheForSdkAuthBackground(member.Id)
	return err
}

func TransferOwnerMember(ctx context.Context, merchantId uint64, memberId uint64) error {
	member := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(member != nil, "member not found")
	if strings.Compare(member.Role, "Owner") == 0 {
		return nil
	}
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
	ReloadMemberCacheForSdkAuthBackground(member.Id)
	ReloadMemberCacheForSdkAuthBackground(memberId)
	return err
}

func AddMerchantMember(ctx context.Context, merchantId uint64, email string, firstName string, lastName string, roleIds []uint64) error {
	one := query.GetMerchantMemberByEmail(ctx, email)
	utility.Assert(one == nil, "email exist")
	for _, roleId := range roleIds {
		role := query.GetRoleById(ctx, roleId)
		utility.Assert(role != nil, "roleId "+strconv.FormatUint(roleId, 10)+" not found")
	}

	one = &entity.MerchantMember{
		MerchantId: merchantId,
		Email:      email,
		CreateTime: gtime.Now().Timestamp(),
		FirstName:  firstName,
		LastName:   lastName,
		Role:       utility.MarshalToJsonString(roleIds),
	}

	result, err := dao.MerchantMember.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	var link = config.GetConfigInstance().Server.GetServerPath()
	if strings.Compare(link, "https://api.unibee.top") == 0 {
		link = "https://merchant.unibee.top"
	} else {
		link = strings.Replace(link, "/api", "", 1)
	}
	merchant := query.GetMerchantById(ctx, one.MerchantId)
	utility.Assert(merchant != nil, "Invalid Merchant")
	{
		ownerEmail := ""
		ownerMember := query.GetMerchantOwnerMember(ctx, merchant.Id)
		if ownerMember != nil {
			ownerEmail = ownerMember.Email
		}

		err = platform.SentPlatformMerchantInviteMember(map[string]string{
			"ownerEmail":  ownerEmail,
			"memberEmail": email,
			"firstName":   one.FirstName,
			"lastName":    one.LastName,
			"link":        link,
		})
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Member(%v)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "AddMerchantMember Error")
	ReloadMemberCacheForSdkAuthBackground(one.Id)
	return nil
}

func FrozenMember(ctx context.Context, memberId uint64) {
	one := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(one != nil, "Member Not Found")
	utility.Assert(one.Status != 2, "Member Already Suspended")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Status:    2,
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Member(%v)", one.Id),
		Content:        "Suspend",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
	ReloadMemberCacheForSdkAuthBackground(one.Id)
}

func ReleaseMember(ctx context.Context, memberId uint64) {
	one := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(one != nil, "member not found")
	_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().Status:    0,
		dao.MerchantMember.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Member(%v)", one.Id),
		Content:        "Resume",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
	ReloadMemberCacheForSdkAuthBackground(one.Id)
}
