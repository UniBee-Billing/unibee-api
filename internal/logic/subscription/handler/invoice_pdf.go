package handler

import (
	"context"
	"fmt"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/oss"
	"go-oversea-pay/internal/logic/subscription/handler/generator"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"golang.org/x/text/currency"
	"golang.org/x/text/number"
	"os"
	"strconv"
	"strings"
)

func GenerateAndUploadInvoicePdf(ctx context.Context, unibInvoice *entity.Invoice) string {

	utility.Assert(unibInvoice.MerchantId > 0, "invalid merchantId")
	utility.Assert(unibInvoice.UserId > 0, "invalid UserId")
	merchantInfo := query.GetMerchantInfoById(ctx, unibInvoice.MerchantId)
	utility.Assert(len(merchantInfo.CompanyLogo) > 0, "invalid CompanyLogo")
	user := query.GetUserAccountById(ctx, uint64(unibInvoice.UserId))

	var savePath = fmt.Sprintf("%s.pdf", unibInvoice.InvoiceId)
	err := createInvoicePdf(ctx, unibInvoice, merchantInfo, user, savePath)
	utility.Assert(err == nil, fmt.Sprintf("createInvoicePdf error:%v", err))

	upload, err := oss.UploadLocalFile(ctx, savePath, unibInvoice.InvoiceId, savePath, "0")
	utility.Assert(err == nil, fmt.Sprintf("UploadLocalFile error:%v", err))

	return upload.Url
}

func createInvoicePdf(ctx context.Context, unibInvoice *entity.Invoice, merchantInfo *entity.MerchantInfo, user *entity.UserAccount, savePath string) error {
	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice: "INVOICE",
		AutoPrint:       true,
		CurrencySymbol:  fmt.Sprintf("%v", currency.Symbol(currency.MustParseISO(strings.ToUpper(unibInvoice.Currency)))),
	})

	doc.SetHeader(&generator.HeaderFooter{
		Text:       "<center>Unibee Billing</center>",
		Pagination: true,
	})

	doc.SetFooter(&generator.HeaderFooter{
		Text:       "<center>Unibee Billing</center>",
		Pagination: true,
	})

	doc.SetRef("testref")
	doc.SetVersion("1.0")

	doc.SetDescription("Subscriptions")
	if unibInvoice.Status == consts.InvoiceStatusProcessing {
		doc.SetNotes("<a href='" + unibInvoice.Link + "'>Invoice Pay Link</a>")
	} else {
		doc.SetNotes("<a href='" + unibInvoice.Link + "'>Invoice Link</a>")
	}

	doc.SetDate(utility.FormatUnixTime(unibInvoice.PeriodStart))
	doc.SetPaymentTerm(utility.FormatUnixTime(unibInvoice.PeriodEnd))

	tempLogoPath := utility.DownloadFile(merchantInfo.CompanyLogo)
	utility.Assert(len(tempLogoPath) > 0, "download Logo error")
	logoBytes, err := os.ReadFile(tempLogoPath)
	if err != nil {
		return err
	}

	doc.SetCompany(&generator.Contact{
		Name: merchantInfo.Name,
		Logo: logoBytes,
		Address: &generator.Address{
			Address: merchantInfo.Location + " " + merchantInfo.Address,
			//PostalCode: "75000",
			City: merchantInfo.Location,
			//Country:    "France",
			//Phone:   merchantInfo.Phone,
			//Email:   merchantInfo.Email,
		},
	})
	var userName = ""
	var userAddress = ""
	if user != nil {
		userName = user.UserName
		userAddress = user.Address
	}

	doc.SetCustomer(&generator.Contact{
		Name: userName,
		Address: &generator.Address{
			Address: userAddress,
			//PostalCode: "29200",
			//City:       "Brest",
			//Country:    "France",
			//Phone: user.Phone,
			//Email: user.Email,
		},
	})

	var lines []*ro.ChannelDetailInvoiceItem
	err = utility.UnmarshalFromJsonString(unibInvoice.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString Logo error:%v", err))

	for i, line := range lines {
		//scale, _ := currency.Cash.Rounding(currency.MustParseISO(strings.ToUpper(unibInvoice.Currency)))
		//dec := fmt.Sprintf("%v", number.Decimal(float64(line.UnitAmountExcludingTax)/100.0, number.Scale(scale)))
		doc.AppendItem(&generator.Item{
			Name:        fmt.Sprintf("%s #%d", line.Description, i),
			Description: fmt.Sprintf("%s-%s", utility.FormatUnixTime(unibInvoice.PeriodStart), utility.FormatUnixTime(unibInvoice.PeriodEnd)),
			UnitCost:    fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
			Quantity:    strconv.FormatInt(line.Quantity, 10),
			Tax: &generator.Tax{
				Percent: utility.ConvertTaxPercentageToPercentageString(unibInvoice.TaxPercentage),
			},
			Discount: &generator.Discount{
				Percent: "0",
				Amount:  "0",
			},
		})
	}

	doc.SetDefaultTax(&generator.Tax{
		Percent: utility.ConvertTaxPercentageToPercentageString(unibInvoice.TaxPercentage),
	})

	// doc.SetDiscount(&generator.Discount{
	// Percent: "90",
	// })
	doc.SetDiscount(&generator.Discount{
		Amount: "0",
	})

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

//func createInvoice(ctx context.Context, creater *creator.Creator, unibInvoice *entity.Invoice) *creator.Invoice {
//	// Create an instance of Logo used as a header for the invoice
//	// If the image is not stored localy, you can use NewImageFromData to generate it from byte array
//	utility.Assert(unibInvoice.MerchantId > 0, "invalid merchantId")
//	utility.Assert(unibInvoice.UserId > 0, "invalid UserId")
//	merchantInfo := query.GetMerchantInfoById(ctx, unibInvoice.MerchantId)
//	utility.Assert(len(merchantInfo.CompanyLogo) > 0, "invalid CompanyLogo")
//	user := query.GetUserAccountById(ctx, uint64(unibInvoice.UserId))
//
//	tempLogoPath := downloadImage(merchantInfo.CompanyLogo)
//	utility.Assert(len(tempLogoPath) > 0, "download Logo error")
//	logo, err := creater.NewImageFromFile(tempLogoPath)
//	checkErr(err)
//
//	// Create a new invoice
//	invoice := creater.NewInvoice()
//
//	// Set invoice logo
//	invoice.SetLogo(logo)
//
//	var paid string
//	if unibInvoice.Status == consts.InvoiceStatusPaid {
//		paid = "YES"
//	} else {
//		paid = "NO"
//	}
//
//	// Set invoice information
//	invoice.SetNumber(unibInvoice.InvoiceId)
//	invoice.SetDate(FormatUnixTime(unibInvoice.PeriodStart))
//	invoice.SetDueDate(FormatUnixTime(unibInvoice.PeriodEnd))
//	invoice.AddInfo("Payment terms", "Due on receipt")
//	invoice.AddInfo("Paid", paid)
//
//	// Set invoice addresses
//	invoice.SetSellerAddress(&creator.InvoiceAddress{
//		Name:    merchantInfo.Name,
//		Street:  merchantInfo.Address,
//		City:    merchantInfo.Location,
//		Zip:     "",
//		Country: "",
//		Phone:   merchantInfo.Phone,
//		Email:   merchantInfo.Email,
//	})
//
//	invoice.SetBuyerAddress(&creator.InvoiceAddress{
//		Name:    user.UserName,
//		Street:  user.Address,
//		City:    "",
//		Zip:     "",
//		Country: "",
//		Phone:   user.Phone,
//		Email:   user.Email,
//	})
//
//	var lines []*ro.ChannelDetailInvoiceItem
//	err = utility.UnmarshalFromJsonString(unibInvoice.Lines, &lines)
//	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString Logo error:%v", err))
//
//	// Add products to invoice
//	for i, line := range lines {
//		invoice.AddLine(
//			fmt.Sprintf("%s #%d\n%s-%s", line.Description, i, FormatUnixTime(unibInvoice.PeriodStart), FormatUnixTime(unibInvoice.PeriodEnd)),
//			strconv.FormatInt(line.Quantity, 10),
//			MustParseCurrencySymbolValue(unibInvoice.Currency, line.UnitAmountExcludingTax),
//			MustParseCurrencySymbolValue(unibInvoice.Currency, line.AmountExcludingTax),
//		)
//	}
//
//	// Set invoice totals
//	invoice.SetSubtotal(MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.SubscriptionAmountExcludingTax))
//	invoice.AddTotalLine(fmt.Sprintf("Tax (%d%%)", unibInvoice.TaxPencentage), MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.TaxAmount))
//	invoice.AddTotalLine("Shipping", MustParseCurrencySymbolValue(unibInvoice.Currency, 0))
//	invoice.SetTotal(MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.TotalAmount))
//
//	// Set invoice content sections
//	invoice.SetNotes("Notes", "Thank you for your business.")
//	invoice.SetTerms("Terms and conditions", "")
//
//	return invoice
//}

func MustParseCurrencySymbolValue(currencyCode string, centAmount int64) string {
	// 将货币代码转换为 currency.Unit 类型的值
	cur := currency.MustParseISO(strings.ToUpper(currencyCode))
	// 将金额从分转换为元
	amountInYuan := float64(centAmount) / 100.0

	// 获取货币的小数位数和舍入规则
	scale, _ := currency.Cash.Rounding(cur)

	// 将数字转换为 number.Formatter 类型的值，指定小数位数和舍入增量
	dec := number.Decimal(amountInYuan, number.Scale(scale))

	// 将货币符号和数字格式化为字符串，并输出到标准输出
	return fmt.Sprintf("%v%v", currency.Symbol(cur), dec)
}
