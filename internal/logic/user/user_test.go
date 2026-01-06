package user

import (
	"context"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/test"
)

func TestUserCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	var one *entity.UserAccount
	var err error
	t.Run("Test for User Create|Login|ChangePassword|Frozen|Release", func(t *testing.T) {
		one, err = CreateUser(ctx, &NewUserInternalReq{
			ExternalUserId: "auto_x_2",
			Email:          "autotestuser@wowow.io",
			FirstName:      "test",
			LastName:       "test",
			Password:       "test123456",
			Phone:          "test",
			Address:        "test",
			UserName:       "test",
			CountryCode:    "CN",
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
		another, err := QueryOrCreateUser(ctx, &NewUserInternalReq{
			ExternalUserId: "auto_x_2",
			Email:          "autotestuser@wowow.io",
			FirstName:      "test",
			LastName:       "test",
			Password:       "test123456",
			Phone:          "test",
			Address:        "test",
			UserName:       "test",
			CountryCode:    "CN",
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

// key:QV6GZORZI54F4A3NGQRKH76K3YBPT6ME
// url:otpauth://totp/Unibee:jack.fu@wowow.io?algorithm=SHA1&digits=6&issuer=Unibee&period=30&secret=QV6GZORZI54F4A3NGQRKH76K3YBPT6ME
func TestTotp(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Unibee",           // Application Name
		AccountName: "jack.fu@wowow.io", // Account
	})
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		return
	}
	fmt.Printf("key:%s\n", key.Secret())
	fmt.Printf("url:%s\n", key.URL())
}

func TestTotpCode(t *testing.T) {
	validateResult, err := totp.ValidateCustom("029844", "QV6GZORZI54F4A3NGQRKH76K3YBPT6ME", time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		return
	}
	fmt.Println(validateResult)
}
