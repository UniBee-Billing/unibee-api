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
	doc.pdf.SetXY(0, 0)
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

	doc.appendLogo()

	x := doc.pdf.GetX()
	y := doc.pdf.GetY()

	doc.appendTitle()
	doc.appendInvoiceHeader()

	y = doc.pdf.GetY()
	if y < 44 {
		y = 44
	}

	doc.pdf.SetXY(x, y)

	y = doc.pdf.GetY() + BaseMargin
	merchantBottom := doc.Company.appendMerchantContactToDoc(doc, y)
	y = doc.appendInvoiceStatus(y)
	customerBottom := doc.Customer.appendCustomerContactToDoc(doc, y)

	if customerBottom > merchantBottom {
		doc.pdf.SetXY(10, customerBottom+BaseMargin)
	} else {
		doc.pdf.SetXY(10, merchantBottom+BaseMargin)
	}

	// Append description
	doc.appendDescription()

	// Append items
	doc.appendItems()

	//// Append Exchange Rate
	//doc.appendExchangeRate()

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

	// Append js to auto print if AutoPrint == true
	if doc.Options.AutoPrint {
		doc.pdf.SetJavascript("print(true);")
	}

	return doc.pdf, nil
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

// appendTitle to document
func (doc *Document) appendTitle() {
	title := doc.Title

	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)
	// Set x y
	doc.pdf.SetXY(120, BaseMarginTop+3)

	//// Draw rect
	//doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	//doc.pdf.Rect(120, BaseMarginTop, 80, 10, "F")

	// Draw text
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 16)
	doc.pdf.CellFormat(80, 0, doc.encodeString(title), "0", 0, "R", false, 0, "")
}

func (doc *Document) appendInvoiceHeader() {
	var x float64 = 120
	var lineBreakHeight float64 = 5
	y := doc.pdf.GetY()
	doc.pdf.SetXY(x, y)

	// Append InvoiceId
	if len(doc.InvoiceId) > 0 {
		doc.pdf.SetFont(doc.Options.Font, "", 11)
		doc.pdf.SetTextColor(
			doc.Options.GreyTextColor[0],
			doc.Options.GreyTextColor[1],
			doc.Options.GreyTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.Options.TextInvoiceIdTitle)), "0", 0, "L", false, 0, "")
		doc.pdf.SetXY(x+40, y)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.InvoiceId)), "0", 0, "R", false, 0, "")
	}

	// Append InvoiceNumber
	y = doc.pdf.GetY() + lineBreakHeight
	doc.pdf.SetXY(x, y)
	doc.pdf.SetFont(doc.Options.Font, "", 11)
	doc.pdf.SetTextColor(
		doc.Options.GreyTextColor[0],
		doc.Options.GreyTextColor[1],
		doc.Options.GreyTextColor[2],
	)
	doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.FitRefundString(doc.Options.TextInvoiceNumberTitle))), "0", 0, "L", false, 0, "")
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)
	doc.pdf.SetXY(x+40, y)
	doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.InvoiceNumber)), "0", 0, "R", false, 0, "")

	// Append InvoiceOriginNumber
	if len(doc.InvoiceOriginNumber) > 0 {
		y = doc.pdf.GetY() + lineBreakHeight
		doc.pdf.SetXY(x, y)
		doc.pdf.SetFont(doc.Options.Font, "", 11)
		doc.pdf.SetTextColor(
			doc.Options.GreyTextColor[0],
			doc.Options.GreyTextColor[1],
			doc.Options.GreyTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.Options.TextOriginInvoiceNumberTitle)), "0", 0, "L", false, 0, "")
		doc.pdf.SetXY(x+40, y)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.InvoiceOriginNumber)), "0", 0, "R", false, 0, "")
	}

	// Append InvoiceType
	if len(doc.InvoiceType) > 0 {
		y = doc.pdf.GetY() + lineBreakHeight
		doc.pdf.SetXY(x, y)
		doc.pdf.SetFont(doc.Options.Font, "", 11)
		doc.pdf.SetTextColor(
			doc.Options.GreyTextColor[0],
			doc.Options.GreyTextColor[1],
			doc.Options.GreyTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("Invoice Type")), "0", 0, "L", false, 0, "")
		doc.pdf.SetXY(x+40, y)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.InvoiceType)), "0", 0, "R", false, 0, "")
	}

	// Append InvoiceDate
	if len(doc.InvoiceDate) > 0 {
		y = doc.pdf.GetY() + lineBreakHeight
		doc.pdf.SetXY(x, y)
		doc.pdf.SetFont(doc.Options.Font, "", 11)
		doc.pdf.SetTextColor(
			doc.Options.GreyTextColor[0],
			doc.Options.GreyTextColor[1],
			doc.Options.GreyTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.FitRefundString(doc.Options.TextInvoiceDateTitle))), "0", 0, "L", false, 0, "")
		doc.pdf.SetXY(x+40, y)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
		doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.InvoiceDate)), "0", 0, "R", false, 0, "")
	}

	// Append paidDate
	paidDate := "-"
	if len(doc.PaidDate) > 0 {
		paidDate = doc.PaidDate
	}
	//dateString := fmt.Sprintf("%s: %s", doc.FitRefundString(doc.Options.TextInvoicePaidDateTitle), paidDate)
	y = doc.pdf.GetY() + lineBreakHeight
	doc.pdf.SetXY(x, y)
	doc.pdf.SetFont(doc.Options.Font, "", 11)
	doc.pdf.SetTextColor(
		doc.Options.GreyTextColor[0],
		doc.Options.GreyTextColor[1],
		doc.Options.GreyTextColor[2],
	)
	doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", doc.FitRefundString(doc.Options.TextInvoicePaidDateTitle))), "0", 0, "L", false, 0, "")
	doc.pdf.SetXY(x+40, y)
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)
	doc.pdf.CellFormat(40, 12, doc.encodeString(fmt.Sprintf("%s", paidDate)), "0", 0, "R", false, 0, "")
}

func (doc *Document) appendInvoiceStatus(y float64) float64 {
	startY := y
	SetBaseTextColor(doc)
	doc.pdf.SetXY(110, startY)
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 10)
	doc.pdf.CellFormat(80, 8, doc.encodeFitRefundString("Invoice status:"), "0", 0, "L", false, 0, "")

	doc.pdf.SetXY(110, startY+4)
	doc.pdf.SetFont(doc.Options.BoldFont, "B", 18)
	if doc.Status == "Paid" {
		doc.pdf.SetTextColor(
			doc.Options.PaidTextColor[0],
			doc.Options.PaidTextColor[1],
			doc.Options.PaidTextColor[2],
		)
	}
	doc.pdf.CellFormat(80, 18, doc.encodeString(doc.Status), "0", 0, "L", false, 0, "")
	SetBaseTextColor(doc)
	return doc.pdf.GetY() + 20
}

func SetGrayTextColor(doc *Document) {
	doc.pdf.SetTextColor(
		doc.Options.GreyTextColor[0],
		doc.Options.GreyTextColor[1],
		doc.Options.GreyTextColor[2],
	)
}

func SetBaseTextColor(doc *Document) {
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)
}

func (doc *Document) appendDescription() {
	if len(doc.Description) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + 10)
		doc.pdf.SetFont(doc.Options.Font, "", 12)
		doc.pdf.CellFormat(190, 5, doc.encodeString(doc.Description), "B", 0, "L", false, 0, "")
	}
}

func (doc *Document) drawsTableTitles(fontSize float64) {
	// Draw table titles
	doc.pdf.SetX(10)
	doc.pdf.SetY(doc.pdf.GetY() + 5)
	doc.pdf.SetFont(doc.Options.BoldFont, "", fontSize)

	// Draw rec
	doc.pdf.SetTextColor(
		doc.Options.GreyTextColor[0],
		doc.Options.GreyTextColor[1],
		doc.Options.GreyTextColor[2],
	)
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rect(BaseMargin, doc.pdf.GetY(), 210-2*BaseMargin, 12, "F")

	// Id
	doc.pdf.SetX(ItemColIdOffset)
	doc.pdf.CellFormat(
		ItemColNameOffset-ItemColIdOffset,
		12,
		doc.encodeString("#"),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Name
	doc.pdf.SetX(ItemColNameOffset)
	doc.pdf.CellFormat(
		ItemColUnitPriceOffset-ItemColNameOffset,
		12,
		doc.encodeString(doc.Options.TextItemsNameTitle),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Unit price
	if ItemColQuantityOffset-ItemColUnitPriceOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetX(ItemColUnitPriceOffset)
		doc.pdf.CellFormat(
			ItemColQuantityOffset-ItemColUnitPriceOffset,
			12,
			doc.encodeString(doc.Options.TextItemsUnitCostTitle),
			"0",
			0,
			"C",
			false,
			0,
			"",
		)
	}

	// Quantity
	if ItemColTotalHTOffset-ItemColQuantityOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetX(ItemColQuantityOffset)
		doc.pdf.CellFormat(
			ItemColTotalHTOffset-ItemColQuantityOffset,
			12,
			doc.encodeString(doc.Options.TextItemsQuantityTitle),
			"0",
			0,
			"C",
			false,
			0,
			"",
		)
	}

	// Total NoTax
	if ItemColDiscountOffset-ItemColTotalHTOffset > 0 && doc.ShowDetailItem {
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
	if ItemColTaxOffset-ItemColDiscountOffset > 0 && doc.ShowDetailItem {
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
	if ItemColDiscountOffset-ItemColTaxOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetX(ItemColTaxOffset)
		doc.pdf.CellFormat(
			ItemColTotalTTCOffset-ItemColDiscountOffset,
			12,
			doc.encodeString(doc.Options.TextItemsTaxTitle),
			"0",
			0,
			"C",
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
		"L",
		false,
		0,
		"",
	)
}

// appendItems to document
func (doc *Document) appendItems() {
	doc.drawsTableTitles(10.0)
	SetBaseTextColor(doc)

	doc.pdf.SetX(10)
	doc.pdf.SetY(doc.pdf.GetY() + 15)
	doc.pdf.SetFont(doc.Options.Font, "", 10.0)

	for i := 0; i < len(doc.Items); i++ {
		item := doc.Items[i]

		// Check item tax
		if item.Tax == nil {
			item.Tax = doc.DefaultTax
		}

		// Append to pdf
		item.appendColTo(doc.Options, i+1, doc)

		if doc.pdf.GetY() > MaxPageHeight {
			// Add page
			doc.pdf.AddPage()
			doc.drawsTableTitles(11.0)
			doc.pdf.SetFont(doc.Options.Font, "", 10)
		}

		doc.pdf.SetX(10)
		doc.pdf.SetY(doc.pdf.GetY() + 5)
	}
}

func (doc *Document) appendExchangeRate() {
	if len(doc.ExchangeRateString) == 0 {
		return
	}

	currentX := doc.pdf.GetX()
	currentY := doc.pdf.GetY()

	doc.pdf.SetFont(doc.Options.Font, "", 8)
	//doc.pdf.SetRightMargin(20)
	doc.pdf.SetY(currentY - 5)
	doc.pdf.SetX(165)

	_, lineHt := doc.pdf.GetFontSize()
	html := doc.pdf.HTMLBasicNew()
	html.Write(lineHt, doc.encodeString(fmt.Sprintf("%s", doc.ExchangeRateString)))

	doc.pdf.SetRightMargin(BaseMargin)
	doc.pdf.SetX(currentX)
	doc.pdf.SetY(currentY)
}

func (doc *Document) appendNotes() {
	if len(doc.Notes) == 0 {
		return
	}

	currentY := doc.pdf.GetY()
	doc.pdf.SetY(currentY + 20)
	doc.pdf.SetFont(doc.Options.Font, "", 11)
	SetGrayTextColor(doc)
	doc.pdf.MultiCell(70, 5, doc.encodeString("Other details"), "0", "L", false)
	doc.pdf.SetY(doc.pdf.GetY() + 3)
	SetBaseTextColor(doc)
	//doc.pdf.SetX(BaseMargin)
	//doc.pdf.SetRightMargin(130)
	//doc.pdf.SetY(currentY + 10)
	//_, lineHt := doc.pdf.GetFontSize()
	//html := doc.pdf.HTMLBasicNew()
	//html.Write(lineHt, doc.encodeString(doc.Notes))
	//doc.pdf.SetRightMargin(BaseMargin)
	doc.pdf.MultiCell(70, 5, doc.encodeString(doc.Notes), "0", "L", false)

	doc.pdf.SetY(currentY)
}

var moneyX = 165.0

// appendTotal to document
func (doc *Document) appendTotal() {
	doc.pdf.SetY(doc.pdf.GetY() + 16)
	doc.pdf.SetFont(doc.Options.Font, "", LargeTextFontSize)
	SetBaseTextColor(doc)
	var lineBreakHeight float64 = 8
	// Draw SUB TOTAL title
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	SetGrayTextColor(doc)
	doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextSubTotal), "0", 0, "R", false, 0, "")

	// Draw SUB TOTAL amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	SetBaseTextColor(doc)
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

	// Draw Promo Credit title
	if len(doc.PromoCreditString) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight)
		doc.pdf.SetX(120)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
		SetGrayTextColor(doc)
		doc.pdf.CellFormat(38, 10, doc.encodeString(doc.PromoCreditTitle), "0", 0, "R", false, 0, "")

		// Draw Promo Credit amount
		doc.pdf.SetX(moneyX)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
		SetBaseTextColor(doc)
		doc.pdf.CellFormat(
			40,
			10,
			doc.encodeString(doc.PromoCreditString),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
	}

	// Draw DISCOUNT TOTAL HT title
	if len(doc.DiscountTotalString) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight)
		doc.pdf.SetX(120)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
		SetGrayTextColor(doc)
		if len(doc.DiscountTitle) > 0 {
			doc.pdf.CellFormat(38, 10, doc.encodeString(doc.DiscountTitle), "0", 0, "R", false, 0, "")
		} else {
			doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextTotalDiscounted), "0", 0, "R", false, 0, "")
		}

		// Draw DISCOUNT TOTAL HT amount
		doc.pdf.SetX(moneyX)
		doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
		SetBaseTextColor(doc)
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
	doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight)
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	SetGrayTextColor(doc)
	if doc.IsRefund {
		doc.pdf.CellFormat(38, 10, doc.encodeString(fmt.Sprintf("VAT Reverse Charge(%s)", doc.TaxPercentageString)), "0", 0, "R", false, 0, "")
	} else {
		doc.pdf.CellFormat(38, 10, doc.encodeString(fmt.Sprintf("%s(%s)", doc.Options.TextTotalTax, doc.TaxPercentageString)), "0", 0, "R", false, 0, "")
	}

	// Draw tax amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	SetBaseTextColor(doc)
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

	// Append Exchange Rate
	if len(doc.OriginalTaxString) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight/2)
		doc.pdf.SetFont(doc.Options.Font, "", BaseTextFontSize)
		doc.pdf.SetX(moneyX)
		//doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		//doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
		doc.pdf.CellFormat(
			40,
			10,
			doc.encodeString(doc.OriginalTaxString),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
		doc.pdf.SetFont(doc.Options.Font, "", LargeTextFontSize)
	}

	// Append Exchange Rate
	if len(doc.ExchangeRateString) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight/2)
		doc.pdf.SetFont(doc.Options.Font, "", BaseTextFontSize)
		doc.pdf.SetX(moneyX)
		//doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
		//doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
		doc.pdf.CellFormat(
			40,
			10,
			doc.encodeString(doc.ExchangeRateString),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
		doc.pdf.SetFont(doc.Options.Font, "", LargeTextFontSize)
	}

	// Draw total with tax title
	doc.pdf.SetY(doc.pdf.GetY() + lineBreakHeight)
	doc.pdf.SetX(120)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(120, doc.pdf.GetY(), 40, 10, "F")
	SetGrayTextColor(doc)
	doc.pdf.CellFormat(38, 10, doc.encodeString(doc.Options.TextTotalWithTax), "0", 0, "R", false, 0, "")

	// Draw total with tax amount
	doc.pdf.SetX(moneyX)
	doc.pdf.SetFillColor(doc.Options.WhiteBgColor[0], doc.Options.WhiteBgColor[1], doc.Options.WhiteBgColor[2])
	doc.pdf.Rect(moneyX-2, doc.pdf.GetY(), 40, 10, "F")
	SetBaseTextColor(doc)
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
