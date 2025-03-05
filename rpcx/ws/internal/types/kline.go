package types

type MarketKline struct {
	Chain        string `json:"chain"`
	TokenAddress string `json:"token_address"`
	Timestamp    int    `json:"timestamp"`
	Interval     string `json:"interval"`
	Open         string `json:"open"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Close        string `json:"close"`
	Swaps        int    `json:"swaps"`
	Amount       string `json:"amount"`
	Volume       string `json:"volume"`
	Buys         int    `json:"buys"`
	BuyAmount    string `json:"buy_amount"`
	BuyVolume    string `json:"buy_volume"`
	Sells        int    `json:"sells"`
	SellAmount   string `json:"sell_amount"`
	SellVolume   string `json:"sell_volume"`
}
