package generator

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

// Item represent a 'product' or a 'service'
type Item struct {
	Name string `json:"name,omitempty" validate:"required"`
	//Description string    `json:"description,omitempty"`
	UnitCostStr  string    `json:"unit_cost_str,omitempty"`
	UnitCost     string    `json:"unit_cost,omitempty"`
	Quantity     string    `json:"quantity,omitempty"`
	AmountString string    `json:"amountString,omitempty"`
	TaxString    string    `json:"taxString,omitempty"`
	Tax          *Tax      `json:"tax,omitempty"`
	Discount     *Discount `json:"discount,omitempty"`

	_unitCost decimal.Decimal
	_quantity decimal.Decimal
}

// Prepare convert strings to decimal
func (i *Item) Prepare() error {
	// Unit cost
	unitCost, err := decimal.NewFromString(i.UnitCost)
	if err != nil {
		return err
	}
	i._unitCost = unitCost

	// Quantity
	quantity, err := decimal.NewFromString(i.Quantity)
	if err != nil {
		return err
	}
	i._quantity = quantity

	//// Tax
	//if i.Tax != nil {
	//	if err := i.Tax.Prepare(); err != nil {
	//		return err
	//	}
	//}
	//
	//// Discount
	//if i.Discount != nil {
	//	if err := i.Discount.Prepare(); err != nil {
	//		return err
	//	}
	//}

	return nil
}

// TotalWithoutTaxAndWithoutDiscount returns the total without tax and without discount
func (i *Item) TotalWithoutTaxAndWithoutDiscount() decimal.Decimal {
	quantity, _ := decimal.NewFromString(i.Quantity)
	price, _ := decimal.NewFromString(i.UnitCost)
	total := price.Mul(quantity)

	return total
}

// TotalWithoutTaxAndWithDiscount returns the total without tax and with discount
func (i *Item) TotalWithoutTaxAndWithDiscount() decimal.Decimal {
	total := i.TotalWithoutTaxAndWithoutDiscount()

	// Check discount
	if i.Discount != nil {
		dType, dNum := i.Discount.getDiscount()

		if dType == DiscountTypeAmount {
			total = total.Sub(dNum)
		} else {
			// Percent
			toSub := total.Mul(dNum.Div(decimal.NewFromFloat(100)))
			total = total.Sub(toSub)
		}
	}

	return total
}

// TotalWithTaxAndDiscount returns the total with tax and discount
func (i *Item) TotalWithTaxAndDiscount() decimal.Decimal {
	return i.TotalWithoutTaxAndWithDiscount().Add(i.TaxWithTotalDiscounted())
}

// TaxWithTotalDiscounted returns the tax with total discounted
func (i *Item) TaxWithTotalDiscounted() decimal.Decimal {
	result := decimal.NewFromFloat(0)

	if i.Tax == nil {
		return result
	}

	totalHT := i.TotalWithoutTaxAndWithDiscount()
	taxType, taxAmount := i.Tax.getTax()

	if taxType == TaxTypeAmount {
		result = taxAmount
	} else {
		divider := decimal.NewFromFloat(100)
		result = totalHT.Mul(taxAmount.Div(divider))
	}

	return result
}

func (i *Item) appendColTo(options *Options, index int, doc *Document) {
	// Get base Y (top of line)
	baseY := doc.pdf.GetY()

	// Name
	doc.pdf.SetX(ItemColNameOffset)
	//doc.pdf.SetFont("SimSun", "", 10)
	doc.pdf.MultiCell(
		ItemColUnitPriceOffset-ItemColNameOffset-4,
		5,
		doc.encodeString(strings.ReplaceAll(i.Name, fmt.Sprintf("#%d", index-1), "")),
		"",
		"L",
		false,
	)

	//// Description
	//if len(i.Description) > 0 {
	//	doc.pdf.SetX(ItemColNameOffset)
	//	doc.pdf.SetY(doc.pdf.GetY() + 1)
	//
	//	doc.pdf.SetFont(doc.Options.Font, "", SmallTextFontSize)
	//	doc.pdf.SetTextColor(
	//		doc.Options.GreyTextColor[0],
	//		doc.Options.GreyTextColor[1],
	//		doc.Options.GreyTextColor[2],
	//	)
	//
	//	doc.pdf.MultiCell(
	//		ItemColUnitPriceOffset-ItemColNameOffset,
	//		3,
	//		doc.encodeString(i.Description),
	//		"",
	//		"",
	//		false,
	//	)
	//
	//	// Reset font
	//	doc.pdf.SetFont(doc.Options.Font, "", BaseTextFontSize)
	//	doc.pdf.SetTextColor(
	//		doc.Options.BaseTextColor[0],
	//		doc.Options.BaseTextColor[1],
	//		doc.Options.BaseTextColor[2],
	//	)
	//}

	// Compute line height
	colHeight := doc.pdf.GetY() - baseY

	doc.pdf.SetY(baseY)
	doc.pdf.SetX(ItemColIdOffset)
	doc.pdf.CellFormat(
		ItemColNameOffset-ItemColIdOffset,
		colHeight,
		doc.encodeString(fmt.Sprintf("%d", index)),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Unit price
	if ItemColQuantityOffset-ItemColUnitPriceOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(ItemColUnitPriceOffset)
		doc.pdf.CellFormat(
			ItemColQuantityOffset-ItemColUnitPriceOffset,
			colHeight,
			i.UnitCostStr,
			"0",
			0,
			"C",
			false,
			0,
			"",
		)
	}

	// Quantity
	if ItemColTaxOffset-ItemColQuantityOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(ItemColQuantityOffset)
		doc.pdf.CellFormat(
			ItemColTaxOffset-ItemColQuantityOffset,
			colHeight,
			i._quantity.String(),
			"0",
			0,
			"C",
			false,
			0,
			"",
		)
	}

	// Total No Tax
	if ItemColTaxOffset-ItemColTotalHTOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(ItemColTotalHTOffset)
		doc.pdf.CellFormat(
			ItemColTaxOffset-ItemColTotalHTOffset,
			colHeight,
			doc.ac.FormatMoneyDecimal(i.TotalWithoutTaxAndWithoutDiscount()),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Discount
	if ItemColTotalTTCOffset-ItemColDiscountOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(ItemColDiscountOffset)
		if i.Discount == nil {
			doc.pdf.CellFormat(
				ItemColTotalTTCOffset-ItemColDiscountOffset,
				colHeight,
				doc.encodeString("--"),
				"0",
				0,
				"",
				false,
				0,
				"",
			)
		} else {
			// If discount
			discountType, discountAmount := i.Discount.getDiscount()
			var discountTitle string
			var discountDesc string

			dCost := i.TotalWithoutTaxAndWithoutDiscount()
			if discountType == DiscountTypePercent {
				discountTitle = fmt.Sprintf("%s %s", discountAmount, doc.encodeString("%"))

				// get amount from percent
				dAmount := dCost.Mul(discountAmount.Div(decimal.NewFromFloat(100)))
				discountDesc = fmt.Sprintf("-%s", doc.ac.FormatMoneyDecimal(dAmount))
			} else {
				discountTitle = fmt.Sprintf("%s %s", discountAmount, doc.encodeString("€"))

				// get percent from amount
				dPerc := discountAmount.Mul(decimal.NewFromFloat(100))
				dPerc = dPerc.Div(dCost)
				discountDesc = fmt.Sprintf("-%s %%", dPerc.StringFixed(2))
			}

			// discount title
			// lastY := doc.pdf.GetY()
			doc.pdf.CellFormat(
				ItemColTotalTTCOffset-ItemColDiscountOffset,
				colHeight/2,
				discountTitle,
				"0",
				0,
				"LB",
				false,
				0,
				"",
			)

			// discount desc
			doc.pdf.SetXY(ItemColDiscountOffset, baseY+(colHeight/2))
			doc.pdf.SetFont(doc.Options.Font, "", SmallTextFontSize)
			doc.pdf.SetTextColor(
				doc.Options.GreyTextColor[0],
				doc.Options.GreyTextColor[1],
				doc.Options.GreyTextColor[2],
			)

			doc.pdf.CellFormat(
				ItemColTotalTTCOffset-ItemColDiscountOffset,
				colHeight/2,
				discountDesc,
				"0",
				0,
				"LT",
				false,
				0,
				"",
			)

			// reset font and y
			doc.pdf.SetFont(doc.Options.Font, "", BaseTextFontSize)
			doc.pdf.SetTextColor(
				doc.Options.BaseTextColor[0],
				doc.Options.BaseTextColor[1],
				doc.Options.BaseTextColor[2],
			)
			doc.pdf.SetY(baseY)
		}
	}

	// Tax
	if ItemColDiscountOffset-ItemColTaxOffset > 0 && doc.ShowDetailItem {
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(ItemColTaxOffset - 8)
		if len(i.TaxString) > 0 {
			// If no tax
			doc.pdf.CellFormat(
				ItemColDiscountOffset-ItemColTaxOffset,
				colHeight,
				i.TaxString,
				"0",
				0,
				"C",
				false,
				0,
				"",
			)
		} else if i.Tax == nil {
			// If no tax
			doc.pdf.CellFormat(
				ItemColDiscountOffset-ItemColTaxOffset,
				colHeight,
				doc.encodeString("--"),
				"0",
				0,
				"",
				false,
				0,
				"",
			)
		} else {
			// If tax
			taxType, taxAmount := i.Tax.getTax()
			var taxTitle string
			var taxDesc string

			if taxType == TaxTypePercent {
				taxTitle = fmt.Sprintf("%s %s", taxAmount, "%")
				// get amount from percent
				dCost := i.TotalWithoutTaxAndWithDiscount()
				dAmount := dCost.Mul(taxAmount.Div(decimal.NewFromFloat(100)))
				taxDesc = doc.ac.FormatMoneyDecimal(dAmount)
			} else {
				taxTitle = fmt.Sprintf("%s %s", doc.ac.Symbol, taxAmount)
				dCost := i.TotalWithoutTaxAndWithDiscount()
				dPerc := taxAmount.Mul(decimal.NewFromFloat(100))
				dPerc = dPerc.Div(dCost)
				// get percent from amount
				taxDesc = fmt.Sprintf("%s %%", dPerc.StringFixed(2))
			}

			// tax title
			// lastY := doc.pdf.GetY()
			doc.pdf.CellFormat(
				ItemColDiscountOffset-ItemColTaxOffset,
				colHeight/2,
				doc.encodeString(taxTitle),
				"0",
				0,
				"LB",
				false,
				0,
				"",
			)

			// tax desc
			doc.pdf.SetXY(ItemColTaxOffset, baseY+(colHeight/2))
			doc.pdf.SetFont(doc.Options.Font, "", SmallTextFontSize)
			doc.pdf.SetTextColor(
				doc.Options.GreyTextColor[0],
				doc.Options.GreyTextColor[1],
				doc.Options.GreyTextColor[2],
			)

			doc.pdf.CellFormat(
				ItemColDiscountOffset-ItemColTaxOffset,
				colHeight/2,
				doc.encodeString(taxDesc),
				"0",
				0,
				"LT",
				false,
				0,
				"",
			)

			// reset font and y
			doc.pdf.SetFont(doc.Options.Font, "", BaseTextFontSize)
			doc.pdf.SetTextColor(
				doc.Options.BaseTextColor[0],
				doc.Options.BaseTextColor[1],
				doc.Options.BaseTextColor[2],
			)
			doc.pdf.SetY(baseY)
		}
	}

	// TOTAL
	doc.pdf.SetY(baseY)
	doc.pdf.SetX(ItemColTotalTTCOffset)
	doc.pdf.CellFormat(
		190-ItemColTotalTTCOffset,
		colHeight,
		i.AmountString,
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Set Y for next line
	doc.pdf.SetY(baseY + colHeight)
}
