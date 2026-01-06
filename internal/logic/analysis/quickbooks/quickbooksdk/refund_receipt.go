package quickbooksdk

import (
	"errors"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type RefundReceipt struct {
	ID                   string          `json:"Id,omitempty"`
	SyncToken            string          `json:"SyncToken,omitempty"`
	GlobalTaxCalculation string          `json:"GlobalTaxCalculation,omitempty"`
	CustomerRef          *ReferenceType  `json:"CustomerRef,omitempty"`
	DepositToAccountRef  *ReferenceType  `json:"DepositToAccountRef,omitempty"` // Which account the money is returned to (e.g., Bank, Crypto)
	PaymentMethodRef     *ReferenceType  `json:"PaymentMethodRef,omitempty"`    // Optional: refund method (e.g., bank transfer)
	TxnDate              string          `json:"TxnDate,omitempty"`             // Refund date
	TotalAmt             decimal.Decimal `json:"TotalAmt,omitempty"`            // Total refund amount
	PrivateNote          string          `json:"PrivateNote,omitempty"`         // Note
	Line                 []Line          `json:"Line,omitempty"`                // Refunded product/service details
	DocNumber            string          `json:"DocNumber,omitempty"`           // Optional: refund document number
	PaymentRefNum        string          `json:"PaymentRefNum,omitempty"`       // Payment reference number (custom number)
}

// CreateRefundReceipt creates the given RefundReceipt on the QuickBooks server, returning
// the resulting RefundReceipt object.
func (c *Client) CreateRefundReceipt(refundReceipt *RefundReceipt) (*RefundReceipt, error) {
	var resp struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err := c.post("refundreceipt", refundReceipt, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.RefundReceipt, nil
}

// DeleteRefundReceipt deletes the refundReceipt
//
// If the refundReceipt was already deleted, QuickBooks returns 400 :(
// The response looks like this:
// {"Fault":{"Error":[{"Message":"Object Not Found","Detail":"Object Not Found : Something you're trying to use has been made inactive. Check the fields with accounts, refundReceipts, items, vendors or employees.","code":"610","element":""}],"type":"ValidationFault"},"time":"2018-03-20T20:15:59.571-07:00"}
//
// This is slightly horrifying and not documented in their API. When this
// happens we just return success; the goal of deleting it has been
// accomplished, just not by us.
func (c *Client) DeleteRefundReceipt(refundReceipt *RefundReceipt) error {
	if refundReceipt.ID == "" || refundReceipt.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("refundreceipt", refundReceipt, nil, map[string]string{"operation": "delete"})
}

// FindRefundReceipts gets the full list of RefundReceipts in the QuickBooks account.
func (c *Client) FindRefundReceipts() ([]RefundReceipt, error) {
	var resp struct {
		QueryResponse struct {
			RefundReceipts []RefundReceipt `json:"RefundReceipt"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM RefundReceipt", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return make([]RefundReceipt, 0), nil
	}

	refundReceipts := make([]RefundReceipt, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM RefundReceipt ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.RefundReceipts == nil {
			return nil, errors.New("no refundReceipts could be found")
		}

		refundReceipts = append(refundReceipts, resp.QueryResponse.RefundReceipts...)
	}

	return refundReceipts, nil
}

// FindRefundReceiptById finds the refundReceipt by the given id
func (c *Client) FindRefundReceiptById(id string) (*RefundReceipt, error) {
	var resp struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err := c.get("refundreceipt/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.RefundReceipt, nil
}

// QueryRefundReceipts accepts an SQL query and returns all refundReceipts found using it
func (c *Client) QueryRefundReceipts(query string) ([]RefundReceipt, error) {
	var resp struct {
		QueryResponse struct {
			RefundReceipts []RefundReceipt `json:"RefundReceipt"`
			StartPosition  int
			MaxResults     int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.RefundReceipts == nil {
		return make([]RefundReceipt, 0), nil
	}

	return resp.QueryResponse.RefundReceipts, nil
}

// SendRefundReceipt sends the refundReceipt to the RefundReceipt.BillEmail if emailAddress is left empty
func (c *Client) SendRefundReceipt(refundReceiptId string, emailAddress string) error {
	queryParameters := make(map[string]string)

	if emailAddress != "" {
		queryParameters["sendTo"] = emailAddress
	}

	return c.post("refundreceipt/"+refundReceiptId+"/send", nil, nil, queryParameters)
}

// UpdateRefundReceipt updates the refundReceipt
func (c *Client) UpdateRefundReceipt(refundReceipt *RefundReceipt) (*RefundReceipt, error) {
	if refundReceipt.ID == "" {
		return nil, errors.New("missing refundReceipt id")
	}

	existingRefundReceipt, err := c.FindRefundReceiptById(refundReceipt.ID)
	if err != nil {
		return nil, err
	}

	refundReceipt.SyncToken = existingRefundReceipt.SyncToken

	payload := struct {
		*RefundReceipt
		Sparse bool `json:"sparse"`
	}{
		RefundReceipt: refundReceipt,
		Sparse:        true,
	}

	var refundReceiptData struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err = c.post("refundreceipt", payload, &refundReceiptData, nil); err != nil {
		return nil, err
	}

	return &refundReceiptData.RefundReceipt, err
}

func (c *Client) CreateOrUpdateRefundReceiptByDocNumber(refundReceipt *RefundReceipt) (*RefundReceipt, error) {
	if refundReceipt.DocNumber == "" {
		return nil, errors.New("missing refundReceipt DocNumber")
	}

	existingRefundReceipts, err := c.QueryRefundReceipts("SELECT * FROM RefundReceipt WHERE DocNumber = '" + strings.Replace(refundReceipt.DocNumber, "'", "''", -1) + "'")
	if err != nil {
		return nil, err
	}

	if existingRefundReceipts != nil && len(existingRefundReceipts) > 0 {
		refundReceipt.ID = existingRefundReceipts[0].ID
		refundReceipt.SyncToken = existingRefundReceipts[0].SyncToken
	}

	payload := struct {
		*RefundReceipt
		Sparse bool `json:"sparse"`
	}{
		RefundReceipt: refundReceipt,
		Sparse:        true,
	}

	var refundReceiptData struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err = c.post("refundreceipt", payload, &refundReceiptData, nil); err != nil {
		return nil, err
	}

	return &refundReceiptData.RefundReceipt, err
}
