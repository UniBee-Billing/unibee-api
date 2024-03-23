package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/merchant/invoice"
	"unibee/test"
)

func TestInvoice(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for invoice create|edit", func(t *testing.T) {
		res, err := CreateInvoice(ctx, &invoice.NewReq{
			UserId:    test.TestUser.Id,
			TaxScale:  0,
			GatewayId: 0,
			Currency:  "USD",
			Name:      "test_invoice",
			Lines:     nil,
			Finish:    false,
		})
		require.Nil(t, err)
		require.NotNil(t, res)
	})
}
