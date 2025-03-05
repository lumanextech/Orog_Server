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

func (c *DB) GetMarketList(ctx context.Context, address string) (*model.Market, error) {
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
