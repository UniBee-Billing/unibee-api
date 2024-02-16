package session

import (
	"context"
	"unibee-api/internal/interface"
	"unibee-api/internal/model"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type sSession struct{}

const (
	sessionKeyUser         = "SessionKeyUser"
	sessionKeyLoginReferer = "SessionKeyReferer"
	sessionKeyNotice       = "SessionKeyNotice"
)

func init() {
	_interface.RegisterSession(New())
}

func New() *sSession {
	return &sSession{}
}

// SetUser 设置用户Session.
func (s *sSession) SetUser(ctx context.Context, user *entity.UserAccount) error {
	return _interface.BizCtx().Get(ctx).Session.Set(sessionKeyUser, user)
}

// GetUser 获取当前登录的用户信息对象，如果用户未登录返回nil。
func (s *sSession) GetUser(ctx context.Context) *entity.UserAccount {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		v, _ := customCtx.Session.Get(sessionKeyUser)
		if !v.IsNil() {
			var user *entity.UserAccount
			_ = v.Struct(&user)
			return user
		}
	}
	return nil
}

// RemoveUser 删除用户Session。
func (s *sSession) RemoveUser(ctx context.Context) error {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(sessionKeyUser)
	}
	return nil
}

// SetLoginReferer 设置LoginReferer.
func (s *sSession) SetLoginReferer(ctx context.Context, referer string) error {
	if s.GetLoginReferer(ctx) == "" {
		customCtx := _interface.BizCtx().Get(ctx)
		if customCtx != nil {
			return customCtx.Session.Set(sessionKeyLoginReferer, referer)
		}
	}
	return nil
}

// GetLoginReferer 获取LoginReferer.
func (s *sSession) GetLoginReferer(ctx context.Context) string {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.MustGet(sessionKeyLoginReferer).String()
	}
	return ""
}

// RemoveLoginReferer 删除LoginReferer.
func (s *sSession) RemoveLoginReferer(ctx context.Context) error {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(sessionKeyLoginReferer)
	}
	return nil
}

// SetNotice 设置Notice
func (s *sSession) SetNotice(ctx context.Context, message *model.SessionNotice) error {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Set(sessionKeyNotice, message)
	}
	return nil
}

// GetNotice 获取Notice
func (s *sSession) GetNotice(ctx context.Context) (*model.SessionNotice, error) {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		var message *model.SessionNotice
		v, err := customCtx.Session.Get(sessionKeyNotice)
		if err != nil {
			return nil, err
		}
		if err = v.Scan(&message); err != nil {
			return nil, err
		}
		return message, nil
	}
	return nil, nil
}

// RemoveNotice 删除Notice
func (s *sSession) RemoveNotice(ctx context.Context) error {
	customCtx := _interface.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(sessionKeyNotice)
	}
	return nil
}
