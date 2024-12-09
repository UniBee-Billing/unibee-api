package credit

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type ConfigListReq struct {
	g.Meta   `path:"/config_list" tags:"Credit" method:"get,post" summary:"Get Credit Config list"`
	Types    []int  `json:"types"                  description:"type list of credit account, 1-main account, 2-promo credit account"`
	Currency string `json:"currency"              description:"currency"`
}

type ConfigListRes struct {
	CreditConfigs []*bean.CreditConfig `json:"creditConfigs" dc:"CreditConfig List Object"`
}

type NewConfigReq struct {
	g.Meta                `path:"/new_config" tags:"Credit" method:"post" summary:"Setup Credit Config"`
	Name                  string                 `json:"name"                  description:"name"`                                                                                                                       // name
	Description           string                 `json:"description"           description:"description"`                                                                                                                // description
	Type                  int                    `json:"type"                  description:"type of credit account, 1-main account, 2-promo credit account" v:"required"`                                                // type of credit account, 1-main account, 2-promo credit account
	Currency              string                 `json:"currency"              description:"currency" v:"required"`                                                                                                      // currency
	ExchangeRate          int64                  `json:"exchangeRate"          description:"keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100"` // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	Recurring             int                    `json:"recurring"             description:"apply to recurring, default no, 0-no,1-yes"`                                                                                 // apply to reucrring, default yes, 0-yes,1-no
	DiscountCodeExclusive int                    `json:"discountCodeExclusive" description:"discount code exclusive when purchase, default no, 0-no, 1-yes"`                                                             // discount code exclusive when purchase, default yes, 0-yes, 1-no
	Logo                  string                 `json:"logo"                  description:"logo image base64, show when user purchase"`                                                                                 // logo image base64, show when user purchase
	LogoUrl               string                 `json:"logoUrl"               description:"logo url, show when user purchase"`                                                                                          // logo url, show when user purchase
	RechargeEnable        *int                   `json:"rechargeEnable"        description:"credit account can be recharged or not, 0-yes, 1-no"`
	PayoutEnable          *int                   `json:"payoutEnable"          description:"credit account can used or payout in purchase or not, 0-yes, 1-no"`
	PreviewDefaultUsed    *int                   `json:"previewDefaultUsed"    description:"is default used when in purchase preview, default no, 0-no, 1-yes"`
	MetaData              map[string]interface{} `json:"metaData"              description:"meta_data(json)"`
}

type NewConfigRes struct {
	CreditConfig *bean.CreditConfig `json:"creditConfig" dc:"CreditConfig Object"`
}

type EditConfigReq struct {
	g.Meta                `path:"/edit_config" tags:"Credit" method:"post" summary:"Setup Credit Config"`
	Id                    uint64                  `json:"id"                    description:"The credit config Id" v:"required"`
	Name                  *string                 `json:"name"                  description:"name"`        // name
	Description           *string                 `json:"description"           description:"description"` // description
	Recurring             *int                    `json:"recurring"             description:"apply to recurring, default no, 0-no,1-yes"`
	DiscountCodeExclusive *int                    `json:"discountCodeExclusive" description:"discount code exclusive when purchase, default no, 0-no, 1-yes"`
	Logo                  *string                 `json:"logo"                  description:"logo image base64, show when user purchase"` // logo image base64, show when user purchase
	LogoUrl               *string                 `json:"logoUrl"               description:"logo url, show when user purchase"`          // logo url, show when user purchase
	RechargeEnable        *int                    `json:"rechargeEnable"        description:"credit account can recharge or not, default yes, 0-yes, 1-no"`
	PayoutEnable          *int                    `json:"payoutEnable"          description:"credit account can payout or not, default yes, 0-yes, 1-no"`
	PreviewDefaultUsed    *int                    `json:"previewDefaultUsed"    description:"is default used when in purchase preview, default no, 0-no, 1-yes"`
	MetaData              *map[string]interface{} `json:"metaData"              description:"meta_data(json)"`
}

type EditConfigRes struct {
	CreditConfig *bean.CreditConfig `json:"creditConfig" dc:"CreditConfig Object"`
}

type DetailReq struct {
	g.Meta `path:"/detail" tags:"Credit" method:"get,post" summary:"Credit Account Detail"`
	Id     uint64 `json:"id"                 dc:"The credit account Id" v:"required"`
}

type DetailRes struct {
	CreditAccount      *detail.CreditAccountDetail `json:"creditAccount" dc:"Credit Account Object"`
	CreditTransactions []*bean.CreditTransaction   `json:"creditTransactions" dc:"Credit Transaction List"`
}

type NewCreditRechargeReq struct {
	g.Meta          `path:"/new_credit_recharge" tags:"Credit" method:"post" summary:"NewCreditRecharge" dc:"New Credit Recharge"`
	UserId          uint64 `json:"userId"  description:"id of user to recharge, either userId&currency or creditAccountId "`
	Currency        string `json:"currency" description:"currency of recharge"`
	CreditAccountId uint64 `json:"creditAccountId"  description:"id of credit account, either userId&currency or creditAccountId "`
	GatewayId       uint64 `json:"gatewayId"  v:"required"`
	RechargeAmount  int64  `json:"rechargeAmount"  v:"required"`
	Name            string `json:"name"             description:"recharge name"`
	Description     string `json:"description"       description:"recharge description"`
	ReturnUrl       string `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl       string `json:"cancelUrl" dc:"CancelUrl"`
}

type NewCreditRechargeRes struct {
	User           *bean.UserAccount    `json:"user"`
	Merchant       *bean.Merchant       `json:"merchant"`
	Gateway        *bean.Gateway        `json:"gateway"`
	CreditAccount  *bean.CreditAccount  `json:"creditAccount"`
	CreditRecharge *bean.CreditRecharge `json:"creditRecharge"`
	Invoice        *bean.Invoice        `json:"invoice"`
	Payment        *bean.Payment        `json:"payment"`
	Link           string               `json:"link"`
	Paid           bool                 `json:"paid" dc:"Paid，true|false"`
}

type CreditAccountListReq struct {
	g.Meta          `path:"/credit_account_list" tags:"Credit" method:"get" summary:"CreditAccountList" dc:"Get Credit Account list"`
	UserId          uint64 `json:"userId"  description:"filter id of user"`
	Email           string `json:"email"  description:"filter email of user"`
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type CreditAccountListRes struct {
	CreditAccounts []*detail.CreditAccountDetail `json:"creditAccounts" dc:"Credit Account List"`
	Total          int                           `json:"total" dc:"Total"`
}

type CreditTransactionListReq struct {
	g.Meta           `path:"/credit_transaction_list" tags:"Credit" method:"get" summary:"CreditTransactionList" dc:"Get Credit Transaction list"`
	AccountType      int    `json:"accountType"  description:"filter type of account, 1-main account, 2-promo credit account" v:"required"`
	UserId           uint64 `json:"userId"  description:"filter id of user"`
	Email            string `json:"email"  description:"filter email of user"`
	Currency         string `json:"currency"  description:"filter currency of account"`
	SortField        string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType         string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	TransactionTypes []int  `json:"transactionTypes" dc:"transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out" `
	Page             int    `json:"page"  dc:"Page, Start 0" `
	Count            int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart  int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd    int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type CreditTransactionListRes struct {
	CreditTransactions []*detail.CreditTransactionDetail `json:"creditTransactions" dc:"Credit Transaction List"`
	Total              int                               `json:"total" dc:"Total"`
}

type PromoCreditIncrementReq struct {
	g.Meta      `path:"/promo_credit_increment" tags:"Credit" method:"post" summary:"PromoCreditIncrement" dc:"Increase user promo credit amount"`
	UserId      uint64 `json:"userId"  description:"filter id of user" v:"required"`
	Currency    string `json:"currency" description:"currency of recharge" v:"required"`
	Amount      uint64 `json:"amount" dc:"The amount to increase, should greater than 0"  v:"required"`
	Name        string `json:"name" description:"name of increase action"`
	Description string `json:"description"  description:"description of increase action"`
}

type PromoCreditIncrementRes struct {
	UserPromoCreditAccount *bean.CreditAccount `json:"UserPromoCreditAccount" dc:"The user promo credit account object" `
}

type PromoCreditDecrementReq struct {
	g.Meta      `path:"/promo_credit_decrement" tags:"Credit" method:"post" summary:"PromoCreditDecrement" dc:"Decrease user promo credit amount, the amount after decreased should greater than 0"`
	UserId      uint64 `json:"userId"  description:"filter id of user" v:"required"`
	Currency    string `json:"currency" description:"currency of recharge" v:"required"`
	Amount      uint64 `json:"amount" dc:"The Amount to decrease, should greater than 0"  v:"required"`
	Name        string `json:"name" description:"name of increase action"`
	Description string `json:"description"  description:"description of increase action"`
}

type PromoCreditDecrementRes struct {
	UserPromoCreditAccount *bean.CreditAccount `json:"UserPromoCreditAccount" dc:"The user promo credit account object" `
}

type EditCreditAccountReq struct {
	g.Meta         `path:"/edit_credit_account" tags:"Credit" method:"post" summary:"EditSingleUserCreditAccountConfig" dc:"Edit Single User Credit Account Config"`
	Id             uint64 `json:"id"  description:"id of credit account" v:"required"`
	RechargeEnable *int   `json:"rechargeEnable"        description:"credit account can be recharged|increment or not, 0-yes, 1-no"`
	PayoutEnable   *int   `json:"payoutEnable"          description:"credit account can used or payout|apply in purchase or not, 0-yes, 1-no"`
}

type EditCreditAccountRes struct {
	UserCreditAccount *bean.CreditAccount `json:"UserCreditAccount" dc:"The user credit account object" `
}
