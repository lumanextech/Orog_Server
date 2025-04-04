// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
// Source: account.proto

package server

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/logic"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"
)

type AccountServer struct {
	svcCtx *svc.ServiceContext
	account.UnimplementedAccountServer
}

func NewAccountServer(svcCtx *svc.ServiceContext) *AccountServer {
	return &AccountServer{
		svcCtx: svcCtx,
	}
}

func (s *AccountServer) Ping(ctx context.Context, in *account.Request) (*account.Response, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}

func (s *AccountServer) Rebate(ctx context.Context, in *account.RebateRequest) (*account.Response, error) {
	l := logic.NewRebateLogic(ctx, s.svcCtx)
	return l.Rebate(in)
}

// 获取Account信息(峰在用)
func (s *AccountServer) GetAccount(ctx context.Context, in *account.AccountRequest) (*account.AccountResponse, error) {
	l := logic.NewGetAccountLogic(ctx, s.svcCtx)
	return l.GetAccount(in)
}

func (s *AccountServer) GetAccountInfo(ctx context.Context, in *account.GetAccountInfoRequest) (*account.GetAccountInfoResponse, error) {
	l := logic.NewGetAccountInfoLogic(ctx, s.svcCtx)
	return l.GetAccountInfo(in)
}

// 检查有没有Account没有创建
func (s *AccountServer) CheckAccount(ctx context.Context, in *account.CheckAccountRequest) (*account.CheckAccountResponse, error) {
	l := logic.NewCheckAccountLogic(ctx, s.svcCtx)
	return l.CheckAccount(in)
}

// 关注token
func (s *AccountServer) AddTokenFollow(ctx context.Context, in *account.AddTokenFollowRequest) (*account.AddTokenFollowResponse, error) {
	l := logic.NewAddTokenFollowLogic(ctx, s.svcCtx)
	return l.AddTokenFollow(in)
}

// 取消关注token
func (s *AccountServer) RemoveTokenFollow(ctx context.Context, in *account.AddTokenFollowRequest) (*account.AddTokenFollowResponse, error) {
	l := logic.NewRemoveTokenFollowLogic(ctx, s.svcCtx)
	return l.RemoveTokenFollow(in)
}

// 更新用户基础信息
func (s *AccountServer) UpdateUserBaseInfo(ctx context.Context, in *account.UpdateUserBaseInfoRequest) (*account.UpdateUserBaseInfoResponse, error) {
	l := logic.NewUpdateUserBaseInfoLogic(ctx, s.svcCtx)
	return l.UpdateUserBaseInfo(in)
}

// 用户是否收藏代币
func (s *AccountServer) IsTokenFollowed(ctx context.Context, in *account.IsTokenFollowedRequest) (*account.IsTokenFollowedResponse, error) {
	l := logic.NewIsTokenFollowedLogic(ctx, s.svcCtx)
	return l.IsTokenFollowed(in)
}

// 获取用户关注的代币
func (s *AccountServer) GetFollowTokenList(ctx context.Context, in *account.GetFollowTokenListRequest) (*account.GetFollowTokenListResponse, error) {
	l := logic.NewGetFollowTokenListLogic(ctx, s.svcCtx)
	return l.GetFollowTokenList(in)
}
