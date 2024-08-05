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

func (c *Client) NewPaymentTokens(ctx context.Context, paymentSource *PaymentSource, requestID string) (*VaultToken, error) {
	type createPaymentTokenRequest struct {
		PaymentSource *PaymentSource `json:"payment_source,omitempty"`
	}

	setupToken := &VaultToken{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v3/vault/payment-tokens"), createPaymentTokenRequest{PaymentSource: paymentSource})
	if err != nil {
		return setupToken, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, setupToken); err != nil {
		return setupToken, err
	}

	return setupToken, nil
}
