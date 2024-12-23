package account

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func InitPromoCreditUserAccount(ctx context.Context, merchantId uint64, userId uint64) {
	list := query.GetCreditConfigList(ctx, merchantId, consts.CreditAccountTypePromo)
	for _, v := range list {
		QueryOrCreateCreditAccount(ctx, userId, v.Currency, consts.CreditAccountTypePromo)
	}
}

func QueryOrCreateCreditAccount(ctx context.Context, userId uint64, currency string, creditType int) *entity.CreditAccount {
	utility.Assert(userId > 0, "Invalid UserId")
	currency = strings.ToUpper(strings.TrimSpace(currency))
	utility.Assert(len(currency) > 0, "invalid currency")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "user not found")
	one := query.GetCreditAccountByUserId(ctx, userId, creditType, currency)
	if one == nil {
		one = &entity.CreditAccount{
			UserId:         userId,
			MerchantId:     user.MerchantId,
			Type:           creditType,
			Currency:       currency,
			Amount:         0,
			CreateTime:     gtime.Now().Timestamp(),
			RechargeEnable: 1,
			PayoutEnable:   1,
		}
		result, err := dao.CreditAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
		utility.AssertError(err, "Server Error")
		id, err := result.LastInsertId()
		utility.AssertError(err, "Server Error")
		one.Id = uint64(id)
	}
	return one
}
