package types

type MarketKlineV2 struct {
	MarketAddress string  `json:"market_address"`
	Interval      string  `json:"interval"`
	Chain         string  `json:"chain"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	Timestamp     int     `json:"time"`
	Volume        float64 `json:"volume"`
}

type MarketTxV2 struct {
	Chain         string  `json:"chain"`
	BaseAmount    float64 `json:"base_amount"`
	BasePrice     float64 `json:"base_price"`
	MakerAddress  string  `json:"maker_address"`
	MarketAddress string  `json:"market_address"`
	MarketType    string  `json:"market_type"`
	QuoteAmount   float64 `json:"quote_amount"`
	QuotePrice    float64 `json:"quote_price"`
	Slot          int     `json:"slot"`
	SwapType      int     `json:"swap_type"`
	Timestamp     int     `json:"timestamp"`
	TxHash        string  `json:"tx_hash"`
	Volume        float64 `json:"volume"`
}
