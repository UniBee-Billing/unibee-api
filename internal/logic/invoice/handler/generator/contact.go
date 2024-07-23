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

	//// Logo
	//if c.Logo != nil {
	//	// Create filename
	//	fileName := b64.StdEncoding.EncodeToString([]byte(c.Name))
	//
	//	// Create reader from logo bytes
	//	ioReader := bytes.NewReader(c.Logo)
	//
	//	// Get image format
	//	_, format, _ := image.DecodeConfig(bytes.NewReader(c.Logo))
	//
	//	// Register image in pdf
	//	imageInfo := doc.pdf.RegisterImageOptionsReader(fileName, fpdf.ImageOptions{
	//		ImageType: format,
	//	}, ioReader)
	//
	//	if imageInfo != nil {
	//		var imageOpt fpdf.ImageOptions
	//		imageOpt.ImageType = format
	//		doc.pdf.ImageOptions(fileName, doc.pdf.GetX(), y, 0, 30, false, imageOpt, 0, "")
	//		doc.pdf.SetY(y + 30)
	//	}
	//}

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

		doc.pdf.Rect(x, doc.pdf.GetY()+9, 70, addrRectHeight, "F")

		// Set address
		doc.pdf.SetFont(doc.Options.Font, "", 10)
		doc.pdf.SetXY(x, doc.pdf.GetY()+8)
		doc.pdf.MultiCell(70, 5, doc.encodeString(c.Address.ToString()), "0", "L", false)
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
		doc.pdf.SetFontSize(BaseTextFontSize)
	}

	return doc.pdf.GetY()
}

func (c *Contact) appendCompanyContactToDoc(doc *Document, y float64) float64 {
	x, _, _, _ := doc.pdf.GetMargins()
	if doc.IsRefund {
		return c.appendContactTODoc(x, y, true, doc, "From:")
	} else {
		return c.appendContactTODoc(x, y, true, doc, "Issued By:")
	}

}

func (c *Contact) appendCustomerContactToDoc(doc *Document, y float64) float64 {
	if doc.IsRefund {
		return c.appendContactTODoc(130, y, true, doc, "To:")
	} else {
		return c.appendContactTODoc(130, y, true, doc, "Invoice To:")
	}
}
