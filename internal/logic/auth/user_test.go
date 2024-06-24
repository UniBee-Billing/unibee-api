package auth

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

func TestUserCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	var one *entity.UserAccount
	var err error
	t.Run("Test for User Create|Login|ChangePassword|Frozen|Release", func(t *testing.T) {
		one, err = CreateUser(ctx, &NewReq{
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
		one, token := PasswordLogin(ctx, one.MerchantId, one.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		ChangeUserPassword(ctx, one.MerchantId, one.Email, "test123456", "test654321")
		one, token = PasswordLogin(ctx, one.MerchantId, one.Email, "test654321")
		require.NotNil(t, one)
		require.NotNil(t, token)
		ChangeUserPasswordWithOutOldVerify(ctx, one.MerchantId, one.Email, "test123456")
		one, token = PasswordLogin(ctx, one.MerchantId, one.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		another, err := QueryOrCreateUser(ctx, &NewReq{
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
		require.NotNil(t, another)
		require.Equal(t, one.Id, another.Id)
		FrozenUser(ctx, int64(one.Id))
		one = query.GetUserAccountById(ctx, one.Id)
		require.NotNil(t, one)
		require.NotNil(t, one.Status == 2)
		ReleaseUser(ctx, int64(one.Id))
		one = query.GetUserAccountById(ctx, one.Id)
		require.NotNil(t, one)
		require.NotNil(t, one.Status == 0)
		list, err := UserList(ctx, &UserListInternalReq{
			MerchantId: test.TestMerchant.Id,
			UserId:     int64(one.Id),
			Email:      "autotestuser@wowow.io",
			SortType:   "desc",
			SortField:  "gmt_create",
			FirstName:  "test",
			LastName:   "test",
			Status:     []int{0, 2},
			Page:       -1,
		})
		require.Nil(t, err)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list.UserAccounts))
		res, err := SearchUser(ctx, test.TestMerchant.Id, "autotestuser@wowow.io")
		if err != nil {
			return
		}
		require.Nil(t, err)
		require.NotNil(t, list)
		require.Equal(t, 1, len(res))
	})
	t.Run("Test For User HardDelete", func(t *testing.T) {
		err = HardDeleteUser(ctx, one.Id)
		require.Nil(t, err)
		one = query.GetUserAccountById(ctx, one.Id)
		require.Nil(t, one)
	})
}
