package generator

// Contact a company information
type Contact struct {
	//Name    string   `json:"name,omitempty" validate:"required,min=1,max=256"`
	Name string `json:"name,omitempty"`
	//Logo    []byte   `json:"logo,omitempty"` // Logo byte array
	Address *Address `json:"address,omitempty"`

	// AdditionalInfo to append after contact information. You can use basic html here (bold, italic tags).
	AdditionalInfo []string `json:"additional_info,omitempty"`
}

// appendContactTODoc append the contact to the document
func (c *Contact) appendContactTODoc(
	x float64,
	y float64,
	fill bool,
	doc *Document,
	title string,
) float64 {
	doc.pdf.SetXY(x, y)

	// Name
	if fill {
		doc.pdf.SetFillColor(
			doc.Options.WhiteBgColor[0],
			doc.Options.WhiteBgColor[1],
			doc.Options.WhiteBgColor[2],
		)
	} else {
		doc.pdf.SetFillColor(255, 255, 255)
	}

	// Reset x
	doc.pdf.SetX(x)

	// Name rect
	doc.pdf.Rect(x, doc.pdf.GetY(), 70, 8, "F")

	// Set name
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 10)
	doc.pdf.Cell(40, 8, doc.encodeString(title))

	doc.pdf.SetFont(doc.Options.Font, "", 10)
	doc.pdf.SetXY(x, doc.pdf.GetY()+8)
	doc.pdf.Cell(40, 8, doc.encodeString(c.Name))

	if c.Address != nil {
		// Address rect
		var addrRectHeight float64 = 17

		if len(c.Address.Address2) > 0 {
			addrRectHeight = addrRectHeight + 5
		}

		if len(c.Address.Country) == 0 {
			addrRectHeight = addrRectHeight - 5
		}

		doc.pdf.Rect(x, doc.pdf.GetY()+9, 90, addrRectHeight, "F")

		// Set address
		doc.pdf.SetFont(doc.Options.Font, "", 10)
		doc.pdf.SetXY(x, doc.pdf.GetY()+7)
		doc.pdf.MultiCell(90, 5, doc.encodeString(c.Address.ToString()), "0", "L", false)
	}

	// Additional info
	if c.AdditionalInfo != nil {
		doc.pdf.SetXY(x, doc.pdf.GetY())
		//doc.pdf.SetFontSize(SmallTextFontSize)
		doc.pdf.SetXY(x, doc.pdf.GetY()+2)

		for _, line := range c.AdditionalInfo {
			doc.pdf.SetXY(x, doc.pdf.GetY())
			doc.pdf.MultiCell(70, 3, doc.encodeString(line), "0", "L", false)
		}

		doc.pdf.SetXY(x, doc.pdf.GetY())
	}

	return doc.pdf.GetY()
}

func (c *Contact) appendMerchantContactToDoc(doc *Document, y float64) float64 {
	x, _, _, _ := doc.pdf.GetMargins()
	if doc.IsRefund {
		return c.appendContactTODoc(x, y, true, doc, "From:")
	} else {
		return c.appendContactTODoc(x, y, true, doc, "Issued by:")
	}

}

func (c *Contact) appendCustomerContactToDoc(doc *Document, y float64) float64 {
	if doc.IsRefund {
		return c.appendContactTODoc(110, y, true, doc, "To:")
	} else {
		return c.appendContactTODoc(110, y, true, doc, "Invoice to:")
	}
}
