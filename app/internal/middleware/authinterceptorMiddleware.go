package middleware

import (
	"net/http"
)

type AuthInterceptorMiddleware struct {
	SecretKey string
}

func NewAuthInterceptorMiddleware(secretKey string) *AuthInterceptorMiddleware {
	return &AuthInterceptorMiddleware{
		SecretKey: secretKey,
	}
}

func (m *AuthInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 Authorization 字段
		//authHeader := r.Header.Get("Authorization")
		//if authHeader == "" {
		//	// 如果没有 Token，直接调用下一个处理函数
		//	next(w, r)
		//	return
		//}
		//
		//// 提取 Token，假设格式为 "Bearer <token>"
		//parts := strings.Split(authHeader, " ")
		//if len(parts) != 2 || parts[0] != "Bearer" {
		//	http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
		//	return
		//}
		//tokenString := parts[1]
		//
		//// 解析 Token
		//claims := &jwt.StandardClaims{}
		//token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		//	return []byte(m.SecretKey), nil
		//})
		//
		//if err != nil || !token.Valid {
		//	http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		//	return
		//}
		//
		//// 将用户信息添加到上下文中
		//ctx := context.WithValue(r.Context(), "user", claims.Subject)
		//
		//// 调用下一个处理函数，并传递更新后的上下文
		//next(w, r.WithContext(ctx))
	}
}
