package merchant

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"encoding/json"
	"unibee-api/api/merchant/auth"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
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
		FirstName, LastName, Email, Password, Phone, Address, UserName string
		MerchantId                                                     int64
	}{}
	err = json.Unmarshal([]byte(userStr.String()), &u)
	// TOO: check err != nil

	user := &entity.MerchantUserAccount{
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Password:   u.Password,
		MerchantId: u.MerchantId,
		// Phone:     u.Phone,
		// Address:   u.Address,
		UserName:   u.UserName,
		CreateTime: gtime.Now().Timestamp(),
	}

	// race condition: email exist checking is too earlier
	result, err := dao.MerchantUserAccount.Ctx(ctx).Data(user).OmitNil().Insert(user)
	// dao.UserAccount.Ctx(ctx).Data(user).OmitNil().Update()
	if err != nil {
		// err = gerror.Newf(`record insert failure %s`, err)
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	user.Id = uint64(id)
	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantAccountById(ctx, user.Id)
	if newOne == nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	// TODO: return &{} is enough, front-end need to re-login, so don't need to return the whole user obj
	return &auth.RegisterVerifyRes{MerchantUser: newOne}, nil
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

/*

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	// utility.Assert(len(req.Phone) > 0, "phone not null")

}

*/
