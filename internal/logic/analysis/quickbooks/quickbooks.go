package quickbooks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"github.com/shopspring/decimal"
	"math"
	"time"
	"unibee/api/bean/detail"
	config2 "unibee/internal/cmd/config"
	log2 "unibee/internal/consumer/webhook/log"
	"unibee/internal/controller/link"
	"unibee/internal/logic/analysis/quickbooks/quickbooksdk"
	detail2 "unibee/internal/logic/invoice/detail"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/query"
	"unibee/utility"
)

// sdk copy from https://github.com/rwestlund/quickbooks-go

type QuickBooksAPIKeys struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func GetCloudQuickBooksPartnerAPIKeys(ctx context.Context) (*QuickBooksAPIKeys, error) {
	apiKeyRes := redismq.Invoke(ctx, &redismq.InvoiceRequest{
		Group:   "GID_UniBee_Cloud",
		Method:  "GetQuickBooksAPIKeys",
		Request: nil,
	}, 0)
	if apiKeyRes == nil {
		return nil, gerror.New("Server Error")
	}
	if !apiKeyRes.Status {
		return nil, gerror.New(fmt.Sprintf("%v", apiKeyRes.Response))
	}
	if apiKeyRes.Response == nil {
		return nil, gerror.New("quickbooks key not found")
	}
	var apiKeys *QuickBooksAPIKeys
	err := redismq.UnmarshalFromJsonString(utility.MarshalToJsonString(apiKeyRes.Response), &apiKeys)
	if err != nil {
		return nil, err
	}
	if apiKeys == nil {
		return nil, gerror.New("quickbooks key not found")
	}
	return apiKeys, nil
}

type MerchantQuickBooksConfig struct {
	SetupReturnUrl string                    `json:"setupReturnUrl"`
	Code           string                    `json:"code"`
	RealmId        string                    `json:"realmId"`
	BearerToken    *quickbooksdk.BearerToken `json:"bearerToken"`
	CompanyName    string                    `json:"companyName"`
	TokenExpired   bool                      `json:"tokenExpired"`
}

const KeyMerchantQuickBooksConfig = "KeyMerchantQuickBooksConfig"

func GetMerchantQuickBooksConfig(ctx context.Context, merchantId uint64) *MerchantQuickBooksConfig {
	config := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantQuickBooksConfig)
	if config != nil && len(config.ConfigValue) > 0 {
		var one *MerchantQuickBooksConfig
		_ = utility.UnmarshalFromJsonString(config.ConfigValue, &one)
		return one
	}
	return &MerchantQuickBooksConfig{
		Code:    "",
		RealmId: "",
	}
}

func SetupMerchantQuickBooksConfig(ctx context.Context, merchantId uint64, config *MerchantQuickBooksConfig) error {
	if merchantId == 0 || config == nil {
		return gerror.New("invalid merchantId or config")
	}
	if len(config.Code) > 0 && len(config.RealmId) > 0 {
		//get accessToken and refreshToken
		apiKeys, err := GetCloudQuickBooksPartnerAPIKeys(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "Get QuickBooks Partner API keys Error: %s", err.Error())
			return err
		}

		qbClient, err := quickbooksdk.NewClient(apiKeys.ClientId, apiKeys.ClientSecret, config.RealmId, config2.GetConfigInstance().IsProd(), "", nil)
		if err != nil {
			g.Log().Errorf(ctx, "NewClient Error: %s", err.Error())
			return err
		}
		// To do first when you receive the authorization code from quickbooks callback
		bearerToken, err := qbClient.RetrieveBearerToken(config.Code, link.GetQuickbooksAuthorizationLink())
		if err != nil {
			g.Log().Errorf(ctx, "RetrieveBearerToken Error: %s", err.Error())
			return err
		}
		config.BearerToken = bearerToken
		err = update.SetMerchantConfig(ctx, merchantId, KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(config))
		if err != nil {
			g.Log().Errorf(ctx, "SetMerchantConfig Error: %s", err.Error())
			return err
		}
		qbClient, err = quickbooksdk.NewClient(apiKeys.ClientId, apiKeys.ClientSecret, config.RealmId, config2.GetConfigInstance().IsProd(), "", config.BearerToken)
		if err != nil {
			g.Log().Errorf(ctx, "NewClient Error: %s", err.Error())
			return err
		}
		info, err := qbClient.FindCompanyInfo()
		if err != nil {
			g.Log().Errorf(ctx, "FindCompanyInfo Error: %s", err.Error())
			return err
		}
		config.CompanyName = info.CompanyName
		err = update.SetMerchantConfig(ctx, merchantId, KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(config))
		if err != nil {
			g.Log().Errorf(ctx, "SetMerchantConfig Error: %s", err.Error())
		}
		return err
	} else {
		return gerror.New("merchantId or config is invalid")
	}
}

func testOrRefreshAccessToken(ctx context.Context, merchantId uint64) *MerchantQuickBooksConfig {
	g.Log().Infof(ctx, "testOrRefreshAccessToken MerchantId:%d", merchantId)
	config := GetMerchantQuickBooksConfig(ctx, merchantId)
	if config == nil || config.BearerToken == nil || config.BearerToken.AccessToken == "" {
		g.Log().Errorf(ctx, "testOrRefreshAccessToken merchant:%d, quickbooks config not setup", merchantId)
		return nil
	}
	apiKeys, err := GetCloudQuickBooksPartnerAPIKeys(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "testOrRefreshAccessToken, Get QuickBooks Partner API keys Error: %s", err.Error())
		return nil
	}
	qbClient, err := quickbooksdk.NewClient(apiKeys.ClientId, apiKeys.ClientSecret, config.RealmId, config2.GetConfigInstance().IsProd(), "", config.BearerToken)
	if err != nil {
		g.Log().Errorf(ctx, "testOrRefreshAccessToken NewClient Error: %s", err.Error())
		return nil
	}
	_, err = qbClient.FindCompanyInfo()
	if err != nil {
		g.Log().Errorf(ctx, "testOrRefreshAccessToken, FindCompanyInfo Error: %s", err.Error())
		token, err := qbClient.RefreshToken(config.BearerToken.RefreshToken)
		if err != nil {
			g.Log().Errorf(ctx, "testOrRefreshAccessToken, RefreshToken Error: %s", err.Error())
			config.TokenExpired = true
			_ = update.SetMerchantConfig(ctx, merchantId, KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(config))
			return nil
		}
		config.BearerToken = token
		_ = update.SetMerchantConfig(ctx, merchantId, KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(config))
	}
	return config
}

func UploadPaidInvoice(ctx context.Context, invoiceId string) {
	if invoiceId == "" {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground invoiceId is empty")
		return
	}
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground invoice not found")
		return
	}
	config := testOrRefreshAccessToken(ctx, one.MerchantId)
	if config == nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground merchant quickbooks config not setup or expired")
		return
	}
	uploadPaidInvoiceBackground(invoiceId, config)
}

func uploadPaidInvoiceBackground(invoiceId string, quickBooksConfig *MerchantQuickBooksConfig) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log2.PrintPanic(ctx, err)
				return
			}
		}()
		invoiceDetail := detail2.InvoiceDetail(ctx, invoiceId)
		if invoiceDetail == nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground invoice invoiceDetail not found")
			return
		}
		apiKeys, err := GetCloudQuickBooksPartnerAPIKeys(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground GetCloudQuickBooksPartnerAPIKeys err:%s", err.Error())
			return
		}
		qbClient, err := quickbooksdk.NewClient(apiKeys.ClientId, apiKeys.ClientSecret, quickBooksConfig.RealmId, config2.GetConfigInstance().IsProd(), "", quickBooksConfig.BearerToken)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground quickbooks client err:%s", err.Error())
			return
		}
		customer, err := qbClient.FindOrCreateCustomer(ctx, invoiceDetail.UserAccount)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateCustomer err:%s", err.Error())
			return
		}
		account, err := qbClient.FindOrCreateBankAccount(invoiceDetail.Gateway)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateBankAccount err:%s", err.Error())
			return
		}
		if invoiceDetail.TotalAmount > 0 && invoiceDetail.Refund == nil {
			UploadToInvoice(ctx, invoiceDetail, qbClient, account, customer)
		} else if invoiceDetail.TotalAmount < 0 && invoiceDetail.Refund != nil {
			UploadToRefundReceipt(ctx, invoiceDetail, qbClient, account, customer)
		}
	}()
}

func UploadToSalesReceipt(ctx context.Context, detail *detail.InvoiceDetail, qbClient *quickbooksdk.Client, account *quickbooksdk.Account, customer *quickbooksdk.Customer) {
	var qbLines []quickbooksdk.SaleReceiptLine

	for _, item := range detail.Lines {
		qbItem, err := qbClient.FindOrCreateItem(ctx, item.Name, account)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
			return
		}
		unitPrice := decimal.NewFromInt(item.UnitAmountExcludingTax).Div(decimal.NewFromInt(100))
		amountExcludingTax := decimal.NewFromInt(item.AmountExcludingTax).Div(decimal.NewFromInt(100))
		quantity := float64(item.Quantity)

		lineItem := quickbooksdk.SaleReceiptLine{
			Amount:     amountExcludingTax,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesReceiptItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbItem.Id,
					Name:  item.Name,
				},
				UnitPrice: unitPrice,
				Qty:       quantity,
			},
		}
		lineItem.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, lineItem)
	}

	qbTaxItem, err := qbClient.FindOrCreateItem(ctx, "UniBee Sales Tax Collected", account)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
		return
	}

	if detail.TaxAmount > 0 {
		taxAmountFloat := decimal.NewFromInt(detail.TaxAmount).Div(decimal.NewFromInt(100))
		taxLine := quickbooksdk.SaleReceiptLine{
			Amount:     taxAmountFloat,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesReceiptItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbTaxItem.Id,
					Name:  qbTaxItem.Name,
				},
				UnitPrice:  taxAmountFloat,
				Qty:        1,
				TaxCodeRef: nil,
			},
		}
		taxLine.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, taxLine)
	}
	salesReceipt := &quickbooksdk.SalesReceipt{
		CustomerRef: &quickbooksdk.ReferenceType{
			Value: customer.Id,
			Name:  customer.DisplayName,
		},
		TxnDate:              time.Unix(detail.PaidTime, 0).Format("2006-01-02"),
		DocNumber:            detail.InvoiceId,
		Line:                 qbLines,
		GlobalTaxCalculation: "NotApplicable",
		TotalAmt:             decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)),
		DepositToAccountRef: &quickbooksdk.ReferenceType{
			Value: account.Id,
			Name:  account.Name,
		},
	}
	createdSalesReceipt, err := qbClient.CreateOrUpdateSaleReceiptsByDocNumber(salesReceipt)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateSaleReceiptsByDocNumber QuickBooks sales receipt for invoice %s: %s", detail.InvoiceId, err.Error())
	}
	if createdSalesReceipt != nil {
		g.Log().Infof(ctx, "uploadPaidInvoiceBackground Successfully submitted or updated invoice %s to QuickBooks as Sales Receipt ID: %s, DocNumber: %s", detail.InvoiceId, createdSalesReceipt.ID, createdSalesReceipt.DocNumber)
	} else {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateSaleReceiptsByDocNumber QuickBooks sales receipt for invoice %s: response is nil", detail.InvoiceId)
	}
}

func UploadToInvoice(ctx context.Context, detail *detail.InvoiceDetail, qbClient *quickbooksdk.Client, account *quickbooksdk.Account, customer *quickbooksdk.Customer) {
	var qbLines []quickbooksdk.Line

	for _, item := range detail.Lines {
		qbItem, err := qbClient.FindOrCreateItem(ctx, item.Name, account)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
			return
		}
		unitPrice := decimal.NewFromInt(item.UnitAmountExcludingTax).Div(decimal.NewFromInt(100))
		amountExcludingTax := decimal.NewFromInt(item.AmountExcludingTax).Div(decimal.NewFromInt(100))
		quantity := float64(item.Quantity)

		lineItem := quickbooksdk.Line{
			Amount:     amountExcludingTax,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbItem.Id,
					Name:  item.Name,
				},
				UnitPrice: unitPrice,
				Qty:       quantity,
			},
		}
		lineItem.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, lineItem)
	}

	qbTaxItem, err := qbClient.FindOrCreateItem(ctx, "UniBee Sales Tax Collected", account)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
		return
	}

	if detail.TaxAmount > 0 {
		taxAmountFloat := decimal.NewFromInt(detail.TaxAmount).Div(decimal.NewFromInt(100))
		taxLine := quickbooksdk.Line{
			Amount:     taxAmountFloat,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbTaxItem.Id,
					Name:  qbTaxItem.Name,
				},
				UnitPrice:  taxAmountFloat,
				Qty:        1,
				TaxCodeRef: nil,
			},
		}
		taxLine.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, taxLine)
	}
	invoice := &quickbooksdk.Invoice{
		CustomerRef: &quickbooksdk.ReferenceType{
			Value: customer.Id,
			Name:  customer.DisplayName,
		},
		TxnDate:              time.Unix(detail.PaidTime, 0).Format("2006-01-02"),
		DocNumber:            detail.InvoiceId,
		Line:                 qbLines,
		GlobalTaxCalculation: "NotApplicable",
		TotalAmt:             decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)),
		DepositToAccountRef: &quickbooksdk.ReferenceType{
			Value: account.Id,
			Name:  account.Name,
		},
	}
	createdInvoice, err := qbClient.CreateOrUpdateInvoiceByDocNumber(invoice)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateInvoiceByDocNumber QuickBooks invoice for invoice %s: %s", detail.InvoiceId, err.Error())
	}
	if createdInvoice != nil {
		g.Log().Infof(ctx, "uploadPaidInvoiceBackground Successfully submitted or updated invoice %s to QuickBooks as Invoice ID: %s, DocNumber: %s", detail.InvoiceId, createdInvoice.ID, createdInvoice.DocNumber)
		createPayment, err := qbClient.CreateOrUpdatePayment(&quickbooksdk.Payment{
			CustomerRef: &quickbooksdk.ReferenceType{
				Value: customer.Id,
				Name:  customer.DisplayName,
			},
			TotalAmt:      decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)),
			TxnDate:       time.Unix(detail.PaidTime, 0).Format("2006-01-02"),
			PaymentRefNum: detail.InvoiceId,
			Line: []quickbooksdk.PaymentLine{
				{
					Amount: decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)),
					LinkedTxn: []quickbooksdk.LinkedTxn{{
						TxnID:   createdInvoice.ID,
						TxnType: "Invoice",
					}},
				},
			},
		})
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdatePayment QuickBooks Payment for invoice %s: %s", detail.InvoiceId, err.Error())
		}
		if createPayment != nil {
			g.Log().Infof(ctx, "uploadPaidInvoiceBackground Successfully submitted or updated CreateOrUpdatePayment %s to QuickBooks as Invoice ID: %s, PaymentNumber: %s", detail.InvoiceId, createdInvoice.ID, createPayment.PaymentRefNum)
		} else {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdatePayment QuickBooks Payment for invoice %s: response is nil", detail.InvoiceId)
		}
	} else {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateInvoiceByDocNumber QuickBooks Invoice for invoice %s: response is nil", detail.InvoiceId)
	}
}

func UploadToRefundReceipt(ctx context.Context, detail *detail.InvoiceDetail, qbClient *quickbooksdk.Client, account *quickbooksdk.Account, customer *quickbooksdk.Customer) {
	var qbLines []quickbooksdk.Line
	var totalAmountExcludingTax = decimal.NewFromFloat(0.0)

	for _, item := range detail.Lines {
		qbItem, err := qbClient.FindOrCreateItem(ctx, item.Name, account)
		if err != nil {
			g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
			return
		}
		unitPrice := decimal.NewFromInt(item.UnitAmountExcludingTax).Abs().Div(decimal.NewFromInt(100))
		amountExcludingTax := decimal.NewFromInt(item.UnitAmountExcludingTax).Mul(decimal.NewFromInt(item.Quantity)).Abs().Div(decimal.NewFromInt(100))
		quantity := math.Abs(float64(item.Quantity))

		lineItem := quickbooksdk.Line{
			Amount:     amountExcludingTax,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbItem.Id,
					Name:  item.Name,
				},
				UnitPrice: unitPrice,
				Qty:       quantity,
			},
		}
		lineItem.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, lineItem)
		totalAmountExcludingTax = totalAmountExcludingTax.Add(unitPrice.Mul(decimal.NewFromFloat(quantity)))
	}

	qbTaxItem, err := qbClient.FindOrCreateItem(ctx, "UniBee Sales Tax Collected", account)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground FindOrCreateItem err:%s", err.Error())
		return
	}

	if detail.TaxAmount != 0 {
		taxAmountFloat := decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)).Abs().Sub(totalAmountExcludingTax)
		taxLine := quickbooksdk.Line{
			Amount:     taxAmountFloat,
			DetailType: "SalesItemLineDetail",
			SalesItemLineDetail: &quickbooksdk.SalesItemLineDetail{
				ItemRef: &quickbooksdk.ReferenceType{
					Value: qbTaxItem.Id,
					Name:  qbTaxItem.Name,
				},
				UnitPrice:  taxAmountFloat,
				Qty:        1,
				TaxCodeRef: nil,
			},
		}
		taxLine.SalesItemLineDetail.TaxCodeRef = &quickbooksdk.ReferenceType{Value: "NON"}
		qbLines = append(qbLines, taxLine)
	}
	salesReceipt := &quickbooksdk.RefundReceipt{
		CustomerRef: &quickbooksdk.ReferenceType{
			Value: customer.Id,
			Name:  customer.DisplayName,
		},
		TxnDate:              time.Unix(detail.PaidTime, 0).Format("2006-01-02"),
		DocNumber:            detail.InvoiceId,
		Line:                 qbLines,
		GlobalTaxCalculation: "NotApplicable",
		TotalAmt:             decimal.NewFromInt(detail.TotalAmount).Div(decimal.NewFromInt(100)).Abs(),
		DepositToAccountRef: &quickbooksdk.ReferenceType{
			Value: account.Id,
			Name:  account.Name,
		},
		PaymentRefNum: detail.Payment.InvoiceId,
	}
	createdSalesReceipt, err := qbClient.CreateOrUpdateRefundReceiptByDocNumber(salesReceipt)
	if err != nil {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateRefundReceiptByDocNumber QuickBooks refund receipt for invoice %s: %s", detail.InvoiceId, err.Error())
	}
	if createdSalesReceipt != nil {
		g.Log().Infof(ctx, "uploadPaidInvoiceBackground Successfully submitted or updated invoice %s to QuickBooks as refund Receipt ID: %s, DocNumber: %s", detail.InvoiceId, createdSalesReceipt.ID, createdSalesReceipt.DocNumber)
	} else {
		g.Log().Errorf(ctx, "uploadPaidInvoiceBackground failed to CreateOrUpdateRefundReceiptByDocNumber QuickBooks refund receipt for invoice %s: response is nil", detail.InvoiceId)
	}
}
