package crypto

import "unibee/internal/consts"

func GetCryptoCurrency() string {
	return "USDT"
}

func GetCryptoAmount(totalAmount int64, taxAmount int64) int64 {
	if consts.GetConfigInstance().IsLocal() || consts.GetConfigInstance().IsStage() {
		return totalAmount
	}
	return totalAmount // todo mark 1:1 now
}
