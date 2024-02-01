package generator

// SetType set type of document
func (d *Document) SetType(docType string) *Document {
	d.Type = docType
	return d
}

// SetHeader set header of document
func (d *Document) SetHeader(header *HeaderFooter) *Document {
	d.Header = header
	return d
}

// SetFooter set footer of document
func (d *Document) SetFooter(footer *HeaderFooter) *Document {
	d.Footer = footer
	return d
}

// SetInvoiceNumber of document
func (d *Document) SetInvoiceNumber(InvoiceNumber string) *Document {
	d.InvoiceNumber = InvoiceNumber
	return d
}

func (d *Document) SetLogo(logo []byte) *Document {
	d.Logo = logo
	return d
}

// SetInvoiceDate of document
func (d *Document) SetInvoiceDate(invoiceDate string) *Document {
	d.InvoiceDate = invoiceDate
	return d
}

// SetStatus of document
func (d *Document) SetStatus(status string) *Document {
	d.Status = status
	return d
}

// SetDescription of document
func (d *Document) SetDescription(desc string) *Document {
	d.Description = desc
	return d
}

// SetNotes of document
func (d *Document) SetNotes(notes string) *Document {
	d.Notes = notes
	return d
}

// SetCompany of document
func (d *Document) SetCompany(company *Contact) *Document {
	d.Company = company
	return d
}

// SetCustomer of document
func (d *Document) SetCustomer(customer *Contact) *Document {
	d.Customer = customer
	return d
}

// AppendItem to document items
func (d *Document) AppendItem(item *Item) *Document {
	d.Items = append(d.Items, item)
	return d
}

// SetPaidDate of document
func (d *Document) SetPaidDate(date string) *Document {
	d.PaidDate = date
	return d
}

// SetPaymentTerm of document
func (d *Document) SetPaymentTerm(term string) *Document {
	d.PaymentTerm = term
	return d
}

// SetDefaultTax of document
func (d *Document) SetDefaultTax(tax *Tax) *Document {
	d.DefaultTax = tax
	return d
}

// SetDiscount of document
func (d *Document) SetDiscount(discount *Discount) *Document {
	d.Discount = discount
	return d
}
