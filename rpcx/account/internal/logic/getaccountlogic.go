package logic

import (
	"context"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountLogic {
	return &GetAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccountLogic) GetAccount(in *account.AccountRequest) (*account.AccountResponse, error) {
	pk := "4iE8n12rpppMo23C6xLXb85m2DRWQ3X9tgXaKs6jDfML4KqZEVVHLUQP6DCQGCMFDdX8TzvGDRxmXgYc2Z2Mt5CW" //okex account
	cryptEn := crypto.
		FromString(pk).
		SetKey(l.svcCtx.Config.PkEncode.Key).
		SetIv(l.svcCtx.Config.PkEncode.Iv).
		Aes().
		CBC().
		PKCS7Padding().
		Encrypt().
		ToBase64String()

	l.Info("cryptEn: ", cryptEn)

	return &account.AccountResponse{
		Id:     0,
		Email:  "18815001004@163.com",
		Mobile: "18815001004",
		Name:   "JZ",
		Pk:     cryptEn,
	}, nil
}
