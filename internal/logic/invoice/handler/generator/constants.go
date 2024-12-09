package generator

const (
	// Invoice define the "invoice" document type
	Invoice string = "INVOICE"

	// Quotation define the "quotation" document type
	Quotation string = "QUOTATION"

	// DeliveryNote define the "delievry note" document type
	DeliveryNote string = "DELIVERY_NOTE"

	scale = 0.35

	// BaseMargin define base margin used in documents
	BaseMargin float64 = 30 * scale
	//BASEMargin         = 10

	// BaseMarginTop define base margin top used in documents
	BaseMarginTop float64 = 40 * scale
	//BaseMarginTop = 20

	// HeaderMarginTop define base header margin top used in documents
	HeaderMarginTop float64 = 5

	// MaxPageHeight define the maximum height for a single page
	MaxPageHeight float64 = 260
)

// Cols offsets
const (
	// ItemColNameOffset ...
	ItemColNameOffset float64 = 15

	// ItemColUnitPriceOffset ...
	//ItemColUnitPriceOffset float64 = 80
	ItemColUnitPriceOffset float64 = 95

	// ItemColQuantityOffset ...
	//ItemColQuantityOffset float64 = 103
	ItemColQuantityOffset float64 = 117

	// ItemColTotalHTOffset ...
	//ItemColTotalHTOffset float64 = 113
	ItemColTotalHTOffset float64 = ItemColTotalTTCOffset

	// ItemColDiscountOffset ...
	//ItemColDiscountOffset float64 = 140
	ItemColDiscountOffset float64 = ItemColTotalTTCOffset

	// ItemColTaxOffset ...
	//ItemColTaxOffset float64 = 157
	ItemColTaxOffset float64 = 140

	// ItemColTotalTTCOffset ...
	ItemColTotalTTCOffset float64 = 165
)

var (
	// BaseTextFontSize define the base font size for text in document
	BaseTextFontSize float64 = 8

	// SmallTextFontSize define the small font size for text in document
	SmallTextFontSize float64 = 7

	// LargeTextFontSize define the large font size for text in document
	LargeTextFontSize float64 = 12
)
