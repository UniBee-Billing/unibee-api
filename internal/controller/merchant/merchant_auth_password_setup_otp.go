package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	auth2 "unibee/internal/logic/member"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) PasswordSetupOtp(ctx context.Context, req *auth.PasswordSetupOtpReq) (res *auth.PasswordSetupOtpRes, err error) {
	setupToken, err := g.Redis().Get(ctx, req.Email+"-MerchantAuth-PasswordSetup-Verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(setupToken != nil, "Setup token invalid")
	utility.Assert((setupToken.String()) == req.SetupToken, "Setup token not match")

	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "User Not Found")
	auth2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, req.Email, req.NewPassword)
	_, err = g.Redis().Del(ctx, req.Email+"-MerchantAuth-PasswordSetup-Verify")
	if err != nil {
		g.Log().Errorf(ctx, "Delete_Setup_Token Error:%s", err.Error())
	}
	return &auth.PasswordSetupOtpRes{}, nil
}
