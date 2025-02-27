// Package generator allows you to easily generate invoices, delivery notes and quotations in GoLang.
package generator

import (
	"context"
	"errors"

	"github.com/creasty/defaults"
	"github.com/go-pdf/fpdf"
	"github.com/leekchan/accounting"
)

var ErrInvalidDocumentType = errors.New("invalid document type")

func New(ctx context.Context, docType string, fontDirStr string, options *Options) (*Document, error) {
	_ = defaults.Set(options)

	if docType != Invoice && docType != Quotation && docType != DeliveryNote {
		return nil, ErrInvalidDocumentType
	}

	doc := &Document{
		Options: options,
		Type:    docType,
	}

	// Prepare pdf
	doc.pdf = fpdf.New("P", "mm", "A4", fontDirStr)
	doc.LoadFonts(ctx, fontDirStr)
	doc.Options.UnicodeTranslateFunc = doc.pdf.UnicodeTranslatorFromDescriptor("")

	// Prepare accounting
	doc.ac = accounting.Accounting{
		Symbol:    doc.Options.CurrencySymbol,
		Precision: doc.Options.CurrencyPrecision,
		Thousand:  doc.Options.CurrencyThousand,
		Decimal:   doc.Options.CurrencyDecimal,
	}

	return doc, nil
}
