package session

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/user/session"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func NewUserSession(ctx context.Context, req *session.NewReq) (res *session.NewRes, err error) {
	one := query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	if one == nil {
		one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	}
	if one == nil {
		// check email not exsit
		emailOne := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
		utility.Assert(emailOne == nil, "email exist")
		one = &entity.UserAccount{
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Password:       utility.PasswordEncrypt(req.Password),
			Email:          req.Email,
			Phone:          req.Phone,
			Address:        req.Address,
			ExternalUserId: req.ExternalUserId,
			MerchantId:     _interface.GetMerchantId(ctx),
			CreateTime:     gtime.Now().Timestamp(),
		}
		result, err := dao.UserAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
		utility.AssertError(err, "Server Error")
		id, err := result.LastInsertId()
		utility.AssertError(err, "Server Error")
		one.Id = uint64(id)
	} else {
		if strings.Compare(one.Email, req.Email) != 0 {
			//email changed, update email
			emailOne := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
			utility.Assert(emailOne == nil || emailOne.Id == one.Id, "email of other externalUserId exist")
			_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().Email:     req.Email,
				dao.UserAccount.Columns().GmtModify: gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
			utility.AssertError(err, "Server Error")
		}
		if strings.Compare(one.ExternalUserId, req.ExternalUserId) != 0 {
			//externalUserId not match, update externalUserId
			otherOne := query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
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
	// create user session
	ss := utility.GenerateRandomAlphanumeric(40)
	_, err = g.Redis().Set(ctx, ss, one.Id)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, ss, 3*60)
	utility.AssertError(err, "Server Error")
	merchantInfo := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(len(merchantInfo.Host) > 0, "user host not set")

	return &session.NewRes{
		UserId:         strconv.FormatUint(one.Id, 10),
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Url:            fmt.Sprintf("%s://%s/session-result?session=%s", consts.GetConfigInstance().Server.GetDomainScheme(), merchantInfo.Host, ss),
	}, nil
}
