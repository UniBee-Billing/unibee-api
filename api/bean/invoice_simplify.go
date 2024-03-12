package bean

type InvoiceSimplify struct {
	InvoiceId                      string                 `json:"invoiceId"`
	InvoiceName                    string                 `json:"invoiceName"`
	TotalAmount                    int64                  `json:"totalAmount"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"`
	Currency                       string                 `json:"currency"`
	TaxAmount                      int64                  `json:"taxAmount"`
	TaxScale                       int64                  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	SubscriptionAmount             int64                  `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemSimplify `json:"lines"`
	PeriodEnd                      int64                  `json:"periodEnd"`
	PeriodStart                    int64                  `json:"periodStart"`
	ProrationDate                  int64                  `json:"prorationDate"`
	ProrationScale                 int64                  `json:"prorationScale"`
}

type InvoiceItemSimplify struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount"`
	Tax                    int64  `json:"tax"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	TaxScale               int64  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Proration              bool   `json:"proration"`
	Quantity               int64  `json:"quantity"`
	PeriodEnd              int64  `json:"periodEnd"`
	PeriodStart            int64  `json:"periodStart"`
}
