package middleware

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	_ "go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/internal/model"
	"go-oversea-pay/internal/query"
	utility "go-oversea-pay/utility"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
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
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// ResponseHandler 返回处理中间件
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
		g.Log().Errorf(r.Context(), "Global_exception panic url: %s params:%s code:%d error:%s", r.GetUrl(), json, err)
		return
	})

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)

	// 如果已经有返回内容，那么该中间件什么也不做
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
		r.Response.ClearBuffer() // 出现 panic 情况框架会自己写入非 json 的返回，需先清除
		r.Response.Status = 200  // 发生错误时候错误码Http 状态吗设置成 200，错误以 Json 形式返回
		message := err.Error()
		if strings.Contains(message, utility.SystemAssertPrefix) || code == gcode.CodeValidationFailed {
			utility.JsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
		} else {
			utility.JsonExit(r, code.Code(), fmt.Sprintf("System Error-%s-%d", _interface.BizCtx().Get(r.Context()).RequestId, code.Code()))
		}
	} else {
		r.Response.Status = 200
		utility.JsonExit(r, code.Code(), "", res)
	}
}

// PreAuth 从 Session 中获取用户
func (s *SMiddleware) PreAuth(r *ghttp.Request) {
	customCtx := _interface.BizCtx().Get(r.Context())
	if userEntity := _interface.Session().GetUser(r.Context()); userEntity != nil {
		customCtx.User = &model.ContextUser{
			Id: userEntity.Id,
			// MobilePhone: userEntity.Mobile,
			// UserName:    userEntity.UserName,
			// AvatarUrl:   userEntity.AvatarUrl,
			// IsAdmin:     false,
		}
	}
	r.Middleware.Next()
}

// PreOpenApiAuth 从 Session 中获取用户 (obsolete)
func (s *SMiddleware) PreOpenApiAuth(r *ghttp.Request) {

	customCtx := _interface.BizCtx().Get(r.Context())
	if userEntity := _interface.Session().GetUser(r.Context()); userEntity != nil {
		customCtx.User = &model.ContextUser{
			Id: userEntity.Id,
			// MobilePhone: userEntity.Mobile,
			// UserName:    userEntity.UserName,
			// AvatarUrl:   userEntity.AvatarUrl,
			// IsAdmin:     false,
		}
	}
	if key := r.GetHeader(consts.ApiKey); len(key) > 0 {
		//openapikey 转化为 api 用户
		customCtx.Data[consts.ApiKey] = key
		customCtx.OpenApiConfig = _interface.OpenApi().GetOpenApiConfig(r.Context(), key)
	}

	r.Middleware.Next()
}

// Auth 前台系统权限控制，用户必须登录才能访问
func (s *SMiddleware) Auth(r *ghttp.Request) {
	user := _interface.Session().GetUser(r.Context())
	if user == nil {
		_ = _interface.Session().SetNotice(r.Context(), &model.SessionNotice{
			Type:    consts.SessionNoticeTypeWarn,
			Content: "未登录或会话已过期，请您登录后再继续",
		})
		// 只有GET请求才支持保存当前URL，以便后续登录后再跳转回来。
		if r.Method == "GET" {
			_ = _interface.Session().SetLoginReferer(r.Context(), r.GetUrl())
		}
		// 根据当前请求方式执行不同的返回数据结构
		if r.IsAjaxRequest() {
			utility.JsonRedirectExit(r, 1, "", s.LoginUrl)
		} else {
			r.Response.RedirectTo(s.LoginUrl)
		}
	}
	r.Middleware.Next()
}

// later define a merchantClaim
type UserClaims struct {
	Id    uint64 `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Token Auth, 和上面的 Auth() 重复了, 但上面的Auth并非用在unibee项目中
var secretKey = []byte("3^&secret-key-for-UniBee*1!8*") // pass this as ENV, user_auth_login.go also uses this
func parseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return parsedAccessToken.Claims.(*UserClaims)
}

func (s *SMiddleware) TokenUserAuth(r *ghttp.Request) {
	if consts.GetConfigInstance().IsServerDev() {
		r.Middleware.Next()
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		g.Log().Errorf(r.Context(), "TokenUserAuth empty token string of auth header")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		g.Log().Errorf(r.Context(), "TokenUserAuth parse error:%s", err.Error())
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}

	if !auth.IsAuthTokenExpired(r.Context(), tokenString) {
		g.Log().Errorf(r.Context(), "TokenUserAuth token invalid")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}

	u := parseAccessToken(tokenString)
	g.Log().Infof(r.Context(), "Parsed User Token: %s, Email: %s, userId: %d", tokenString, u.Email, u.ID)

	userAccount := query.GetUserAccountById(r.Context(), u.Id)
	if userAccount == nil {
		g.Log().Errorf(r.Context(), "TokenUserAuth user not found")
		utility.JsonRedirectExit(r, 61, "user not found", s.LoginUrl)
		r.Exit()
	}

	customCtx := _interface.BizCtx().Get(r.Context())
	customCtx.User = &model.ContextUser{
		Id:    u.Id,
		Token: tokenString,
		Email: u.Email,
	}

	r.Middleware.Next()
}

func (s *SMiddleware) TokenMerchantAuth(r *ghttp.Request) {
	if consts.GetConfigInstance().IsServerDev() {
		r.Middleware.Next()
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		g.Log().Errorf(r.Context(), "TokenMerchantAuth empty token string of auth header")
		utility.JsonRedirectExit(r, 61, "invalid merchant token", s.LoginUrl)
		r.Exit()
	}
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		g.Log().Errorf(r.Context(), "TokenMerchantAuth parse error:%s", err.Error())
		utility.JsonRedirectExit(r, 61, "invalid merchant token", s.LoginUrl)
		r.Exit()
	}

	if !auth.IsAuthTokenExpired(r.Context(), tokenString) {
		g.Log().Errorf(r.Context(), "TokenMerchantAuth token invalid")
		utility.JsonRedirectExit(r, 61, "invalid merchant token", s.LoginUrl)
		r.Exit()
	}

	u := parseAccessToken(tokenString)

	g.Log().Infof(r.Context(), "Parsed Merchant Token: %s, Email: %s, merchantUserId: %d", tokenString, u.Email, u.ID)

	merchantAccount := query.GetMerchantAccountById(r.Context(), u.Id)
	if merchantAccount == nil {
		g.Log().Errorf(r.Context(), "TokenMerchantAuth merchant user not found")
		utility.JsonRedirectExit(r, 61, "merchant user not found", s.LoginUrl)
		r.Exit()
	}
	//有接口调用，顺延 5 分钟失效时间
	auth.SetAuthTokenNewTTL(r.Context(), tokenString, 5*60)

	customCtx := _interface.BizCtx().Get(r.Context())
	customCtx.MerchantUser = &model.ContextMerchantUser{
		Id:         u.Id,
		MerchantId: uint64(merchantAccount.MerchantId),
		Token:      tokenString,
		Email:      u.Email,
	}

	r.Middleware.Next()
}

func (s *SMiddleware) ApiAuth(r *ghttp.Request) {
	openApiConfig := _interface.BizCtx().Get(r.Context()).OpenApiConfig
	if openApiConfig == nil {
		if key := _interface.BizCtx().Get(r.Context()).Data[consts.ApiKey]; key == nil {
			utility.Json(r, 401, "key require in header")
		} else {
			utility.Json(r, 401, "invalid key")
		}
	}
	r.Middleware.Next()
}
