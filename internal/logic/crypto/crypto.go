package crypto

func GetCryptoCurrency() string {
	return "USDT"
}

func GetCryptoAmount(amount int64, taxAmount int64) int64 {
	return amount + taxAmount // todo mark 1:1 now
}
