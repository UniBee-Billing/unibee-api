package crypto

func GetCryptoCurrency() string {
	return "USDT"
}

func GetCryptoAmount(totalAmount int64, taxAmount int64) int64 {
	return totalAmount + taxAmount // todo mark 1:1 now
}
