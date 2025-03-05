package server

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"github.com/simance-ai/smdx/rpcx/ws/internal/config"
	"github.com/simance-ai/smdx/rpcx/ws/internal/svc"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type WebSocketHandler struct {
	upgrade       websocket.Upgrader
	websocketCtx  *svc.WebsocketContext
	conf          config.Config
	heartBeatPool *ants.Pool
}

func NewWebSocketHandler(conf config.Config, websocketContext *svc.WebsocketContext) *WebSocketHandler {

	pool, err := ants.NewPool(100000)
	if err != nil {
		log.Fatalln("NewWebSocketHandler ants.NewPool err:", err)
	}

	return &WebSocketHandler{
		upgrade: websocket.Upgrader{
			ReadBufferSize:    conf.WebSocket.ReadBufferSize,
			WriteBufferSize:   conf.WebSocket.WriteBufferSize,
			HandshakeTimeout:  time.Duration(conf.WebSocket.WebsocketHandshakeTimeout),
			EnableCompression: conf.WebSocket.EnableCompression == 1,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		websocketCtx:  websocketContext,
		conf:          conf,
		heartBeatPool: pool,
	}
}

func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error("ServeHTTP panic: %v\n", err)
			debug.PrintStack()
		}
	}()

	conn, err := wsh.upgrade.Upgrade(w, r, nil)
	if err != nil {
		logx.Error("ServeHTTP when upgrading connection to websocket", err)
		return
	}

	var userId int64
	values := r.Header.Values("fomo-go")
	if len(values) > 0 {
		token := values[0]
		if token != "" {
			userInfo, err := util.ParseAccountAuthClaims(token, wsh.conf.App.Secret)
			if err != nil {
				logx.Error("ServeHTTP ParseAccountAuthClaims: ", err)
				_ = util.WriteFailedJson(conn, "", types.ActionConnect, "", err.Error())
				return
			}
			userId = userInfo.UserId
		}
	}

	clientID := wsh.HandleClientConn(r.Context(), strconv.FormatInt(userId, 10), conn)
	if clientID == "" {
		logx.Error("ServeHTTP clientID is empty!")
		return
	}

	err = wsh.heartBeatPool.Submit(func() {
		wsh.websocketCtx.ListenHeartbeatCheck(context.Background(), clientID)
	})
	if err != nil {
		logx.Error("ServeHTTP heartBeatPool.Submit err:", err)
		return
	}

	_ = util.WriteOkJson(conn, clientID, types.ActionConnect, "", nil)

	defer func() {
		err := conn.Close()
		if err != nil {
			logx.Error("ServeHTTP Close: ", err)
		}

		wsh.websocketCtx.OnDisconnect(r.Context(), clientID)
	}()

	//read listen
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logx.Debug("ServeHTTP ReadMessage: ", err)
			wsh.HandleClientDisconnect(r.Context(), clientID)
			break
		} else {
			err = wsh.HandleClientMsg(r.Context(), clientID, msgType, msg)
			if err != nil {
				logx.Error("ServeHTTP HandleClientMsg err: ", err)
				_ = util.WriteFailedJson(conn, clientID, "", "", err.Error())
			}
		}
	}
}

func (wsh *WebSocketHandler) HandleClientConn(ctx context.Context, id string, conn *websocket.Conn) string {
	clientID := uuid.New().String()

	wsh.websocketCtx.OnConnect(ctx, clientID, id, conn)

	return clientID
}

func (wsh *WebSocketHandler) HandleClientDisconnect(ctx context.Context, clientID string) {
	// 使用 Load 和 Delete 方法，不需要额外的锁定操作
	wsh.websocketCtx.OnDisconnect(ctx, clientID)
}

func (wsh *WebSocketHandler) HandleClientMsg(ctx context.Context, clientID string, messageType int, message []byte) error {
	return wsh.websocketCtx.OnHandleMessage(ctx, clientID, messageType, message)
}
