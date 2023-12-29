package auth

import (
	"context"

	"encoding/json"
	v1 "go-oversea-pay/api/auth/v1"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) RegisterVerify(ctx context.Context, req *v1.RegisterVerifyReq) (res *v1.RegisterVerifyRes, err error) {
		// utility.Assert(len(req.Phone) > 0, "phone not null")
		verificationCode, err := g.Redis().Get(ctx, req.Email + "-verify")
		if err != nil {
			// return nil, gerror.New("internal error")
			return nil, gerror.NewCode(gcode.New(500, "server error", nil))
		}
		if verificationCode == nil {
			return nil, gerror.NewCode(gcode.New(400, "invalid code", nil))
		}

		if (verificationCode.String()) != req.VerificationCode {
			return nil, gerror.NewCode(gcode.New(401, "invalid code", nil))
		}

		userStr, err := g.Redis().Get(ctx, req.Email)
		if err != nil {
			return nil, gerror.NewCode(gcode.New(500, "server error", nil))
		}
		if userStr == nil {
			return nil, gerror.NewCode(gcode.New(401, "invalid code", nil))
		}

		u := struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName string
		}{}
		err = json.Unmarshal([]byte(userStr.String()), &u)

		user := &entity.UserAccount{
			FirstName: u.FirstName,
			LastName: u.LastName,
			Email:    u.Email,
			Password: u.Password,
			Phone:   u.Phone,
			Address: u.Address,
			UserName: u.UserName,
		}
		
		result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitEmpty().Insert(user)
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
		return &v1.RegisterVerifyRes{User: newOne}, nil
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
