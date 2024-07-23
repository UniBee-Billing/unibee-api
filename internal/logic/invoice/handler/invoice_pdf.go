package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
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
		return "", gerror.New("File service need setup")
	}
	upload, err := oss.UploadLocalFile(ctx, filePath, invoiceId, filePath, "0")
	if err != nil {
		g.Log().Errorf(ctx, fmt.Sprintf("UploadInvoicePdf error:%v", err))
		return "", err
	}
	return upload.Url, nil
}

func createInvoicePdf(one *entity.Invoice, merchantInfo *entity.Merchant, user *entity.UserAccount, savePath string) error {
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("createInvoicePdf Unmarshal Metadata error:%s", err.Error())
		}
	}

	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(one.Currency))))
	doc, _ := generator2.New(generator2.Invoice, &generator2.Options{
		AutoPrint:      true,
		CurrencySymbol: symbol,
	})

	doc.SetFooter(&generator2.HeaderFooter{
		Text:       fmt.Sprintf("PDF Generated on %s %s", time.Now().Format(time.RFC850), one.CountryCode),
		Pagination: true,
	})

	doc.SetInvoiceNumber(one.InvoiceId)
	doc.SetInvoiceDate(one.GmtCreate.Layout("2006-01-02"))
	//doc.Description = "Test Description"
	if len(one.RefundId) > 0 {
		doc.IsRefund = true
		doc.SetOriginInvoiceNumber(one.SendNote)
	}

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
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString error:%v", err))

	if len(one.RefundId) > 0 {
		doc.SetIsRefund(true)
		for i, line := range lines {
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", line.Description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				AmountString: fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.AmountExcludingTax, one.Currency)),
			})
		}
	} else {
		for i, line := range lines {
			doc.AppendItem(&generator2.Item{
				Name:         fmt.Sprintf("%s #%d", line.Description, i),
				UnitCost:     fmt.Sprintf("%f", float64(line.UnitAmountExcludingTax)/100.0),
				Quantity:     strconv.FormatInt(line.Quantity, 10),
				AmountString: fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(line.UnitAmountExcludingTax*line.Quantity, one.Currency)),
			})
		}
	}
	doc.SetDefaultTax(&generator2.Tax{
		Percent: utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
	})
	doc.SubTotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.SubscriptionAmountExcludingTax, one.Currency))
	if one.DiscountAmount > 0 {
		if len(one.DiscountCode) > 0 {
			doc.DiscountTitle = fmt.Sprintf("TOTAL DISCOUNTED(%s)", one.DiscountCode)
		}
		doc.DiscountTotalString = fmt.Sprintf("%s -%s", symbol, utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency))
	}
	doc.TotalString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency))
	doc.TaxString = fmt.Sprintf("%s%s", symbol, utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency))
	doc.TaxPercentageString = fmt.Sprintf("%s%s", utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage), "%")

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
