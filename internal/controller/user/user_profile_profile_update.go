package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"go-oversea-pay/api/user/profile"
	dao "go-oversea-pay/internal/dao/oversea_pay"
)

func (c *ControllerProfile) ProfileUpdate(ctx context.Context, req *profile.ProfileUpdateReq) (res *profile.ProfileUpdateRes, err error) {

	update, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		// dao.UserAccount.Columns().Email :        req.Email,
		dao.UserAccount.Columns().LastName:        req.LastName,
		dao.UserAccount.Columns().FirstName:       req.FirstName,
		dao.UserAccount.Columns().Address:         req.Address,
		dao.UserAccount.Columns().CompanyName:     req.CompanyName,
		dao.UserAccount.Columns().VATNumber:       req.VATNumber,
		dao.UserAccount.Columns().Phone:           req.Phone,
		dao.UserAccount.Columns().Telegram:        req.Telegram,
		dao.UserAccount.Columns().WhatsAPP:        req.WhatsApp,
		dao.UserAccount.Columns().WeChat:          req.WeChat,
		dao.UserAccount.Columns().LinkedIn:        req.LinkedIn,
		dao.UserAccount.Columns().Facebook:        req.Facebook,
		dao.UserAccount.Columns().TikTok:          req.TikTok,
		dao.UserAccount.Columns().OtherSocialInfo: req.OtherSocialInfo,
		dao.UserAccount.Columns().PaymentMethod:   req.PaymentMethod,
		dao.UserAccount.Columns().CountryCode:     req.CountryCode,
		dao.UserAccount.Columns().CountryName:     req.CountryName,
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, req.Id).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	// TODO: return the updated user account
	return &profile.ProfileUpdateRes{}, nil
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
