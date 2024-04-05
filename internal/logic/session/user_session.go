package session

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/api/merchant/session"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/auth"
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
	return one
}

func NewUserSession(ctx context.Context, merchantId uint64, req *session.NewReq) (res *session.NewRes, err error) {
	one, err := auth.QueryOrCreateUser(ctx, &auth.NewReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Phone:          req.Phone,
		Address:        req.Address,
		MerchantId:     merchantId,
	})
	utility.AssertError(err, "Server Error")
	utility.Assert(one != nil, "Server Error")
	// create user session
	ss := utility.GenerateRandomAlphanumeric(40)
	_, err = g.Redis().Set(ctx, ss, one.Id)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, ss, config.GetConfigInstance().Auth.Login.Expire*60)
	utility.AssertError(err, "Server Error")
	merchantInfo := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(len(merchantInfo.Host) > 0, "user host not set")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")

	return &session.NewRes{
		UserId:         strconv.FormatUint(one.Id, 10),
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Url:            fmt.Sprintf("%s://%s/session-result?session=%s", config.GetConfigInstance().Server.GetDomainScheme(), merchantInfo.Host, ss),
		ClientToken:    token,
		ClientSession:  ss,
	}, nil
}
