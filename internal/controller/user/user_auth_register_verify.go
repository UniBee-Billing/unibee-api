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

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	// utility.Assert(len(req.Phone) > 0, "phone not null")
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-verify")
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "Invalid Code")
	//if verificationCode == nil {
	//	return nil, gerror.NewCode(gcode.New(400, "invalid code", nil))
	//}

	utility.Assert((verificationCode.String()) == req.VerificationCode, "Invalid Code")
	//if (verificationCode.String()) != req.VerificationCode {
	//	return nil, gerror.NewCode(gcode.New(401, "invalid code", nil))
	//}

	userStr, err := g.Redis().Get(ctx, req.Email)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(userStr != nil, "Invalid Code")
	//if userStr == nil {
	//	return nil, gerror.NewCode(gcode.New(401, "invalid code", nil))
	//}

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
		MerchantId:  _interface.BizCtx().Get(ctx).MerchantId,
		CreateTime:  gtime.Now().Timestamp(),
	}

	result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitNil().Insert(user)
	// dao.UserAccount.Ctx(ctx).Data(user).OmitEmpty().Update()
	if err != nil {
		// err = gerror.Newf(`record insert failure %s`, err)
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	user.Id = uint64(id)
	var newOne *entity.UserAccount
	newOne = query.GetUserAccountById(ctx, user.Id)
	if newOne == nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	newOne.Password = ""
	return &auth.RegisterVerifyRes{User: newOne}, nil
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
