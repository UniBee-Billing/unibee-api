package generator

// UnicodeTranslateFunc ...
type UnicodeTranslateFunc func(string) string

// Options for Document
type Options struct {
	AutoPrint bool `json:"auto_print,omitempty"`

	CurrencySymbol    string `default:"€ " json:"currency_symbol,omitempty"`
	CurrencyPrecision int    `default:"2" json:"currency_precision,omitempty"`
	CurrencyDecimal   string `default:"." json:"currency_decimal,omitempty"`
	CurrencyThousand  string `default:" " json:"currency_thousand,omitempty"`

	TextTypeInvoice      string `default:"INVOICE" json:"text_type_invoice,omitempty"`
	TextTypeQuotation    string `default:"QUOTATION" json:"text_type_quotation,omitempty"`
	TextTypeDeliveryNote string `default:"DELIVERY NOTE" json:"text_type_delivery_note,omitempty"`

	TextInvoiceIdTitle           string `default:"Invoice id" json:"invoice_id_title,omitempty"`
	TextInvoiceNumberTitle       string `default:"Invoice number" json:"invoice_number_title,omitempty"`
	TextOriginInvoiceNumberTitle string `default:"Original invoice number" json:"text_origin_number_title,omitempty"`
	TextInvoiceDateTitle         string `default:"Invoice date" json:"text_data_title,omitempty"`
	TextInvoicePaidDateTitle     string `default:"Invoice payment date" json:"text_paid_date_title,omitempty"`
	TextPaymentTermTitle         string `default:"Payment term" json:"text_payment_term_title,omitempty"`

	TextItemsNameTitle     string `default:"Description" json:"text_items_name_title,omitempty"`
	TextItemsUnitCostTitle string `default:"Unit price" json:"text_items_unit_cost_title,omitempty"`
	TextItemsQuantityTitle string `default:"Quantity" json:"text_items_quantity_title,omitempty"`
	TextItemsTotalHTTitle  string `default:"Total no tax" json:"text_items_total_ht_title,omitempty"`
	TextItemsTaxTitle      string `default:"VAT" json:"text_items_tax_title,omitempty"`
	TextItemsDiscountTitle string `default:"Discount" json:"text_items_discount_title,omitempty"`
	TextItemsTotalTTCTitle string `default:"Total" json:"text_items_total_ttc_title,omitempty"`

	TextSubTotal        string `default:"SUBTOTAL" json:"text_sub_total,omitempty"`
	TextTotalDiscounted string `default:"TOTAL DISCOUNTED" json:"text_total_discounted,omitempty"`
	TextTotalTax        string `default:"VAT" json:"text_total_tax,omitempty"`
	TextTotalWithTax    string `default:"TOTAL" json:"text_total_with_tax,omitempty"`

	BaseTextColor []int `default:"[51,51,51]" json:"base_text_color,omitempty"`
	PaidTextColor []int `default:"[73,167,101]" json:"paid_text_color,omitempty"`
	GreyTextColor []int `default:"[153,153,153]" json:"grey_text_color,omitempty"`
	WhiteBgColor  []int `default:"[255,255,255]" json:"grey_bg_color,omitempty"`
	//WhiteBgColor  []int `default:"[255,255,255]" json:"white_bg_color,omitempty"`
	//DarkBgColor []int `default:"[247,247,247]" json:"dark_bg_color,omitempty"`
	DarkBgColor []int `default:"[242,242,242]" json:"dark_bg_color,omitempty"`
	DeepBgColor []int `default:"[200,200,200]" json:"deep_bg_color,omitempty"`

	Font     string `default:"Helvetica"`
	BoldFont string `default:"Helvetica"`

	UnicodeTranslateFunc UnicodeTranslateFunc
}
