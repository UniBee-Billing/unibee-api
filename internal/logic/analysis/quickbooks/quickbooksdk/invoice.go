// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooksdk

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

// Invoice represents a QuickBooks Invoice object.
type Invoice struct {
	ID                           string           `json:"Id,omitempty"`
	SyncToken                    string           `json:"SyncToken,omitempty"`
	MetaData                     *MetaData        `json:"MetaData,omitempty"`
	CustomField                  []CustomField    `json:"CustomField,omitempty"`
	DocNumber                    string           `json:"DocNumber,omitempty"`
	TxnDate                      string           `json:"TxnDate,omitempty"`
	DepartmentRef                *ReferenceType   `json:"DepartmentRef,omitempty"`
	PrivateNote                  string           `json:"PrivateNote,omitempty"`
	LinkedTxn                    []LinkedTxn      `json:"LinkedTxn,omitempty"`
	Line                         []Line           `json:"Line,omitempty"`
	TxnTaxDetail                 *TxnTaxDetail    `json:"TxnTaxDetail,omitempty"`
	CustomerRef                  *ReferenceType   `json:"CustomerRef,omitempty"`
	CustomerMemo                 *MemoRef         `json:"CustomerMemo,omitempty"`
	BillAddr                     *PhysicalAddress `json:"BillAddr,omitempty"`
	ShipAddr                     *PhysicalAddress `json:"ShipAddr,omitempty"`
	ClassRef                     *ReferenceType   `json:"ClassRef,omitempty"`
	SalesTermRef                 *ReferenceType   `json:"SalesTermRef,omitempty"`
	DueDate                      Date             `json:"DueDate,omitempty"`
	GlobalTaxCalculation         string           `json:"GlobalTaxCalculation,omitempty"`
	ShipMethodRef                *ReferenceType   `json:"ShipMethodRef,omitempty"`
	ShipDate                     Date             `json:"ShipDate,omitempty"`
	TrackingNum                  string           `json:"TrackingNum,omitempty"`
	TotalAmt                     decimal.Decimal  `json:"TotalAmt,omitempty"`
	CurrencyRef                  *ReferenceType   `json:"CurrencyRef,omitempty"`
	ExchangeRate                 json.Number      `json:"ExchangeRate,omitempty"`
	HomeAmtTotal                 json.Number      `json:"HomeAmtTotal,omitempty"`
	HomeBalance                  json.Number      `json:"HomeBalance,omitempty"`
	ApplyTaxAfterDiscount        bool             `json:"ApplyTaxAfterDiscount,omitempty"`
	PrintStatus                  string           `json:"PrintStatus,omitempty"`
	EmailStatus                  string           `json:"EmailStatus,omitempty"`
	BillEmail                    *EmailAddress    `json:"BillEmail,omitempty"`
	BillEmailCC                  *EmailAddress    `json:"BillEmailCc,omitempty"`
	BillEmailBCC                 *EmailAddress    `json:"BillEmailBcc,omitempty"`
	DeliveryInfo                 *DeliveryInfo    `json:"DeliveryInfo,omitempty"`
	Balance                      json.Number      `json:"Balance,omitempty"`
	TxnSource                    string           `json:"TxnSource,omitempty"`
	AllowOnlineCreditCardPayment bool             `json:"AllowOnlineCreditCardPayment,omitempty"`
	AllowOnlineACHPayment        bool             `json:"AllowOnlineACHPayment,omitempty"`
	Deposit                      json.Number      `json:"Deposit,omitempty"`
	DepositToAccountRef          *ReferenceType   `json:"DepositToAccountRef,omitempty"`
}

type DeliveryInfo struct {
	DeliveryType string
	DeliveryTime Date
}

type LinkedTxn struct {
	TxnID   string `json:"TxnId"`
	TxnType string `json:"TxnType"`
}

type TxnTaxDetail struct {
	TxnTaxCodeRef ReferenceType `json:",omitempty"`
	TotalTax      json.Number   `json:",omitempty"`
	TaxLine       []Line        `json:",omitempty"`
}

type AccountBasedExpenseLineDetail struct {
	AccountRef ReferenceType
	TaxAmount  json.Number `json:",omitempty"`
	// TaxInclusiveAmt json.Number              `json:",omitempty"`
	// ClassRef        ReferenceType `json:",omitempty"`
	// TaxCodeRef      ReferenceType `json:",omitempty"`
	// MarkupInfo MarkupInfo `json:",omitempty"`
	// BillableStatus BillableStatusEnum       `json:",omitempty"`
	// CustomerRef    ReferenceType `json:",omitempty"`
}

type Line struct {
	Id                            string                         `json:"Id,omitempty"`
	LineNum                       int                            `json:"LineNum,omitempty"`
	Description                   string                         `json:"Description,omitempty"`
	Amount                        decimal.Decimal                `json:"Amount,omitempty"`
	DetailType                    string                         `json:"DetailType,omitempty"`
	AccountBasedExpenseLineDetail *AccountBasedExpenseLineDetail `json:"AccountBasedExpenseLineDetail,omitempty"`
	SalesItemLineDetail           *SalesItemLineDetail           `json:"SalesItemLineDetail,omitempty"`
	DiscountLineDetail            *DiscountLineDetail            `json:"DiscountLineDetail,omitempty"`
	TaxLineDetail                 *TaxLineDetail                 `json:"TaxLineDetail,omitempty"`
}

// TaxLineDetail ...
type TaxLineDetail struct {
	PercentBased     bool           `json:",omitempty"`
	NetAmountTaxable json.Number    `json:",omitempty"`
	TaxPercent       json.Number    `json:",omitempty"`
	TaxRateRef       *ReferenceType `json:",omitempty"`
}

// SalesItemLineDetail ...
type SalesItemLineDetail struct {
	ItemRef         *ReferenceType  `json:",omitempty"`
	ClassRef        *ReferenceType  `json:",omitempty"`
	UnitPrice       decimal.Decimal `json:",omitempty"`
	Qty             float64         `json:",omitempty"`
	ItemAccountRef  *ReferenceType  `json:",omitempty"`
	TaxCodeRef      *ReferenceType  `json:",omitempty"`
	ServiceDate     Date            `json:",omitempty"`
	TaxInclusiveAmt json.Number     `json:",omitempty"`
	DiscountRate    json.Number     `json:",omitempty"`
	DiscountAmt     json.Number     `json:",omitempty"`
}

// DiscountLineDetail ...
type DiscountLineDetail struct {
	PercentBased    bool
	DiscountPercent float32 `json:",omitempty"`
}

// CreateInvoice creates the given Invoice on the QuickBooks server, returning
// the resulting Invoice object.
func (c *Client) CreateInvoice(invoice *Invoice) (*Invoice, error) {
	var resp struct {
		Invoice Invoice
		Time    Date
	}

	if err := c.post("invoice", invoice, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Invoice, nil
}

// DeleteInvoice deletes the invoice
//
// If the invoice was already deleted, QuickBooks returns 400 :(
// The response looks like this:
// {"Fault":{"Error":[{"Message":"Object Not Found","Detail":"Object Not Found : Something you're trying to use has been made inactive. Check the fields with accounts, invoices, items, vendors or employees.","code":"610","element":""}],"type":"ValidationFault"},"time":"2018-03-20T20:15:59.571-07:00"}
//
// This is slightly horrifying and not documented in their API. When this
// happens we just return success; the goal of deleting it has been
// accomplished, just not by us.
func (c *Client) DeleteInvoice(invoice *Invoice) error {
	if invoice.ID == "" || invoice.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("invoice", invoice, nil, map[string]string{"operation": "delete"})
}

// FindInvoices gets the full list of Invoices in the QuickBooks account.
func (c *Client) FindInvoices() ([]Invoice, error) {
	var resp struct {
		QueryResponse struct {
			Invoices      []Invoice `json:"Invoice"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Invoice", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no invoices could be found")
	}

	invoices := make([]Invoice, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Invoice ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Invoices == nil {
			return nil, errors.New("no invoices could be found")
		}

		invoices = append(invoices, resp.QueryResponse.Invoices...)
	}

	return invoices, nil
}

// FindInvoiceById finds the invoice by the given id
func (c *Client) FindInvoiceById(id string) (*Invoice, error) {
	var resp struct {
		Invoice Invoice
		Time    Date
	}

	if err := c.get("invoice/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Invoice, nil
}

// QueryInvoices accepts an SQL query and returns all invoices found using it
func (c *Client) QueryInvoices(query string) ([]Invoice, error) {
	var resp struct {
		QueryResponse struct {
			Invoices      []Invoice `json:"Invoice"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Invoices == nil {
		return make([]Invoice, 0), nil
	}

	return resp.QueryResponse.Invoices, nil
}

// SendInvoice sends the invoice to the Invoice.BillEmail if emailAddress is left empty
func (c *Client) SendInvoice(invoiceId string, emailAddress string) error {
	queryParameters := make(map[string]string)

	if emailAddress != "" {
		queryParameters["sendTo"] = emailAddress
	}

	return c.post("invoice/"+invoiceId+"/send", nil, nil, queryParameters)
}

// UpdateInvoice updates the invoice
func (c *Client) UpdateInvoice(invoice *Invoice) (*Invoice, error) {
	if invoice.ID == "" {
		return nil, errors.New("missing invoice id")
	}

	existingInvoice, err := c.FindInvoiceById(invoice.ID)
	if err != nil {
		return nil, err
	}

	invoice.SyncToken = existingInvoice.SyncToken

	payload := struct {
		*Invoice
		Sparse bool `json:"sparse"`
	}{
		Invoice: invoice,
		Sparse:  true,
	}

	var invoiceData struct {
		Invoice Invoice
		Time    Date
	}

	if err = c.post("invoice", payload, &invoiceData, nil); err != nil {
		return nil, err
	}

	return &invoiceData.Invoice, err
}

func (c *Client) CreateOrUpdateInvoiceByDocNumber(invoice *Invoice) (*Invoice, error) {
	if invoice.DocNumber == "" {
		return nil, errors.New("missing invoice DocNumber")
	}

	existingInvoices, err := c.QueryInvoices("SELECT * FROM Invoice WHERE DocNumber = '" + strings.Replace(invoice.DocNumber, "'", "''", -1) + "'")
	if err != nil {
		return nil, err
	}

	if existingInvoices != nil && len(existingInvoices) > 0 {
		invoice.ID = existingInvoices[0].ID
		invoice.SyncToken = existingInvoices[0].SyncToken
	}

	payload := struct {
		*Invoice
		Sparse bool `json:"sparse"`
	}{
		Invoice: invoice,
		Sparse:  true,
	}

	var invoiceData struct {
		Invoice Invoice
		Time    Date
	}

	if err = c.post("invoice", payload, &invoiceData, nil); err != nil {
		return nil, err
	}

	return &invoiceData.Invoice, err
}

func (c *Client) VoidInvoice(invoice Invoice) error {
	if invoice.ID == "" {
		return errors.New("missing invoice id")
	}

	existingInvoice, err := c.FindInvoiceById(invoice.ID)
	if err != nil {
		return err
	}

	invoice.SyncToken = existingInvoice.SyncToken

	return c.post("invoice", invoice, nil, map[string]string{"operation": "void"})
}
