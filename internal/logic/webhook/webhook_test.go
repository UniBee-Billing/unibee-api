package webhook

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

func TestWebhook(t *testing.T) {
	ctx := context.Background()
	var one *entity.MerchantWebhook
	var err error
	t.Run("Test for webhook endpoint Create|Edit|Delete", func(t *testing.T) {
		list := MerchantWebhookEndpointList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 0, len(list))
		one, err = NewMerchantWebhookEndpoint(ctx, test.TestMerchant.Id, "http://test.endpoint.unibee.api", []string{})
		require.Nil(t, err)
		one = query.GetMerchantWebhook(ctx, one.Id)
		require.NotNil(t, one)
		require.Equal(t, "http://test.endpoint.unibee.api", one.WebhookUrl)
		require.Equal(t, 0, len(one.WebhookEvents))
		list = MerchantWebhookEndpointList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		err = UpdateMerchantWebhookEndpoint(ctx, test.TestMerchant.Id, one.Id, "http://test2.endpoint.unibee.api", []string{"subscription.created"})
		require.Nil(t, err)
		one = query.GetMerchantWebhook(ctx, one.Id)
		require.NotNil(t, one)
		require.Equal(t, "http://test2.endpoint.unibee.api", one.WebhookUrl)
		require.Equal(t, "subscription.created", one.WebhookEvents)

		//log list
		logList := MerchantWebhookEndpointLogList(ctx, &EndpointLogListInternalReq{
			MerchantId: one.MerchantId,
			EndpointId: one.Id,
			Page:       -1,
		})
		require.NotNil(t, logList)
		require.Equal(t, 0, len(logList))

		list = MerchantWebhookEndpointList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		require.Equal(t, "http://test2.endpoint.unibee.api", list[0].WebhookUrl)
		require.Equal(t, 1, len(list[0].WebhookEvents))
		err = DeleteMerchantWebhookEndpoint(ctx, one.MerchantId, one.Id)
		require.Nil(t, err)
		one, err = NewMerchantWebhookEndpoint(ctx, test.TestMerchant.Id, "http://test2.endpoint.unibee.api", []string{})
		require.Nil(t, err)
		list = MerchantWebhookEndpointList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		require.Equal(t, "http://test2.endpoint.unibee.api", list[0].WebhookUrl)
		require.Equal(t, 0, len(list[0].WebhookEvents))
	})
	t.Run("Test for webhook HardDelete", func(t *testing.T) {
		err = HardDeleteMerchantWebhookEndpoint(ctx, one.MerchantId, one.Id)
		require.Nil(t, err)
	})
}
