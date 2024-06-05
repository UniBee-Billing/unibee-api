package paypal

import (
	"context"
	"fmt"
)

func (c *Client) GetPaymentMethodTokens(ctx context.Context, customerId string) (*PaymentMethodToken, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v3/vault/payment-tokens?customer_id=%s", c.APIBase, customerId), nil)
	if err != nil {
		return nil, err
	}

	response := &PaymentMethodToken{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}
