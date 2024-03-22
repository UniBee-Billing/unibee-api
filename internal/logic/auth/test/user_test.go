package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/logic/auth"
	"unibee/internal/query"
	"unibee/test"
)

func TestUserCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for User Create|Login|ChangePassword|Delete", func(t *testing.T) {
		one, err := auth.CreateUser(ctx, &auth.NewReq{
			ExternalUserId: "auto_x_2",
			Email:          "autotestuser@wowow.io",
			FirstName:      "test",
			LastName:       "test",
			Password:       "test123456",
			Phone:          "test",
			Address:        "test",
			UserName:       "test",
			CountryCode:    "CN",
			CountryName:    "China",
			MerchantId:     test.TestMerchant.Id,
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		one = query.GetUserAccountById(ctx, one.Id)
		require.NotNil(t, one)
		one, token := auth.PasswordLogin(ctx, one.MerchantId, one.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		auth.ChangeUserPassword(ctx, one.MerchantId, one.Email, "test123456", "test654321")
		one, token = auth.PasswordLogin(ctx, one.MerchantId, one.Email, "test654321")
		require.NotNil(t, one)
		require.NotNil(t, token)
		auth.ChangeUserPasswordWithOutOldVerify(ctx, one.MerchantId, one.Email, "test123456")
		one, token = auth.PasswordLogin(ctx, one.MerchantId, one.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		err = auth.HardDeleteUser(ctx, one.Id)
		require.Nil(t, err)
		one = query.GetUserAccountById(ctx, one.Id)
		require.Nil(t, one)
	})
}
