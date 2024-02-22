package middleware

import (
	"fmt"
	"net/url"
	"strings"
	"unibee-api/internal/consts"
	_ "unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/jwt"
	"unibee-api/internal/model"
	"unibee-api/internal/query"
	utility "unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SMiddleware struct {
	LoginUrl string // 登录路由地址
}

func init() {
	_interface.RegisterMiddleware(New())
}

func New() *SMiddleware {
	return &SMiddleware{
		LoginUrl: "/login",
	}
}

func (s *SMiddleware) CORS(r *ghttp.Request) {
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Host:%s", r.GetHost())
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Origin:%s", r.GetHeader("Origin"))
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Referer:%s", r.GetHeader("Referer"))
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header User-Agent:%s", r.GetHeader("User-Agent"))
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *SMiddleware) ResponseHandler(r *ghttp.Request) {
	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	customCtx.RequestId = utility.CreateRequestId()
	_interface.BizCtx().Init(r, customCtx)
	r.Assigns(g.Map{
		consts.ContextKey: customCtx,
	})

	utility.Try(r.Middleware.Next, func(err interface{}) {
		json, _ := r.GetJson()
		g.Log().Errorf(r.Context(), "Global_Exception Panic Url: %s Params:%s Error:%v", r.GetUrl(), json, err)
		return
	})

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)

	if err == nil && r.Response.BufferLength() > 0 {
		return
	}

	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		json, _ := r.GetJson()
		g.Log().Errorf(r.Context(), "Global_exception requestId:%s url: %s params:%s code:%d error:%s", _interface.BizCtx().Get(r.Context()).RequestId, r.GetUrl(), json, code.Code(), err.Error())
		r.Response.ClearBuffer() // inner panic will contain json data，need clean
		r.Response.Status = 200  // error reply in json code, http code always 200
		message := err.Error()
		if strings.Contains(message, utility.SystemAssertPrefix) || code == gcode.CodeValidationFailed {
			utility.JsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
		} else {
			utility.JsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", _interface.BizCtx().Get(r.Context()).RequestId, code.Code()))
		}
	} else {
		r.Response.Status = 200
		utility.JsonExit(r, code.Code(), "", res)
	}
}

func (s *SMiddleware) UserPortalPreAuth(r *ghttp.Request) {
	customCtx := _interface.BizCtx().Get(r.Context())
	list := query.GetActiveMerchantInfoList(r.Context())
	if len(list) == 0 {
		g.Log().Infof(r.Context(), "UserPortalPreAuth Merchant Need Init")
		utility.JsonRedirectExit(r, 61, "Merchant Not Found", s.LoginUrl)
		r.Exit()
	} else if len(list) == 1 {
		//SingleMerchant
		customCtx.MerchantId = list[0].Id
		r.Middleware.Next()
		return
	} else {
		if consts.GetConfigInstance().IsServerDev() || consts.GetConfigInstance().IsLocal() {
			customCtx.MerchantId = 15621
			r.Middleware.Next()
			return
		}
		one := query.GetMerchantInfoByHost(r.Context(), r.GetHost())
		if one == nil {
			//try match merchant from origin
			origin := r.GetHeader("Origin")
			if len(origin) > 0 {
				g.Log().Infof(r.Context(), "UserPortalPreAuth Try Extract Domain From Origin:%s", origin)
				parsedURL, err := url.Parse(origin)
				if err == nil {
					// Extract the host (domain) from the parsed URL
					domain := parsedURL.Hostname()
					one = query.GetMerchantInfoByHost(r.Context(), domain)
				}
			}
		}
		if one == nil {
			g.Log().Infof(r.Context(), "UserPortalPreAuth Merchant Not Found For Host:%s", r.GetHost())
			utility.JsonRedirectExit(r, 61, "Merchant Not Found", s.LoginUrl)
			r.Exit()
		} else {
			customCtx.MerchantId = one.Id
			r.Middleware.Next()
			return
		}
	}
}

func (s *SMiddleware) TokenAuth(r *ghttp.Request) {
	customCtx := _interface.BizCtx().Get(r.Context())
	if consts.GetConfigInstance().IsServerDev() {
		customCtx.MerchantId = 15621
		r.Middleware.Next()
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		g.Log().Infof(r.Context(), "TokenAuth empty token string of auth header")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // remove Bearer
	}
	if jwt.IsPortalToken(tokenString) {
		// Portal Call
		if !jwt.IsAuthTokenExpired(r.Context(), tokenString) {
			g.Log().Infof(r.Context(), "TokenAuth Invalid Token:%s", tokenString)
			utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
			r.Exit()
		}

		token := jwt.ParsePortalToken(tokenString)
		g.Log().Infof(r.Context(), "Parsed Token: %s, URL: %s", utility.MarshalToJsonString(token), r.GetUrl())

		if token.TokenType == jwt.TOKENTYPEUSER {
			userAccount := query.GetUserAccountById(r.Context(), token.Id)
			if userAccount == nil {
				g.Log().Infof(r.Context(), "TokenAuth user not found :%v", utility.MarshalToJsonString(token))
				utility.JsonRedirectExit(r, 61, "user not found", s.LoginUrl)
				r.Exit()
			}
			customCtx.User = &model.ContextUser{
				Id:         token.Id,
				Token:      tokenString,
				MerchantId: userAccount.MerchantId,
				Email:      token.Email,
			}
			customCtx.MerchantId = userAccount.MerchantId
		} else if token.TokenType == jwt.TOKENTYPEMERCHANTUSER {
			merchantAccount := query.GetMerchantUserAccountById(r.Context(), token.Id)
			if merchantAccount == nil {
				g.Log().Infof(r.Context(), "TokenMerchantAuth merchant user not found token:%s", utility.MarshalToJsonString(token))
				utility.JsonRedirectExit(r, 61, "merchant user not found", s.LoginUrl)
				r.Exit()
			}

			customCtx.MerchantUser = &model.ContextMerchantUser{
				Id:         token.Id,
				MerchantId: merchantAccount.MerchantId,
				Token:      tokenString,
				Email:      token.Email,
			}
			customCtx.MerchantId = merchantAccount.MerchantId
		} else {
			g.Log().Infof(r.Context(), "TokenAuth invalid token type token:%v", utility.MarshalToJsonString(token))
			utility.JsonRedirectExit(r, 61, "invalid token type", s.LoginUrl)
			r.Exit()
		}
		//Reset Expire Time
		jwt.ResetAuthTokenTTL(r.Context(), tokenString)
	} else {
		// Api Call
		merchantInfo := query.GetMerchantInfoByApiKey(r.Context(), tokenString)
		utility.Assert(merchantInfo != nil, "invalid api key")
		customCtx.MerchantId = merchantInfo.Id
	}

	r.Middleware.Next()
}
