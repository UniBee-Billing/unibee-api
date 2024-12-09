package credit

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/credit/transaction"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type TaskCreditTransactionExport struct {
}

func (t TaskCreditTransactionExport) TaskName() string {
	return "CreditTransactionExport"
}

func (t TaskCreditTransactionExport) Header() interface{} {
	return ExportCreditTransactionEntity{}
}

func (t TaskCreditTransactionExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
		return mainList, nil
	}
	//merchant := query.GetMerchantById(ctx, task.MerchantId)
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &transaction.CreditTransactionListInternalReq{
		MerchantId: task.MerchantId,
		Page:       page,
		Count:      count,
	}
	timeZone := 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(float64); ok {
			timeZone = int(value)
			if timeZone > 0 {
				timeZoneStr = fmt.Sprintf("UTC+%d", timeZone)
			} else if timeZone < 0 {
				timeZoneStr = fmt.Sprintf("UTC%d", timeZone)
			}
		}
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = uint64(value)
		}
		if value, ok := payload["accountType"].(float64); ok {
			req.AccountType = int(value)
		}
		if value, ok := payload["transactionTypes"].([]interface{}); ok {
			req.TransactionTypes = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
	}
	req.SkipTotal = true
	result, _ := transaction.CreditTransactionList(ctx, req)
	if result != nil && result.CreditTransactions != nil {
		for _, one := range result.CreditTransactions {
			if one.User == nil {
				one.User = &bean.UserAccount{}
			}
			accountType := "Credit"
			if one.AccountType == 2 {
				accountType = "Promo Credit"
			}
			mainList = append(mainList, &ExportCreditTransactionEntity{
				Id:                 fmt.Sprintf("%v", one.Id),
				UserId:             fmt.Sprintf("%d", one.User.Id),
				Email:              one.User.Email,
				TransactionType:    consts.CreditTransactionTypeToEnum(one.TransactionType).Description(),
				TransactionId:      one.TransactionId,
				Currency:           one.Currency,
				InvoiceId:          one.InvoiceId,
				ChangedAmount:      utility.ConvertCentToDollarStr(one.DeltaAmount, one.Currency),
				CreditAmountBefore: utility.ConvertCentToDollarStr(one.CreditAmountBefore, one.Currency),
				CreditAmountAfter:  utility.ConvertCentToDollarStr(one.CreditAmountAfter, one.Currency),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime + int64(timeZone*3600)),
				Name:               one.Name,
				Description:        one.Description,
				AccountType:        accountType,
				TimeZone:           timeZoneStr,
			})
		}
	}
	return mainList, nil
}

type ExportCreditTransactionEntity struct {
	Id                 string      `json:"Id"    comment:"" group:"Transaction"`
	UserId             string      `json:"UserId"              comment:"The unique id of user" group:"Transaction"`
	Email              string      `json:"Email"               comment:"The email of user" group:"Transaction"`
	TransactionType    string      `json:"TransactionType"    comment:"" group:"Transaction"`
	TransactionId      string      `json:"TransactionId"       comment:"" group:"Transaction"`
	Currency           string      `json:"Currency" comment:"The currency of invoice" group:"Transaction"`
	InvoiceId          string      `json:"InvoiceId"  comment:"The invoice id of transaction, pure digital" group:"Transaction"`
	ChangedAmount      string      `json:"ChangedAmount" comment:"The amount changed" group:"Transaction"`
	CreditAmountBefore string      `json:"CreditAmountBefore" comment:"The amount before transaction" group:"Transaction"`
	CreditAmountAfter  string      `json:"CreditAmountAfter" comment:"The amount after transaction" group:"Transaction"`
	CreateTime         *gtime.Time `json:"CreateTime"  layout:"2006-01-02 15:04:05"   comment:"The create time of invoice" group:"Transaction"`
	Name               string      `json:"Name" comment:"The name of transaction"  group:"Transaction"`
	Description        string      `json:"Description" comment:"The description of transaction"  group:"Transaction"`
	AccountType        string      `json:"AccountType" comment:"The type of transaction account"  group:"Transaction"`
	TimeZone           string      `json:"TimeZone"         comment:"" group:"Transaction"`
}
