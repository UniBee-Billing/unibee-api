package admin

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/credit"
	"unibee/internal/logic/credit/config"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type CreditAccountAdminChangeInternalReq struct {
	UserId        uint64 `json:"userId"`
	MerchantId    uint64 `json:"merchantId"`
	CreditType    int    `json:"creditType"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	AdminMemberId uint64 `json:"adminMemberId"`
}

type CreditAccountAdminChangeInternalRes struct {
	User          *bean.UserAccount   `json:"user"`
	Merchant      *bean.Merchant      `json:"merchant"`
	CreditAccount *bean.CreditAccount `json:"creditAccount"`
}

func CreditAccountAdminChange(ctx context.Context, req *CreditAccountAdminChangeInternalReq) (*CreditAccountAdminChangeInternalRes, error) {
	utility.Assert(req != nil, "request is nil")
	utility.Assert(req.UserId > 0, "invalid userId")
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.Amount != 0, "invalid amount")
	utility.Assert(req.Currency != "", "invalid currency")
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.MerchantId == req.MerchantId, "merchant not match")
	creditAccount := credit.QueryOrCreateCreditAccount(ctx, req.UserId, req.Currency, req.CreditType)
	utility.Assert(creditAccount != nil, "Credit account create failed")
	utility.AssertError(config.CheckCreditConfig(ctx, req.MerchantId, creditAccount.Type, req.Currency), "Invalid Credit Config")
	if req.Amount < 0 {
		utility.Assert(creditAccount.Amount >= -req.Amount, "no enough amount to decrement")
	}
	err := dao.CreditRecharge.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		creditAccount = query.GetCreditAccountById(ctx, creditAccount.Id)
		transactionId := utility.CreateEventId()
		trans := &entity.CreditTransaction{
			UserId:             req.UserId,
			CreditId:           creditAccount.Id,
			Currency:           strings.ToUpper(req.Currency),
			TransactionId:      transactionId,
			TransactionType:    consts.CreditTransactionAdminChange,
			CreditAmountAfter:  creditAccount.Amount + req.Amount,
			CreditAmountBefore: creditAccount.Amount,
			DeltaAmount:        req.Amount,
			BizId:              transactionId,
			Name:               req.Name,
			Description:        req.Description,
			CreateTime:         gtime.Now().Timestamp(),
			MerchantId:         req.MerchantId,
			AccountType:        creditAccount.Type,
			AdminMemberId:      req.AdminMemberId,
		}
		_, err := dao.CreditTransaction.Ctx(ctx).Data(trans).OmitNil().Insert(trans)
		if err != nil {
			return err
		}

		//op := fmt.Sprintf("%s+%d", dao.CreditAccount.Columns().Amount, req.Amount)
		//if req.Amount < 0 {
		//	op = fmt.Sprintf("%s-%d", dao.CreditAccount.Columns().Amount, -req.Amount)
		//}
		// append the credit amount
		//update, err := dao.CreditAccount.Ctx(ctx).Data(g.Map{
		//	dao.CreditAccount.Columns().Amount:    op,
		//	dao.CreditAccount.Columns().GmtModify: gtime.Now(),
		//}).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Update()
		if req.Amount > 0 {
			update, err := dao.CreditAccount.Ctx(ctx).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Increment(dao.CreditAccount.Columns().Amount, req.Amount)
			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("update credit account failed, err: %v", err.Error()))
				return err
			}
			if update == nil {
				return gerror.New("update credit amount err")
			}
			affected, err := update.RowsAffected()
			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("update credit account failed, err: %v", err.Error()))
				return err
			}
			if affected != 1 {
				return gerror.New("update credit amount err")
			}
		} else {
			update, err := dao.CreditAccount.Ctx(ctx).Where(dao.CreditAccount.Columns().Id, creditAccount.Id).Decrement(dao.CreditAccount.Columns().Amount, -req.Amount)
			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("update credit account failed, err: %v", err.Error()))
				return err
			}
			if update == nil {
				return gerror.New("update credit amount err")
			}
			affected, err := update.RowsAffected()
			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("update credit account failed, err: %v", err.Error()))
				return err
			}
			if affected != 1 {
				return gerror.New("update credit amount err")
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	creditAccount = credit.QueryOrCreateCreditAccount(ctx, req.UserId, req.Currency, req.CreditType)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     req.MerchantId,
		Target:         fmt.Sprintf("CreditAccount(%d)", creditAccount.Id),
		Content:        "AdminChange",
		UserId:         creditAccount.UserId,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &CreditAccountAdminChangeInternalRes{
		User:          bean.SimplifyUserAccount(user),
		CreditAccount: bean.SimplifyCreditAccount(creditAccount),
	}, nil
}
