package dbx

import (
	"fmt"
	"strings"
	"time"
)

//func FindByPageWithJoin(queryParam *Query, orderField, Direction string, page, size int) (result []MarketJoinedData, count int64, err error) {
//	offset := (page - 1) * size
//	query := fmt.Sprintf(`SELECT *
//        FROM market_real_time_data
//        INNER JOIN market ON market_real_time_data.address = market.address
//        ORDER BY %s %s
//        LIMIT %d OFFSET %d`, orderField, Direction, size, offset)
//
//	println("query:", query)
//	var results []MarketJoinedData
//	err = queryParam.db.Raw(query).Scan(&results).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return results, int64(len(results)), nil
//}

func FindByPageWithJoinWithMarketType(queryParam *Query, orderField, direction string, page, size int, marketType string) (result []MarketJoinedData, count int64, err error) {
	offset := (page - 1) * size

	// 查询分页数据
	dataQuery := fmt.Sprintf(`
        SELECT market_real_time_data.*, market.*
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address
        WHERE market.market_type = '%s'
        ORDER BY %s %s
        LIMIT %d OFFSET %d`, marketType, orderField, direction, size, offset)

	// 查询总记录数
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address 
        WHERE market.market_type = '%s'`, marketType)

	// 执行分页查询
	var results []MarketJoinedData
	err = queryParam.db.Raw(dataQuery).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// 执行总数查询
	err = queryParam.db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func FindByPageWithJoin(queryParam *Query, orderField, direction string, page, size int) (result []MarketJoinedData, count int64, err error) {
	offset := (page - 1) * size

	// 查询分页数据
	dataQuery := fmt.Sprintf(`
        SELECT market_real_time_data.*, market.*
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address
        ORDER BY %s %s
        LIMIT %d OFFSET %d`, orderField, direction, size, offset)

	// 查询总记录数
	countQuery := `
        SELECT COUNT(*)
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address`

	// 执行分页查询
	var results []MarketJoinedData
	err = queryParam.db.Raw(dataQuery).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// 执行总数查询
	err = queryParam.db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func FindByPageWithJoinFollow(queryParam *Query, orderField, direction string, page, size int, tokenAddress []string) (result []MarketJoinedData, count int64, err error) {
	offset := (page - 1) * size

	// 构建 WHERE 子句
	var addressCondition string
	if len(tokenAddress) > 0 {
		quotedAddresses := make([]string, len(tokenAddress))
		for i, addr := range tokenAddress {
			quotedAddresses[i] = fmt.Sprintf("'%s'", addr)
		}
		addressCondition = fmt.Sprintf("WHERE market.quote_token_mint_address IN (%s)", strings.Join(quotedAddresses, ","))
	} else {
		// 如果 tokenAddress 为空，添加一个永不成立的条件
		addressCondition = "WHERE 1 = 0"
	}

	// 查询分页数据
	dataQuery := fmt.Sprintf(`
        SELECT market_real_time_data.*, market.*
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address
        %s
        ORDER BY %s %s
        LIMIT %d OFFSET %d`, addressCondition, orderField, direction, size, offset)

	// 查询总记录数
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM market_real_time_data 
        INNER JOIN market ON market_real_time_data.address = market.address
        %s`, addressCondition)

	// 执行分页查询
	var results []MarketJoinedData
	err = queryParam.db.Raw(dataQuery).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// 执行总数查询
	err = queryParam.db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func SearchCurrencyByPage(queryParam *Query, orderField string, page, size int) (result []MarketJoinedData, count int64, err error) {
	offset := (page - 1) * size
	var condStr string

	if len(orderField) == 44 {
		// 如果长度为 44，认为是地址精确匹配
		condStr = fmt.Sprintf("market.address = '%s' OR market.quote_token_mint_address = '%s'", orderField, orderField)
	} else if len(orderField) == 0 {
		condStr = "1 = 1"
		page = 0
		size = 20
	} else {
		// 否则认为是代币符号模糊匹配
		condStr = fmt.Sprintf("market.quote_symbol ILIKE '%%%s%%'", orderField)
	}

	// 分页查询
	dataQuery := fmt.Sprintf(`
		SELECT *
		FROM market_real_time_data 
		INNER JOIN market ON market_real_time_data.address = market.address
		WHERE %s 
		ORDER BY market_real_time_data.sell_volume_24h + market_real_time_data.buy_volume_24h DESC
		LIMIT %d OFFSET %d`, condStr, size, offset)

	// 总数查询
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM market_real_time_data 
		INNER JOIN market ON market_real_time_data.address = market.address
		WHERE %s`, condStr)

	// 执行分页查询
	var results []MarketJoinedData
	err = queryParam.db.Raw(dataQuery).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// 执行总数查询
	err = queryParam.db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

//func SearchCurrencyByPage(queryParam *Query, orderField string, page, size int) (result []MarketJoinedData, count int64, err error) {
//	offset := (page - 1) * size
//	var condStr string
//
//	if len(orderField) == 44 {
//		// 如果长度为 44，认为是地址精确匹配
//		condStr = fmt.Sprintf("market.address = '%s'", orderField)
//	} else {
//		// 否则认为是代币符号模糊匹配
//		condStr = fmt.Sprintf("market.quote_symbol LIKE '%%%s%%'", orderField)
//	}
//
//	query := fmt.Sprintf(`SELECT *
//	FROM market_real_time_data
//	INNER JOIN market ON market_real_time_data.address = market.address
//	WHERE %s
//	ORDER BY market_real_time_data.sell_volume_24h + market_real_time_data.buy_volume_24h DESC
//	LIMIT %d OFFSET %d`, condStr, size, offset)
//	println("query:", query)
//	var results []MarketJoinedData
//	err = queryParam.db.Raw(query).Scan(&results).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return results, int64(len(results)), nil
//}

type MarketJoinedData struct {
	// market
	ID                    int64     `gorm:"column:id" json:"id"`                                             // 主键
	CreatedTimestamp      time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`               // 创建时间戳
	OpenTimestamp         time.Time `gorm:"column:open_timestamp" json:"open_timestamp"`                     // 开放时间戳
	MarketType            string    `gorm:"column:market_type" json:"market_type"`                           // 市场类型
	BaseSymbol            string    `gorm:"column:base_symbol" json:"base_symbol"`                           // 基础代币符号
	BaseTokenAddress      string    `gorm:"column:base_token_address" json:"base_token_address"`             // 基础代币地址
	BaseTokenMintAddress  string    `gorm:"column:base_token_mint_address" json:"base_token_mint_address"`   // 基础代币铸造地址
	BaseTokenDecimals     int32     `gorm:"column:base_token_decimals" json:"base_token_decimals"`           // 基础代币小数位
	QuoteSymbol           string    `gorm:"column:quote_symbol" json:"quote_symbol"`                         // 报价代币符号
	QuoteTokenAddress     string    `gorm:"column:quote_token_address" json:"quote_token_address"`           // 报价代币地址
	QuoteTokenMintAddress string    `gorm:"column:quote_token_mint_address" json:"quote_token_mint_address"` // 报价代币铸造地址
	QuoteTokenDecimals    int32     `gorm:"column:quote_token_decimals" json:"quote_token_decimals"`         // 报价代币小数位
	DevAddress            string    `gorm:"column:dev_address" json:"dev_address"`                           // 开发者地址
	IsPump                bool      `gorm:"column:is_pump" json:"is_pump"`                                   // 是否为 Pump 市场
	InitBaseVault         float64   `gorm:"column:init_base_vault" json:"init_base_vault"`                   // 初始基础代币仓库
	InitQuoteVault        float64   `gorm:"column:init_quote_vault" json:"init_quote_vault"`                 // 初始报价代币仓库
	LogoURL               string    `gorm:"column:logo_url" json:"logo_url"`                                 // Logo 链接
	OfficialWebsite       string    `gorm:"column:official_website" json:"official_website"`                 // 官方网站
	Twitter               string    `gorm:"column:twitter" json:"twitter"`                                   // Twitter 链接
	Facebook              string    `gorm:"column:facebook" json:"facebook"`                                 // Facebook 链接
	Reddit                string    `gorm:"column:reddit" json:"reddit"`                                     // Reddit 链接
	GitHub                string    `gorm:"column:github" json:"github"`                                     // GitHub 链接
	Whitepaper            string    `gorm:"column:whitepaper" json:"whitepaper"`                             // 白皮书链接
	QuoteSupply           float64   `gorm:"column:quote_supply" json:"quote_supply"`                         // 报价代币供应量
	// market_real_time_data
	Address        string  `gorm:"column:address" json:"address"`
	QuotePrice     float64 `gorm:"column:quote_price" json:"quote_price"`
	BasePrice      float64 `gorm:"column:base_price" json:"base_price"`
	QuoteVault     float64 `gorm:"column:quote_vault" json:"quote_vault"`
	BaseVault      float64 `gorm:"column:base_vault" json:"base_vault"`
	Liquidity      float64 `gorm:"column:liquidity" json:"liquidity"`
	InitLiquidity  float64 `gorm:"column:init_liquidity" json:"init_liquidity"`
	MarketCap      float64 `gorm:"column:market_cap" json:"market_cap"`
	HolderCount    int64   `gorm:"column:holder_count" json:"holder_count"`
	Swaps          int64   `gorm:"column:swaps" json:"swaps"`
	Sells          int64   `gorm:"column:sells" json:"sells"`
	Buys           int64   `gorm:"column:buys" json:"buys"`
	Volume         float64 `gorm:"column:volume" json:"volume"`
	PriceChange1m  float64 `gorm:"column:price_change_1m" json:"price_change_1m"`
	PriceChange5m  float64 `gorm:"column:price_change_5m" json:"price_change_5m"`
	PriceChange30m float64 `gorm:"column:price_change_30m" json:"price_change_30m"`
	PriceChange1h  float64 `gorm:"column:price_change_1h" json:"price_change_1h"`
	PriceChange6h  float64 `gorm:"column:price_change_6h" json:"price_change_6h"`
	PriceChange24h float64 `gorm:"column:price_change_24h" json:"price_change_24h"`
	SellVolume1m   float64 `gorm:"column:sell_volume_1m" json:"sell_volume_1m"`
	SellVolume5m   float64 `gorm:"column:sell_volume_5m" json:"sell_volume_5m"`
	SellVolume30m  float64 `gorm:"column:sell_volume_30m" json:"sell_volume_30m"`
	SellVolume1h   float64 `gorm:"column:sell_volume_1h" json:"sell_volume_1h"`
	SellVolume6h   float64 `gorm:"column:sell_volume_6h" json:"sell_volume_6h"`
	SellVolume24h  float64 `gorm:"column:sell_volume_24h" json:"sell_volume_24h"`
	BuyVolume1m    float64 `gorm:"column:buy_volume_1m" json:"buy_volume_1m"`
	BuyVolume5m    float64 `gorm:"column:buy_volume_5m" json:"buy_volume_5m"`
	BuyVolume30m   float64 `gorm:"column:buy_volume_30m" json:"buy_volume_30m"`
	BuyVolume1h    float64 `gorm:"column:buy_volume_1h" json:"buy_volume_1h"`
	BuyVolume6h    float64 `gorm:"column:buy_volume_6h" json:"buy_volume_6h"`
	BuyVolume24h   float64 `gorm:"column:buy_volume_24h" json:"buy_volume_24h"`
	BuyCount1m     float64 `gorm:"column:buy_count_1m" json:"buy_count_1m"`
	BuyCount5m     float64 `gorm:"column:buy_count_5m" json:"buy_count_5m"`
	BuyCount30m    float64 `gorm:"column:buy_count_30m" json:"buy_count_30m"`
	BuyCount1h     float64 `gorm:"column:buy_count_1h" json:"buy_count_1h"`
	BuyCount6h     float64 `gorm:"column:buy_count_6h" json:"buy_count_6h"`
	BuyCount24h    float64 `gorm:"column:buy_count_24h" json:"buy_count_24h"`
	SellCount1m    float64 `gorm:"column:sell_count_1m" json:"sell_count_1m"`
	SellCount5m    float64 `gorm:"column:sell_count_5m" json:"sell_count_5m"`
	SellCount30m   float64 `gorm:"column:sell_count_30m" json:"sell_count_30m"`
	SellCount1h    float64 `gorm:"column:sell_count_1h" json:"sell_count_1h"`
	SellCount6h    float64 `gorm:"column:sell_count_6h" json:"sell_count_6h"`
	SellCount24h   float64 `gorm:"column:sell_count_24h" json:"sell_count_24h"`
}
