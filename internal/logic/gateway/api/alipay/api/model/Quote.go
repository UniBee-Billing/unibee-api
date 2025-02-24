package model

type Quote struct {
	QuoteId           string  `json:"quoteId,omitempty"`
	QuoteCurrencyPair string  `json:"quoteCurrencyPair,omitempty"`
	QuotePrice        float64 `json:"quotePrice,omitempty"`
	QuoteStartTime    string  `json:"quoteStartTime,omitempty"`
	QuoteExpiryTime   string  `json:"quoteExpiryTime,omitempty"`
	Guaranteed        bool    `json:"guaranteed,omitempty"`
}
