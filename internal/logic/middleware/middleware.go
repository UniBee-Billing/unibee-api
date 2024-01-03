package middleware

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	_ "go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/model"
	utility "go-oversea-pay/utility"

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
		g.Log().Errorf(r.Context(), "Global_exception err url: %s params:%s code:%d error:%s", r.GetUrl(), json, code.Code(), err.Error())
		//if r.IsAjaxRequest() {
		r.Response.ClearBuffer() // 出现 panic 情况框架会自己写入非 json 的返回，需先清除
		utility.JsonExit(r, code.Code(), err.Error())
		//} else {
		//interface.View().Render500(r.Context(), model.Vie w{
		//	Error: err.Error(),
		//})
		//}
	} else {
		//if r.IsAjaxRequest() {
		utility.JsonExit(r, code.Code(), "", res)
		//} else {
		// 什么都不做，业务API自行处理模板渲染的成功逻辑。
		//}
	}
}

// PreAuth 从 Session 中获取用户
func (s *SMiddleware) PreAuth(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	_interface.BizCtx().Init(r, customCtx)
	if userEntity := _interface.Session().GetUser(r.Context()); userEntity != nil {
		customCtx.User = &model.ContextUser{
			Id:          userEntity.Id,
			MobilePhone: userEntity.Mobile,
			UserName:    userEntity.UserName,
			AvatarUrl:   userEntity.AvatarUrl,
			IsAdmin:     false,
		}
	}
	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		consts.ContextKey: customCtx,
	})
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// PreOpenApiAuth 从 Session 中获取用户
func (s *SMiddleware) PreOpenApiAuth(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	_interface.BizCtx().Init(r, customCtx)
	if userEntity := _interface.Session().GetUser(r.Context()); userEntity != nil {
		customCtx.User = &model.ContextUser{
			Id:          userEntity.Id,
			MobilePhone: userEntity.Mobile,
			UserName:    userEntity.UserName,
			AvatarUrl:   userEntity.AvatarUrl,
			IsAdmin:     false,
		}
	}
	if key := r.GetHeader(consts.ApiKey); len(key) > 0 {
		//openapikey 转化未api 用户
		customCtx.Data[consts.ApiKey] = key
		customCtx.OpenApiConfig = _interface.OpenApi().GetOpenApiConfig(r.Context(), key)
	}
	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		consts.ContextKey: customCtx,
	})
	// 执行下一步请求逻辑
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
	// UserId
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

func (s *SMiddleware) TokenAuth(r *ghttp.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		fmt.Println("empty token string of auth header")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}
	// fmt.Println("token str: ", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("parse error")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}

	if !token.Valid {
		fmt.Println("token invalid")
		utility.JsonRedirectExit(r, 61, "invalid token", s.LoginUrl)
		r.Exit()
	}

	u := parseAccessToken(tokenString)
	fmt.Println("parsed token: ", u.Email)

	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	_interface.BizCtx().Init(r, customCtx)
	// if := _interface.Session().GetUser(r.Context()); userEntity != nil {
	customCtx.User = &model.ContextUser{
		/*
			Id:          userEntity.Id,
			MobilePhone: userEntity.Mobile,
			UserName:    userEntity.UserName,
			AvatarUrl:   userEntity.AvatarUrl,
			IsAdmin:     false,
		*/
		Email: u.Email,
	}
	// }
	// if key := r.GetHeader(consts.ApiKey); len(key) > 0 {
	//openapikey 转化为api 用户
	// customCtx.Data[consts.ApiKey] = key
	// customCtx.OpenApiConfig = _interface.OpenApi().GetOpenApiConfig(r.Context(), key)
	// }
	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		consts.ContextKey: customCtx,
	})

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
