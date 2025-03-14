package prepare

import (
	"context"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type GatewayBank struct {
	AccountHolder string `json:"accountHolder"   dc:"The AccountHolder of wire transfer " v:"required" `
	BIC           string `json:"bic"   dc:"The BIC of wire transfer " v:"required" `
	IBAN          string `json:"iban"   dc:"The IBAN of wire transfer " v:"required" `
	Address       string `json:"address"   dc:"The address of wire transfer " v:"required" `
}

func CreateTestGateway(ctx context.Context, merchantId uint64) *entity.MerchantGateway {
	one := query.GetGatewayByGatewayName(ctx, merchantId, "autotest")
	if one != nil {
		return one
	}
	if config.GetConfigInstance().IsProd() {
		err := dao.MerchantGateway.Ctx(ctx).
			Where(dao.MerchantGateway.Columns().GatewayName, "autotest").
			Where(dao.MerchantGateway.Columns().GatewayKey, "autotest").
			Where(dao.MerchantGateway.Columns().GatewaySecret, "autotest").
			OmitEmpty().
			Scan(&one)
		utility.AssertError(err, "system error")
		utility.Assert(one == nil, "same gateway exist")
	}
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   "autotest",
		Name:          "autotest",
		GatewayKey:    "autotest",
		GatewaySecret: "autotest",
		GatewayType:   consts.GatewayTypeCard,
		Logo:          "autotest",
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	one = query.GetGatewayByGatewayName(ctx, merchantId, "autotest")
	utility.Assert(one != nil, "autotest gateway error")
	return one
}

func CreateTestCryptoGateway(ctx context.Context, merchantId uint64) *entity.MerchantGateway {
	one := query.GetGatewayByGatewayName(ctx, merchantId, "autotest_crypto")
	if one != nil {
		return one
	}
	if config.GetConfigInstance().IsProd() {
		err := dao.MerchantGateway.Ctx(ctx).
			Where(dao.MerchantGateway.Columns().GatewayName, "autotest_crypto").
			Where(dao.MerchantGateway.Columns().GatewayKey, "autotest_crypto").
			Where(dao.MerchantGateway.Columns().GatewaySecret, "autotest_crypto").
			OmitEmpty().
			Scan(&one)
		utility.AssertError(err, "system error")
		utility.Assert(one == nil, "same gateway exist")
	}
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   "autotest_crypto",
		Name:          "autotest_crypto",
		GatewayKey:    "autotest_crypto",
		GatewaySecret: "autotest_crypto",
		GatewayType:   consts.GatewayTypeCrypto,
		Logo:          "autotest_crypto",
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	one = query.GetGatewayByGatewayName(ctx, merchantId, "autotest_crypto")
	utility.Assert(one != nil, "autotest gateway error")
	return one
}

func CreateTestWireTransferGateway(ctx context.Context, merchantId uint64) *entity.MerchantGateway {
	one := query.GetGatewayByGatewayName(ctx, merchantId, "wire_transfer")
	if one != nil {
		return one
	}
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   "wire_transfer",
		Name:          "autotest_wire_transfer",
		GatewayKey:    "autotest_wire_transfer",
		GatewaySecret: "autotest_wire_transfer",
		GatewayType:   consts.GatewayTypeWireTransfer,
		Logo:          "autotest_wire_transfer",
		Currency:      "USD",
		MinimumAmount: 10,
		BankData: utility.MarshalToJsonString(&GatewayBank{
			AccountHolder: "testAccountHolder",
			BIC:           "testBic",
			IBAN:          "testIBAN",
			Address:       "testAddress",
		}),
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	one = query.GetGatewayByGatewayName(ctx, merchantId, "wire_transfer")
	utility.Assert(one != nil, "autotest gateway error")
	return one
}
