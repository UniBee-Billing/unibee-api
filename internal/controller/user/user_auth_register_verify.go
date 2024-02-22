package user

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/api/user/auth"
	dao "unibee-api/internal/dao/oversea_pay"
	_interface "unibee-api/internal/interface"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify")
	utility.AssertError(err, "Server Error")
	utility.Assert(verificationCode != nil, "Invalid Code")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "Invalid Code")

	userStr, err := g.Redis().Get(ctx, CacheKeyUserRegisterPrefix+req.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(userStr != nil, "Invalid Code")
	u := struct {
		FirstName, LastName, Email, Password, Phone, Address, UserName, CountryCode, CountryName string
	}{}
	err = json.Unmarshal([]byte(userStr.String()), &u)

	user := &entity.UserAccount{
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		CountryCode: u.CountryCode,
		CountryName: u.CountryName,
		Password:    u.Password,
		Phone:       u.Phone,
		Address:     u.Address,
		UserName:    u.UserName,
		MerchantId:  _interface.GetMerchantId(ctx),
		CreateTime:  gtime.Now().Timestamp(),
	}

	result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitNil().Insert(user)
	utility.AssertError(err, "Server Error")
	id, _ := result.LastInsertId()
	user.Id = uint64(id)
	var newOne *entity.UserAccount
	newOne = query.GetUserAccountById(ctx, user.Id)
	utility.Assert(newOne != nil, "Server Error")
	newOne.Password = ""
	return &auth.RegisterVerifyRes{User: newOne}, nil
}
