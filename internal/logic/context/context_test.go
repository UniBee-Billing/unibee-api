package context

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	"unibee/internal/model"
	"unibee/test"
	"unibee/utility"
)

func TestContext(t *testing.T) {
	//request, err := http.NewRequestWithContext(context.Background(), "Get", "http://api.unibee.top", nil)
	//require.Nil(t, err)
	//var r = &ghttp.Request{
	//	Request: request,
	//}
	ctx := context.WithValue(context.Background(), consts.ContextKey, &model.Context{
		Data: make(g.Map),
	})
	t.Run("Test for Request context ", func(t *testing.T) {
		//customCtx := &model.Context{
		//	//Session: r.Session,
		//	Data: make(g.Map),
		//}
		//customCtx.RequestId = utility.CreateRequestId()
		_interface.Context().Get(ctx).RequestId = utility.CreateRequestId()
		require.NotNil(t, _interface.Context().Get(ctx).RequestId)
		require.Equal(t, uint64(0), _interface.Context().Get(ctx).MerchantId)
		require.Nil(t, _interface.Context().Get(ctx).User)
		require.Nil(t, _interface.Context().Get(ctx).MerchantMember)
		_interface.Context().SetUser(ctx, &model.ContextUser{
			Id:         test.TestUser.Id,
			MerchantId: test.TestMerchant.Id,
			Email:      test.TestUser.Email,
		})
		_interface.Context().SetMerchantMember(ctx, &model.ContextMerchantMember{
			Id:         test.TestMerchantMember.Id,
			MerchantId: test.TestMerchant.Id,
			Email:      test.TestMerchantMember.Email,
			IsOwner:    true,
		})
		_interface.Context().Get(ctx).MerchantId = test.TestMerchant.Id
		require.Equal(t, test.TestMerchant.Id, _interface.Context().Get(ctx).MerchantId)
		require.Equal(t, test.TestUser.Id, _interface.Context().Get(ctx).User.Id)
		require.Equal(t, test.TestUser.Email, _interface.Context().Get(ctx).User.Email)
		require.Equal(t, test.TestMerchantMember.Id, _interface.Context().Get(ctx).MerchantMember.Id)
		require.Equal(t, test.TestMerchantMember.Email, _interface.Context().Get(ctx).MerchantMember.Email)
	})
}
