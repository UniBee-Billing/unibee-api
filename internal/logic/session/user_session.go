package session

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee-api/api/session/user"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	auth2 "unibee-api/internal/logic/auth"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func NewUserSession(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error) {
	one := query.GetUserAccountByEmail(ctx, req.Email)
	if one == nil {
		one = &entity.UserAccount{
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Email:          req.Email,
			Phone:          req.Phone,
			Address:        req.Address,
			ExternalUserId: req.ExternalUserId,
			CreateTime:     gtime.Now().Timestamp(),
		}
		result, err := dao.UserAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			// err = gerror.Newf(`record insert failure %s`, err)
			return nil, gerror.NewCode(gcode.New(500, "server error", nil))
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(id)
	}
	// create user session
	session := utility.CreateSessionId(strconv.FormatUint(one.Id, 10))
	_, err = g.Redis().Set(ctx, session, one.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, session, 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	return &user.NewRes{
		UserId:         strconv.FormatUint(one.Id, 10),
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Url:            fmt.Sprintf("%s/session/redirect/%s/forward", consts.GetConfigInstance().Server.DomainPath, session),
	}, nil
}

func UserSessionRedirectEntrance(r *ghttp.Request) {
	session := r.Get("session").String()
	utility.Assert(len(session) > 0, "Session Is Nil")
	id, err := g.Redis().Get(r.Context(), session)
	utility.AssertError(err, "Get Session")
	utility.Assert(id != nil, "Session Expired")
	utility.Assert(len(id.String()) > 0, "Invalid Session")
	userId, err := strconv.Atoi(id.String())
	if err != nil {
		g.Log().Errorf(r.Context(), "UserSessionRedirectEntrance panic url: %s id: %s err:%s", r.GetUrl(), id, err)
		return
	}
	one := query.GetUserAccountById(r.Context(), uint64(userId))
	utility.Assert(one != nil, "Invalid Session, User Not Found")
	token, err := auth2.CreateToken(one.Email, one.Id)
	utility.AssertError(err, "Generate Session")
	utility.Assert(auth2.PutAuthTokenToCache(r.Context(), token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	// todo mark Redirect to UserPortal
}
