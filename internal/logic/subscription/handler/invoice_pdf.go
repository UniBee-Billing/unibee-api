package handler

import (
	"context"
	"fmt"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/oss"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"golang.org/x/text/currency"
	"golang.org/x/text/number"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//https://unidoc.io/post/simple-invoices/

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Fail: %v\n", err)
	}
}

func GenerateAndUploadInvoicePdf(ctx context.Context, unibInvoice *entity.Invoice) string {
	creater := creator.New()
	// Create a new PDF page and select it for editing
	creater.NewPage()
	// Create new invoice and populate it with data
	invoice := createInvoice(ctx, creater, unibInvoice)
	// Write invoice to page
	err := creater.Draw(invoice)
	utility.Assert(err == nil, fmt.Sprintf("GenerateInvoicePdf Draw error:%v", err))
	// Write output file.
	// Alternative is writing to a Writer interface by using c.Write
	var savePath = fmt.Sprintf("%s.pdf", unibInvoice.InvoiceId)
	err = creater.WriteToFile(savePath)
	utility.Assert(err == nil, fmt.Sprintf("GenerateInvoicePdf WriteToFile error:%v", err))

	upload, err := oss.UploadLocalFile(ctx, savePath, unibInvoice.InvoiceId, savePath, "0")
	if err != nil {
		return ""
	}

	return upload.Url
}

func createInvoice(ctx context.Context, creater *creator.Creator, unibInvoice *entity.Invoice) *creator.Invoice {
	// Create an instance of Logo used as a header for the invoice
	// If the image is not stored localy, you can use NewImageFromData to generate it from byte array
	utility.Assert(unibInvoice.MerchantId > 0, "invalid merchantId")
	utility.Assert(unibInvoice.UserId > 0, "invalid UserId")
	merchantInfo := query.GetMerchantInfoById(ctx, unibInvoice.MerchantId)
	utility.Assert(len(merchantInfo.CompanyLogo) > 0, "invalid CompanyLogo")
	user := query.GetUserAccountById(ctx, uint64(unibInvoice.UserId))

	tempLogoPath := downloadImage(merchantInfo.CompanyLogo)
	utility.Assert(len(tempLogoPath) > 0, "download Logo error")
	logo, err := creater.NewImageFromFile(tempLogoPath)
	checkErr(err)

	// Create a new invoice
	invoice := creater.NewInvoice()

	// Set invoice logo
	invoice.SetLogo(logo)

	var paid string
	if unibInvoice.Status == consts.InvoiceStatusPaid {
		paid = "YES"
	} else {
		paid = "NO"
	}

	// Set invoice information
	invoice.SetNumber(unibInvoice.InvoiceId)
	invoice.SetDate(FormatUnixTime(unibInvoice.PeriodStart))
	invoice.SetDueDate(FormatUnixTime(unibInvoice.PeriodEnd))
	invoice.AddInfo("Payment terms", "Due on receipt")
	invoice.AddInfo("Paid", paid)

	// Set invoice addresses
	invoice.SetSellerAddress(&creator.InvoiceAddress{
		Name:    merchantInfo.Name,
		Street:  merchantInfo.Address,
		City:    merchantInfo.Location,
		Zip:     "",
		Country: "",
		Phone:   merchantInfo.Phone,
		Email:   merchantInfo.Email,
	})

	invoice.SetBuyerAddress(&creator.InvoiceAddress{
		Name:    user.UserName,
		Street:  user.Address,
		City:    "",
		Zip:     "",
		Country: "",
		Phone:   user.Phone,
		Email:   user.Email,
	})

	var lines []*ro.ChannelDetailInvoiceItem
	err = utility.UnmarshalFromJsonString(unibInvoice.Lines, &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString Logo error:%v", err))

	// Add products to invoice
	for i, line := range lines {
		invoice.AddLine(
			fmt.Sprintf("%s #%d\n%s-%s", line.Description, i, FormatUnixTime(unibInvoice.PeriodStart), FormatUnixTime(unibInvoice.PeriodEnd)),
			strconv.FormatInt(line.Quantity, 10),
			MustParseCurrencySymbolValue(unibInvoice.Currency, line.UnitAmountExcludingTax),
			MustParseCurrencySymbolValue(unibInvoice.Currency, line.AmountExcludingTax),
		)
	}

	// Set invoice totals
	invoice.SetSubtotal(MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.SubscriptionAmountExcludingTax))
	invoice.AddTotalLine(fmt.Sprintf("Tax (%d%%)", unibInvoice.TaxPencentage), MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.TaxAmount))
	invoice.AddTotalLine("Shipping", MustParseCurrencySymbolValue(unibInvoice.Currency, 0))
	invoice.SetTotal(MustParseCurrencySymbolValue(unibInvoice.Currency, unibInvoice.TotalAmount))

	// Set invoice content sections
	invoice.SetNotes("Notes", "Thank you for your business.")
	invoice.SetTerms("Terms and conditions", "")

	return invoice
}

func customizeInvoice(i *creator.Invoice) {
	// Load custom font
	fontHelvetica := model.NewStandard14FontMustCompile(model.HelveticaName)

	// Create colors from RGB
	lightBlue := creator.ColorRGBFrom8bit(217, 240, 250)
	red := creator.ColorRGBFrom8bit(225, 0, 0)

	// Set invoice title text style
	i.SetTitleStyle(creator.TextStyle{
		Color:    red,
		Font:     fontHelvetica,
		FontSize: 32,
	})

	// Set invoice address heading style
	i.SetAddressHeadingStyle(creator.TextStyle{
		Font:     fontHelvetica,
		Color:    red,
		FontSize: 16,
	})

	// Set columns and rows styling
	//  Line formatting can be changed immediately after adding a line
	for cn, col := range i.Columns() {
		col.BackgroundColor = lightBlue
		col.BorderColor = lightBlue
		col.TextStyle.FontSize = 9
		col.Alignment = creator.CellHorizontalAlignmentCenter

		for _, line := range i.Lines() {
			line[cn].BorderColor = lightBlue
			line[cn].TextStyle.FontSize = 9
			line[cn].Alignment = creator.CellHorizontalAlignmentCenter
		}
	}

	// Change Total text syle
	titleCell, contentCell := i.Total()
	titleCell.BackgroundColor = lightBlue
	titleCell.BorderColor = lightBlue
	contentCell.BackgroundColor = lightBlue
	contentCell.BorderColor = lightBlue

	// Set Note text style
	i.SetNoteHeadingStyle(creator.TextStyle{
		Color:    red,
		Font:     fontHelvetica,
		FontSize: 16,
	})
}

func downloadImage(url string) string {
	// 获取图片文件名
	fileName := filepath.Base(url)

	currentDir, err := os.Getwd()
	// 构建本地文件路径
	localFilePath := filepath.Join(currentDir, fileName)
	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 发起 HTTP 请求获取图片
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// 检查 HTTP 状态码
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error:", response.Status)
		return ""
	}

	// 将图片内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return ""
	}
	fmt.Println("Image downloaded successfully:", localFilePath)
	return localFilePath
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

func FormatUnixTime(unixTime int64) string {
	// Convert Unix time to time.Time
	timeValue := time.Unix(unixTime, 0)

	// Format time using a layout
	return timeValue.Format("2006-01-02 15:04:05 MST")
}
