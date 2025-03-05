package sol_scan

type BaseResponse struct {
	Success bool `json:"success"`
}

type TokenMarketInfo struct {
	BaseResponse

	Data struct {
		PoolAddress   string  `json:"pool_address"`
		ProgramID     string  `json:"program_id"`
		Token1        string  `json:"token1"`
		Token2        string  `json:"token2"`
		Token1Account string  `json:"token1_account"`
		Token2Account string  `json:"token2_account"`
		Token1Amount  float64 `json:"token1_amount"`
		Token2Amount  float64 `json:"token2_amount"`
	} `json:"data"`
}

type TokenPriceResp struct {
	BaseResponse

	Data []struct {
		Date  int     `json:"date"`
		Price float64 `json:"price"`
	} `json:"data"`
}

type TokenMetadataResp struct {
	BaseResponse

	Data struct {
		Supply         string  `json:"supply"`
		Address        string  `json:"address"`
		Name           string  `json:"name"`
		Symbol         string  `json:"symbol"`
		Icon           string  `json:"icon"`
		Decimals       int     `json:"decimals"`
		Holder         int     `json:"holder"`
		Creator        string  `json:"creator"`
		CreateTx       string  `json:"create_tx"`
		CreatedTime    int     `json:"created_time"`
		FirstMintTx    string  `json:"first_mint_tx"`
		FirstMintTime  int     `json:"first_mint_time"`
		Price          float64 `json:"price"`
		Volume24H      int64   `json:"volume_24h"`
		MarketCap      int     `json:"market_cap"`
		MarketCapRank  int     `json:"market_cap_rank"`
		PriceChange24H float64 `json:"price_change_24h"`
	} `json:"data"`
}
