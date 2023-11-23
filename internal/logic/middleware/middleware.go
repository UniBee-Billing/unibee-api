package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/consts"
	_ "go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/model"
	"go-oversea-pay/internal/service"
	response "go-oversea-pay/utility"
)

type SMiddleware struct {
	LoginUrl string // 登录路由地址
}

func init() {
	service.RegisterMiddleware(New())
}

func New() *SMiddleware {
	return &SMiddleware{
		LoginUrl: "/login",
	}
}

// ResponseHandler 返回处理中间件
func (s *SMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有返回内容，那么该中间件什么也不做
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)
	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		//if r.IsAjaxRequest() {
		response.JsonExit(r, code.Code(), err.Error())
		//} else {
		//service.View().Render500(r.Context(), model.Vie w{
		//	Error: err.Error(),
		//})
		//}
	} else {
		//if r.IsAjaxRequest() {
		response.JsonExit(r, code.Code(), "", res)
		//} else {
		// 什么都不做，业务API自行处理模板渲染的成功逻辑。
		//}
	}
}

// Ctx 自定义上下文对象
func (s *SMiddleware) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	service.BizCtx().Init(r, customCtx)
	if userEntity := service.Session().GetUser(r.Context()); userEntity != nil {
		customCtx.User = &model.ContextUser{
			Id:          userEntity.Id,
			MobilePhone: userEntity.MobilePhone,
			UserName:    userEntity.UserName,
			AvatarUrl:   userEntity.AvatarUrl,
			IsAdmin:     false,
		}
	}
	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		"Context": customCtx,
	})
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// Auth 前台系统权限控制，用户必须登录才能访问
func (s *SMiddleware) Auth(r *ghttp.Request) {
	user := service.Session().GetUser(r.Context())
	if user == nil {
		_ = service.Session().SetNotice(r.Context(), &model.SessionNotice{
			Type:    consts.SessionNoticeTypeWarn,
			Content: "未登录或会话已过期，请您登录后再继续",
		})
		// 只有GET请求才支持保存当前URL，以便后续登录后再跳转回来。
		if r.Method == "GET" {
			_ = service.Session().SetLoginReferer(r.Context(), r.GetUrl())
		}
		// 根据当前请求方式执行不同的返回数据结构
		if r.IsAjaxRequest() {
			response.JsonRedirectExit(r, 1, "", s.LoginUrl)
		} else {
			r.Response.RedirectTo(s.LoginUrl)
		}
	}
	r.Middleware.Next()
}
