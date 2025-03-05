package types

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"sync"
)

/*
WebSocketClient
底层连接
*/
type WebSocketClient struct {
	ID            string
	Conn          *websocket.Conn
	LastHeartBeat int64
	Lock          sync.Mutex
	BindUid       string
	JoinGroup     []string
}

type AccountAuthClaims struct {
	ChainId int64  `json:"chain_id"`
	UserId  int64  `json:"user_id"`
	Address string `json:"address"`
	Env     string `json:"env"`
	jwt.StandardClaims
}
