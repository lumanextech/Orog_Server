package sol_scan

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetTokenMarketInfo 获取代币市场信息
// @api https://pro-api.solscan.io/pro-api-docs/v2.0/reference/v2-token-market-info
func (c *Client) GetTokenMarketInfo(ctx context.Context, address string) (*TokenMarketInfo, error) {
	resp, err := c.get("/v2.0/token/market/info", map[string]string{
		"address": address,
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetTokenMarketInfo.error: %s %v", resp.Status(), resp.String())
	}

	result := new(TokenMarketInfo)
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("GetTokenMarketInfo.error: %v", err)
	}

	c.Debugf("GetTokenMarketInfo.result: %v", result)

	return result, nil
}

// GetTokenPrice 获取代币价格
// @api https://pro-api.solscan.io/pro-api-docs/v2.0/reference/v2-token-price
func (c *Client) GetTokenPrice(ctx context.Context, address string) (*TokenPriceResp, error) {
	resp, err := c.get("/v2.0/token/price", map[string]string{
		"address": address,
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetTokenPrice.error: %s %v", resp.Status(), resp.String())
	}

	result := new(TokenPriceResp)
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("GetTokenPrice.error: %v", err)
	}

	return result, nil
}

// GetTokenMetadata 获取代币元数据
// @api https://pro-api.solscan.io/pro-api-docs/v2.0/reference/v2-token-meta
func (c *Client) GetTokenMetadata(ctx context.Context, address string) (*TokenMetadataResp, error) {
	resp, err := c.get("/v2.0/token/metadata", map[string]string{
		"address": address,
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetTokenMetadata.error: %s %v", resp.Status(), resp.String())
	}

	result := new(TokenMetadataResp)
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("GetTokenMetadata.error: %v", err)
	}

	return result, nil
}
