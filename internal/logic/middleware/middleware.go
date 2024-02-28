package middleware

import (
	"fmt"
	"net/url"
	"strings"
	"unibee/internal/consts"
	_ "unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/jwt"
	"unibee/internal/model"
	"unibee/internal/query"
	utility "unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SMiddleware struct {
	LoginUrl string
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

		message := err.Error()
		if strings.Contains(message, "Session Expired") {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				utility.OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), "Session Expired")
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				utility.JsonRedirectExit(r, 61, "Session Expired", s.LoginUrl)
			}
		} else if strings.Contains(message, utility.SystemAssertPrefix) || code == gcode.CodeValidationFailed {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				utility.OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				utility.JsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			}
		} else {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				utility.OpenApiJsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", _interface.BizCtx().Get(r.Context()).RequestId, code.Code()))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				utility.JsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", _interface.BizCtx().Get(r.Context()).RequestId, code.Code()))
			}
		}
	} else {
		r.Response.Status = 200
		if customCtx.IsOpenApiCall {
			utility.OpenApiJsonExit(r, code.Code(), "", res)
		} else {
			utility.JsonExit(r, code.Code(), "", res)
		}
	}
}

func (s *SMiddleware) OpenApiDetach(r *ghttp.Request) {
	customCtx := _interface.BizCtx().Get(r.Context())
	userAgent := r.Header.Get("User-Agent")
	if len(userAgent) > 0 && strings.Contains(userAgent, "OpenAPI") {
		customCtx.IsOpenApiCall = true
	}
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 0 && strings.HasPrefix(tokenString, "Bearer ") {
		customCtx.IsOpenApiCall = true
	}
	r.Middleware.Next()
}

func (s *SMiddleware) UserPortalPreAuth(r *ghttp.Request) {
	customCtx := _interface.BizCtx().Get(r.Context())
	list := query.GetActiveMerchantInfoList(r.Context())
	userAgent := r.Header.Get("User-Agent")
	if len(userAgent) > 0 && strings.Contains(userAgent, "OpenAPI") {
		customCtx.IsOpenApiCall = true
	}
	tokenString := r.Header.Get("Authorization")
	if customCtx.IsOpenApiCall == true || (len(tokenString) > 0 && strings.HasPrefix(tokenString, "Bearer ")) {
		g.Log().Infof(r.Context(), "UserPortal Api Not Support OpenApi Call")
		utility.JsonRedirectExit(r, 61, "UserPortal Api Not Support OpenApi Call", s.LoginUrl)
		r.Exit()
	}
	if len(list) == 0 {
		g.Log().Infof(r.Context(), "UserPortalPreAuth Merchant Need Init")
		utility.JsonRedirectExit(r, 61, "Merchant Not Found", s.LoginUrl)
		r.Exit()
	} else if len(list) == 1 {
		//SingleMerchant
		g.Log().Infof(r.Context(), "UserPortalPreAuth SingleMerchant")
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
			//try match merchant from Http Origin
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
			utility.JsonRedirectExit(r, 61, "Merchant Not Ready", s.LoginUrl)
			r.Exit()
		} else {
			g.Log().Infof(r.Context(), "UserPortalPreAuth Checked Merchant:%d", one.Id)
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
	userAgent := r.Header.Get("User-Agent")
	if len(userAgent) > 0 && strings.Contains(userAgent, "OpenAPI") {
		customCtx.IsOpenApiCall = true
	}
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		g.Log().Infof(r.Context(), "TokenAuth empty token string of auth header")
		if customCtx.IsOpenApiCall {
			r.Response.Status = 401
			utility.OpenApiJsonExit(r, 61, "invalid token")
		} else {
			utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		}
		r.Exit()
	}
	if strings.HasPrefix(tokenString, "Bearer ") {
		customCtx.IsOpenApiCall = true
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // remove Bearer
	}
	if !customCtx.IsOpenApiCall || jwt.IsPortalToken(tokenString) {
		// Portal Call
		if !jwt.IsAuthTokenExpired(r.Context(), tokenString) {
			g.Log().Infof(r.Context(), "TokenAuth Invalid Token:%s", tokenString)
			utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
			r.Exit()
		}

		token := jwt.ParsePortalToken(tokenString)
		g.Log().Infof(r.Context(), "TokenAuth Parsed Token: %s, URL: %s", utility.MarshalToJsonString(token), r.GetUrl())

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
				g.Log().Infof(r.Context(), "TokenAuth merchant user not found token:%s", utility.MarshalToJsonString(token))
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
		customCtx.IsOpenApiCall = true
		merchantInfo := query.GetMerchantInfoByApiKey(r.Context(), tokenString)
		if merchantInfo == nil {
			r.Response.Status = 401
			utility.OpenApiJsonExit(r, 61, "invalid token")
		} else {
			customCtx.MerchantId = merchantInfo.Id
		}
	}

	r.Middleware.Next()
}
