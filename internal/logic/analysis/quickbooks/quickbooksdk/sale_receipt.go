package quickbooksdk

import (
	"errors"
	"github.com/shopspring/decimal"
	"strings"
)

// https://developer.intuit.com/app/developer/qbo/docs/api/accounting/all-entities/salesreceipt#the-salesreceipt-object
type SalesReceipt struct {
	ID                   string            `json:"Id,omitempty"`
	SyncToken            string            `json:"SyncToken,omitempty"`
	MetaData             *MetaData         `json:"MetaData,omitempty"`
	DocNumber            string            `json:"DocNumber,omitempty"`
	TxnDate              string            `json:"TxnDate,omitempty"` // e.g., "2023-01-01"
	Line                 []SaleReceiptLine `json:"Line,omitempty"`
	CustomerRef          *ReferenceType    `json:"CustomerRef,omitempty"`
	BillEmail            *EmailAddress     `json:"BillEmail,omitempty"`
	TotalAmt             decimal.Decimal   `json:"TotalAmt,omitempty"`
	PrivateNote          string            `json:"PrivateNote,omitempty"`
	DepositToAccountRef  *ReferenceType    `json:"DepositToAccountRef,omitempty"`
	GlobalTaxCalculation string            `json:"GlobalTaxCalculation,omitempty"`
}

type SaleReceiptLine struct {
	ID                  string                      `json:"Id,omitempty"`
	LineNum             int                         `json:"LineNum,omitempty"`
	Description         string                      `json:"Description,omitempty"`
	Amount              decimal.Decimal             `json:"Amount"`
	DetailType          string                      `json:"DetailType"`
	SalesItemLineDetail *SalesReceiptItemLineDetail `json:"SalesItemLineDetail,omitempty"`
}

type SalesReceiptItemLineDetail struct {
	ItemRef    *ReferenceType  `json:"ItemRef,omitempty"`
	UnitPrice  decimal.Decimal `json:"UnitPrice,omitempty"`
	Qty        float64         `json:"Qty,omitempty"`
	TaxCodeRef *ReferenceType  `json:"TaxCodeRef,omitempty"`
}

func (c *Client) CreateSalesReceipt(receipt *SalesReceipt) (*SalesReceipt, error) {
	var resp struct {
		SalesReceipt SalesReceipt
		Time         Date
	}
	if err := c.post("salesreceipt", receipt, &resp, nil); err != nil {
		return nil, err
	}
	return &resp.SalesReceipt, nil
}

func (c *Client) FindSaleReceiptById(id string) (*SalesReceipt, error) {
	var resp struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err := c.get("salesreceipt/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.SalesReceipt, nil
}

func (c *Client) QuerySaleReceipts(query string) ([]SalesReceipt, error) {
	var resp struct {
		QueryResponse struct {
			SalesReceipts []SalesReceipt `json:"SalesReceipt"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.SalesReceipts == nil {
		return make([]SalesReceipt, 0), nil
	}

	return resp.QueryResponse.SalesReceipts, nil
}

// UpdateInvoice updates the invoice
func (c *Client) UpdateSaleReceipts(saleReceipt *SalesReceipt) (*SalesReceipt, error) {
	if saleReceipt.ID == "" {
		return nil, errors.New("missing salesReceipt id")
	}

	existingSaleReceipt, err := c.FindSaleReceiptById(saleReceipt.ID)
	if err != nil {
		return nil, err
	}

	saleReceipt.SyncToken = existingSaleReceipt.SyncToken

	payload := struct {
		*SalesReceipt
		Sparse bool `json:"sparse"`
	}{
		SalesReceipt: saleReceipt,
		Sparse:       true,
	}

	var salesReceiptData struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err = c.post("salesreceipt", payload, &salesReceiptData, nil); err != nil {
		return nil, err
	}

	return &salesReceiptData.SalesReceipt, err
}

// CreateOrUpdateSaleReceiptsByDocNumber create or updates the invoice
func (c *Client) CreateOrUpdateSaleReceiptsByDocNumber(saleReceipt *SalesReceipt) (*SalesReceipt, error) {
	if saleReceipt.DocNumber == "" {
		return nil, errors.New("missing salesReceipt DocNumber")
	}
	existingSaleReceipts, err := c.QuerySaleReceipts("SELECT * FROM SalesReceipt WHERE DocNumber = '" + strings.Replace(saleReceipt.DocNumber, "'", "''", -1) + "'")
	if err != nil {
		return nil, err
	}

	if existingSaleReceipts != nil && len(existingSaleReceipts) > 0 {
		saleReceipt.ID = existingSaleReceipts[0].ID
		saleReceipt.SyncToken = existingSaleReceipts[0].SyncToken
	}

	payload := struct {
		*SalesReceipt
		Sparse bool `json:"sparse"`
	}{
		SalesReceipt: saleReceipt,
		Sparse:       true,
	}

	var salesReceiptData struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err = c.post("salesreceipt", payload, &salesReceiptData, nil); err != nil {
		return nil, err
	}

	return &salesReceiptData.SalesReceipt, err
}

func (c *Client) VoidSalesReceipt(salesReceipt SalesReceipt) error {
	if salesReceipt.ID == "" {
		return errors.New("missing salesReceipt id")
	}

	existingSaleReceipt, err := c.FindSaleReceiptById(salesReceipt.ID)
	if err != nil {
		return err
	}

	salesReceipt.SyncToken = existingSaleReceipt.SyncToken

	return c.post("salesreceipt", salesReceipt, nil, map[string]string{"operation": "void"})
}
