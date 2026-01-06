package invoice

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/invoice/service"
	preload2 "unibee/internal/logic/preload"
	"unibee/internal/logic/subscription/config"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskInvoiceV2Export struct {
}

func (t TaskInvoiceV2Export) TaskName() string {
	return fmt.Sprintf("InvoiceExport")
}

func (t TaskInvoiceV2Export) Header() interface{} {
	return ExportInvoiceEntity{}
}

func (t TaskInvoiceV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.InvoiceListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZoneStr = value
				timeZone = zone
			}
		}
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = uint64(value)
		}
		if value, ok := payload["gatewayIds"].([]interface{}); ok {
			req.GatewayIds = export.JsonArrayTypeConvertInt64(ctx, value)
		}
		if value, ok := payload["firstName"].(string); ok {
			req.FirstName = value
		}
		if value, ok := payload["lastName"].(string); ok {
			req.LastName = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		//if value, ok := payload["deleteInclude"].(bool); ok {
		//	req.DeleteInclude = value
		//}
		if value, ok := payload["type"].(float64); ok {
			req.Type = unibee.Int(int(value))
		}
		if value, ok := payload["sendEmail"].(string); ok {
			req.SendEmail = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["amountStart"].(float64); ok {
			req.AmountStart = unibee.Int64(int64(value))
		}
		if value, ok := payload["amountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeStart"].(float64); ok {
			req.ReportTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["reportTimeEnd"].(float64); ok {
			req.ReportTimeEnd = int64(value) - timeZone
		}
	}
	req.SkipTotal = true
	invoices := invoiceList(ctx, req)
	preload := preload2.InvoiceListPreload(ctx, invoices)
	for _, one := range invoices {
		var invoiceGateway = ""
		if preload.Gateways[one.GatewayId] != nil {
			invoiceGateway = preload.Gateways[one.GatewayId].GatewayName
		}
		userAccount := &bean.UserAccount{}
		if preload.Users[one.UserId] != nil {
			userAccount = bean.SimplifyUserAccount(preload.Users[one.UserId])
		}
		userSnapshot := &bean.UserAccount{}
		if len(one.Data) > 0 {
			var userSnapShotEntity *entity.UserAccount
			err = gjson.Unmarshal([]byte(one.Data), &userSnapShotEntity)
			if err != nil {
				fmt.Printf("UserSnapshot Unmarshal Metadata error:%s", err.Error())
				userSnapshot = userAccount
			} else {
				userSnapshot = bean.SimplifyUserAccount(userSnapShotEntity)
			}
		} else {
			userSnapshot = userAccount
		}
		subscription := &entity.Subscription{}
		if preload.Subscriptions[one.SubscriptionId] != nil {
			subscription = preload.Subscriptions[one.SubscriptionId]
		}
		payment := &bean.Payment{}
		if preload.Payments[one.PaymentId] != nil {
			payment = bean.SimplifyPayment(preload.Payments[one.PaymentId])
		}
		invoiceType := "Tax invoice"
		OriginInvoiceId := ""
		var refund *bean.Refund
		if len(one.RefundId) > 0 {
			invoiceType = "Credit Note"
			OriginInvoiceId = payment.InvoiceId
			if preload.Refunds[one.RefundId] != nil {
				refund = bean.SimplifyRefund(preload.Refunds[one.RefundId])
			}
		}
		userType := "Individual"
		if userAccount.Type == 2 {
			userType = "Business"
		}
		var dueTime int64 = 0
		if one.FinishTime > 0 {
			dueTime = one.FinishTime + one.DayUtilDue*86400
		}
		var billingPeriod = ""
		if one.BizType == consts.BizTypeSubscription {
			if subscription != nil && subscription.PlanId > 0 {
				plan := query.GetPlanById(ctx, subscription.PlanId)
				if plan != nil {
					if plan.IntervalCount <= 1 {
						billingPeriod = plan.IntervalUnit
					} else {
						billingPeriod = fmt.Sprintf("%d x %s", plan.IntervalCount, plan.IntervalUnit)
					}
				}
			}
		} else {
			billingPeriod = "one time purchase"
		}
		countryName := ""
		IsEu := ""
		// Use preloaded VAT data to avoid N+1 queries
		key := fmt.Sprintf("%d_%s", one.MerchantId, one.CountryCode)
		vatCountryRate := preload.VatCountryRates[key]
		if vatCountryRate != nil {
			countryName = vatCountryRate.CountryName
			if vatCountryRate.IsEU {
				IsEu = "EU"
			} else {
				IsEu = "Non-EU"
			}
		}
		var productIdStr string
		var planIdStr string
		//var planNameStr string
		//var planInternalNameStr string
		//var planIntervalUnitStr string
		//var planIntervalCountStr string
		var lines []*bean.InvoiceItemSimplify
		err = utility.UnmarshalFromJsonString(one.Lines, &lines)
		for _, line := range lines {
			line.Currency = one.Currency
			line.TaxPercentage = one.TaxPercentage
		}
		if err != nil {
			fmt.Printf("ConvertInvoiceLines err:%s", err)
		}
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
			}
		}
		_, planSnapshot := detail.GetInvoicePlanSnapshot(ctx, one, metadata, lines)
		if planSnapshot != nil && planSnapshot.Plan != nil {
			planIdStr = fmt.Sprintf("%v", planSnapshot.Plan.Id)
			productIdStr = fmt.Sprintf("%v", &planSnapshot.Plan.ProductId)
			//planNameStr = onePlan.PlanName
			//planInternalNameStr = onePlan.InternalName
			//planIntervalUnitStr = onePlan.IntervalUnit
			//planIntervalCountStr = fmt.Sprintf("%d", onePlan.IntervalCount)

		}
		var transactionType = "payment"
		var transactionId = one.PaymentId
		var externalTransactionId = payment.GatewayPaymentId
		if refund != nil {
			transactionType = "refund"
			transactionId = one.RefundId
			externalTransactionId = refund.GatewayRefundId
		}
		var exchangeCurrency = ""
		var exchangeAmount = ""
		var exchangeRate = ""
		if refund == nil && payment.GatewayCurrencyExchange != nil {
			exchangeCurrency = payment.GatewayCurrencyExchange.ToCurrency
			exchangeAmount = utility.ConvertCentToDollarStr(payment.GatewayCurrencyExchange.ExchangeAmount, payment.GatewayCurrencyExchange.ToCurrency)
			exchangeRate = fmt.Sprintf("%.2f", payment.GatewayCurrencyExchange.ExchangeRate)
		} else if refund != nil && refund.GatewayCurrencyExchange != nil {
			exchangeCurrency = refund.GatewayCurrencyExchange.ToCurrency
			exchangeAmount = utility.ConvertCentToDollarStr(refund.GatewayCurrencyExchange.ExchangeAmount, refund.GatewayCurrencyExchange.ToCurrency)
			exchangeRate = fmt.Sprintf("%.2f", refund.GatewayCurrencyExchange.ExchangeRate)
		}
		var promoCreditChanged int64 = 0
		var promoCreditDiscountAmount int64 = 0
		if preload.PromoCreditTransactions[one.InvoiceId] != nil {
			promoCreditChanged = preload.PromoCreditTransactions[one.InvoiceId].DeltaAmount
			promoCreditDiscountAmount = one.PromoCreditDiscountAmount
		}
		mainList = append(mainList, &ExportInvoiceEntity{
			InvoiceId:                      one.InvoiceId,
			InvoiceNumber:                  fmt.Sprintf("%s%s", api.GatewayShortNameMapping[invoiceGateway], one.InvoiceId),
			UserId:                         fmt.Sprintf("%v", one.UserId),
			ExternalUserId:                 fmt.Sprintf("%v", userAccount.ExternalUserId),
			FirstName:                      userSnapshot.FirstName,
			LastName:                       userSnapshot.LastName,
			FullName:                       fmt.Sprintf("%s %s", userSnapshot.FirstName, userSnapshot.LastName),
			UserType:                       userType,
			Email:                          userSnapshot.Email,
			City:                           userSnapshot.City,
			Address:                        userSnapshot.Address,
			InvoiceName:                    one.InvoiceName,
			ProductId:                      productIdStr,
			ProductName:                    one.ProductName,
			TaxCode:                        one.CountryCode,
			CountryCode:                    one.CountryCode,
			PostCode:                       userSnapshot.ZipCode,
			VatNumber:                      one.VatNumber,
			CountryName:                    countryName,
			IsEU:                           IsEu,
			InvoiceType:                    invoiceType,
			Gateway:                        invoiceGateway,
			CompanyName:                    userSnapshot.CompanyName,
			MerchantName:                   merchant.CompanyName,
			DiscountCode:                   one.DiscountCode,
			OriginalAmount:                 utility.ConvertCentToDollarStr(one.TotalAmount+one.DiscountAmount+one.PromoCreditDiscountAmount, one.Currency),
			TotalAmount:                    utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
			DiscountAmount:                 utility.ConvertCentToDollarStr(one.DiscountAmount, one.Currency),
			TotalAmountExcludingTax:        utility.ConvertCentToDollarStr(one.TotalAmountExcludingTax, one.Currency),
			Currency:                       one.Currency,
			TaxAmount:                      utility.ConvertCentToDollarStr(one.TaxAmount, one.Currency),
			TaxPercentage:                  utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
			SubscriptionAmount:             utility.ConvertCentToDollarStr(one.SubscriptionAmount, one.Currency),
			SubscriptionAmountExcludingTax: utility.ConvertCentToDollarStr(one.SubscriptionAmountExcludingTax, one.Currency),
			PeriodEnd:                      gtime.NewFromTimeStamp(one.PeriodEnd + timeZone),
			PeriodStart:                    gtime.NewFromTimeStamp(one.PeriodStart + timeZone),
			FinishTime:                     gtime.NewFromTimeStamp(one.FinishTime + timeZone),
			DueDate:                        gtime.NewFromTimeStamp(dueTime + timeZone),
			CreateTime:                     gtime.NewFromTimeStamp(one.CreateTime + timeZone),
			PaidTime:                       gtime.NewFromTimeStamp(payment.PaidTime + timeZone),
			Status:                         consts.InvoiceStatusToEnum(one.Status).Description(),
			PaymentId:                      one.PaymentId,
			RefundId:                       one.RefundId,
			SubscriptionId:                 one.SubscriptionId,
			PlanId:                         planIdStr,
			TrialEnd:                       gtime.NewFromTimeStamp(one.TrialEnd + timeZone),
			BillingCycleAnchor:             gtime.NewFromTimeStamp(one.BillingCycleAnchor + timeZone),
			OriginalInvoiceId:              OriginInvoiceId,
			BillingPeriod:                  billingPeriod,
			TransactionType:                transactionType,
			TransactionId:                  transactionId,
			ExternalTransactionId:          externalTransactionId,
			TimeZone:                       timeZoneStr,
			ExchangeAmount:                 exchangeAmount,
			ExchangeCurrency:               exchangeCurrency,
			ExchangeRate:                   exchangeRate,
			PromoCreditChanged:             utility.ConvertCreditAmountToDollarStr(promoCreditChanged, one.Currency, consts.CreditAccountTypePromo),
			PromoCreditDiscountAmount:      utility.ConvertCentToDollarStr(promoCreditDiscountAmount, one.Currency),
		})
	}

	return mainList, nil
}

func invoiceList(ctx context.Context, req *service.InvoiceListInternalReq) []*entity.Invoice {
	var mainList = make([]*entity.Invoice, 0)
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	var isDeletes = []int{0}
	if req.DeleteInclude {
		isDeletes = append(isDeletes, 1)
	}
	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("invoice_id|gmt_create|gmt_modify|period_end|total_amount", req.SortField), "sortField should one of invoice_id|gmt_create|period_end|total_amount")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	query := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
		Where(dao.Invoice.Columns().Currency, strings.ToUpper(req.Currency))
	if !config.GetMerchantSubscriptionConfig(ctx, req.MerchantId).ShowZeroInvoice {
		query = query.WhereNot(dao.Invoice.Columns().TotalAmount, 0)
	}

	if len(req.SendEmail) > 0 {
		query = query.WhereLike(dao.Invoice.Columns().SendEmail, "%"+req.SendEmail+"%")
	}
	if req.UserId > 0 {
		query = query.Where(dao.Invoice.Columns().UserId, req.UserId)
	}
	if req.GatewayIds != nil && len(req.GatewayIds) > 0 {
		query = query.WhereIn(dao.Invoice.Columns().GatewayId, req.GatewayIds)
	}
	if req.Status != nil && len(req.Status) > 0 {
		query = query.WhereIn(dao.Invoice.Columns().Status, req.Status)
	}
	if req.AmountStart != nil && req.AmountEnd != nil {
		utility.Assert(*req.AmountStart <= *req.AmountEnd, "amountStart should lower than amountEnd")
	}
	if req.AmountStart != nil {
		query = query.WhereGTE(dao.Invoice.Columns().TotalAmount, &req.AmountStart)
	}
	if req.AmountEnd != nil {
		query = query.WhereLTE(dao.Invoice.Columns().TotalAmount, &req.AmountEnd)
	}
	if req.Type != nil {
		utility.Assert(*req.Type == 0 || *req.Type == 1, "type should be 0 or 1")
		if *req.Type == 0 {
			query = query.WhereNull(dao.Invoice.Columns().RefundId)
		} else if *req.Type == 1 {
			query = query.WhereNotNull(dao.Invoice.Columns().RefundId)
		}
	}
	if len(req.FirstName) > 0 || len(req.LastName) > 0 {
		var userIdList = make([]uint64, 0)
		var list []*entity.UserAccount
		userQuery := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().MerchantId, req.MerchantId)
		if len(req.FirstName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().FirstName, "%"+req.FirstName+"%")
		}
		if len(req.LastName) > 0 {
			userQuery = userQuery.WhereLike(dao.UserAccount.Columns().LastName, "%"+req.LastName+"%")
		}
		_ = userQuery.Where(dao.UserAccount.Columns().IsDeleted, 0).Scan(&list)
		for _, user := range list {
			userIdList = append(userIdList, user.Id)
		}
		if len(userIdList) == 0 {
			return mainList
		}
		query = query.WhereIn(dao.Invoice.Columns().UserId, userIdList)

	}
	if req.CreateTimeStart > 0 {
		query = query.WhereGTE(dao.Invoice.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		query = query.WhereLTE(dao.Invoice.Columns().CreateTime, req.CreateTimeEnd)
	}
	if req.ReportTimeStart > 0 {
		query = query.Where(query.Builder().WhereOrGTE(dao.Invoice.Columns().CreateTime, req.ReportTimeStart).
			WhereOrGTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeStart)))
	}
	if req.ReportTimeEnd > 0 {
		query = query.Where(query.Builder().WhereOrLTE(dao.Invoice.Columns().CreateTime, req.ReportTimeEnd).
			WhereOrLTE(dao.Invoice.Columns().GmtModify, gtime.New(req.ReportTimeEnd)))
	}
	query = query.WhereIn(dao.Invoice.Columns().IsDeleted, isDeletes).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	var err error
	if req.SkipTotal {
		err = query.Scan(&mainList)
	} else {
		err = query.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return mainList
	}

	return mainList
}
