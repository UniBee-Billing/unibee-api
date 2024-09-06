package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"golang.org/x/text/currency"
	"golang.org/x/text/number"
	"os"
	"strconv"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	generator2 "unibee/internal/logic/invoice/handler/generator"
	entity "unibee/internal/model/entity/default"
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

	err := createInvoicePdf(unibInvoice, merchantInfo, user, query.GetGatewayById(ctx, unibInvoice.GatewayId), savePath)
	utility.Assert(err == nil, fmt.Sprintf("createInvoicePdf error:%v", err))
	return savePath
}

//func UploadInvoicePdf(ctx context.Context, invoiceId string, filePath string) (string, error) {
//	//if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
//	//	len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
//	//	len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
//	//	len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
//	//	g.Log().Errorf(ctx, "UploadInvoicePdf error:Oss service not setup")
//	//	return "", gerror.New("File service need setup")
//	//}
//	upload, err := oss.UploadLocalFile(ctx, filePath, invoiceId, filePath, "0")
//	if err != nil {
//		g.Log().Errorf(ctx, fmt.Sprintf("UploadInvoicePdf error:%v", err))
//		return "", err
//	}
//	return upload.Url, nil
//}

func createInvoicePdf(one *entity.Invoice, merchantInfo *entity.Merchant, user *entity.UserAccount, gateway *entity.MerchantGateway, savePath string) error {
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("createInvoicePdf Unmarshal Metadata error:%s", err.Error())
		}
	}

	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(one.Currency))))
	doc, _ := generator2.New(generator2.Invoice, "/usr/share/fonts", &generator2.Options{
		AutoPrint:      true,
		CurrencySymbol: symbol,
	})
	//doc.Pdf().AddUTF8Font("SimSun", "", "simsun.ttf")
	doc.SetFooter(&generator2.HeaderFooter{
		Text:       fmt.Sprintf("PDF Generated on %s                                                    -%s-", time.Now().Format(time.RFC850), one.CountryCode),
		Pagination: true,
	})

	var invoiceGateway = ""
	if gateway != nil {
		invoiceGateway = gateway.GatewayName
	}
	//doc.SetInvoiceId(one.InvoiceId)
	doc.SetInvoiceNumber(fmt.Sprintf("%s%s", api.GatewayShortNameMapping[invoiceGateway], one.InvoiceId))
	doc.SetInvoiceDate(one.GmtCreate.Layout("2006-01-02"))

	if one.Status == consts.InvoiceStatusProcessing {
		doc.SetStatus("Process")
	} else if one.Status == consts.InvoiceStatusPaid {
		if len(one.RefundId) > 0 {
			doc.SetStatus("Refunded")
		} else {
			doc.SetStatus("Paid")
		}
	} else if one.Status == consts.InvoiceStatusCancelled {
		doc.SetStatus("Cancelled")
	} else if one.Status == consts.InvoiceStatusFailed {
		doc.SetStatus("Failed")
	} else if one.Status == consts.InvoiceStatusReversed {
		doc.SetStatus("Reversed")
	}

	doc.SetPaidDate(one.GmtModify.Layout("2006-01-02"))

	if len(merchantInfo.CompanyLogo) > 0 {
		tempLogoPath := utility.DownloadFile(merchantInfo.CompanyLogo)
		utility.Assert(len(tempLogoPath) > 0, "download Logo error")
		logoBytes, err := os.ReadFile(tempLogoPath)
		if err != nil {
			return err
		}
		doc.SetLogo(logoBytes)
	}

	//Localized currency
	localizedCurrency := metadata["LocalizedCurrency"]
	localizedCurrencyRate := metadata["LocalizedExchangeRate"]
	localizedSymbol := ""
	localizedCurrencyStr := ""
	localizedExchangeRateDescription := metadata["LocalizedExchangeRateDescription"]
	var localizedExchangeRate float64
	localized := false
	if localizedCurrencyRate != nil && localizedCurrency != nil {
		if rate, ok := localizedCurrencyRate.(float64); ok {
			localizedCurrencyStr = strings.ToUpper(fmt.Sprintf("%v", localizedCurrency))
			iso, err := currency.ParseISO(strings.ToUpper(localizedCurrencyStr))
			if err == nil {
				localizedExchangeRate = rate
				localized = true
				localizedSymbol = fmt.Sprintf("%v ", currency.NarrowSymbol(iso))
			} else {
				fmt.Printf("Invoice PDF Localized failed %s\n", err.Error())
			}
		}
	}

	doc.ShowDetailItem = true
	showDetailItem := metadata["ShowDetailItem"]
	if showDetailItem != nil {
		if _showDetailItem, ok := showDetailItem.(bool); ok {
			if _showDetailItem {
				doc.ShowDetailItem = _showDetailItem
			}
		}
	}

	var vatNumber = metadata["IssueVatNumber"]
	var regNumber = metadata["IssueRegNumber"]
	var companyName = metadata["IssueCompanyName"]
	var address = metadata["IssueAddress"]
	if vatNumber == nil {
		vatNumber = ""
	}
	if regNumber == nil {
		regNumber = ""
	}
	if companyName == nil {
		companyName = merchantInfo.CompanyName
	}
	if address == nil {
		address = merchantInfo.Address
	}
	doc.SetCompany(&generator2.Contact{
		Name: fmt.Sprintf("%s", companyName),
		Address: &generator2.Address{
			Address:   fmt.Sprintf("%s", address),
			VatNumber: fmt.Sprintf("%s", vatNumber),
			RegNumber: fmt.Sprintf("%s", regNumber),
		},
	})
	var userName = ""
	var userAddress = ""
	var userCity = ""
	var userPostalCode = ""
	var userCountry = ""
	var userRegNumber = ""
	if user != nil {
		if user.Type == 1 {
			if len(user.FirstName) == 0 || len(user.LastName) == 0 {
				userName = user.Email
			} else {
				userName = fmt.Sprintf("%s %s(%s)", user.FirstName, user.LastName, user.Email)
			}
		} else {
			if len(user.CompanyName) == 0 {
				userName = user.Email
			} else {
				userName = fmt.Sprintf("%s(%s)", user.CompanyName, user.Email)
			}
		}
		userAddress = user.Address
		userPostalCode = user.ZipCode
		userCountry = user.CountryName
		userCity = user.City
		userRegNumber = user.RegistrationNumber
	}

	doc.SetCustomer(&generator2.Contact{
		Name: userName,
		Address: &generator2.Address{
			RegNumber: userRegNumber,
			//VatNumber:  one.VatNumber,
			Country:    userCountry,
			City:       userCity,
			PostalCode: userPostalCode,
			Address:    userAddress,
		},
	})

	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString error:%v", err))

	doc.TaxPercentageString = fmt.Sprintf("%s%s", utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage), "%")
	if len(one.RefundId) > 0 {
		doc.IsRefund = true
		doc.SetOriginInvoiceNumber(one.SendNote)
		doc.Title = "TAX CREDIT NOTE"
		doc.Notes = one.CreateFrom
		if len(one.VatNumber) > 0 {
			//doc.Customer.AdditionalInfo = []string{"VAT reverse charge"}
			doc.Customer.AdditionalInfo = []string{fmt.Sprintf("VAT Number:%s", one.VatNumber)}
		}
		doc.SetIsRefund(true)
		for i, line := range lines {
			amountString := fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.Amount, one.Currency))
			taxString := fmt.Sprintf("%s %s%s", doc.TaxPercentageString, symbol, utility.ConvertCentToDollarStr(line.Tax, one.Currency))
			//if localized {
			//	amountString = fmt.Sprintf("%s | %s%s", amountString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(line.Amount)*localizedExchangeRate), localizedCurrencyStr))
			//	taxString = fmt.Sprintf("%s | %s%s", taxString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(line.Tax)*localizedExchangeRate), localizedCurrencyStr))
			//}
			description := line.Description
			if len(line.PdfDescription) > 0 {
				description = line.PdfDescription
			}
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				UnitCostStr:  fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.UnitAmountExcludingTax, one.Currency)),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				TaxString:    taxString,
				AmountString: amountString,
			})
		}
	} else {
		doc.Title = "TAX INVOICE"
		if len(one.VatNumber) > 0 {
			doc.Customer.AdditionalInfo = []string{fmt.Sprintf("VAT Number:%s", one.VatNumber)}
		}
		for i, line := range lines {
			amountString := fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.Amount, one.Currency))
			taxString := fmt.Sprintf("%s %s%s", doc.TaxPercentageString, symbol, utility.ConvertCentToDollarStr(line.Tax, one.Currency))
			//if localized {
			//	amountString = fmt.Sprintf("%s | %s%s", amountString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(line.Amount)*localizedExchangeRate), localizedCurrencyStr))
			//	taxString = fmt.Sprintf("%s | %s%s", taxString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(line.Tax)*localizedExchangeRate), localizedCurrencyStr))
			//}
			description := line.Description
			if len(line.PdfDescription) > 0 {
				description = line.PdfDescription
			}
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				UnitCostStr:  fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.UnitAmountExcludingTax, one.Currency)),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				TaxString:    taxString,
				AmountString: amountString,
			})
		}
	}
	doc.SetDefaultTax(&generator2.Tax{
		Percent: utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
	})
	doc.SubTotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.SubscriptionAmountExcludingTax, one.Currency))
	if one.DiscountAmount > 0 {
		if len(one.DiscountCode) > 0 {
			doc.DiscountTitle = fmt.Sprintf("TOTAL DISCOUNTED(code: %s)", one.DiscountCode)
		}
		doc.DiscountTotalString = fmt.Sprintf("%s -%s", symbol, utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency))
	}
	doc.TotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency))
	doc.TaxString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency))

	if localized {
		if localizedExchangeRateDescription != nil {
			doc.ExchangeRateString = fmt.Sprintf("%s", localizedExchangeRateDescription)
		} else {
			doc.ExchangeRateString = fmt.Sprintf("* %s1 = %s%s", symbol, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(100)*localizedExchangeRate), localizedCurrencyStr))
		}
		//doc.SubTotalString = fmt.Sprintf("%s | %s%s", doc.SubTotalString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(one.SubscriptionAmountExcludingTax)*localizedExchangeRate), localizedCurrencyStr))
		//doc.TotalString = fmt.Sprintf("%s | %s%s", doc.TotalString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(one.TotalAmount)*localizedExchangeRate), localizedCurrencyStr))
		doc.TaxString = fmt.Sprintf("%s | %s%s", doc.TaxString, localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(one.TaxAmount)*localizedExchangeRate), localizedCurrencyStr))
		//if one.DiscountAmount > 0 {
		//	doc.DiscountTotalString = fmt.Sprintf("%s -%s  |  %s -%s", symbol, utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency), localizedSymbol, utility.ConvertCentToDollarStr(int64(float64(one.DiscountAmount)*localizedExchangeRate), localizedCurrencyStr))
		//}
	}

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
