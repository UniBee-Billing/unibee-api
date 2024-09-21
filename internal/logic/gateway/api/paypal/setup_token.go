package paypal

import (
	"context"
	"fmt"
)

func (c *Client) NewSetupTokens(ctx context.Context, customer *Customer, paymentSource *PaymentSource, requestID string) (*VaultToken, error) {
	type createSetupTokenRequest struct {
		PaymentSource *PaymentSource `json:"payment_source,omitempty"`
		Customer      *Customer      `json:"customer,omitempty"`
	}

	setupToken := &VaultToken{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v3/vault/setup-tokens"), createSetupTokenRequest{PaymentSource: paymentSource, Customer: customer})
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
