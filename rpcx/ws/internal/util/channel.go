package util

import (
	"fmt"
	"strings"
)

func GetMarketTxActivityRedisKey(action string, channel string, chain string, marketAddress string) (string, error) {
	if action == "" {
		return "", fmt.Errorf("action is required")
	}
	if channel == "" {
		return "", fmt.Errorf("channel is required")
	}
	if chain == "" {
		return "", fmt.Errorf("chain is required")
	}
	if marketAddress == "" {
		return "", fmt.Errorf("marketAddress is required")
	}
	return strings.ToLower(fmt.Sprintf("ws:%v_%v_%v_%v", action, channel, chain, marketAddress)), nil
}

func GetSmartTradeActivityRedisKey(action string, channel string, userId string) (string, error) {
	if action == "" {
		return "", fmt.Errorf("action is required")
	}
	if channel == "" {
		return "", fmt.Errorf("channel is required")
	}
	if userId == "" {
		return "", fmt.Errorf("marketAddress is required")
	}
	return strings.ToLower(fmt.Sprintf("ws:%v_%v_%v", action, channel, userId)), nil
}

func GetMarketKlineRedisKey(action string, channel string, chain string, marketAddress string, interval string) (string, error) {
	if action == "" {
		return "", fmt.Errorf("action is required")
	}
	if channel == "" {
		return "", fmt.Errorf("channel is required")
	}
	if chain == "" {
		return "", fmt.Errorf("chain is required")
	}
	if marketAddress == "" {
		return "", fmt.Errorf("marketAddress is required")
	}
	if interval == "" {
		return "", fmt.Errorf("interval is required")
	}
	return strings.ToLower(fmt.Sprintf("ws:%v_%v_%v_%v_%v", action, channel, chain, marketAddress, interval)), nil
}

func GetClientIDKey(clientID string) (string, error) {
	if clientID == "" {
		return "", fmt.Errorf("clientID is required")
	}
	return strings.ToLower(fmt.Sprintf("ws:client_id:%v", clientID)), nil
}
