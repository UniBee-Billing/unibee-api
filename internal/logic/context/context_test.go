package context

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	_interface "unibee/internal/interface"
	"unibee/internal/model"
	"unibee/test"
	"unibee/utility"
)

func TestContext(t *testing.T) {
	request, err := http.NewRequestWithContext(context.Background(), "Get", "http://api.unibee.top", nil)
	require.Nil(t, err)
	var r = &ghttp.Request{
		Request: request,
	}
	t.Run("Test for Request context ", func(t *testing.T) {
		customCtx := &model.Context{
			Session: r.Session,
			Data:    make(g.Map),
		}
		customCtx.RequestId = utility.CreateRequestId()
		_interface.Context().Init(r, customCtx)
		require.NotNil(t, _interface.Context().Get(r.Context()).RequestId)
		require.Equal(t, uint64(0), _interface.Context().Get(r.Context()).MerchantId)
		require.Nil(t, _interface.Context().Get(r.Context()).User)
		require.Nil(t, _interface.Context().Get(r.Context()).MerchantMember)
		_interface.Context().SetUser(r.Context(), &model.ContextUser{
			Id:         test.TestUser.Id,
			MerchantId: test.TestMerchant.Id,
			Email:      test.TestUser.Email,
		})
		_interface.Context().SetMerchantMember(r.Context(), &model.ContextMerchantMember{
			Id:         test.TestMerchantMember.Id,
			MerchantId: test.TestMerchant.Id,
			Email:      test.TestMerchantMember.Email,
			IsOwner:    true,
		})
		_interface.Context().SetData(r.Context(), g.Map{})
		_interface.Context().Get(r.Context()).MerchantId = test.TestMerchant.Id
		require.Equal(t, test.TestMerchant.Id, _interface.Context().Get(r.Context()).MerchantId)
		require.Equal(t, test.TestUser.Id, _interface.Context().Get(r.Context()).User.Id)
		require.Equal(t, test.TestUser.Email, _interface.Context().Get(r.Context()).User.Email)
		require.Equal(t, test.TestMerchantMember.Id, _interface.Context().Get(r.Context()).MerchantMember.Id)
		require.Equal(t, test.TestMerchantMember.Email, _interface.Context().Get(r.Context()).MerchantMember.Email)
	})
}
