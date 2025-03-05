package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
)

// GetMarketKline1m get market kline 1m from redis or db
func (c *DB) GetMarketKline1m(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline1m, error) {
	marketKline1mDB := dbx.Use(c.pgDB).MarketKline1m
	timeDuration := time.Minute
	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline1m, truncateTime.Unix())
	// get cache
	var marketList []*model.MarketKline1m
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1m.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline1mDB.WithContext(ctx).Where(
			marketKline1mDB.MarketAddress.Eq(marketAddress),
			marketKline1mDB.Timestamp.Lte(truncateTime),
		).Order(marketKline1mDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1m.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline1m.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline5m get market kline 5m from redis or db
func (c *DB) GetMarketKline5m(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline5m, error) {
	marketKline5mDB := dbx.Use(c.pgDB).MarketKline5m
	timeDuration := time.Minute * 5
	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline5m, truncateTime.Unix())

	var marketList []*model.MarketKline5m
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline5m.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline5mDB.WithContext(ctx).Where(
			marketKline5mDB.MarketAddress.Eq(marketAddress),
			marketKline5mDB.Timestamp.Lte(truncateTime),
		).Order(marketKline5mDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline5m.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline5m.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline15m get market kline 15m from redis or db
func (c *DB) GetMarketKline15m(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline15m, error) {
	marketKline15mDB := dbx.Use(c.pgDB).MarketKline15m
	timeDuration := time.Minute * 15

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline15m, truncateTime.Unix())

	var marketList []*model.MarketKline15m
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline15m.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline15mDB.WithContext(ctx).Where(
			marketKline15mDB.MarketAddress.Eq(marketAddress),
			marketKline15mDB.Timestamp.Lte(truncateTime),
		).Order(marketKline15mDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline15m.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline15m.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline15m.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline15m.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline30m get market kline 30m from redis or db
func (c *DB) GetMarketKline30m(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline30m, error) {
	marketKline30mDB := dbx.Use(c.pgDB).MarketKline30m
	timeDuration := time.Minute * 30
	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline30m, truncateTime.Unix())

	var marketList []*model.MarketKline30m
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline30m.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline30mDB.WithContext(ctx).Where(
			marketKline30mDB.MarketAddress.Eq(marketAddress),
			marketKline30mDB.Timestamp.Lte(truncateTime),
		).Order(marketKline30mDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline30m.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline30m.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline30m.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline30m.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline1h get market kline 1h from redis or db
func (c *DB) GetMarketKline1h(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline1h, error) {
	marketKline1hDB := dbx.Use(c.pgDB).MarketKline1h
	timeDuration := time.Hour

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline1h, truncateTime.Unix())

	var marketList []*model.MarketKline1h
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1h.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline1hDB.WithContext(ctx).Where(
			marketKline1hDB.MarketAddress.Eq(marketAddress),
			marketKline1hDB.Timestamp.Lte(truncateTime),
		).Order(marketKline1hDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1h.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1h.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1h.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline1h.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline4h get market kline 4h from redis or db
func (c *DB) GetMarketKline4h(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline4h, error) {
	marketKline4hDB := dbx.Use(c.pgDB).MarketKline4h
	timeDuration := time.Hour * 4

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline4h, truncateTime.Unix())

	var marketList []*model.MarketKline4h
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline4h.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline4hDB.WithContext(ctx).Where(
			marketKline4hDB.MarketAddress.Eq(marketAddress),
			marketKline4hDB.Timestamp.Lte(truncateTime),
		).Order(marketKline4hDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline4h.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline4h.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline4h.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline4h.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline6h get market kline 6h from redis or db
func (c *DB) GetMarketKline6h(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline6h, error) {
	marketKline6hDB := dbx.Use(c.pgDB).MarketKline6h
	timeDuration := time.Hour * 6

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline6h, truncateTime.Unix())

	var marketList []*model.MarketKline6h
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline6h.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline6hDB.WithContext(ctx).Where(
			marketKline6hDB.MarketAddress.Eq(marketAddress),
			marketKline6hDB.Timestamp.Lte(truncateTime),
		).Order(marketKline6hDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline6h.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline6h.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline6h.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline6h.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline12h get market kline 12h from redis or db
func (c *DB) GetMarketKline12h(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline12h, error) {
	marketKline12hDB := dbx.Use(c.pgDB).MarketKline12h
	timeDuration := time.Hour * 12

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline12h, truncateTime.Unix())

	var marketList []*model.MarketKline12h
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline12h.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline12hDB.WithContext(ctx).Where(
			marketKline12hDB.MarketAddress.Eq(marketAddress),
			marketKline12hDB.Timestamp.Lte(truncateTime),
		).Order(marketKline12hDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline12h.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline12h.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline12h.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline12h.redis.Get: %w", err)
	}
	return marketList, nil
}

// GetMarketKline1d get market kline 1d from redis or db
func (c *DB) GetMarketKline1d(ctx context.Context, marketAddress string, truncateTime time.Time) ([]*model.MarketKline1d, error) {
	marketKline1dDB := dbx.Use(c.pgDB).MarketKline1d
	timeDuration := time.Hour * 24

	cacheKey := fmt.Sprintf(MarketKlineInternalEndTimestampRedisKey, marketAddress, common.MarketKline1d, truncateTime.Unix())

	var marketList []*model.MarketKline1d
	cacheResult, err := c.redisClient.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		err = json.Unmarshal([]byte(cacheResult), &marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1d.json.Unmarshal: %w", err)
		}
	case errors.Is(err, redis.Nil):
		marketList, err = marketKline1dDB.WithContext(ctx).Where(
			marketKline1dDB.MarketAddress.Eq(marketAddress),
			marketKline1dDB.Timestamp.Lte(truncateTime),
		).Order(marketKline1dDB.Timestamp.Asc()).Limit(DefaultLimit).Find()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1d.db.Where: %w", err)
		}

		marketListBytes, err := json.Marshal(marketList)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1d.json.Marshal: %w", err)
		}

		_, err = c.redisClient.SetNX(ctx, cacheKey, marketListBytes, timeDuration).Result()
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline1d.redis.SetNX: %w", err)
		}
	default:
		return nil, fmt.Errorf("GetMarketKline1d.redis.Get: %w", err)
	}
	return marketList, nil
}
