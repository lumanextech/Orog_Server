package svc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/alibaba/tair-go/tair"
	"github.com/redis/go-redis/v9"

	"github.com/simance-ai/smdx/rpcx/ws/internal/config"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"

	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

var fireMessagePool, _ = ants.NewPool(50000, ants.WithPreAlloc(true)) //理论5万用户同时推送

// WebsocketContext 可以用channel替代sync.Map
type WebsocketContext struct {
	clients       sync.Map //clientID-Conn
	users         sync.Map //clientID-UserID
	heartbeatTime time.Duration

	redisChannel *RedisChannel
	ctx          context.Context
}

func NewWebsocketContext(c config.Config) *WebsocketContext {
	if len(c.Cache.Redis) <= 0 {
		log.Fatalf("redis config is empty")
	}

	redisClusterConf := c.Cache.Redis[0]
	var tlsConfig *tls.Config // Explicitly define the type
	if redisClusterConf.Tls {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: redisClusterConf.Tls,
		}
	}
	redisClient := tair.NewTairClusterClient(&tair.TairClusterOptions{
		ClusterOptions: &redis.ClusterOptions{
			Addrs:     []string{redisClusterConf.Host},
			Password:  redisClusterConf.Pass,
			TLSConfig: tlsConfig,
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	successStatusErr := redisClient.Ping(ctx).Err()
	if successStatusErr != nil {
		log.Fatalf("redis connect error: %v", successStatusErr)
	}

	return &WebsocketContext{
		users:         sync.Map{},
		clients:       sync.Map{},
		redisChannel:  NewRedisChannel(redisClient),
		ctx:           context.Background(),
		heartbeatTime: types.DefaultHeartbeatTime,
	}
}

/*
ListenHeartbeatCheck
  - @description: 客户端心跳检测，超时即断开连接（主要是为了降低服务端承载压力）
  - @param {string} clientID
  - @return {*}
*/
func (wc *WebsocketContext) ListenHeartbeatCheck(ctx context.Context, clientID string) {
	for {
		time.Sleep(5 * time.Second)

		count := 0
		wc.clients.Range(func(key, value interface{}) bool {
			client, ok := value.(*types.WebSocketClient)
			if ok {
				logx.Infow("clients",
					logx.Field("key", key),
					logx.Field("clientId", client.ID),
					logx.Field("heartBeat", time.Unix(client.LastHeartBeat, 0)),
					logx.Field("fireMessagePoolLen", fireMessagePool.Running()),
				)
				count++
			}
			return true
		})

		logx.Debug("clientLen: ", count)

		clientInterface, exists := wc.clients.Load(clientID)
		if !exists {
			return
		}

		client, ok := clientInterface.(*types.WebSocketClient)
		if !ok {
			logx.Error("clientInterface .*types.WebSocketClient failed")
			return
		}

		if time.Now().Sub(time.Unix(client.LastHeartBeat, 0)).Milliseconds() >=
			wc.heartbeatTime.Milliseconds() {
			logx.Debugf("Client: %v %v", clientID, "heartbeat timeout")
			wc.OnDisconnect(ctx, clientID)
			break
		}
	}
}

/*
OnConnect
  - @description: 客户端连接
  - @param {string} clientID
  - @param {types.WebSocketClient} conn
  - @return {*}
*/
func (wc *WebsocketContext) OnConnect(ctx context.Context, clientID string, uuid string, conn *websocket.Conn) {
	if uuid == "0" {
		uuid = ""
	}

	client := &types.WebSocketClient{
		ID:            clientID,
		Conn:          conn,
		BindUid:       uuid,
		LastHeartBeat: time.Now().Unix(),
	}

	wc.clients.Store(clientID, client)
	wc.users.Store(uuid, clientID)

	err := wc.redisChannel.AddClientID(ctx, clientID)
	if err != nil {
		logx.Error("OnConnect AddClientID: ", err)
	}
}

/*
OnDisconnect
  - @description: 客户端断开
  - @param {string} clientID
  - @return {*}
*/
func (wc *WebsocketContext) OnDisconnect(ctx context.Context, clientID string) {
	value, ok := wc.clients.Load(clientID)
	if ok {
		client, ok := value.(*types.WebSocketClient)
		if ok {
			client.Conn.Close()
		} else {
			logx.Error("OnDisconnect to types.WebSocketClient: ", clientID)
		}

		wc.clients.Delete(clientID)
		wc.users.Delete(client.BindUid)
	}

	err := wc.redisChannel.DelClientID(ctx, clientID, "")
	if err != nil {
		logx.Error("OnDisconnect DelClientID: ", err)
	}
}

/*
OnHandleMessage
  - @description: 客户端心跳保活
  - @param {string} clientID
  - @return {*}
*/
func (wc *WebsocketContext) OnHandleMessage(ctx context.Context, clientID string, messageType int, message []byte) error {

	value, ok := wc.clients.Load(clientID)
	if !ok {
		wc.OnDisconnect(ctx, clientID)
		return errors.New("clientID not exist")
	}

	client, ok := value.(*types.WebSocketClient)
	if !ok {
		wc.OnDisconnect(ctx, clientID)
		return errors.New("clientID not match")
	}

	if client.Conn == nil {
		wc.OnDisconnect(ctx, clientID)
		return errors.New("clientID conn not exist")
	}

	switch messageType {
	case websocket.CloseMessage:
		logx.Debug("CloseMessage..")
		wc.OnDisconnect(ctx, clientID)
	case websocket.PingMessage:
		wc.OnConnect(ctx, clientID, client.BindUid, client.Conn)
		if err := client.Conn.WriteMessage(websocket.PongMessage, []byte(types.Pong)); err != nil {
			wc.OnDisconnect(ctx, clientID)
			return err
		}
		return nil
	case websocket.TextMessage:
		//请求
		var subRequest types.BaseRequest
		err := json.Unmarshal(message, &subRequest)
		if err != nil {
			return err
		}

		switch subRequest.Action {
		case types.ActionPing:
			//保活
			wc.OnConnect(ctx, clientID, client.BindUid, client.Conn)
			err = util.WriteOkJson(client.Conn, clientID, types.Pong, "", nil)
			if err != nil {
				wc.OnDisconnect(ctx, clientID)
				return err
			}
		case types.ActionSubscribe:
			channel := subRequest.Channel
			switch channel {
			case types.ChannelMarketKline:
				err = wc.redisChannel.AddMarketKlineOP(ctx, clientID, subRequest)
				if err != nil {
					return err
				}
				err = util.WriteOkJson(client.Conn, clientID, subRequest.Action, channel, nil)
				if err != nil {
					return err
				}
			case types.ChannelMarketTxActivity:
				err = wc.redisChannel.AddMarketTxActivityOP(ctx, clientID, subRequest)
				if err != nil {
					return err
				}
				err = util.WriteOkJson(client.Conn, clientID, subRequest.Action, channel, nil)
				if err != nil {
					return err
				}
			case types.ChannelSmartTradeActivity:
				if client.BindUid == "" {
					return errors.New("smart_trade_activity subscribe need dx token")
				}

				logx.Info("smart_trade_activity subscribe need dx token", client.BindUid)

				err = wc.redisChannel.AddSmartTradeActivityOP(ctx, clientID, subRequest, client.BindUid)
				if err != nil {
					return err
				}
				err = util.WriteOkJson(client.Conn, clientID, subRequest.Action, channel, nil)
				if err != nil {
					return err
				}
			default:
				return errors.New("channelType not support now")
			}
			return nil
		case types.ActionUnSubscribe:
			err = wc.redisChannel.DelMarketKlineByClientIDOP(ctx, clientID)
			if err != nil {
				return err
			}
			err = util.WriteOkJson(client.Conn, clientID, subRequest.Action, types.ChannelMarketKline, nil)
			if err != nil {
				return err
			}
		default:
			return errors.New("action not support")
		}
	default:
		return errors.New("not support message type")
	}
	return nil
}

/*
FireJsonMessage
  - @description: 发送消息
  - @param {string} clientID
  - @return {*}
*/
func (wc *WebsocketContext) FireJsonMessage(ctx context.Context, channelKey string, channel string, data interface{}) error {

	clientIds, err := wc.redisChannel.GetClientIDsByKey(ctx, channelKey)
	if err != nil {
		return fmt.Errorf("redisChannel.GetClientIDsByKey: %w", err)
	}

	//logx.Debugf("FireJsonMessage clientIds: %v %v", clientIds, channelKey)

	for clientId, _ := range clientIds {
		value, ok := wc.clients.Load(clientId)
		if !ok {
			logx.Error("FireJsonMessage clientID not exist: ", clientId)
			wc.OnDisconnect(ctx, clientId)
			continue
		}

		wsClient, ok := value.(*types.WebSocketClient)
		if !ok {
			logx.Error("FireJsonMessage clientID not *types.WebSocketClient: ", clientId)
			wc.OnDisconnect(ctx, clientId)
			continue
		}

		err = fireMessagePool.Submit(func() {

			err = util.WriteOkJsonConcurrent(wsClient, clientId, types.ActionSubscribe, channel, data)
			if err != nil {
				logx.Errorf("FireJsonMessage WriteMessage err: %v", err)
				wc.OnDisconnect(ctx, clientId)
			}

		})

		if err != nil {
			return fmt.Errorf("fireMessagePool.Submit: %w", err)
		}
	}

	return nil
}
