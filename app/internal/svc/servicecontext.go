package svc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/simance-ai/smdx/app/internal/config"
	"github.com/simance-ai/smdx/rpcx/account/accountclient"
	"github.com/simance-ai/smdx/rpcx/chains/sol/solclient"
	"github.com/simance-ai/smdx/rpcx/order/orderclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	SolClient solclient.Sol

	AccountClient accountclient.Account

	OrderClient orderclient.Order
}

func (c ServiceContext) AuthInterceptor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 Authorization 字段,未传入直接跳过
		token := r.Header.Get("Authorization")
		if token == "" {
			next(w, r)
			return
		}

		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:] // 去掉 "Bearer " 前缀，获取实际的 JWT
		}
		if len(token) < 128 {
			next.ServeHTTP(w, r)
			return
		}

		address, err := JwtTokenHandler(c.Config.Auth.AccessSecret, token)
		if err != nil {
			fmt.Printf("middleware wrong: %v\n", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "payload", address)
		// 将上下文传递给下一个处理器
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// JwtTokenHandler 模拟解析 Token 的方法
func JwtTokenHandler(sk, tokenString string) (string, error) {
	// 替换为实际的 secret key
	secretKey := []byte(sk)

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	// 从 claims 中提取 "address" 字段
	address, ok := claims["payload"].(string)
	if !ok {
		return "", fmt.Errorf("address not found in token")
	}

	return address, nil
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		SolClient:     solclient.NewSol(zrpc.MustNewClient(c.SolClientConf)),
		AccountClient: accountclient.NewAccount(zrpc.MustNewClient(c.AccountClientConf)),
		OrderClient:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderClientConf)),
	}
}
