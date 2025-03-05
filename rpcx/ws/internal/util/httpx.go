package util

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
)

func WriteOkJsonConcurrent(wc *types.WebSocketClient, clientId string, action string, channel string, data interface{}) error {
	wc.Lock.Lock()
	defer wc.Lock.Unlock()

	reps := types.BaseResponse{
		Id:      clientId,
		Code:    0,
		Message: "success",
		Action:  action,
		Channel: channel,
		Data:    data,
	}
	marshal, err := json.Marshal(reps)
	if err != nil {
		return err
	}
	return wc.Conn.WriteMessage(websocket.TextMessage, marshal)
}

func WriteOkJson(conn *websocket.Conn, clientId string, action string, channel string, data interface{}) error {
	reps := types.BaseResponse{
		Id:      clientId,
		Code:    0,
		Message: "success",
		Action:  action,
		Channel: channel,
		Data:    data,
	}
	marshal, err := json.Marshal(reps)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, marshal)
}

func WriteFailedJson(conn *websocket.Conn, clientId string, action string, channel string, errMsg string) error {
	reps := types.BaseResponse{
		Id:      clientId,
		Code:    -1,
		Message: errMsg,
		Action:  action,
		Channel: channel,
	}
	marshal, err := json.Marshal(reps)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, marshal)
}
