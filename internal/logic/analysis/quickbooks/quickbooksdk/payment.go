package quickbooksdk

import (
	"errors"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type Payment struct {
	Id                  string          `json:"Id,omitempty"`                  // Payment ID (for updates)
	SyncToken           string          `json:"SyncToken,omitempty"`           // Required for updates
	MetaData            *MetaData       `json:"MetaData,omitempty"`            // Returned after creation, includes creation and update times
	TxnDate             string          `json:"TxnDate,omitempty"`             // Payment date, format: YYYY-MM-DD
	TotalAmt            decimal.Decimal `json:"TotalAmt"`                      // Total amount (required)
	UnappliedAmt        float64         `json:"UnappliedAmt,omitempty"`        // Unapplied amount (auto-returned, cannot submit)
	PaymentRefNum       string          `json:"PaymentRefNum,omitempty"`       // Payment reference number (custom number)
	ProcessPayment      bool            `json:"ProcessPayment,omitempty"`      // Whether to process payment (default false)
	CustomerRef         *ReferenceType  `json:"CustomerRef"`                   // Customer (required)
	DepositToAccountRef *ReferenceType  `json:"DepositToAccountRef,omitempty"` // Optional, account for funds deposit
	Line                []PaymentLine   `json:"Line"`                          // Each payment line (for invoice association)
	Domain              string          `json:"domain,omitempty"`
}

type PaymentLine struct {
	Amount    decimal.Decimal `json:",omitempty"`
	LinkedTxn []LinkedTxn     `json:",omitempty"`
}

// CreatePayment creates the given payment within QuickBooks.
func (c *Client) CreatePayment(payment *Payment) (*Payment, error) {
	var resp struct {
		Payment Payment
		Time    Date
	}

	if err := c.post("payment", payment, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Payment, nil
}

// DeletePayment deletes the given payment from QuickBooks.
func (c *Client) DeletePayment(payment *Payment) error {
	if payment.Id == "" || payment.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("payment", payment, nil, map[string]string{"operation": "delete"})
}

// FindPayments gets the full list of Payments in the QuickBooks account.
func (c *Client) FindPayments() ([]Payment, error) {
	var resp struct {
		QueryResponse struct {
			Payments      []Payment `json:"Payment"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Payment", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no payments could be found")
	}

	payments := make([]Payment, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Payment ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Payments == nil {
			return nil, errors.New("no payments could be found")
		}

		payments = append(payments, resp.QueryResponse.Payments...)
	}

	return payments, nil
}

// FindPaymentById returns an payment with a given Id.
func (c *Client) FindPaymentById(id string) (*Payment, error) {
	var resp struct {
		Payment Payment
		Time    Date
	}

	if err := c.get("payment/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Payment, nil
}

// QueryPayments accepts a SQL query and returns all payments found using it.
func (c *Client) QueryPayments(query string) ([]Payment, error) {
	var resp struct {
		QueryResponse struct {
			Payments      []Payment `json:"Payment"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Payments == nil {
		return make([]Payment, 0), nil
	}

	return resp.QueryResponse.Payments, nil
}

// UpdatePayment updates the given payment in QuickBooks.
func (c *Client) UpdatePayment(payment *Payment) (*Payment, error) {
	if payment.Id == "" {
		return nil, errors.New("missing payment id")
	}

	existingPayment, err := c.FindPaymentById(payment.Id)
	if err != nil {
		return nil, err
	}

	payment.SyncToken = existingPayment.SyncToken

	payload := struct {
		*Payment
		Sparse bool `json:"sparse"`
	}{
		Payment: payment,
		Sparse:  true,
	}

	var paymentData struct {
		Payment Payment
		Time    Date
	}

	if err = c.post("payment", payload, &paymentData, nil); err != nil {
		return nil, err
	}

	return &paymentData.Payment, err
}

func (c *Client) CreateOrUpdatePayment(payment *Payment) (*Payment, error) {
	if payment.PaymentRefNum == "" {
		return nil, errors.New("missing payment PaymentRefNum")
	}

	existingPayments, err := c.QueryPayments("SELECT * FROM Payment WHERE PaymentRefNum = '" + strings.Replace(payment.PaymentRefNum, "'", "''", -1) + "'")
	if err != nil {
		return nil, err
	}

	if existingPayments != nil && len(existingPayments) > 0 {
		payment.Id = existingPayments[0].Id
		payment.SyncToken = existingPayments[0].SyncToken
	}

	payload := struct {
		*Payment
		Sparse bool `json:"sparse"`
	}{
		Payment: payment,
		Sparse:  true,
	}

	var paymentData struct {
		Payment Payment
		Time    Date
	}

	if err = c.post("payment", payload, &paymentData, nil); err != nil {
		return nil, err
	}

	return &paymentData.Payment, err
}

// VoidPayment voids the given payment in QuickBooks.
func (c *Client) VoidPayment(payment Payment) error {
	if payment.Id == "" {
		return errors.New("missing payment id")
	}

	existingPayment, err := c.FindPaymentById(payment.Id)
	if err != nil {
		return err
	}

	payment.SyncToken = existingPayment.SyncToken

	return c.post("payment", payment, nil, map[string]string{"operation": "update", "include": "void"})
}
