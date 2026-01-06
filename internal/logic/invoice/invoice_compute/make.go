package invoice_compute

import "unibee/api/bean"

type InvoiceSimplifyInternalReq struct {
	Id             uint64                            `json:"id"`
	MerchantId     uint64                            `json:"merchantId"`
	InvoiceId      string                            `json:"invoiceId"`
	InvoiceName    string                            `json:"invoiceName"`
	DiscountCode   string                            `json:"discountCode"`
	Currency       string                            `json:"currency"`
	TaxPercentage  int64                             `json:"taxPercentage"`
	Lines          []*InvoiceItemSimplifyInternalReq `json:"lines"`
	PeriodEnd      int64                             `json:"periodEnd"`
	PeriodStart    int64                             `json:"periodStart"`
	ProrationDate  int64                             `json:"prorationDate"`
	ProrationScale int64                             `json:"prorationScale"`
	FinishTime     int64                             `json:"finishTime"`
	SendStatus     int                               `json:"sendStatus"`
	DayUtilDue     int64                             `json:"dayUtilDue"`
	TimeNow        int64                             `json:"timeNow"`
}

type InvoiceItemSimplifyInternalReq struct {
	UnitAmountExcludingTax int64      `json:"unitAmountExcludingTax"`
	Quantity               int64      `json:"quantity"`
	Name                   string     `json:"name"`
	Description            string     `json:"description"`
	Plan                   *bean.Plan `json:"plan"`
}
