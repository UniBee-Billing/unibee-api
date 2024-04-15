package generator

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"github.com/go-pdf/fpdf"
	"image"
)

// Build pdf document from data provided
func (doc *Document) Build() (*fpdf.Fpdf, error) {
	// Validate document data
	if err := doc.Validate(); err != nil {
		return nil, err
	}

	// Build base doc
	doc.pdf.SetMargins(BaseMargin, BaseMarginTop, BaseMargin)
	doc.pdf.SetXY(10, 10)
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)

	// Set header
	if doc.Header != nil {
		if err := doc.Header.applyHeader(doc); err != nil {
			return nil, err
		}
	}

	// Set footer
	if doc.Footer != nil {
		if err := doc.Footer.applyFooter(doc); err != nil {
			return nil, err
		}
	}

	// Add first page
	doc.pdf.AddPage()

	// Load font
	doc.pdf.SetFont(doc.Options.Font, "", 12)

	//// Append document title
	//doc.appendTitle()

	doc.appendLogo()

	// Append document metas (ref & version)
	doc.appendMetas()

	y := doc.pdf.GetY()
	// Append company contact to doc
	companyBottom := doc.Company.appendCompanyContactToDoc(doc, y)

	// Append customer contact to doc
	customerBottom := doc.Customer.appendCustomerContactToDoc(doc, y)

	if customerBottom > companyBottom {
		doc.pdf.SetXY(10, customerBottom)
	} else {
		doc.pdf.SetXY(10, companyBottom)
	}

	// Append description
	doc.appendDescription()

	// Append items
	doc.appendItems()

	// Check page height (total bloc height = 30, 45 when doc discount)
	offset := doc.pdf.GetY() + 30
	if doc.Discount != nil {
		offset += 15
	}
	if offset > MaxPageHeight {
		doc.pdf.AddPage()
	}

	// Append notes
	doc.appendNotes()

	// Append total
	doc.appendTotal()

	// Append payment term
	doc.appendPaymentTerm()

	// Append js to autoprint if AutoPrint == true
	if doc.Options.AutoPrint {
		doc.pdf.SetJavascript("print(true);")
	}

	return doc.pdf, nil
}

// appendTitle to document
func (doc *Document) appendTitle() {
	title := doc.typeAsString()

	// Set x y
	doc.pdf.SetXY(120, BaseMarginTop)

	// Draw rect
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rect(120, BaseMarginTop, 80, 10, "F")

	// Draw text
	doc.pdf.SetFont(doc.Options.Font, "", 14)
	doc.pdf.CellFormat(80, 10, doc.encodeString(title), "0", 0, "C", false, 0, "")
}

func (doc *Document) appendLogo() float64 {
	if doc.Logo != nil {
		// Create filename
		fileName := b64.StdEncoding.EncodeToString([]byte("unibee"))

		// Create reader from logo bytes
		ioReader := bytes.NewReader(doc.Logo)

		// Get image format
		_, format, _ := image.DecodeConfig(bytes.NewReader(doc.Logo))

		// Register image in pdf
		imageInfo := doc.pdf.RegisterImageOptionsReader(fileName, fpdf.ImageOptions{
			ImageType: format,
		}, ioReader)

		if imageInfo != nil {
			var imageOpt fpdf.ImageOptions
			imageOpt.ImageType = format
			doc.pdf.ImageOptions(fileName, doc.pdf.GetX()+1, doc.pdf.GetY(), 0, 30, false, imageOpt, 0, "")
			doc.pdf.SetY(doc.pdf.GetY() + 30)
		}
		return 30.0
	}
	_, y, _, _ := doc.pdf.GetMargins()
	return y
}

// appendMetas to document
func (doc *Document) appendMetas() {
	// Append ref
	x, _, _, _ := doc.pdf.GetMargins()
	refString := fmt.Sprintf("%s: %s", doc.Options.TextInvoiceNumberTitle, doc.InvoiceNumber)
	//y := doc.pdf.GetY()
	startY := doc.pdf.GetY() + 8
	doc.pdf.SetXY(x, startY)
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 10)
	doc.pdf.CellFormat(80, 12, doc.encodeString(refString), "0", 0, "L", false, 0, "")

	// Append version
	if len(doc.InvoiceDate) > 0 {
		dataString := fmt.Sprintf("%s: %s", doc.Options.TextInvoiceDateTitle, doc.InvoiceDate)
		doc.pdf.SetXY(x, doc.pdf.GetY()+8)
		doc.pdf.SetFont(doc.Options.Font, "", 10)
		doc.pdf.CellFormat(80, 12, doc.encodeString(dataString), "0", 0, "L", false, 0, "")
	}

	// Append paidDate
	paidDate := "-"
	if len(doc.PaidDate) > 0 {
		paidDate = doc.PaidDate
	}
	dateString := fmt.Sprintf("%s: %s", doc.Options.TextInvoicePaidDateTitle, paidDate)
	doc.pdf.SetXY(x, doc.pdf.GetY()+8)
	doc.pdf.SetFont(doc.Options.Font, "", 10)
	doc.pdf.CellFormat(80, 12, doc.encodeString(dateString), "0", 0, "L", false, 0, "")

	returnY := doc.pdf.GetY() + 15

	doc.pdf.SetXY(130, startY)
	doc.pdf.SetFont(doc.Options.BoldFont, "", 10)
	doc.pdf.CellFormat(80, 12, doc.encodeString("Invoice Status:"), "0", 0, "L", false, 0, "")

	doc.pdf.SetXY(130, startY+10)
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 18)
	if doc.Status == "Paid" {
		doc.pdf.SetTextColor(
			doc.Options.PaidTextColor[0],
			doc.Options.PaidTextColor[1],
			doc.Options.PaidTextColor[2],
		)
	} else {
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
	}
	doc.pdf.CellFormat(80, 12, doc.encodeString(doc.Status), "0", 0, "L", false, 0, "")
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)
	doc.pdf.SetXY(x, returnY)
}

// appendDescription to document
func (doc *Document) appendDescription() {
	if len(doc.Description) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + 10)
		doc.pdf.SetFont(doc.Options.Font, "", 12)
		doc.pdf.MultiCell(190, 5, doc.encodeString(doc.Description), "B", "L", false)
	}
}

// drawsTableTitles in document
func (doc *Document) drawsTableTitles(fontSize float64) {
	// Draw table titles
	doc.pdf.SetX(10)
	doc.pdf.SetY(doc.pdf.GetY() + 5)
	doc.pdf.SetFont(doc.Options.BoldFont, "B", fontSize)

	// Draw rec
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rect(10, doc.pdf.GetY(), 190, 12, "F")

	// Name
	doc.pdf.SetX(ItemColNameOffset)
	doc.pdf.CellFormat(
		ItemColUnitPriceOffset-ItemColNameOffset,
		12,
		doc.encodeString(doc.Options.TextItemsNameTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Unit price
	if ItemColQuantityOffset-ItemColUnitPriceOffset > 0 {
		doc.pdf.SetX(ItemColUnitPriceOffset)
		doc.pdf.CellFormat(
			ItemColQuantityOffset-ItemColUnitPriceOffset,
			12,
			doc.encodeString(doc.Options.TextItemsUnitCostTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Quantity
	if ItemColTotalHTOffset-ItemColQuantityOffset > 0 {
		doc.pdf.SetX(ItemColQuantityOffset)
		doc.pdf.CellFormat(
			ItemColTotalHTOffset-ItemColQuantityOffset,
			12,
			doc.encodeString(doc.Options.TextItemsQuantityTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Total NoTax
	if ItemColDiscountOffset-ItemColTotalHTOffset > 0 {
		doc.pdf.SetX(ItemColTotalHTOffset)
		doc.pdf.CellFormat(
			ItemColDiscountOffset-ItemColTotalHTOffset,
			12,
			doc.encodeString(doc.Options.TextItemsTotalHTTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Discount
	if ItemColTaxOffset-ItemColDiscountOffset > 0 {
		doc.pdf.SetX(ItemColDiscountOffset)
		doc.pdf.CellFormat(
			ItemColTaxOffset-ItemColDiscountOffset,
			12,
			doc.encodeString(doc.Options.TextItemsDiscountTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Tax
	if ItemColTotalTTCOffset-ItemColDiscountOffset > 0 {
		doc.pdf.SetX(ItemColTaxOffset)
		doc.pdf.CellFormat(
			ItemColTotalTTCOffset-ItemColDiscountOffset,
			12,
			doc.encodeString(doc.Options.TextItemsTaxTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// TOTAL TTC
	doc.pdf.SetX(ItemColTotalTTCOffset)
	doc.pdf.CellFormat(
		190-ItemColTotalTTCOffset,
		12,
		doc.encodeString(doc.Options.TextItemsTotalTTCTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)
}

// appendItems to document
func (doc *Document) appendItems() {
	doc.drawsTableTitles(11.0)

	doc.pdf.SetX(10)
	doc.pdf.SetY(doc.pdf.GetY() + 18)
	doc.pdf.SetFont(doc.Options.Font, "", 10.0)

	for i := 0; i < len(doc.Items); i++ {
		item := doc.Items[i]

		// Check item tax
		if item.Tax == nil {
			item.Tax = doc.DefaultTax
		}

		// Append to pdf
		item.appendColTo(doc.Options, doc)

		if doc.pdf.GetY() > MaxPageHeight {
			// Add page
			doc.pdf.AddPage()
			doc.drawsTableTitles(12.0)
			doc.pdf.SetFont(doc.Options.Font, "", 12)
		}

		doc.pdf.SetX(10)
		doc.pdf.SetY(doc.pdf.GetY() + 8)
	}
}

// appendNotes to document
func (doc *Document) appendNotes() {
	if len(doc.Notes) == 0 {
		return
	}

	currentY := doc.pdf.GetY()

	doc.pdf.SetFont(doc.Options.Font, "", 12)
	doc.pdf.SetX(BaseMargin)
	doc.pdf.SetRightMargin(100)
	doc.pdf.SetY(currentY + 10)

	_, lineHt := doc.pdf.GetFontSize()
	html := doc.pdf.HTMLBasicNew()
	html.Write(lineHt, doc.encodeString(doc.Notes))

	doc.pdf.SetRightMargin(BaseMargin)
	doc.pdf.SetY(currentY)
}

var moneyX = 180.0

// appendTotal to document
func (doc *Document) appendTotal() {
	doc.pdf.SetY(doc.pdf.GetY() + 10)
	doc.pdf.SetFont(doc.Options.Font, "", LargeTextFontSize)
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)

	// Draw SUB TOTAL HT title
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextSubTotal), "0", 0, "R", false, 0, "")

	// Draw SUB TOTAL HT amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.SubTotalString),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Draw DISCOUNT TOTAL HT title
	if len(doc.DiscountTotalString) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + 10)
		doc.pdf.SetX(120)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
		if len(doc.DiscountTitle) > 0 {
			doc.pdf.CellFormat(38, 10, doc.encodeString(doc.DiscountTitle), "0", 0, "R", false, 0, "")
		} else {
			doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextTotalDiscounted), "0", 0, "R", false, 0, "")
		}

		// Draw DISCOUNT TOTAL HT amount
		doc.pdf.SetX(moneyX)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
		doc.pdf.CellFormat(
			40,
			10,
			doc.encodeString(doc.DiscountTotalString),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
	}

	// Draw tax title
	doc.pdf.SetY(doc.pdf.GetY() + 10)
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(38, 10, doc.encodeString(fmt.Sprintf("%s(%s)", doc.Options.TextTotalTax, doc.TaxPercentageString)), "0", 0, "R", false, 0, "")

	// Draw tax amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.TaxString),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Draw total with tax title
	doc.pdf.SetY(doc.pdf.GetY() + 10)
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextTotalWithTax), "0", 0, "R", false, 0, "")

	// Draw total with tax amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	doc.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.TotalString),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)
}

// appendPaymentTerm to document
func (doc *Document) appendPaymentTerm() {
	if len(doc.PaymentTerm) > 0 {
		paymentTermString := fmt.Sprintf(
			"%s",
			doc.PaymentTerm,
		)
		doc.pdf.SetY(doc.pdf.GetY() + 30)

		doc.pdf.SetX(0)
		doc.pdf.SetFont(doc.Options.BoldFont, "B", 10)
		//doc.pdf.SetFillColor(doc.Options.DeepBgColor[0], doc.Options.DeepBgColor[1], doc.Options.DeepBgColor[2])
		//doc.pdf.Rect(0, doc.pdf.GetY(), 210, 10, "F")
		doc.pdf.CellFormat(210, 10, doc.encodeString(paymentTermString), "0", 0, "C", false, 0, "")
	}

}
