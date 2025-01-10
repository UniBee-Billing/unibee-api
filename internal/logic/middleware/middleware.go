package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"net/url"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/cmd/i18n"
	"unibee/internal/consts"
	_ "unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/interface/context"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/merchant"
	"unibee/internal/model"
	"unibee/internal/query"
	"unibee/utility"

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
	context.Context().Init(r, customCtx)
	r.Assigns(g.Map{
		consts.ContextKey: customCtx,
	})

	// Setup System Default Language
	r.SetCtx(gi18n.WithLanguage(r.Context(), "en"))
	customCtx.Language = "en"
	lang := ""
	if r.Get("lang") != nil {
		lang = r.Get("lang").String()
	}
	if len(lang) == 0 {
		lang = r.GetHeader("lang")
	}
	if len(lang) > 0 && i18n.IsLangAvailable(lang) {
		r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
		customCtx.Language = lang
	}

	customCtx.UserAgent = r.Header.Get("User-Agent")
	if len(customCtx.UserAgent) > 0 && strings.Contains(customCtx.UserAgent, "OpenAPI") {
		customCtx.IsOpenApiCall = true
	}
	customCtx.Authorization = r.Header.Get("Authorization")
	customCtx.TokenString = customCtx.Authorization
	if len(customCtx.TokenString) > 0 && strings.HasPrefix(customCtx.TokenString, "Bearer ") && !jwt.IsPortalToken(customCtx.TokenString) {
		customCtx.IsOpenApiCall = true
		customCtx.TokenString = strings.Replace(customCtx.TokenString, "Bearer ", "", 1) // remove Bearer
	}
	g.Log().Info(r.Context(), fmt.Sprintf("[Request][%s][%s][%s][%s] IsOpenApi:%v Token:%s Body:%s", customCtx.Language, customCtx.RequestId, r.Method, r.GetUrl(), customCtx.IsOpenApiCall, customCtx.TokenString, r.GetBodyString()))

	utility.Try(r.Middleware.Next, func(err interface{}) {
		g.Log().Errorf(r.Context(), "[Request][%s][%s][%s] Global_Exception Panic Body:%s Error:%v", customCtx.RequestId, r.Method, r.GetUrl(), r.GetBodyString(), err)
		return
	})
	g.Log().Info(r.Context(), fmt.Sprintf("[Request][%s][%s][%s] MerchantId:%d", customCtx.RequestId, r.Method, r.GetUrl(), customCtx.MerchantId))

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
		g.Log().Errorf(r.Context(), "Global_exception requestId:%s url: %s params:%s code:%d error:%s", context.Context().Get(r.Context()).RequestId, r.GetUrl(), json, code.Code(), err.Error())
		r.Response.ClearBuffer() // inner panic will contain json data，need clean

		message := err.Error()
		if strings.Contains(message, "Session Expired") {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				_interface.OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), "Session Expired")
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				_interface.JsonRedirectExit(r, 61, "Session Expired", s.LoginUrl)
			}
		} else if strings.Contains(message, utility.SystemAssertPrefix) || code == gcode.CodeValidationFailed {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				_interface.OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				_interface.JsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			}
		} else {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				_interface.OpenApiJsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", context.Context().Get(r.Context()).RequestId, code.Code()))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				_interface.JsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", context.Context().Get(r.Context()).RequestId, code.Code()))
			}
		}
	} else {
		r.Response.Status = 200
		if customCtx.IsOpenApiCall {
			_interface.OpenApiJsonExit(r, code.Code(), "", res)
		} else {
			_interface.JsonExit(r, code.Code(), "", res)
		}
	}
}

func (s *SMiddleware) MerchantHandler(r *ghttp.Request) {
	customCtx := context.Context().Get(r.Context())
	if len(customCtx.TokenString) == 0 {
		g.Log().Infof(r.Context(), "MerchantHandler empty token string of auth header")
		if customCtx.IsOpenApiCall {
			r.Response.Status = 401
			_interface.OpenApiJsonExit(r, 61, "invalid token")
		} else {
			_interface.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		}
		r.Exit()
	}
	if !customCtx.IsOpenApiCall {
		// Admin Portal Call
		customCtx.IsAdminPortalCall = true
		if !jwt.IsAuthTokenAvailable(r.Context(), customCtx.TokenString) {
			g.Log().Infof(r.Context(), "MerchantHandler Invalid Token:%s", customCtx.TokenString)
			_interface.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
			r.Exit()
		}

		customCtx.Token = jwt.ParsePortalToken(customCtx.TokenString)
		g.Log().Debugf(r.Context(), "MerchantHandler Parsed Token: %s, URL: %s", utility.MarshalToJsonString(customCtx.Token), r.GetUrl())

		if customCtx.Token.TokenType == jwt.TOKENTYPEMERCHANTMember {
			member := query.GetMerchantMemberById(r.Context(), customCtx.Token.Id)
			permissionKey := jwt.GetMemberPermissionKey(r.Context(), member)
			if member == nil {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member not found token:%s", utility.MarshalToJsonString(customCtx.Token))
				_interface.JsonRedirectExit(r, 61, "merchant user not found", s.LoginUrl)
				r.Exit()
			} else if member.Status == 2 {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member has suspend :%v", utility.MarshalToJsonString(customCtx.Token))
				_interface.JsonRedirectExit(r, 61, "Your account has been suspended. Please contact billing admin for further assistance.", s.LoginUrl)
				r.Exit()
			} else if strings.Compare(permissionKey, customCtx.Token.PermissionKey) != 0 && !strings.Contains(r.GetUrl(), "logout") {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member permission has change, need reLogin")
				_interface.JsonRedirectExit(r, 62, "Your permission has changed. Please reLogin.", s.LoginUrl)
				r.Exit()
			}

			customCtx.MerchantMember = &model.ContextMerchantMember{
				Id:         customCtx.Token.Id,
				MerchantId: customCtx.Token.MerchantId,
				Token:      customCtx.TokenString,
				Email:      customCtx.Token.Email,
				IsOwner:    strings.Compare(strings.Trim(member.Role, " "), "Owner") == 0,
			}
			customCtx.MerchantId = customCtx.Token.MerchantId
			doubleRequestLimit(strconv.FormatUint(customCtx.MerchantMember.Id, 10), r)
			lang := ""
			if r.Get("lang") != nil {
				lang = r.Get("lang").String()
			}
			if len(lang) == 0 {
				lang = r.GetHeader("lang")
			}
			if len(lang) > 0 && i18n.IsLangAvailable(lang) {
				r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
			}
		} else {
			g.Log().Infof(r.Context(), "MerchantHandler invalid token type token:%v", utility.MarshalToJsonString(customCtx.Token))
			_interface.JsonRedirectExit(r, 61, "invalid token type", s.LoginUrl)
			r.Exit()
		}
		//Reset Expire Time
		jwt.ResetAuthTokenTTL(r.Context(), customCtx.TokenString)
	} else {
		// Api Call
		customCtx.IsOpenApiCall = true
		merchantInfo := query.GetMerchantByApiKey(r.Context(), customCtx.TokenString)
		if merchantInfo == nil {
			merchantInfo = merchant.GetMerchantByOpenApiKeyFromCache(r.Context(), customCtx.TokenString)
		}
		if merchantInfo == nil {
			r.Response.Status = 401
			_interface.OpenApiJsonExit(r, 61, "invalid token")
		} else {
			customCtx.MerchantId = merchantInfo.Id
			customCtx.OpenApiKey = customCtx.TokenString
		}
		lang := ""
		if r.Get("lang") != nil {
			lang = r.Get("lang").String()
		}
		if len(lang) == 0 {
			lang = r.GetHeader("lang")
		}
		if len(lang) > 0 && i18n.IsLangAvailable(lang) {
			r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
		}
	}
	r.Middleware.Next()
}

func (s *SMiddleware) UserPortalMerchantRouterHandler(r *ghttp.Request) {
	customCtx := context.Context().Get(r.Context())
	list := query.GetActiveMerchantList(r.Context())
	if customCtx.IsOpenApiCall == true {
		g.Log().Infof(r.Context(), "UserPortal Api Not Support OpenApi Call")
		_interface.JsonRedirectExit(r, 61, "UserPortal Api Not Support OpenApi Call", s.LoginUrl)
		r.Exit()
	}
	if len(list) == 0 {
		g.Log().Infof(r.Context(), "UserPortalMerchantRouterHandler Merchant Need Init")
		_interface.JsonRedirectExit(r, 61, "Merchant Not Found", s.LoginUrl)
		r.Exit()
	} else if len(list) == 1 {
		//SingleMerchant
		g.Log().Infof(r.Context(), "UserPortalMerchantRouterHandler SingleMerchant")
		customCtx.MerchantId = list[0].Id
		r.Middleware.Next()
		return
	} else {
		if config.GetConfigInstance().IsServerDev() || config.GetConfigInstance().IsLocal() {
			customCtx.MerchantId = consts.CloudModeManagerMerchantId
			r.Middleware.Next()
			return
		}
		host := r.GetHost()
		if !config.GetConfigInstance().IsProd() && (host == "127.0.0.1" || host == "localhost") {
			host = "user.unibee.top"
		}
		one := query.GetMerchantByHost(r.Context(), host)
		if one == nil {
			//try match merchant from Http Origin
			origin := r.GetHeader("Origin")
			if len(origin) > 0 {
				g.Log().Infof(r.Context(), "UserPortalMerchantRouterHandler Try Extract Domain From Origin:%s", origin)
				parsedURL, err := url.Parse(origin)
				if err == nil {
					// Extract the host (domain) from the parsed URL
					domain := parsedURL.Hostname()
					if !config.GetConfigInstance().IsProd() && (domain == "127.0.0.1" || domain == "localhost") {
						domain = "user.unibee.top"
					}
					one = query.GetMerchantByHost(r.Context(), domain)
				}
			}
		}
		if one == nil {
			g.Log().Infof(r.Context(), "UserPortalMerchantRouterHandler Merchant Not Found For Host:%s", r.GetHost())
			_interface.JsonRedirectExit(r, 61, "Merchant Not Ready", s.LoginUrl)
			r.Exit()
		} else {
			g.Log().Infof(r.Context(), "UserPortalMerchantRouterHandler Checked Merchant:%d", one.Id)
			customCtx.MerchantId = one.Id
			r.Middleware.Next()
			return
		}
	}
}

func (s *SMiddleware) UserPortalApiHandler(r *ghttp.Request) {
	customCtx := context.Context().Get(r.Context())
	if len(customCtx.TokenString) == 0 {
		g.Log().Infof(r.Context(), "MerchantHandler empty token string of auth header")
		if customCtx.IsOpenApiCall {
			r.Response.Status = 401
			_interface.OpenApiJsonExit(r, 61, "invalid token")
		} else {
			_interface.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		}
		r.Exit()
	}
	if !customCtx.IsOpenApiCall {
		// Merchant Portal Call
		if !jwt.IsAuthTokenAvailable(r.Context(), customCtx.TokenString) {
			g.Log().Infof(r.Context(), "MerchantHandler Invalid Token:%s", customCtx.TokenString)
			_interface.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
			r.Exit()
		}

		customCtx.Token = jwt.ParsePortalToken(customCtx.TokenString)
		g.Log().Debugf(r.Context(), "MerchantHandler Parsed Token: %s, URL: %s", utility.MarshalToJsonString(customCtx.Token), r.GetUrl())
		if customCtx.Token.TokenType == jwt.TOKENTYPEUSER {
			userAccount := query.GetUserAccountById(r.Context(), customCtx.Token.Id)
			if userAccount == nil {
				g.Log().Infof(r.Context(), "MerchantHandler user not found :%v", utility.MarshalToJsonString(customCtx.Token))
				_interface.JsonRedirectExit(r, 61, "account not found", s.LoginUrl)
				r.Exit()
				return
			} else if userAccount.Status == 2 {
				g.Log().Infof(r.Context(), "MerchantHandler user has suspend :%v", utility.MarshalToJsonString(customCtx.Token))
				_interface.JsonRedirectExit(r, 61, "Your account has been suspended. Please contact billing admin for further assistance.", s.LoginUrl)
				r.Exit()
				return
			}
			customCtx.User = &model.ContextUser{
				Id:         customCtx.Token.Id,
				Token:      customCtx.TokenString,
				MerchantId: customCtx.Token.MerchantId,
				Email:      userAccount.Email,
				Lang:       userAccount.Language,
			}
			customCtx.MerchantId = customCtx.Token.MerchantId
			doubleRequestLimit(strconv.FormatUint(customCtx.User.Id, 10), r)
			//UserPortalTrack
			segment.TrackSegmentEventBackground(r.Context(), userAccount.MerchantId, userAccount, r.URL.Path, nil)
			lang := ""
			if r.Get("lang") != nil {
				lang = r.Get("lang").String()
			}
			if len(lang) == 0 {
				lang = r.GetHeader("lang")
			}
			if len(lang) > 0 && i18n.IsLangAvailable(lang) {
				r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
			} else if customCtx.User != nil && len(customCtx.User.Lang) > 0 && i18n.IsLangAvailable(customCtx.User.Lang) {
				r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(customCtx.User.Lang))))
			}
		} else {
			g.Log().Infof(r.Context(), "MerchantHandler invalid token type token:%v", utility.MarshalToJsonString(customCtx.Token))
			_interface.JsonRedirectExit(r, 61, "invalid token type", s.LoginUrl)
			r.Exit()
		}
		//Reset Expire Time
		jwt.ResetAuthTokenTTL(r.Context(), customCtx.TokenString)
	} else {
		g.Log().Infof(r.Context(), "UserPortal Api Not Support OpenApi Call")
		_interface.JsonRedirectExit(r, 61, "UserPortal Api Not Support OpenApi Call", s.LoginUrl)
		r.Exit()
	}
	r.Middleware.Next()
}

func doubleRequestLimit(id string, r *ghttp.Request) {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
		if strings.HasSuffix(r.GetUrl(), "detail") || strings.HasSuffix(r.GetUrl(), "list") || strings.HasSuffix(r.GetUrl(), "get") {
			return
		}
		md5 := utility.MD5(fmt.Sprintf("%s%s%s", id, r.GetUrl(), r.GetBodyString()))
		if !utility.TryLock(r.Context(), md5, 2) {
			utility.Assert(false, i18n.LocalizationFormat(r.Context(), "{#ClickTooFast}"))
		}
	}
}
