package handler

import (
	"context"
	"fmt"
	"unibee-api/internal/consts"
	"unibee-api/internal/logic/gateway/ro"
	generator2 "unibee-api/internal/logic/invoice/handler/generator"
	"unibee-api/internal/logic/oss"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
	"golang.org/x/text/currency"
	"golang.org/x/text/number"
	"os"
	"strconv"
	"strings"
	"time"
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
	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(unibInvoice.Currency))))
	doc, _ := generator2.New(generator2.Invoice, &generator2.Options{
		//TextTypeInvoice: "INVOICE",
		AutoPrint:      true,
		CurrencySymbol: symbol,
	})

	//doc.SetHeader(&generator2.HeaderFooter{
	//	Text:       "<center>UniBee Billing</center>",
	//	Pagination: true,
	//})
	//
	doc.SetFooter(&generator2.HeaderFooter{
		Text:       fmt.Sprintf("PDF Generated on %s", time.Now().Format(time.RFC850)),
		Pagination: true,
	})

	doc.SetInvoiceNumber(unibInvoice.InvoiceId)
	doc.SetInvoiceDate(unibInvoice.GmtCreate.Layout("2006-01-01"))

	//doc.SetDescription("Subscriptions")
	if unibInvoice.Status == consts.InvoiceStatusProcessing {
		doc.SetStatus("Process")
		//doc.SetNotes("<a href='" + unibInvoice.Link + "'>Processing</a>")
	} else if unibInvoice.Status == consts.InvoiceStatusPaid {
		doc.SetStatus("Paid")
		//doc.SetNotes("<a href='" + unibInvoice.Link + "'>Invoice Link</a>")
	} else if unibInvoice.Status == consts.InvoiceStatusCancelled {
		doc.SetStatus("Cancelled")
		//doc.SetNotes("<a href='" + unibInvoice.Link + "'>Invoice Link</a>")
	} else if unibInvoice.Status == consts.InvoiceStatusFailed {
		doc.SetStatus("Failed")
		//doc.SetNotes("<a href='" + unibInvoice.Link + "'>Invoice Link</a>")
	}

	doc.SetPaidDate(unibInvoice.GmtModify.Layout("2006-01-02"))

	tempLogoPath := utility.DownloadFile(merchantInfo.CompanyLogo)
	utility.Assert(len(tempLogoPath) > 0, "download Logo error")
	logoBytes, err := os.ReadFile(tempLogoPath)
	if err != nil {
		return err
	}

	doc.SetLogo(logoBytes)

	doc.SetCompany(&generator2.Contact{
		Name: merchantInfo.Name,
		//Logo: logoBytes,
		Address: &generator2.Address{
			Address: merchantInfo.Location + " " + merchantInfo.Address,
			//PostalCode: "75000",
			//City: merchantInfo.Location,
			//Country:    "France",
			//Phone:   merchantInfo.Phone,
			//Email:   merchantInfo.Email,
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
			//PostalCode: "29200",
			//City:       "Brest",
			//Country:    "France",
			//Phone: user.Phone,
			//Email: user.Email,
		},
	})

	var lines []*ro.InvoiceItemDetailRo
	err = utility.UnmarshalFromJsonString(unibInvoice.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString Logo error:%v", err))

	for i, line := range lines {
		//scale, _ := currency.Cash.Rounding(currency.MustParseISO(strings.ToUpper(unibInvoice.Currency)))
		//dec := fmt.Sprintf("%v", number.Decimal(float64(line.UnitAmountExcludingTax)/100.0, number.Scale(scale)))
		doc.AppendItem(&generator2.Item{
			Name: fmt.Sprintf("%s #%d", line.Description, i),
			//Description: fmt.Sprintf("%s-%s", utility.FormatUnixTime(unibInvoice.PeriodStart), utility.FormatUnixTime(unibInvoice.PeriodEnd)),
			UnitCost: fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
			Quantity: strconv.FormatInt(line.Quantity, 10),
			//Tax: &generator2.Tax{
			//	Percent: utility.ConvertTaxScaleToPercentageString(unibInvoice.TaxScale),
			//},
			//Discount: &generator2.Discount{
			//	Percent: "0",
			//	Amount:  "0",
			//},
		})
	}

	doc.SetDefaultTax(&generator2.Tax{
		Percent: utility.ConvertTaxScaleToPercentageString(unibInvoice.TaxScale),
	})
	doc.SubTotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.SubscriptionAmountExcludingTax, unibInvoice.Currency))
	doc.TaxString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.TaxAmount, unibInvoice.Currency))
	doc.TotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(unibInvoice.TotalAmount, unibInvoice.Currency))
	doc.TaxPercentageString = fmt.Sprintf("%s%s", utility.ConvertTaxScaleToPercentageString(unibInvoice.TaxScale), "%")

	// doc.SetDiscount(&generator.Discount{
	// Percent: "90",
	// })
	//doc.SetDiscount(&generator2.Discount{
	//	Amount: "0",
	//})

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
