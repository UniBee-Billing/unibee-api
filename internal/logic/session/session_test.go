package session

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/merchant/session"
	"unibee/test"
)

func TestUserSessionTransfer(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for UserSession Create|Transfer", func(t *testing.T) {
		userSession, err := NewUserSession(ctx, test.TestMerchant.Id, &session.NewReq{
			ExternalUserId: "auto_x",
			Email:          "testuser@wowow.io",
		})
		require.Nil(t, err)
		require.NotNil(t, userSession)
		require.NotNil(t, userSession.ClientSession)
		one := UserSessionTransfer(ctx, userSession.ClientSession)
		require.Equal(t, one.Email, "testuser@wowow.io")
	})
}
