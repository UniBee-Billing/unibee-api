package merchant_config

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/test"
)

func TestMerchantConfig(t *testing.T) {
	ctx := context.Background()
	var one *entity.MerchantConfig
	var err error
	t.Run("Test for Merchant Config Set|Get", func(t *testing.T) {
		one = GetMerchantConfig(ctx, test.TestMerchant.Id, "test_config_key")
		require.Equal(t, true, one == nil || len(one.ConfigValue) == 0)
		err = SetMerchantConfig(ctx, test.TestMerchant.Id, "test_config_key", "test")
		require.Nil(t, err)
		one = GetMerchantConfig(ctx, test.TestMerchant.Id, "test_config_key")
		require.Equal(t, true, one != nil && len(one.ConfigValue) > 0)
	})
	t.Run("Test for Clean Merchant Config", func(t *testing.T) {
		err = SetMerchantConfig(ctx, test.TestMerchant.Id, "test_config_key", "")
		require.Nil(t, err)
		one = GetMerchantConfig(ctx, test.TestMerchant.Id, "test_config_key")
		require.Equal(t, true, one == nil || len(one.ConfigValue) == 0)
	})
}
