package user

import (
	"context"
	"fmt"
	"unibee-api/time"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"unibee-api/api/user/profile"
	dao "unibee-api/internal/dao/oversea_pay"
)

func (c *ControllerProfile) ProfileUpdate(ctx context.Context, req *profile.ProfileUpdateReq) (res *profile.ProfileUpdateRes, err error) {
	// timezone check
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}
	_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
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
		dao.UserAccount.Columns().TimeZone:        req.TimeZone,
		dao.UserAccount.Columns().CountryCode:     req.CountryCode,
		dao.UserAccount.Columns().CountryName:     req.CountryName,
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, req.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	return &profile.ProfileUpdateRes{}, nil
}
