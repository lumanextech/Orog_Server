package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibaba/tair-go/tair"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedisChannel struct {
	redisClient *tair.TairClusterClient
}

// NewRedisChannel 处理订阅topic逻辑
func NewRedisChannel(redisClient *tair.TairClusterClient) *RedisChannel {
	return &RedisChannel{redisClient: redisClient}
}

func (rtl *RedisChannel) AddMarketKlineOP(ctx context.Context, clientID string, req types.BaseRequest) error {
	if clientID == "" {
		return errors.New("clientID is empty")
	}
	if req.Action != types.ActionSubscribe {
		return errors.New("SubChannelTokenKline action err")
	}
	if req.Channel != types.ChannelMarketKline {
		return fmt.Errorf("SubChannelMarketKline channel err")
	}
	if !common.CheckIsSupportChain(req.Data.Chain) {
		return fmt.Errorf("SubChannelMarketKline chain err")
	}
	if !common.CheckIsSupportMarketKline(req.Data.Interval) {
		return fmt.Errorf("SubChannelMarketKline interval err")
	}
	if req.Data.Address == "" {
		return fmt.Errorf("SubChannelMarketKline address err")
	}

	channelKey, err := util.GetMarketKlineRedisKey(req.Action, req.Channel, req.Data.Chain, req.Data.Address, req.Data.Interval)
	if err != nil {
		return err
	}

	count, err := rtl.getClientIDCount(ctx, clientID, types.ChannelMarketKline)
	if err != nil {
		return err
	}
	if count > types.MaxTokenKlineChannelCount {
		return fmt.Errorf("Channel-MarketKline count too many, please unsubscribe it")
	}

	err = rtl.redisClient.HSet(ctx, channelKey, clientID, types.ChannelMarketKline).Err()
	if err != nil {
		return err
	}

	//topic channelKey expire in 24 hour times
	err = rtl.redisClient.Expire(ctx, channelKey, types.DefaultChannelExpireTime).Err()
	if err != nil {
		return err
	}

	err = rtl.addKeyToClientID(ctx, clientID, channelKey, req.Channel)
	if err != nil {
		return err
	}

	return nil
}

func (rtl *RedisChannel) AddMarketTxActivityOP(ctx context.Context, clientID string, req types.BaseRequest) error {
	if clientID == "" {
		return errors.New("clientID is empty")
	}
	if req.Action != types.ActionSubscribe {
		return errors.New("SubChannelMarketKline action err")
	}
	if req.Channel != types.ChannelMarketTxActivity {
		return fmt.Errorf("SubChannelMarketKline channel err")
	}
	if !common.CheckIsSupportChain(req.Data.Chain) {
		return fmt.Errorf("SubChannelMarketKline chain err")
	}
	if req.Data.Address == "" {
		return fmt.Errorf("SubChannelMarketKline address err")
	}

	key, err := util.GetMarketTxActivityRedisKey(req.Action, req.Channel, req.Data.Chain, req.Data.Address)
	if err != nil {
		return err
	}

	err = rtl.redisClient.HSet(ctx, key, clientID, req.Channel).Err()
	if err != nil {
		return err
	}

	//topic key expire in 24 hour times
	err = rtl.redisClient.Expire(ctx, key, types.DefaultChannelExpireTime).Err()
	if err != nil {
		return err
	}

	err = rtl.addKeyToClientID(ctx, clientID, key, types.ChannelMarketTxActivity)
	if err != nil {
		return err
	}

	return nil
}

func (rtl *RedisChannel) AddSmartTradeActivityOP(ctx context.Context, clientID string, req types.BaseRequest, userId string) error {
	if clientID == "" {
		return errors.New("clientID is empty")
	}
	if req.Action != types.ActionSubscribe {
		return errors.New("SubChannelMarketKline action err")
	}
	if req.Channel != types.ChannelSmartTradeActivity {
		return fmt.Errorf("SubChannelMarketKline channel err")
	}

	key, err := util.GetSmartTradeActivityRedisKey(req.Action, req.Channel, userId)
	if err != nil {
		return err
	}

	//用户订阅了 SmartTradeActivity Channel
	err = rtl.redisClient.HSet(ctx, key, clientID, req.Channel).Err()
	if err != nil {
		return err
	}

	//topic key expire in 24 hour times
	err = rtl.redisClient.Expire(ctx, key, types.DefaultChannelExpireTime).Err()
	if err != nil {
		return err
	}

	err = rtl.addKeyToClientID(ctx, clientID, key, types.ChannelSmartTradeActivity)
	if err != nil {
		return err
	}

	return nil
}

func (rtl *RedisChannel) DelMarketKlineByClientIDOP(ctx context.Context, clientID string) error {
	if clientID == "" {
		return errors.New("clientID is empty")
	}

	kValue, err := rtl.redisClient.HGetAll(ctx, clientID).Result()
	if err != nil {
		return err
	}

	//k: channelKey, v: channel
	for channelKey, channel := range kValue {
		if channel == types.ChannelMarketKline {
			err = rtl.redisClient.HDel(ctx, clientID, channelKey).Err()
			if err != nil {
				return err
			}

			err = rtl.redisClient.HDel(ctx, channelKey, clientID).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rtl *RedisChannel) AddClientID(ctx context.Context, clientID string) error {
	clientIDKey, err := util.GetClientIDKey(clientID)
	if err != nil {
		return err
	}

	err = rtl.redisClient.Expire(ctx, clientIDKey, types.DefaultHeartbeatTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rtl *RedisChannel) DelClientID(ctx context.Context, clientID string, forceDelChannelKey string) error {
	clientIDKey, err := util.GetClientIDKey(clientID)
	if err != nil {
		return err
	}

	mapKeys, err := rtl.redisClient.HGetAll(ctx, clientIDKey).Result()
	if err != nil {
		return err
	}

	success := true
	for channelKey, _ := range mapKeys {
		err = rtl.delClientIDFromKey(ctx, channelKey, clientID)
		if err != nil {
			success = false
			logx.Error("DelClientID delClientIDFromChannelKey: ", err)
		}
	}

	if success {
		err = rtl.redisClient.Del(ctx, clientIDKey).Err()
		if err != nil {
			return err
		}
	}

	if forceDelChannelKey != "" {
		err = rtl.delClientIDFromKey(ctx, forceDelChannelKey, clientID)
		if err != nil {
			logx.Error("DelClientID forceDelChannelKey: ", err)
		}
	}

	return nil
}

func (rtl *RedisChannel) GetClientIDsByKey(ctx context.Context, channelKey string) (map[string]string, error) {
	result, err := rtl.redisClient.HGetAll(ctx, channelKey).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rtl *RedisChannel) getClientIDCount(ctx context.Context, clientID string, channel string) (int, error) {
	clientIDKey, err := util.GetClientIDKey(clientID)
	if err != nil {
		return 0, err
	}

	kValue, err := rtl.redisClient.HGetAll(ctx, clientIDKey).Result()
	if err != nil {
		return 0, err
	}

	var count int
	for _, channelV := range kValue {
		if channelV == channel {
			count++
		}
	}
	return count, nil
}

func (rtl *RedisChannel) addKeyToClientID(ctx context.Context, clientID string, channelKey string, channel string) error {
	clientIDKey, err := util.GetClientIDKey(clientID)
	if err != nil {
		return err
	}

	err = rtl.redisClient.HSet(ctx, clientIDKey, channelKey, channel).Err()
	if err != nil {
		err = rtl.redisClient.HDel(ctx, channelKey, clientID).Err()
		if err != nil {
			return err
		}
		return err
	}

	err = rtl.redisClient.Expire(ctx, clientID, types.DefaultHeartbeatTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rtl *RedisChannel) delClientIDFromKey(ctx context.Context, channelKey string, clientID string) error {
	err := rtl.redisClient.HDel(ctx, channelKey, clientID).Err()
	if err != nil {
		return err
	}
	return nil
}
