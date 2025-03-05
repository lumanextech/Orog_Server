package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
	"gorm.io/gorm"
)

// GetMarketInfo get market info from redis or db
// @return marketInfo: market info
// @return err: error ErrMarketNotFound if market info not found
func (c *DB) GetMarketInfo(ctx context.Context, address string) (*model.Market, error) {
	marketDB := dbx.Use(c.pgDB).Market
	var marketInfo *model.Market
	marketInfoResult, err := c.redisClient.Get(ctx, fmt.Sprintf(MarketAddressInfoRedisKey, address)).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(marketInfoResult), &marketInfo)
		if err != nil {
			return nil, ErrMarketNotFound
		}

		return marketInfo, nil
	case errors.Is(err, redis.Nil):
		//from db
		marketInfo, err = marketDB.WithContext(ctx).ReadDB().Where(marketDB.Address.Eq(address)).First()
		switch {
		case err == nil:
			//set to redis
			marketInfoBytes, err := json.Marshal(marketInfo)
			if err != nil {
				return nil, err
			}

			//nx: not exist and than set
			err = c.redisClient.SetNX(ctx, fmt.Sprintf(MarketAddressInfoRedisKey, address), marketInfoBytes, DefaultExpireDuration).Err()
			if err != nil {
				c.Errorf("GetMarketInfo.SetNX MarketAddressInfoRedisKey error: %v", err)
			}
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrMarketNotFound
		default:
			return nil, err
		}

		return marketInfo, nil
	default:
		return nil, err
	}
}

// SetMarketInfo set market info to redis
// @return err: error
func (c *DB) SetMarketInfo(ctx context.Context, marketInfo *model.Market, updateDB bool) error {
	marketDB := dbx.Use(c.pgDB).Market

	if marketInfo == nil {
		return fmt.Errorf("SetMarketInfo.marketInfo is nil")
	}

	marketInfoBytes, err := json.Marshal(marketInfo)
	if err != nil {
		return err
	}

	err = c.redisClient.SetNX(ctx, fmt.Sprintf(MarketAddressInfoRedisKey, marketInfo.Address), marketInfoBytes, DefaultExpireDuration).Err()
	if err != nil {
		return fmt.Errorf("SetMarketInfo.SetNX MarketAddressInfoRedisKey error: %v", err)
	}

	//update db
	if updateDB {
		_, err = marketDB.WithContext(ctx).Where(marketDB.Address.Eq(marketInfo.Address)).Updates(marketInfo)
		if err != nil {
			return fmt.Errorf("SetMarketInfo.Updates error: %v", err)
		}
	}

	return nil
}
