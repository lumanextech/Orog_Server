package account_v1

import (
	"context"
	"errors"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/account/account"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const KnownMessage = "Hello OROG" // 假设您有一个已知的消息

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VerifySignature verifies the signature using the provided public key and message.
func VerifySignature(signature []byte, publicKey string, message []byte) error {
	// Convert []byte to solana.Signature
	var signatureArray solana.Signature
	copy(signatureArray[:], signature)

	// 假设您有一个已知的公钥
	pubKey, err := solana.PublicKeyFromBase58(publicKey)
	if err != nil {
		return errors.New("无效的公钥格式")
	}

	// 验证签名
	if !pubKey.Verify(message, signatureArray) {
		return errors.New("签名验证失败")
	}

	return nil
}

func getJwtToken(secretKey string, iat, seconds int64, payload string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func CreateToken(secretKey, payload string, duration int64) (string, error) {
	currentTime := time.Now().Unix()
	return getJwtToken(secretKey, currentTime, duration, payload)
}

func (l *GetTokenLogic) GetToken(req *types.GetTokenRequest) (resp *types.GetTokenResponse, err error) {
	// 使用常量 KnownMessage
	message := []byte(KnownMessage)
	// 调用 VerifySignature 函数
	err = VerifySignature(req.Signature, req.PublicKey, message)
	if err != nil {
		return nil, err
	}
	//查询该地址是否已经创建过account，没有则创建account
	_, err = l.svcCtx.AccountClient.CheckAccount(l.ctx, &account.CheckAccountRequest{
		Address: req.PublicKey,
	})
	// 创建token
	tokenString, err := CreateToken("orog2025", req.PublicKey, 604800)
	if err != nil {
		return nil, err
	}
	// 创建响应
	resp = &types.GetTokenResponse{
		Token: tokenString,
	}
	return resp, nil
}
