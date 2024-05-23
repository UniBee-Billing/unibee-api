package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/text/currency"
	"golang.org/x/text/number"
	"os"
	"strconv"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	generator2 "unibee/internal/logic/invoice/handler/generator"
	"unibee/internal/logic/oss"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func GenerateInvoicePdf(ctx context.Context, unibInvoice *entity.Invoice) string {
	utility.Assert(unibInvoice.MerchantId > 0, "invalid merchantId")
	utility.Assert(unibInvoice.UserId > 0, "invalid UserId")
	merchantInfo := query.GetMerchantById(ctx, unibInvoice.MerchantId)
	//utility.Assert(len(merchantInfo.CompanyLogo) > 0, "invalid CompanyLogo")
	user := query.GetUserAccountById(ctx, unibInvoice.UserId)
	var savePath = fmt.Sprintf("%s.pdf", unibInvoice.InvoiceId)
	err := createInvoicePdf(unibInvoice, merchantInfo, user, savePath)
	utility.Assert(err == nil, fmt.Sprintf("createInvoicePdf error:%v", err))
	return savePath
}

func UploadInvoicePdf(ctx context.Context, invoiceId string, filePath string) (string, error) {
	if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
		len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
		len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
		len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
		g.Log().Errorf(ctx, "UploadInvoicePdf error:Oss service not setup")
		return "", gerror.New("Oss service not setup")
	}
	upload, err := oss.UploadLocalFile(ctx, filePath, invoiceId, filePath, "0")
	if err != nil {
		g.Log().Errorf(ctx, fmt.Sprintf("UploadInvoicePdf error:%v", err))
		return "", err
	}
	return upload.Url, nil
}

func createInvoicePdf(unibInvoice *entity.Invoice, merchantInfo *entity.Merchant, user *entity.UserAccount, savePath string) error {
	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(unibInvoice.Currency))))
	doc, _ := generator2.New(generator2.Invoice, &generator2.Options{
		AutoPrint:      true,
		CurrencySymbol: symbol,
	})

	doc.SetFooter(&generator2.HeaderFooter{
		Text:       fmt.Sprintf("PDF Generated on %s", time.Now().Format(time.RFC850)),
		Pagination: true,
	})

	doc.SetInvoiceNumber(unibInvoice.InvoiceId)
	doc.SetInvoiceDate(unibInvoice.GmtCreate.Layout("2006-01-02"))
	if len(unibInvoice.RefundId) > 0 {
		doc.SetOriginInvoiceNumber(unibInvoice.SendNote)
	}

	if unibInvoice.Status == consts.InvoiceStatusProcessing {
		doc.SetStatus("Process")
	} else if unibInvoice.Status == consts.InvoiceStatusPaid {
		doc.SetStatus("Paid")
	} else if unibInvoice.Status == consts.InvoiceStatusCancelled {
		doc.SetStatus("Cancelled")
	} else if unibInvoice.Status == consts.InvoiceStatusFailed {
		doc.SetStatus("Failed")
	} else if unibInvoice.Status == consts.InvoiceStatusReversed {
		doc.SetStatus("Reversed")
	}

	doc.SetPaidDate(unibInvoice.GmtModify.Layout("2006-01-02"))

	if len(merchantInfo.CompanyLogo) > 0 {
		tempLogoPath := utility.DownloadFile(merchantInfo.CompanyLogo)
		utility.Assert(len(tempLogoPath) > 0, "download Logo error")
		logoBytes, err := os.ReadFile(tempLogoPath)
		if err != nil {
			return err
		}
		doc.SetLogo(logoBytes)
	}

	doc.SetCompany(&generator2.Contact{
		Name: merchantInfo.CompanyName,
		Address: &generator2.Address{
			Address: merchantInfo.Location + " " + merchantInfo.Address,
		},
	})
	var userName = ""
	var userAddress = ""
	if user != nil {
		userName = user.FirstName + " " + user.LastName
		userAddress = user.Address
	}

	doc.SetCustomer(&generator2.Contact{
		Name: userName,
		Address: &generator2.Address{
			Address: userAddress,
		},
	})

	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(unibInvoice.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString error:%v", err))

	if len(unibInvoice.RefundId) > 0 {
		for i, line := range lines {
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", line.Description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				AmountString: fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.AmountExcludingTax, unibInvoice.Currency)),
			})
		}
	} else {
		for i, line := range lines {
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", line.Description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				AmountString: fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.UnitAmountExcludingTax*line.Quantity, unibInvoice.Currency)),
			})
		}
	}
	doc.SetDefaultTax(&generator2.Tax{
		Percent: utility.ConvertTaxPercentageToPercentageString(unibInvoice.TaxPercentage),
	})
	doc.SubTotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.SubscriptionAmountExcludingTax, unibInvoice.Currency))
	if unibInvoice.DiscountAmount > 0 {
		if len(unibInvoice.DiscountCode) > 0 {
			doc.DiscountTitle = fmt.Sprintf("TOTAL DISCOUNTED(%s)", unibInvoice.DiscountCode)
		}
		doc.DiscountTotalString = fmt.Sprintf("%s -%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.DiscountAmount, unibInvoice.Currency))
	}
	doc.TotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.TotalAmount, unibInvoice.Currency))
	doc.TaxString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.TaxAmount, unibInvoice.Currency))
	doc.TaxPercentageString = fmt.Sprintf("%s%s", utility.ConvertTaxPercentageToPercentageString(unibInvoice.TaxPercentage), "%")

	pdf, err := doc.Build()
	if err != nil {
		return err
	}

	err = pdf.OutputFileAndClose(savePath)

	if err != nil {
		return err
	}
	return nil
}

func MustParseCurrencySymbolValue(currencyCode string, centAmount int64) string {
	cur := currency.MustParseISO(strings.ToUpper(currencyCode))
	amountInYuan := float64(centAmount) / 100.0
	scale, _ := currency.Cash.Rounding(cur)
	dec := number.Decimal(amountInYuan, number.Scale(scale))
	return fmt.Sprintf("%v%v", currency.Symbol(cur), dec)
}
