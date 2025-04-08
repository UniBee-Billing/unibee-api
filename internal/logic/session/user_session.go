package session

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/api/merchant/session"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/user"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func UserSessionTransfer(ctx context.Context, session string) (*entity.UserAccount, string, string) {
	utility.Assert(len(session) > 0, "Session Is Nil")
	id, err := g.Redis().Get(ctx, session)
	utility.AssertError(err, "System Error")
	utility.Assert(id != nil && !id.IsNil() && !id.IsEmpty(), "Session Expired")
	utility.Assert(len(id.String()) > 0, "Invalid Session")
	userId, err := strconv.Atoi(id.String())
	utility.AssertError(err, "System Error")
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "Invalid Session, User Not Found")
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	var returnUrl = ""
	returnData, err := g.Redis().Get(ctx, session+"_returnUrl")
	if err == nil {
		returnUrl = returnData.String()
	}
	var cancelUrl = ""
	cancelUrlData, err := g.Redis().Get(ctx, session+"_cancelUrl")
	if err == nil {
		cancelUrl = cancelUrlData.String()
	}
	return one, returnUrl, cancelUrl
}

func NewUserSession(ctx context.Context, merchantId uint64, userId uint64, returnUrl string, cancelUrl string) (clientToken string, clientSession string, err error) {
	merchantInfo := query.GetMerchantById(ctx, merchantId)
	one := query.GetUserAccountById(ctx, userId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.AssertError(err, "Server Error")
	utility.Assert(one != nil, "Server Error")
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	// create user session
	ss := utility.GenerateRandomAlphanumeric(40)
	_, err = g.Redis().Set(ctx, ss, one.Id)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, ss, config.GetConfigInstance().Auth.Login.Expire*60)
	utility.AssertError(err, "Server Error")
	if len(returnUrl) > 0 {
		_, err = g.Redis().Set(ctx, ss+"_returnUrl", returnUrl)
		utility.AssertError(err, "Server Error")
		_, err = g.Redis().Expire(ctx, ss+"_returnUrl", config.GetConfigInstance().Auth.Login.Expire*60)
		utility.AssertError(err, "Server Error")
	}
	if len(cancelUrl) > 0 {
		_, err = g.Redis().Set(ctx, ss+"_cancelUrl", cancelUrl)
		utility.AssertError(err, "Server Error")
		_, err = g.Redis().Expire(ctx, ss+"_cancelUrl", config.GetConfigInstance().Auth.Login.Expire*60)
		utility.AssertError(err, "Server Error")
	}

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, one.Email, one.Language)
	fmt.Println("logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")

	return token, ss, nil
}

func NewUserPortalSession(ctx context.Context, merchantId uint64, req *session.NewReq) (res *session.NewRes, err error) {
	one, err := user.QueryOrCreateUser(ctx, &user.NewUserInternalReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Phone:          req.Phone,
		Address:        req.Address,
		MerchantId:     merchantId,
	})
	merchantInfo := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.AssertError(err, "Server Error")
	utility.Assert(one != nil, "Server Error")
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	// create user session
	ss := utility.GenerateRandomAlphanumeric(40)
	_, err = g.Redis().Set(ctx, ss, one.Id)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, ss, config.GetConfigInstance().Auth.Login.Expire*60)
	utility.AssertError(err, "Server Error")
	if len(req.ReturnUrl) > 0 {
		_, err = g.Redis().Set(ctx, ss+"_returnUrl", req.ReturnUrl)
		utility.AssertError(err, "Server Error")
		_, err = g.Redis().Expire(ctx, ss+"_returnUrl", config.GetConfigInstance().Auth.Login.Expire*60)
		utility.AssertError(err, "Server Error")
	}
	if len(req.CancelUrl) > 0 {
		_, err = g.Redis().Set(ctx, ss+"_cancelUrl", req.CancelUrl)
		utility.AssertError(err, "Server Error")
		_, err = g.Redis().Expire(ctx, ss+"_cancelUrl", config.GetConfigInstance().Auth.Login.Expire*60)
		utility.AssertError(err, "Server Error")
	}

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email, one.Language)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	url := ""
	if len(merchantInfo.Host) > 0 {
		url = fmt.Sprintf("%s://%s/session-result?session=%s", config.GetConfigInstance().Server.GetDomainScheme(), merchantInfo.Host, ss)
	}
	//utility.Assert(len(merchantInfo.Host) > 0, "user portal host not set")
	return &session.NewRes{
		UserId:         strconv.FormatUint(one.Id, 10),
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Url:            url,
		ClientToken:    token,
		ClientSession:  ss,
	}, nil
}
