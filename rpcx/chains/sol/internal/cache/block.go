package cache

import (
	"context"
	"fmt"
	"strconv"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// AddBlockHeightBasePrice
func (c *DB) AddBlockHeightBasePrice(ctx context.Context, height int64, basePrice float64) error {
	return c.redisClient.ZAdd(ctx, BlockHeightBasePriceRedisZKey, redis.Z{
		Score:  float64(height),
		Member: basePrice,
	}).Err()
}

// GetBlockHeightBasePrice
func (c *DB) GetBlockHeightBasePrice(ctx context.Context, height int64) (float64, error) {
	prices, err := c.redisClient.ZRevRange(ctx, BlockHeightBasePriceRedisZKey, 1, 10).Result()
	if err != nil {
		return 0, err
	}

	if len(prices) <= 0 {
		//get from solscan

		return 210, nil
	}

	var priceAll float64
	for _, price := range prices {
		priceF, err := decimal.NewFromString(price)
		if err != nil {
			return 0, fmt.Errorf("helpGetBasePriceByBlock priceF not string")
		}
		priceAll = priceAll + priceF.InexactFloat64()
	}
	return priceAll / float64(len(prices)), nil
}

// AddBlockedHeight
func (c *DB) AddBlockedHeight(ctx context.Context, height int64) error {
	_, err := c.redisClient.ZAdd(ctx, BlockHeightRedisZKey, redis.Z{
		Score:  float64(height),
		Member: height,
	}).Result()
	return err
}

// GetBlockedHeight
func (c *DB) GetBlockedHeight(ctx context.Context) int64 {
	heightList, err := c.redisClient.ZRevRange(ctx, BlockHeightRedisZKey, 0, 1).Result()
	if err != nil {
		return 0
	}

	if len(heightList) <= 0 {
		latestHeight, err := c.solClient.GetSlot(ctx, solrpc.CommitmentConfirmed)
		if err != nil {
			c.Errorf("GetBlockedHeight GetSlot: %v", err)
			return 0
		}
		return int64(latestHeight)
	}

	height := heightList[0]

	heightInt, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		return 0
	}

	return heightInt
}
