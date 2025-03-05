package types

import "time"

var (
	DefaultChannelExpireTime  = 24 * time.Hour   //24小时 会断开连接 客户端需要重新连接一次
	MaxTokenKlineChannelCount = 2                //最多订阅2个TokenKline类型
	DefaultHeartbeatTime      = 15 * time.Minute //心跳检测断开 15分钟
)

var Pong = "pong"

var (
	SuccessCode = int32(1000)
)

const (
	ActionConnect     = "connect"
	ActionPing        = "ping"
	ActionSubscribe   = "subscribe"
	ActionUnSubscribe = "unsubscribe"
)

const (
	ChannelMarketKline        = "market_kline"
	ChannelMarketTxActivity   = "market_tx_activity"
	ChannelSmartTradeActivity = "smart_trade_activity"
)

type BaseRequest struct {
	Action  string `json:"action"`  //"subscribe/unsubscribe/ping"
	Id      string `json:"id"`      //clientID
	Channel string `json:"channel"` //market_kline/market_tx_activity/smart_trade_activity
	Data    struct {
		Chain    string `json:"chain"`    //chainName
		Address  string `json:"address"`  //market_address地址
		Interval string `json:"interval"` //订阅时长market_kline(1s 1s 1m 5m 15m 30m 1h 4h 6h 12h 1d)
	} `json:"data"`
}

type BaseResponse struct {
	Id      string      `json:"id"`
	Code    int64       `json:"code"`    //1000成功 其他错误
	Message string      `json:"message"` //success成功/其他错误内容
	Action  string      `json:"action"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
