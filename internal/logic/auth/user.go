package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/jwt"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func UserSessionTransfer(ctx context.Context, session string) *entity.UserAccount {
	utility.Assert(len(session) > 0, "Session Is Nil")
	id, err := g.Redis().Get(ctx, session)
	utility.AssertError(err, "System Error")
	utility.Assert(id != nil && !id.IsNil() && !id.IsEmpty(), "Session Expired")
	utility.Assert(len(id.String()) > 0, "Invalid Session")
	userId, err := strconv.Atoi(id.String())
	utility.AssertError(err, "System Error")
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "Invalid Session, User Not Found")
	one.Password = ""
	return one
}

func ChangeUserPasswordWithOutOldVerify(ctx context.Context, merchantId uint64, email string, newPassword string) {
	one := query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "user not found")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

func ChangeUserPassword(ctx context.Context, merchantId uint64, email string, oldPassword string, newPassword string) {
	one := query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "user not found")
	utility.Assert(utility.ComparePasswords(one.Password, oldPassword), "password not match")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

func FrozenUser(ctx context.Context, userId int64) {
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "user not found")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Status:    2,
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

func ReleaseUser(ctx context.Context, userId int64) {
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "user not found")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Status:    0,
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "server error")
}

type NewReq struct {
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	Password       string `json:"password" dc:"Password"`
	Phone          string `json:"phone" dc:"Phone" `
	Address        string `json:"address" dc:"Address"`
	UserName       string `json:"userName" dc:"UserName"`
	CountryCode    string `json:"countryCode" dc:"CountryCode"`
	CountryName    string `json:"countryName" dc:"CountryName"`
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId"`
}

func PasswordLogin(ctx context.Context, merchantId uint64, email string, password string) (one *entity.UserAccount, token string) {
	one = query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "Email Not Found")
	utility.Assert(one.Status == 0, "Account Status Abnormal")
	utility.Assert(utility.ComparePasswords(one.Password, password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, one.Email)
	fmt.Println("logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	return one, token
}

func CreateUser(ctx context.Context, req *NewReq) (one *entity.UserAccount, err error) {
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(req != nil, "Server Error")
	if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	}
	if one == nil {
		one = query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	}
	utility.Assert(one == nil, "email or externalUserId exist")
	emailOne := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	utility.Assert(emailOne == nil, "email exist")
	one = &entity.UserAccount{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       utility.PasswordEncrypt(req.Password),
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		ExternalUserId: req.ExternalUserId,
		CountryCode:    req.CountryCode,
		CountryName:    req.CountryName,
		UserName:       req.UserName,
		MerchantId:     req.MerchantId,
		CreateTime:     gtime.Now().Timestamp(),
	}
	result, err := dao.UserAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "Server Error")
	id, err := result.LastInsertId()
	utility.AssertError(err, "Server Error")
	one.Id = uint64(id)
	return one, nil
}

func QueryOrCreateUser(ctx context.Context, req *NewReq) (one *entity.UserAccount, err error) {
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(req != nil, "Server Error")
	if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	}
	if one == nil {
		one = query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	}
	if one == nil {
		// check email not exist
		one, err = CreateUser(ctx, req)
		utility.AssertError(err, "Server Error")
	} else {
		if strings.Compare(one.Email, req.Email) != 0 {
			//email changed, update email
			emailOne := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
			utility.Assert(emailOne == nil || emailOne.Id == one.Id, "email of other externalUserId exist")
			_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().Email:     req.Email,
				dao.UserAccount.Columns().GmtModify: gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
			utility.AssertError(err, "Server Error")
		}
		if strings.Compare(one.ExternalUserId, req.ExternalUserId) != 0 {
			//externalUserId not match, update externalUserId
			otherOne := query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
			utility.Assert(otherOne == nil || otherOne.Id == one.Id, "externalUserId of other email exist")
			_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().ExternalUserId: req.ExternalUserId,
				dao.UserAccount.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
			utility.AssertError(err, "Server Error")
		}
		utility.Assert(one.Status == 0, "account status abnormal")
		_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().Address:   req.Address,
			dao.UserAccount.Columns().Phone:     req.Phone,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
		utility.AssertError(err, "Server Error")
	}
	return
}

func HardDeleteUser(ctx context.Context, userId uint64) error {
	_, err := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().Id, userId).Delete()
	if err != nil {
		return err
	}
	return nil
}
